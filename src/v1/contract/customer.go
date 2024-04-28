package contract

type CustomerRequest struct {
	CustomerName string `json:"customer_name" validate:"required"`
	Address      string `json:"address" validate:"required"`
}
