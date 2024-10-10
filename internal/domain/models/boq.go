package models

import (
	"time"

	"github.com/google/uuid"
)

type BOQ struct {
	BID          uuid.UUID `db:"b_id"`
	ProjectID    uuid.UUID `db:"project_id"`
	Status       string    `db:"status"`
	CompleteStep int       `db:"complete_step"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
