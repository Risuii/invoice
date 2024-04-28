package Invoices

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Risuii/invoice/src/entity"
	"github.com/Risuii/invoice/src/v1/contract"
	"github.com/google/uuid"
	"github.com/mariomac/gostream/stream"

	frsAtomic "github.com/Risuii/frs-lib/atomic"
	frsUtils "github.com/Risuii/frs-lib/utils"
	errorss "github.com/Risuii/invoice/src/errors"
)

type UUIDGenerator interface {
	New() uuid.UUID
}

type DefaultUUIDGenerator struct{}

func (g DefaultUUIDGenerator) New() uuid.UUID {
	return uuid.New()
}

type Invoiceservice struct {
	InvoicesRepo  InvoicesRepository
	CustomerRepo  CustomerRepository
	ItemRepo      ItemRepository
	AtomicSession frsAtomic.AtomicSessionProvider
	UUIDGen       UUIDGenerator
}

func InitInvoiceservice(InvoicesRepo InvoicesRepository, customerRepo CustomerRepository, item ItemRepository, aSession frsAtomic.AtomicSessionProvider, uuid UUIDGenerator) *Invoiceservice {
	return &Invoiceservice{
		InvoicesRepo:  InvoicesRepo,
		CustomerRepo:  customerRepo,
		ItemRepo:      item,
		AtomicSession: aSession,
		UUIDGen:       uuid,
	}
}

func incrementInvoiceID(lastInvoiceID string) string {
	var lastID int
	_, err := fmt.Sscanf(lastInvoiceID, "%d", &lastID)
	if err != nil {
		log.Fatal(err)
	}

	newID := lastID + 1

	return fmt.Sprintf("%04d", newID)
}

func findItemID(arr1, arr2 []uuid.UUID) []uuid.UUID {
	// Membuat map untuk menyimpan frekuensi UUID pada array pertama
	uuidCount := make(map[uuid.UUID]int)

	// Menghitung frekuensi UUID pada array pertama
	for _, id := range arr1 {
		uuidCount[id]++
	}

	// Mengurangkan frekuensi UUID pada array kedua
	for _, id := range arr2 {
		uuidCount[id]--
	}

	// Membuat slice untuk menyimpan UUID yang tidak ada di kedua array
	missingUUIDs := make([]uuid.UUID, 0)

	// Menemukan UUID yang frekuensinya tidak nol (artinya hanya ada di salah satu array)
	for id, count := range uuidCount {
		if count != 0 {
			missingUUIDs = append(missingUUIDs, id)
		}
	}

	return missingUUIDs
}

