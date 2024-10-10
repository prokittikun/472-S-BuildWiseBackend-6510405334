package models

import (
	"time"

	"github.com/google/uuid"
)

type Quotation struct {
	OfferID    uuid.UUID `db:"offer_id"`
	ProjectID  uuid.UUID `db:"project_id"`
	SalePrice  float64   `db:"sale_price"`
	ValidUntil time.Time `db:"valid_until"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
