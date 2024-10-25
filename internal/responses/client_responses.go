package responses

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ClientResponse struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Tel       string          `json:"tel"`
	Address   json.RawMessage `json:"address"`
	TaxID     string          `json:"tax_id"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type ClientListResponse struct {
	Clients []ClientResponse `json:"clients"`
	Total   int64            `json:"total"`
}