func (ts *Invoiceservice) Create(ctx context.Context, request contract.InvoiceRequest) (contract.InvcResponse, error) {
	var res contract.InvcResponse
	var newInvoiceID string

	invoiceID, _ := ts.InvoicesRepo.GetLatestInvoiceID(ctx)

	if invoiceID == "" {
		newInvoiceID = incrementInvoiceID("0000")
	} else {
		newInvoiceID = incrementInvoiceID(invoiceID)
	}

	layout := "02-01-2006"
	uuidForCustomer := ts.UUIDGen.New()
	newIssueDate, _ := time.Parse(layout, request.IssueDate)

	newDueDate, _ := time.Parse(layout, request.DueDate)

	err := frsAtomic.Atomic(ctx, ts.AtomicSession, func(ctx context.Context) error {

		insertDataCustomer := entity.Customer{
			CustomerData: entity.CustomerData{
				CustomerID: uuidForCustomer,
				Name:       request.CustomerRequest.CustomerName,
				Address:    request.CustomerRequest.Address,
			},
		}

		insertDataInvoice := entity.Invoices{
			InvoicesData: entity.InvoicesData{
				InvoiceID:  newInvoiceID,
				IssueDate:  newIssueDate,
				Subject:    request.Subject,
				TotalItems: len(request.ItemRequest),
				CustomerID: insertDataCustomer.CustomerID,
				DueDate:    newDueDate,
				SubTotal:   request.SubTotal,
				Tax:        request.Tax,
				GrandTotal: request.GrandTotal,
				Status:     "Unpaid",
			},
		}

		items := stream.Map(stream.OfSlice(request.ItemRequest), func(i contract.ItemRequest) *entity.Item {
			return &entity.Item{
				ItemData: entity.ItemData{
					InvoiceID: insertDataInvoice.InvoiceID,
					ItemID:    ts.UUIDGen.New(),
					Name:      i.Name,
					Type:      i.Type,
					Quantity:  i.Quantity,
					UnitPrice: i.UnitPrice,
					Amount:    i.Amount,
				},
			}
		}).ToSlice()

		err := ts.CustomerRepo.Create(ctx, &insertDataCustomer)
		if err != nil {
			log.Println("create customer err: ", err)
			return err
		}

		invoiceData, err := ts.InvoicesRepo.Create(ctx, &insertDataInvoice)
		if err != nil {
			log.Println("create invoice err: ", err)
			return err
		}

		err = ts.ItemRepo.Create(ctx, items)
		if err != nil {
			log.Println("create item err: ", err)
			return err
		}

		res = contract.InvcResponse{
			InvoiceID: invoiceData.InvoiceID,
		}

		return nil
	})

	if err != nil {
		log.Println("err: ", err)
		return res, err
	}

	return res, nil
}

func (ts *Invoiceservice) GetList(ctx context.Context, params contract.GetListParam) (contract.ListInvoiceResponse, error) {
	var response contract.ListInvoiceResponse

	Invoices, err := ts.InvoicesRepo.GetList(ctx, params)
	if err != nil {
		log.Println("getList err: ", err)
		return response, err
	}

	count, err := ts.InvoicesRepo.GetInvoicesCount(ctx, params)
	if err != nil {
		log.Println("InvoicesCount err: ", err)
		return response, err
	}

	pagination := frsUtils.GetPaginationData(params.Page, params.Limit, int(count))

	responseInvoicesList := stream.Map(stream.OfSlice(Invoices), func(t *entity.Invoices) *contract.Invoice {
		return &contract.Invoice{
			InvoiceID:    t.InvoiceID,
			IssueDate:    t.IssueDate.Format("02-01-2006"),
			Subject:      t.Subject,
			TotalItem:    t.TotalItems,
			CustomerName: t.CustomerName,
			DueDate:      t.DueDate.Format("02-01-2006"),
			Status:       t.Status,
			SubTotal:     t.SubTotal,
			Tax:          t.Tax,
			GrandTotal:   t.GrandTotal,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		}
	}).ToSlice()

	response = contract.ListInvoiceResponse{
		Data:       responseInvoicesList,
		Pagination: pagination,
	}

	return response, nil
}

func (ts *Invoiceservice) GetDetail(ctx context.Context, id string) (contract.InvoiceResponse, error) {
	var res contract.InvoiceResponse

	dataInvoices, err := ts.InvoicesRepo.Get(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return res, errorss.ErrInvoiceIdNotFound
		}
		log.Println(err)
		return res, err
	}

	customerID := dataInvoices.CustomerID.String()

	dataCustomer, err := ts.CustomerRepo.Get(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return res, errorss.ErrCustomerIdNotFound
		}
		log.Println(err)
		return res, err
	}

	dataItems, err := ts.ItemRepo.GetByInvoiceID(ctx, dataInvoices.InvoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return res, errorss.ErrInvoiceIdNotFound
		}
		log.Println(err)
		return res, err
	}

	items := stream.Map(stream.OfSlice(dataItems), func(i *entity.Item) contract.ItemResponse {
		return contract.ItemResponse{
			ItemID:    i.ItemID,
			Name:      i.Name,
			Quantity:  i.Quantity,
			UnitPrice: i.UnitPrice,
			Amount:    i.Amount,
		}
	}).ToSlice()

	res = contract.InvoiceResponse{
		InvoiceID:    dataInvoices.InvoiceID,
		IssueDate:    dataInvoices.IssueDate.Format("02-01-2006"),
		Subject:      dataInvoices.Subject,
		TotalItem:    dataInvoices.TotalItems,
		Items:        items,
		CustomerName: dataCustomer.Name,
		DueDate:      dataInvoices.DueDate.Format("02-01-2006"),
		Status:       dataInvoices.Status,
		SubTotal:     dataInvoices.SubTotal,
		Tax:          dataInvoices.Tax,
		GrandTotal:   dataInvoices.GrandTotal,
	}

	return res, nil
}

