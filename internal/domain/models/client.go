package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Client struct {
	ClientID uuid.UUID       `db:"client_id"`
	Name     string          `db:"name"`
	Email    string          `db:"email"`
	Tel      string          `db:"tel"`
	Address  json.RawMessage `db:"address"`
	TaxID    string          `db:"tax_id"`
}
