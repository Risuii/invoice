package contract

import "github.com/google/uuid"

type Item struct {
	InvoiceID string    `json:"invoice_id"`
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Quantity  float64   `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Amount    float64   `json:"amount"`
}

type ItemResponse struct {
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Quantity  float64   `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Amount    float64   `json:"amount"`
}

type ItemRequest struct {
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Quantity  float64   `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Amount    float64   `json:"amount"`
}
