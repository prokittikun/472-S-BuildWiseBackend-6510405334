package requests

import "encoding/json"

type CreateSupplierRequest struct {
	Name    string          `json:"name" validate:"required"`
	Email   string          `json:"email" validate:"required"`
	Tel     string          `json:"tel" validate:"required"`
	Address json.RawMessage `json:"address" validate:"required"`
}

type UpdateSupplierRequest struct {
	Name    string          `json:"name" validate:"required"`
	Email   string          `json:"email" validate:"required"`
	Tel     string          `json:"tel" validate:"required"`
	Address json.RawMessage `json:"address" validate:"required"`
}
