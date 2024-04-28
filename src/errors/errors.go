package errors

import (
	i18n_err "github.com/Risuii/frs-lib/i18n/errors"
)

var (
	ErrDuplicateInvoices  = i18n_err.NewI18nError("err_Invoices_duplicate")
	ErrCustomerIdNotFound = i18n_err.NewI18nError("err_customer_id_not_found")
	ErrInvoiceIdNotFound  = i18n_err.NewI18nError("err_invoice_id_not_found")
)
