package Invoices

import (
	"context"

	"github.com/Risuii/invoice/src/entity"
	"github.com/Risuii/invoice/src/v1/contract"
	"github.com/google/uuid"
)

type InvoicesRepository interface {
	Create(ctx context.Context, data *entity.Invoices) (contract.InvoiceResponseDB, error)
	GetList(ctx context.Context, params contract.GetListParam) ([]*entity.Invoices, error)
	GetInvoicesCount(ctx context.Context, param contract.GetListParam) (int64, error)
	Get(ctx context.Context, id string) (entity.Invoices, error)
	GetLatestInvoiceID(ctx context.Context) (string, error)
	Update(ctx context.Context, data *entity.Invoices) error
}

type CustomerRepository interface {
	Create(ctx context.Context, data *entity.Customer) error
	Get(ctx context.Context, id string) (entity.Customer, error)
	Update(ctx context.Context, data *entity.Customer) error
}

type ItemRepository interface {
	Create(ctx context.Context, data []*entity.Item) error
	GetByInvoiceID(ctx context.Context, invID string) ([]*entity.Item, error)
	Update(ctx context.Context, data []*entity.Item) error
	Delete(ctx context.Context, ids []uuid.UUID) error
}
