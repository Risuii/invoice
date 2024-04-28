package handler

import (
	"context"

	"github.com/Risuii/invoice/src/v1/contract"
)

type InvoiceService interface {
	GetList(ctx context.Context, params contract.GetListParam) (contract.ListInvoiceResponse, error)
	GetDetail(ctx context.Context, id string) (contract.InvoiceResponse, error)
	Create(ctx context.Context, request contract.InvoiceRequest) (contract.InvcResponse, error)
	Update(ctx context.Context, request contract.InvoiceRequest, id string) (contract.InvcResponse, error)
}
