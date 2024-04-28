package Invoices

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"testing"
	"time"

	atomic "github.com/Risuii/frs-lib/atomic"
	"github.com/Risuii/invoice/src/entity"
	errorss "github.com/Risuii/invoice/src/errors"
	"github.com/Risuii/invoice/src/v1/contract"
	"github.com/go-faker/faker/v4"
	"github.com/go-playground/assert"
	"github.com/google/uuid"
	"github.com/mariomac/gostream/stream"

	"go.uber.org/mock/gomock"

	mock_atomic "github.com/Risuii/frs-lib/atomic/mock"
	frsUtils "github.com/Risuii/frs-lib/utils"
	mock_Invoices "github.com/Risuii/invoice/src/v1/service/mock/invoice"
)

type FixedUUIDGenerator struct{}

func (g FixedUUIDGenerator) New() uuid.UUID {
	return uuid.MustParse("00000000-0000-0000-0000-000000000000") // Use a fixed UUID for testing
}

func TestInvoiceService_GetList(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	type (
		getListInvoices struct {
			Invoices []*entity.Invoices
			err      error
		}

		getInvoicesCount struct {
			count int64
			err   error
		}

		given struct {
			params           contract.GetListParam
			getListInvoices  getListInvoices
			getInvoicesCount getInvoicesCount
		}

		expected struct {
			res contract.ListInvoiceResponse
			err error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	currentDate := time.Now()
	var mockInvoices []*entity.Invoices
	mockID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	mockParams := contract.GetListParam{
		Page:  1,
		Limit: 10,
	}

	sizeDataSet := 5

	for i := 0; i < sizeDataSet; i++ {
		mockInvoices = append(mockInvoices, &entity.Invoices{
			ModelID: entity.ModelID{
				Id: 1,
			},
			ModelLogTime: entity.ModelLogTime{
				CreatedAt: currentDate,
				UpdatedAt: currentDate,
			},
			InvoicesData: entity.InvoicesData{
				InvoiceID:    faker.Name(),
				IssueDate:    currentDate,
				Subject:      faker.Name(),
				TotalItems:   3,
				CustomerID:   mockID,
				DueDate:      currentDate,
				Status:       faker.Name(),
				CustomerName: faker.Name(),
			},
		})
	}

	mockInvoicesResp := stream.Map(stream.OfSlice(mockInvoices), func(t *entity.Invoices) *contract.Invoice {
		return &contract.Invoice{
			InvoiceID:    t.InvoiceID,
			IssueDate:    t.IssueDate.Format("02-01-2006"),
			Subject:      t.Subject,
			TotalItem:    t.TotalItems,
			CustomerName: t.CustomerName,
			DueDate:      t.DueDate.Format("02-01-2006"),
			Status:       t.Status,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		}
	}).ToSlice()

	testCases := []testCase{
		{
			name: "error sql",
			given: given{
				params: mockParams,
				getListInvoices: getListInvoices{
					Invoices: mockInvoices,
					err:      errors.New("error"),
				},
			},
			expected: expected{
				res: contract.ListInvoiceResponse{},
				err: errors.New("error"),
			},
		},
		{
			name: "error count",
			given: given{
				params: mockParams,
				getListInvoices: getListInvoices{
					Invoices: mockInvoices,
					err:      nil,
				},
				getInvoicesCount: getInvoicesCount{
					err: errors.New("error"),
				},
			},
			expected: expected{
				res: contract.ListInvoiceResponse{},
				err: errors.New("error"),
			},
		},
		{
			name: "success",
			given: given{
				params: mockParams,
				getListInvoices: getListInvoices{
					Invoices: mockInvoices,
					err:      nil,
				},
				getInvoicesCount: getInvoicesCount{
					count: int64(2 * len(mockInvoices)),
					err:   nil,
				},
			},
			expected: expected{
				res: contract.ListInvoiceResponse{
					Data: mockInvoicesResp,
					Pagination: &frsUtils.Pagination{
						Page:      1,
						TotalPage: 1,
						TotalData: int64(2 * len(mockInvoices)),
					},
				},
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		mockInvoicesRepo := mock_Invoices.NewMockInvoicesRepository(mockCtrl)
		mockCustomerRepo := mock_Invoices.NewMockCustomerRepository(mockCtrl)
		mockItemRepo := mock_Invoices.NewMockItemRepository(mockCtrl)
		mockAsession := mock_atomic.NewMockAtomicSessionProvider(mockCtrl)

		func() {
			mockInvoicesRepo.EXPECT().GetList(gomock.Any(), testCase.given.params).
				Return(testCase.given.getListInvoices.Invoices, testCase.given.getListInvoices.err).
				Times(1)

			if testCase.given.getListInvoices.err == nil {
				mockInvoicesRepo.EXPECT().GetInvoicesCount(gomock.Any(), testCase.given.params).
					Return(testCase.given.getInvoicesCount.count, testCase.given.getInvoicesCount.err).
					Times(1)
			}
		}()

		Invoices := InitInvoiceservice(mockInvoicesRepo, mockCustomerRepo, mockItemRepo, mockAsession, FixedUUIDGenerator{})
		got, actualErr := Invoices.GetList(context.Background(), testCase.given.params)
		assert.Equal(t, testCase.expected.res, got)
		assert.Equal(t, testCase.expected.err, actualErr)
	}
}

func TestInvoiceService_GetDetail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	type (
		getInvoices struct {
			Invoices entity.Invoices
			err      error
		}

		getCustomer struct {
			customer entity.Customer
			err      error
		}

		getItems struct {
			items []*entity.Item
			err   error
		}

		given struct {
			id          string
			uuid        uuid.UUID
			getInvoices getInvoices
			getCustomer getCustomer
			getItems    getItems
		}

		expected struct {
			res contract.InvoiceResponse
			err error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	currentDate := time.Now()
	mockID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	var mockItem []*entity.Item

	sizeDataSet := 5

	for i := 0; i < sizeDataSet; i++ {
		mockItem = append(mockItem, &entity.Item{
			ModelID: entity.ModelID{
				Id: 1,
			},
			ModelLogTime: entity.ModelLogTime{
				CreatedAt: currentDate,
				UpdatedAt: currentDate,
			},
			ItemData: entity.ItemData{
				InvoiceID: faker.Name(),
				Name:      faker.Name(),
				Type:      faker.Name(),
				Quantity:  5,
				UnitPrice: 5,
				Amount:    5,
			},
		})
	}

	mockItems := stream.Map(stream.OfSlice(mockItem), func(i *entity.Item) contract.ItemResponse {
		return contract.ItemResponse{
			Name:      i.Name,
			Quantity:  i.Quantity,
			UnitPrice: i.UnitPrice,
			Amount:    i.Amount,
		}
	}).ToSlice()

	testCases := []testCase{
		{
			name: "error sql no rows Invoices",
			given: given{
				getInvoices: getInvoices{
					err: sql.ErrNoRows,
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{},
				err: errorss.ErrInvoiceIdNotFound,
			},
		},
		{
			name: "error Invoices",
			given: given{
				getInvoices: getInvoices{
					err: errors.New("error"),
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{},
				err: errors.New("error"),
			},
		},
		{
			name: "error sql no rows customer",
			given: given{
				getInvoices: getInvoices{
					Invoices: entity.Invoices{
						InvoicesData: entity.InvoicesData{
							CustomerID: mockID,
						},
					},
					err: nil,
				},
				getCustomer: getCustomer{
					err: sql.ErrNoRows,
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{},
				err: errorss.ErrCustomerIdNotFound,
			},
		},
		{
			name: "error customer",
			given: given{
				getInvoices: getInvoices{
					Invoices: entity.Invoices{
						InvoicesData: entity.InvoicesData{
							CustomerID: mockID,
						},
					},
					err: nil,
				},
				getCustomer: getCustomer{
					err: errors.New("error"),
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{},
				err: errors.New("error"),
			},
		},
		{
			name: "error sql no rows items",
			given: given{
				getInvoices: getInvoices{
					Invoices: entity.Invoices{
						InvoicesData: entity.InvoicesData{
							CustomerID: mockID,
							InvoiceID:  "test-invoiceID",
						},
					},
					err: nil,
				},
				getCustomer: getCustomer{
					customer: entity.Customer{},
					err:      nil,
				},
				getItems: getItems{
					err: sql.ErrNoRows,
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{},
				err: errorss.ErrInvoiceIdNotFound,
			},
		},
		{
			name: "error items",
			given: given{
				getInvoices: getInvoices{
					Invoices: entity.Invoices{
						InvoicesData: entity.InvoicesData{
							CustomerID: mockID,
							InvoiceID:  "test-invoiceID",
						},
					},
					err: nil,
				},
				getCustomer: getCustomer{
					customer: entity.Customer{},
					err:      nil,
				},
				getItems: getItems{
					err: errors.New("error"),
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{},
				err: errors.New("error"),
			},
		},
		{
			name: "success",
			given: given{
				getInvoices: getInvoices{
					Invoices: entity.Invoices{
						InvoicesData: entity.InvoicesData{
							CustomerID: mockID,
							InvoiceID:  "",
						},
					},
					err: nil,
				},
				getCustomer: getCustomer{
					customer: entity.Customer{},
					err:      nil,
				},
				getItems: getItems{
					items: mockItem,
					err:   nil,
				},
			},
			expected: expected{
				res: contract.InvoiceResponse{
					IssueDate: "01-01-0001",
					Items:     mockItems,
					DueDate:   "01-01-0001",
				},
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		mockInvoicesRepo := mock_Invoices.NewMockInvoicesRepository(mockCtrl)
		mockCustomerRepo := mock_Invoices.NewMockCustomerRepository(mockCtrl)
		mockItemRepo := mock_Invoices.NewMockItemRepository(mockCtrl)
		mockAsession := mock_atomic.NewMockAtomicSessionProvider(mockCtrl)

		func() {
			mockInvoicesRepo.EXPECT().Get(gomock.Any(), testCase.given.id).
				Return(testCase.given.getInvoices.Invoices, testCase.given.getInvoices.err).
				Times(1)
			if testCase.given.getInvoices.err == nil {
				mockCustomerRepo.EXPECT().Get(gomock.Any(), testCase.given.uuid.String()).
					Return(testCase.given.getCustomer.customer, testCase.given.getCustomer.err).
					Times(1)

				if testCase.given.getCustomer.err == nil {
					mockItemRepo.EXPECT().GetByInvoiceID(gomock.Any(), testCase.given.getInvoices.Invoices.InvoiceID).
						Return(testCase.given.getItems.items, testCase.given.getItems.err).
						Times(1)
				}
			}
		}()

		Invoices := InitInvoiceservice(mockInvoicesRepo, mockCustomerRepo, mockItemRepo, mockAsession, FixedUUIDGenerator{})
		got, actualErr := Invoices.GetDetail(context.Background(), testCase.given.id)
		log.Println(got)
		assert.Equal(t, testCase.expected.res, got)
		assert.Equal(t, testCase.expected.err, actualErr)
	}
}

func TestInvoiceService_Create(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	type (
		getLatestInvoiceID struct {
			invoiceID string
			err       error
		}

		createCustomer struct {
			err error
		}

		createInvoice struct {
			invoice contract.InvoiceResponseDB
			err     error
		}

		createItem struct {
			err error
		}

		given struct {
			req                contract.InvoiceRequest
			dataCustomer       entity.Customer
			dataInvoices       entity.Invoices
			dataItem           []*entity.Item
			getLatestInvoiceID getLatestInvoiceID
			createCustomer     createCustomer
			createInvoice      createInvoice
			createItem         createItem
		}

		expected struct {
			res contract.InvcResponse
			err error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	mockID := uuid.MustParse("00000000-0000-0000-0000-000000000000")

	mockInvoiceRequest := contract.InvoiceRequest{
		Subject:    faker.Name(),
		SubTotal:   1,
		Tax:        1,
		GrandTotal: 1,
		CustomerRequest: contract.CustomerRequest{
			CustomerName: faker.Name(),
			Address:      faker.Name(),
		},
		ItemRequest: []contract.ItemRequest{},
	}

	mockInsertDataCustomer := entity.Customer{
		CustomerData: entity.CustomerData{
			CustomerID: mockID,
			Name:       mockInvoiceRequest.CustomerRequest.CustomerName,
			Address:    mockInvoiceRequest.CustomerRequest.Address,
		},
	}

	mockInsertDataInvoice := entity.Invoices{
		ModelID:      entity.ModelID{},
		ModelLogTime: entity.ModelLogTime{},
		InvoicesData: entity.InvoicesData{
			InvoiceID:  "0001",
			Subject:    mockInvoiceRequest.Subject,
			TotalItems: len(mockInvoiceRequest.ItemRequest),
			CustomerID: mockInsertDataCustomer.CustomerID,
			SubTotal:   mockInvoiceRequest.SubTotal,
			Tax:        mockInvoiceRequest.Tax,
			GrandTotal: mockInvoiceRequest.GrandTotal,
			Status:     "Unpaid",
		},
	}

	mockItemResp := stream.Map(stream.OfSlice(mockInvoiceRequest.ItemRequest), func(t contract.ItemRequest) *entity.Item {
		return &entity.Item{
			ItemData: entity.ItemData{
				ItemID:    mockID,
				Name:      t.Name,
				Type:      t.Type,
				Quantity:  t.Quantity,
				UnitPrice: t.UnitPrice,
				Amount:    t.Amount,
			},
		}
	}).ToSlice()

	testCases := []testCase{
		{
			name: "err create customer",
			given: given{
				createCustomer: createCustomer{
					err: errors.New("error internal server"),
				},
			},
			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "err create invoice",
			given: given{
				req:          mockInvoiceRequest,
				dataInvoices: mockInsertDataInvoice,
				dataCustomer: mockInsertDataCustomer,
				createCustomer: createCustomer{
					err: nil,
				},
				createInvoice: createInvoice{
					err: errors.New("error internal server"),
				},
			},
			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "err create item",
			given: given{
				req:          mockInvoiceRequest,
				dataInvoices: mockInsertDataInvoice,
				dataItem:     mockItemResp,
				dataCustomer: mockInsertDataCustomer,
				createCustomer: createCustomer{
					err: nil,
				},
				createInvoice: createInvoice{
					invoice: contract.InvoiceResponseDB{
						InvoiceID:  faker.Name(),
						CustomerID: faker.Name(),
					},
					err: nil,
				},
				createItem: createItem{
					err: errors.New("error internal server"),
				},
			},
			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "success",
			given: given{
				req:          mockInvoiceRequest,
				dataInvoices: mockInsertDataInvoice,
				dataItem:     mockItemResp,
				dataCustomer: mockInsertDataCustomer,
				createCustomer: createCustomer{
					err: nil,
				},
				createInvoice: createInvoice{
					invoice: contract.InvoiceResponseDB{
						InvoiceID:  "test-id",
						CustomerID: "test-id",
					},
					err: nil,
				},
				createItem: createItem{
					err: nil,
				},
			},
			expected: expected{
				res: contract.InvcResponse{
					InvoiceID: "test-id",
				},
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockInvoicesRepo := mock_Invoices.NewMockInvoicesRepository(mockCtrl)
			mockCustomerRepo := mock_Invoices.NewMockCustomerRepository(mockCtrl)
			mockItemRepo := mock_Invoices.NewMockItemRepository(mockCtrl)

			mockAsession := mock_atomic.NewMockAtomicSessionProvider(mockCtrl)
			mockAtomicSession := mock_atomic.NewMockAtomicSession(mockCtrl)
			mockAtomicSessionCtx := atomic.NewAtomicSessionContext(context.Background(), mockAtomicSession)

			func() {
				mockInvoicesRepo.EXPECT().GetLatestInvoiceID(gomock.Any()).
					Return(testCase.given.getLatestInvoiceID.invoiceID, testCase.given.getLatestInvoiceID.err).
					Times(1)

				mockAsession.EXPECT().BeginSession(gomock.Any()).
					Return(mockAtomicSessionCtx, nil).
					Times(1)

				mockCustomerRepo.EXPECT().Create(gomock.Any(), &testCase.given.dataCustomer).
					Return(testCase.given.createCustomer.err).
					Times(1)

				testCase.given.dataInvoices.Status = "Unpaid"

				mockInvoicesRepo.EXPECT().Create(gomock.Any(), &testCase.given.dataInvoices).
					Return(testCase.given.createInvoice.invoice, testCase.given.createInvoice.err).
					Times(1)

				mockItemRepo.EXPECT().Create(gomock.Any(), testCase.given.dataItem).
					Return(testCase.given.createItem.err).
					Times(1)

				mockAtomicSession.EXPECT().Rollback(gomock.Any()).Times(1)

				mockAtomicSession.EXPECT().Commit(gomock.Any()).Times(1)
			}()

			Invoices := InitInvoiceservice(mockInvoicesRepo, mockCustomerRepo, mockItemRepo, mockAsession, FixedUUIDGenerator{})
			got, actualErr := Invoices.Create(context.Background(), testCase.given.req)
			log.Println(got)
			assert.Equal(t, testCase.expected.res, got)
			assert.Equal(t, testCase.expected.err, actualErr)
		})
	}
}

func TestInvoiceService_Update(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockCtrl.Finish()

	type (
		getDataInvoice struct {
			dataInvoice entity.Invoices
			err         error
		}

		getDataCustomer struct {
			dataCustomer entity.Customer
			err          error
		}

		getDataItems struct {
			dataItem []*entity.Item
			err      error
		}

		deleteDataItems struct {
			err error
		}

		updateCustomer struct {
			err error
		}

		updateInvoice struct {
			err error
		}

		updateItem struct {
			err error
		}

		given struct {
			req             contract.InvoiceRequest
			id              string
			getDataInvoice  getDataInvoice
			getDataCustomer getDataCustomer
			getDataItems    getDataItems
			deleteDataItems deleteDataItems
			updateCustomer  updateCustomer
			updateInvoice   updateInvoice
			updateItem      updateItem
		}

		expected struct {
			res contract.InvcResponse
			err error
		}

		testCase struct {
			name     string
			given    given
			expected expected
		}
	)

	mockID := uuid.MustParse("00000000-0000-0000-0000-000000000000")
	var mockUUID []uuid.UUID

	mockInvoiceRequest := contract.InvoiceRequest{
		Subject:    "test-subject",
		SubTotal:   0,
		Tax:        0,
		GrandTotal: 0,
		CustomerRequest: contract.CustomerRequest{
			CustomerName: "test-name",
			Address:      "test-address",
		},
		ItemRequest: []contract.ItemRequest{},
	}

	mockEntityInvoice := entity.Invoices{
		ModelID:      entity.ModelID{},
		ModelLogTime: entity.ModelLogTime{},
		InvoicesData: entity.InvoicesData{
			InvoiceID:  "0001",
			Subject:    "test-subject",
			TotalItems: len(mockInvoiceRequest.ItemRequest),
			CustomerID: mockID,
			SubTotal:   0,
			Tax:        0,
			GrandTotal: 0,
			Status:     "Unpaid",
		},
	}

	mockEntityCustomer := entity.Customer{
		CustomerData: entity.CustomerData{
			CustomerID: mockID,
			Name:       "test-name",
			Address:    "test-address",
		},
	}

	mockEntityItem := []*entity.Item{}

	testCases := []testCase{
		{
			name: "error invoice id not found",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					err: sql.ErrNoRows,
				},
			},

			expected: expected{
				err: errorss.ErrInvoiceIdNotFound,
			},
		},
		{
			name: "error get data invoices",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "error customer id not found",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					err: sql.ErrNoRows,
				},
			},

			expected: expected{
				err: errorss.ErrCustomerIdNotFound,
			},
		},
		{
			name: "error get customer data",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "error invoice id not found in items",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					err: sql.ErrNoRows,
				},
			},

			expected: expected{
				err: errorss.ErrInvoiceIdNotFound,
			},
		},
		{
			name: "error get items data",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "error delete item data",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					dataItem: mockEntityItem,
					err:      nil,
				},
				deleteDataItems: deleteDataItems{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "error update customer",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					dataItem: mockEntityItem,
					err:      nil,
				},
				deleteDataItems: deleteDataItems{
					err: nil,
				},
				updateCustomer: updateCustomer{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "error update invoice",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					dataItem: mockEntityItem,
					err:      nil,
				},
				deleteDataItems: deleteDataItems{
					err: nil,
				},
				updateCustomer: updateCustomer{
					err: nil,
				},
				updateInvoice: updateInvoice{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "error update item",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					dataItem: mockEntityItem,
					err:      nil,
				},
				deleteDataItems: deleteDataItems{
					err: nil,
				},
				updateCustomer: updateCustomer{
					err: nil,
				},
				updateInvoice: updateInvoice{
					err: nil,
				},
				updateItem: updateItem{
					err: errors.New("error internal server"),
				},
			},

			expected: expected{
				err: errors.New("error internal server"),
			},
		},
		{
			name: "success",
			given: given{
				id:  "test-id",
				req: mockInvoiceRequest,
				getDataInvoice: getDataInvoice{
					dataInvoice: mockEntityInvoice,
					err:         nil,
				},
				getDataCustomer: getDataCustomer{
					dataCustomer: mockEntityCustomer,
					err:          nil,
				},
				getDataItems: getDataItems{
					dataItem: mockEntityItem,
					err:      nil,
				},
				deleteDataItems: deleteDataItems{
					err: nil,
				},
				updateCustomer: updateCustomer{
					err: nil,
				},
				updateInvoice: updateInvoice{
					err: nil,
				},
				updateItem: updateItem{
					err: nil,
				},
			},

			expected: expected{
				res: contract.InvcResponse{
					InvoiceID: mockEntityInvoice.InvoiceID,
				},
				err: nil,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			mockInvoicesRepo := mock_Invoices.NewMockInvoicesRepository(mockCtrl)
			mockCustomerRepo := mock_Invoices.NewMockCustomerRepository(mockCtrl)
			mockItemRepo := mock_Invoices.NewMockItemRepository(mockCtrl)

			mockAsession := mock_atomic.NewMockAtomicSessionProvider(mockCtrl)
			mockAtomicSession := mock_atomic.NewMockAtomicSession(mockCtrl)
			mockAtomicSessionCtx := atomic.NewAtomicSessionContext(context.Background(), mockAtomicSession)

			func() {
				mockInvoicesRepo.EXPECT().Get(gomock.Any(), testCase.given.id).
					Return(testCase.given.getDataInvoice.dataInvoice, testCase.given.getDataInvoice.err).
					Times(1)

				customerID := testCase.given.getDataInvoice.dataInvoice.CustomerID.String()

				mockCustomerRepo.EXPECT().Get(gomock.Any(), customerID).
					Return(testCase.given.getDataCustomer.dataCustomer, testCase.given.getDataCustomer.err).
					Times(1)

				mockItemRepo.EXPECT().GetByInvoiceID(gomock.Any(), testCase.given.getDataInvoice.dataInvoice.InvoiceID).
					Return(testCase.given.getDataItems.dataItem, testCase.given.getDataItems.err).
					Times(1)

				mockAsession.EXPECT().BeginSession(gomock.Any()).
					Return(mockAtomicSessionCtx, nil).
					Times(1)

				mockItemRepo.EXPECT().Delete(gomock.Any(), mockUUID).
					Return(testCase.given.deleteDataItems.err).
					Times(1)

				mockCustomerRepo.EXPECT().Update(gomock.Any(), &mockEntityCustomer).
					Return(testCase.given.updateCustomer.err).
					Times(1)

				mockInvoicesRepo.EXPECT().Update(gomock.Any(), &mockEntityInvoice).
					Return(testCase.given.updateInvoice.err).
					Times(1)

				mockItemRepo.EXPECT().Update(gomock.Any(), mockEntityItem).
					Return(testCase.given.updateItem.err).
					Times(1)

				mockAtomicSession.EXPECT().Rollback(gomock.Any()).Times(1)

				mockAtomicSession.EXPECT().Commit(gomock.Any()).Times(1)

			}()

			Invoices := InitInvoiceservice(mockInvoicesRepo, mockCustomerRepo, mockItemRepo, mockAsession, FixedUUIDGenerator{})
			got, actualErr := Invoices.Update(context.Background(), testCase.given.req, testCase.given.id)
			log.Println(got)
			assert.Equal(t, testCase.expected.res, got)
			assert.Equal(t, testCase.expected.err, actualErr)
		})
	}
}
