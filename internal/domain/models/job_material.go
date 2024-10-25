package models

import (
	"github.com/google/uuid"
)

type JobMaterial struct {
	JobID      uuid.UUID `db:"job_id"`
	MaterialID string    `db:"material_id"`
	Quantity   int32     `db:"quantity"`
}
