package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Supplier struct {
	SupplierID uuid.UUID       `db:"supplier_id"`
	Name       string          `db:"name"`
	Email      string          `db:"email"`
	Tel        string          `db:"tel"`
	Address    json.RawMessage `db:"address"`
}
