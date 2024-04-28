package entity

import "github.com/google/uuid"

type Item struct {
	ModelID
	ModelLogTime
	ItemData
}

type ItemData struct {
	InvoiceID string    `db:"invoice_id"`
	ItemID    uuid.UUID `db:"item_id"`
	Name      string    `db:"name"`
	Type      string    `db:"type"`
	Quantity  float64   `db:"quantity"`
	UnitPrice float64   `db:"unit_price"`
	Amount    float64   `db:"amount"`
}
