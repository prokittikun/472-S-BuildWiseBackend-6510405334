package models

import (
	"github.com/google/uuid"
)

type BOQJob struct {
	BOQID        uuid.UUID `db:"boq_id"`
	JobID        uuid.UUID `db:"job_id"`
	Quantity     int       `db:"quantity"`
	LaborCost    float64   `db:"labor_cost"`
	SellingPrice float64   `db:"selling_price"`
}
