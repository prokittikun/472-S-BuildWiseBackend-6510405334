package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type QuotationStatus string

const (
	QuotationStatusDraft    QuotationStatus = "draft"
	QuotationStatusApproved QuotationStatus = "approved"
)

type Quotation struct {
	QuotationID   uuid.UUID       `db:"quotation_id"`
	ProjectID     uuid.UUID       `db:"project_id"`
	ValidDate     sql.NullTime    `db:"valid_date"`
	Status        QuotationStatus `db:"status"`
	FinalAmount   sql.NullFloat64 `db:"final_amount"`
	TaxPercentage sql.NullFloat64 `db:"tax_percentage"`
}
