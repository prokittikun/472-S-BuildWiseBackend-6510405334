package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type Job struct {
	JobID       uuid.UUID      `db:"job_id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	Unit        string         `db:"unit"`
}
