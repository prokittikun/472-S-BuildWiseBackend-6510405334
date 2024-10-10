package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ProjectID   uuid.UUID `db:"project_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Status      string    `db:"status"`
	ContractURL string    `db:"contract_url"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date"`
	QuotationID uuid.UUID `db:"quotation_id"`
	ContractID  uuid.UUID `db:"contract_id"`
	InvoiceID   uuid.UUID `db:"invoice_id"`
	BID         uuid.UUID `db:"b_id"`
	ClientID    uuid.UUID `db:"client_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
