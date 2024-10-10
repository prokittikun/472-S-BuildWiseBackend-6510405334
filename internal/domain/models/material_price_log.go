package models

import (
	"time"

	"github.com/google/uuid"
)

type MaterialPriceLog struct {
	MPLID          uuid.UUID `db:"mpl_id"`
	MaterialName   string    `db:"material_name"`
	ActualPrice    float64   `db:"actual_price"`
	SalePrice      float64   `db:"sale_price"`
	EstimatedPrice float64   `db:"estimated_price"`
	POID           uuid.UUID `db:"po_id"`
	CreatedAt      time.Time `db:"created_at"`
}
