package requests

import "encoding/json"

type CreateClientRequest struct {
	Name    string          `json:"name" validate:"required"`
	Email   string          `json:"email" validate:"required,email"`
	Tel     string          `json:"tel" validate:"required,len=10"`
	Address json.RawMessage `json:"address" validate:"required"`
	TaxID   string          `json:"tax_id" validate:"required,len=13"`
}

type UpdateClientRequest struct {
	Name    string          `json:"name" validate:"required"`
	Email   string          `json:"email" validate:"required,email"`
	Tel     string          `json:"tel" validate:"required,len=10"`
	Address json.RawMessage `json:"address" validate:"required"`
	TaxID   string          `json:"tax_id" validate:"required,len=13"`
}