func (ts *Invoiceservice) Update(ctx context.Context, request contract.InvoiceRequest, id string) (contract.InvcResponse, error) {
	var res contract.InvcResponse

	dataInvoices, err := ts.InvoicesRepo.Get(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return res, errorss.ErrInvoiceIdNotFound
		}
		log.Println(err)
		return res, err
	}

	customerID := dataInvoices.CustomerID.String()

	dataCustomer, err := ts.CustomerRepo.Get(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return res, errorss.ErrCustomerIdNotFound
		}
		log.Println(err)
		return res, err
	}

	dataItems, err := ts.ItemRepo.GetByInvoiceID(ctx, dataInvoices.InvoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return res, errorss.ErrInvoiceIdNotFound
		}
		log.Println(err)
		return res, err
	}

	layout := "02-01-2006"
	newIssueDate, _ := time.Parse(layout, request.IssueDate)
	newDueDate, _ := time.Parse(layout, request.DueDate)

	err = frsAtomic.Atomic(ctx, ts.AtomicSession, func(ctx context.Context) error {

		// invoice
		dataInvoices.IssueDate = newIssueDate
		dataInvoices.Subject = request.Subject
		dataInvoices.TotalItems = len(request.ItemRequest)
		dataInvoices.DueDate = newDueDate
		dataInvoices.SubTotal = request.SubTotal
		dataInvoices.Tax = request.Tax
		dataInvoices.GrandTotal = request.GrandTotal

		// customer
		dataCustomer.Name = request.CustomerRequest.CustomerName
		dataCustomer.Address = request.CustomerRequest.Address

		var data []uuid.UUID
		var dataRequest []uuid.UUID
		for _, v := range dataItems {
			data = append(data, v.ItemID)
		}
		for _, v := range request.ItemRequest {
			dataRequest = append(dataRequest, v.ItemID)
		}

		var newID []uuid.UUID
		if len(data) > len(dataRequest) {
			missID := findItemID(data, dataRequest)
			newID = append(newID, missID...)
		}

		err := ts.ItemRepo.Delete(ctx, newID)
		if err != nil {
			log.Println("delete id err: ", err)
			return err
		}

		// item
		for i, v := range request.ItemRequest {
			dataItems[i] = &entity.Item{
				ItemData: entity.ItemData{
					InvoiceID: dataInvoices.InvoiceID,
					ItemID:    data[i],
					Name:      v.Name,
					Type:      v.Type,
					Quantity:  v.Quantity,
					UnitPrice: v.UnitPrice,
					Amount:    v.Amount,
				},
			}
		}

		err = ts.CustomerRepo.Update(ctx, &dataCustomer)
		if err != nil {
			log.Println("update customer err: ", err)
			return err
		}

		err = ts.InvoicesRepo.Update(ctx, &dataInvoices)
		if err != nil {
			log.Println("update invoice err: ", err)
			return err
		}

		err = ts.ItemRepo.Update(ctx, dataItems)
		if err != nil {
			log.Println("update item err: ", err)
			return err
		}

		res = contract.InvcResponse{
			InvoiceID: dataInvoices.InvoiceID,
		}

		return nil
	})

	if err != nil {
		log.Println("err: ", err)
		return res, err
	}

	return res, nil
}
