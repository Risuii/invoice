package contract

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	frsUtils "github.com/Risuii/frs-lib/utils"
	"github.com/go-playground/validator/v10"
)

type Invoice struct {
	InvoiceID    string    `json:"invoice_id"`
	IssueDate    string    `json:"issue_date"`
	Subject      string    `json:"subject"`
	TotalItem    int       `json:"total_item"`
	CustomerName string    `json:"customer_name"`
	DueDate      string    `json:"due_date"`
	Status       string    `json:"status"`
	SubTotal     float64   `json:"sub_total"`
	Tax          float64   `json:"tax"`
	GrandTotal   float64   `json:"grand_total"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ListInvoiceResponse struct {
	Data       []*Invoice
	Pagination *frsUtils.Pagination
}

type InvoiceResponse struct {
	InvoiceID    string         `json:"invoice_id"`
	IssueDate    string         `json:"issue_date"`
	Subject      string         `json:"subject"`
	TotalItem    int            `json:"total_item"`
	Items        []ItemResponse `json:"item"`
	CustomerName string         `json:"customer_name"`
	DueDate      string         `json:"due_date"`
	Status       string         `json:"status"`
	SubTotal     float64        `json:"sub_total"`
	Tax          float64        `json:"tax"`
	GrandTotal   float64        `json:"grand_total"`
}

type InvoiceRequest struct {
	Subject         string          `json:"subject" validate:"required"`
	IssueDate       string          `json:"issue_date" validate:"required"`
	DueDate         string          `json:"due_date" validate:"required"`
	SubTotal        float64         `json:"sub_total"`
	Tax             float64         `json:"tax"`
	GrandTotal      float64         `json:"grand_total"`
	CustomerRequest CustomerRequest `json:"customer_request"`
	ItemRequest     []ItemRequest   `json:"item_request"`
}

type InvcResponse struct {
	InvoiceID string `json:"invoice_id"`
}

type InvoiceResponseDB struct {
	InvoiceID  string `db:"invoice_id"`
	CustomerID string `db:"customer_id"`
}

func checkSpecialCharacter(inputString string) bool {
	for _, char := range inputString {
		if unicode.IsPunct(char) && char != '-' && char != '_' {
			return true
		}
	}
	return false
}

func BuildAndValidateInvoiceRequest(r *http.Request) (InvoiceRequest, error) {
	var payload InvoiceRequest

	bodyByte, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("read request body err: ", err)
		return payload, err
	}

	if err := json.Unmarshal(bodyByte, &payload); err != nil {
		log.Println("unmarshal request body err: ", err)
		return payload, err
	}

	payload.Subject = strings.ToLower(payload.Subject)
	payload.CustomerRequest.CustomerName = strings.ToLower(payload.CustomerRequest.CustomerName)

	isSpecial := checkSpecialCharacter(payload.Subject)
	if isSpecial {
		return payload, errors.New("there is special characters")
	}

	validator := validator.New()

	if err := validator.Struct(payload); err != nil {
		log.Println("validate request body err: ", err)
		return payload, err
	}

	return payload, nil
}
