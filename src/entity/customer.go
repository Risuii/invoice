package entity

import "github.com/google/uuid"

type Customer struct {
	ModelID
	ModelLogTime
	CustomerData
}

type CustomerData struct {
	CustomerID uuid.UUID `db:"customer_id"`
	Name       string    `db:"name"`
	Address    string    `db:"address"`
}
