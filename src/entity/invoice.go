package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invoices struct {
	ModelID
	ModelLogTime
	InvoicesData
}

type InvoicesData struct {
	InvoiceID    string    `db:"invoice_id"`
	IssueDate    time.Time `db:"issue_date"`
	Subject      string    `db:"subject"`
	TotalItems   int       `db:"total_items"`
	CustomerID   uuid.UUID `db:"customer_id"`
	DueDate      time.Time `db:"due_date"`
	Status       string    `db:"status"`
	SubTotal     float64   `db:"sub_total"`
	Tax          float64   `db:"tax"`
	GrandTotal   float64   `db:"grand_total"`
	CustomerName string    `db:"customer_name"`
}
