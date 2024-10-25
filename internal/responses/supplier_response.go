package responses

import (
	"encoding/json"

	"github.com/google/uuid"
)

type SupplierResponse struct {
	ID      uuid.UUID       `json:"id"`
	Name    string          `json:"name"`
	Email   string          `json:"email"`
	Tel     string          `json:"tel"`
	Address json.RawMessage `json:"address"`
}

type SupplierListResponse struct {
	Suppliers []SupplierResponse `json:"suppliers"`
	Total     int64              `json:"total"`
}
