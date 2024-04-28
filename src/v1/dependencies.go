package v1

import (
	"context"
	"log"

	"github.com/Risuii/invoice/src/app"
	"github.com/google/uuid"

	frsAtomicSQLX "github.com/Risuii/frs-lib/atomic/sqlx"
	customerRepo "github.com/Risuii/invoice/src/repository/customers"
	InvoicesRepo "github.com/Risuii/invoice/src/repository/invoice"
	itemsRepo "github.com/Risuii/invoice/src/repository/items"
	Invoicesvc "github.com/Risuii/invoice/src/v1/service/invoice"
)

type repositories struct {
	AtomicSessionProvider frsAtomicSQLX.SqlxAtomicSessionProvider
	InvoicesRepo          *InvoicesRepo.InvoicesRepository
	CustomersRepo         *customerRepo.CustomersRepository
	ItemsRepo             *itemsRepo.ItemsRepository
}

type services struct {
	Invoicesvc *Invoicesvc.Invoiceservice
}

type Dependency struct {
	Repositories *repositories
	Services     *services
}

type UUIDGeneratorImplementation struct{}

func (g UUIDGeneratorImplementation) New() uuid.UUID {
	return uuid.New()
}

func initRepositories(ctx context.Context) *repositories {
	var r repositories
	var err error

	r.AtomicSessionProvider = *frsAtomicSQLX.NewSqlxAtomicSessionProvider(app.DB())

	r.InvoicesRepo, err = InvoicesRepo.InitInvoicesRepository(ctx, app.DB(), app.Cache())
	if err != nil {
		log.Fatal("init Invoices repo err: ", err)
	}

	r.CustomersRepo, err = customerRepo.InitCustomersRepository(ctx, app.DB(), app.Cache())
	if err != nil {
		log.Fatal("init customers repo err: ", err)
	}

	r.ItemsRepo, err = itemsRepo.InitItemsRepository(ctx, app.DB(), app.Cache())
	if err != nil {
		log.Fatal("init items repo err: ", err)
	}

	return &r
}

func initServices(ctx context.Context, r *repositories) *services {

	uuidGen := UUIDGeneratorImplementation{}

	return &services{
		Invoicesvc: Invoicesvc.InitInvoiceservice(r.InvoicesRepo, r.CustomersRepo, r.ItemsRepo, &r.AtomicSessionProvider, uuidGen),
	}
}

func Dependencies(ctx context.Context) *Dependency {
	repositories := initRepositories(ctx)
	services := initServices(ctx, repositories)

	return &Dependency{
		Repositories: repositories,
		Services:     services,
	}
}
