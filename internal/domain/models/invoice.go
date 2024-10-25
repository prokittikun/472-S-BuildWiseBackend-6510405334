package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	InvoiceID uuid.UUID      `db:"invoice_id"`
	ProjectID uuid.UUID      `db:"project_id"`
	FileURL   sql.NullString `db:"file_url"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}
