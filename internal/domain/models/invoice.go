package models

import (
	"time"

	"github.com/google/uuid"
)

type Invoice struct {
	InvoiceID uuid.UUID `db:"invoice_id"`
	ProjectID uuid.UUID `db:"project_id"`
	FileURL   string    `db:"file_url"`
	Amount    float64   `db:"amount"`
	Status    string    `db:"status"`
	DueDate   time.Time `db:"due_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
