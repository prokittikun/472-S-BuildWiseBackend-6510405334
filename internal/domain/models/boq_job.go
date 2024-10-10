package models

import (
	"github.com/google/uuid"
)

type BOQJob struct {
	BID       uuid.UUID `db:"b_id"`
	JobID     uuid.UUID `db:"job_id"`
	Unit      int       `db:"unit"`
	LaborCost float64   `db:"labor_cost"`
}
