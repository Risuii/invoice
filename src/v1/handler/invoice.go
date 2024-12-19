package handler

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Risuii/invoice/src/app"
	errorss "github.com/Risuii/invoice/src/errors"
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
			case errorss.ErrInvoiceIdNotFound,
				errorss.ErrCustomerIdNotFound:
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
			case errorss.ErrInvoiceIdNotFound,
				errorss.ErrCustomerIdNotFound:
				response.JSONUnprocessableEntity(r.Context(), w, err)
			default:
				response.JSONInternalErrorResponse(r.Context(), w)
			}
			return
		}

		response.JSONSuccessResponse(r.Context(), w, data)
	}
}

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Image       string  `json:"image"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type Pagination struct {
	Page      int `json:"page"`
	TotalPage int `json:"total_page"`
	TotalData int `json:"total_data"`
}

type Data struct {
	Movies     []Movie    `json:"Data"`
	Pagination Pagination `json:"Pagination"`
}

type APIResponse struct {
	Data Data `json:"data"`
}

func TestApi(svc InvoiceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoint := "http://localhost:3001/Movies"

		ResponseFromAnotherEndPoint, err := fetchFromEndpoint(endpoint)
		if err != nil {
			log.Println(err)
			response.JSONTimeOut(r.Context(), w, errorss.ErrTimeOut)
			return
		}

		var res APIResponse
		err = json.Unmarshal([]byte(ResponseFromAnotherEndPoint), &res)
		if err != nil {
			http.Error(w, "Error parsing response", http.StatusInternalServerError)
			return
		}

		response.JSONSuccessResponse(r.Context(), w, res)
	}
}

func fetchFromEndpoint(url string) (string, error) {
	var netErr net.Error
	delayTime := app.Config().DelayTime
	client := &http.Client{
		Timeout: time.Duration(delayTime) * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		if errors.As(err, &netErr) && netErr.Timeout() {
			return "", errorss.ErrTimeOut
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("received non-200 response code")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
