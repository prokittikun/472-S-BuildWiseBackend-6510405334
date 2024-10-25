package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type MaterialPriceLog struct {
	MplID          uuid.UUID       `db:"mpl_id"`
	MaterialID     string          `db:"material_id"`
	BOQID          uuid.UUID       `db:"boq_id"`
	SupplierID     uuid.UUID       `db:"supplier_id"`
	ActualPrice    sql.NullFloat64 `db:"actual_price"`
	EstimatedPrice sql.NullFloat64 `db:"estimated_price"`
	UpdatedAt      sql.NullTime    `db:"updated_at"`
}
