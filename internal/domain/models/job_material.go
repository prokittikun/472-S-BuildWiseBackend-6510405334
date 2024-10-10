package models

import (
	"github.com/google/uuid"
)

type JobMaterial struct {
	JobID         uuid.UUID `db:"job_id"`
	MaterialName  string    `db:"material_name"`
	Quantity      float64   `db:"quantity"`
	Type          string    `db:"type"`
	UnitOfMeasure string    `db:"unit_of_measure"`
}
