package contract

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GetListParam struct {
	Page      int    `json:"page" db:"page"`
	Limit     int    `json:"limit" db:"limit"`
	Offset    int    `json:"offset" db:"offset"`
	Keyword   string `json:"keyword" db:"keyword"`
	InvoiceID string `json:"invoice_id" db:"invoice_id"`
	IssueDate string `json:"issue_date" db:"issue_date"`
	Subject   string `json:"subject" db:"subject"`
	TotalItem int    `json:"total_item" db:"total_item"`
	Customer  string `json:"customer" db:"customer"`
	DueDate   string `json:"due_date" db:"due_date"`
	Status    string `json:"status" db:"status"`
}

// ValidateQuery return common converted parameter from query parameter for get list data
// common query parameter is keyword, page, limit, and offset
// page is number page where the data is now, keyword is for search data by string keyword,
// limit is limit data loaded per page, offset is number data skiped when loaded data
// data page and limit from query parameter is always number in string
// its need to converted to int, it will return error if page and limit is not a number
func ValidateAndBuildRequest(r *http.Request) (getListParam *GetListParam, err error) {
	// default value for page and limit
	page, limit := 1, 10
	var item int

	// get data from query parameter
	queryParams := r.URL.Query()
	limitQuery := queryParams.Get("limit")
	pageQuery := queryParams.Get("page")
	keyword := queryParams.Get("keyword")
	InvoiceID := queryParams.Get("invoice_id")
	IssueDate := queryParams.Get("issue_date")
	Subject := queryParams.Get("subject")
	TotalItem := queryParams.Get("total_item")
	Customer := queryParams.Get("customer")
	DueDate := queryParams.Get("due_date")
	Status := queryParams.Get("status")

	// query param validation
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			return
		}
	}

	if limitQuery != "" {
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			return
		}
	}

	if TotalItem != "" {
		item, err = strconv.Atoi(TotalItem)
		if err != nil {
			return
		}
	}

	// offset for OFFSET in get list query
	offset := (page - 1) * limit
	getListParam = &GetListParam{
		Page:      page,
		Limit:     limit,
		Offset:    offset,
		Keyword:   keyword,
		InvoiceID: InvoiceID,
		IssueDate: IssueDate,
		Subject:   Subject,
		TotalItem: item,
		Customer:  Customer,
		DueDate:   DueDate,
		Status:    Status,
	}

	return
}

func ValidateIDParamRequest(r *http.Request) (id string, err error) {
	idParam := chi.URLParam(r, "id")

	return idParam, nil
}
