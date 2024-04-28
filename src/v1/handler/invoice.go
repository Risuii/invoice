package handler

import (
	"log"
	"net/http"

	"github.com/Risuii/invoice/src/errors"
	"github.com/Risuii/invoice/src/middleware/response"
	"github.com/Risuii/invoice/src/v1/contract"
)

func CreateInvoiceHandler(svc InvoiceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		invoiceRequest, err := contract.BuildAndValidateInvoiceRequest(r)
		if err != nil {
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		res, err := svc.Create(r.Context(), invoiceRequest)
		if err != nil {
			log.Println(err)
			response.JSONInternalErrorResponse(r.Context(), w)
			return
		}

		response.JSONSuccessResponse(r.Context(), w, res)
	}
}

func UpdateInvoiceHandler(svc InvoiceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := contract.ValidateIDParamRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		invoiceRequest, err := contract.BuildAndValidateInvoiceRequest(r)
		if err != nil {
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		res, err := svc.Update(r.Context(), invoiceRequest, id)
		if err != nil {
			log.Println(err)
			switch err {
			case errors.ErrInvoiceIdNotFound,
				errors.ErrCustomerIdNotFound:
				response.JSONUnprocessableEntity(r.Context(), w, err)
			default:
				response.JSONInternalErrorResponse(r.Context(), w)
			}
			return
		}

		response.JSONSuccessResponse(r.Context(), w, res)
	}
}

func GetListInvoicesHandler(svc InvoiceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params, err := contract.ValidateAndBuildRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		data, err := svc.GetList(r.Context(), *params)
		if err != nil {
			log.Println(err)
			response.JSONInternalErrorResponse(r.Context(), w)
			return
		}

		response.JSONSuccessResponse(r.Context(), w, data)
	}
}

func GetDetailInvoicesHandler(svc InvoiceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := contract.ValidateIDParamRequest(r)
		if err != nil {
			log.Println(err)
			response.JSONBadRequestResponse(r.Context(), w)
			return
		}

		data, err := svc.GetDetail(r.Context(), id)
		if err != nil {
			log.Println(err)
			switch err {
			case errors.ErrInvoiceIdNotFound,
				errors.ErrCustomerIdNotFound:
				response.JSONUnprocessableEntity(r.Context(), w, err)
			default:
				response.JSONInternalErrorResponse(r.Context(), w)
			}
			return
		}

		response.JSONSuccessResponse(r.Context(), w, data)
	}
}
