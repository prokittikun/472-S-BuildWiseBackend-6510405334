package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type GeneralCost struct {
	GID           uuid.UUID       `db:"g_id"`
	BOQID         uuid.UUID       `db:"boq_id"`
	TypeName      string          `db:"type_name"`
	ActualCost    sql.NullFloat64 `db:"actual_cost"`
	EstimatedCost sql.NullFloat64 `db:"estimated_cost"`
}
