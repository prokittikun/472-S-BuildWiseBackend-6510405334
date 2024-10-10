package models

import (
	"time"

	"github.com/google/uuid"
)

type GeneralCost struct {
	GID           uuid.UUID `db:"g_id"`
	BID           uuid.UUID `db:"b_id"`
	TypeName      string    `db:"type_name"`
	SaleCost      float64   `db:"sale_cost"`
	EstimatedCost float64   `db:"estimated_cost"`
	Description   string    `db:"description"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
