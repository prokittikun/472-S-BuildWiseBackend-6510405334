package models

import (
	"database/sql"
	"encoding/json"
	"time"

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

type QuotationJob struct {
	QuotationID    uuid.UUID       `db:"quotation_id"`
	Status         string          `db:"status"`
	ValidDate      time.Time       `db:"valid_date"`
	TaxPercentage  float64         `db:"tax_percentage"`
	JobName        string          `db:"name"`
	Unit           string          `db:"unit"`
	Quantity       float64         `db:"quantity"`
	LaborCost      float64         `db:"labor_cost"`
	TotalLaborCost float64         `db:"total_labor_cost"`
	EstimatedPrice sql.NullFloat64 `db:"estimated_price"`
	TotalEstPrice  sql.NullFloat64 `db:"total_estimated_price"`
	Total          sql.NullFloat64 `db:"total"`
	SellingPrice   sql.NullFloat64 `db:"selling_price"`
}

type QuotationGeneralCost struct {
	BOQID         uuid.UUID       `db:"boq_id"`
	GID           uuid.UUID       `db:"g_id"`
	TypeName      string          `db:"type_name"`
	EstimatedCost sql.NullFloat64 `db:"estimated_cost"`
}

type QuotationExportData struct {
	ProjectID     uuid.UUID       `db:"project_id"`
	ProjectName   string          `db:"name"`
	Description   string          `db:"description"`
	Address       json.RawMessage `db:"address"`
	ClientName    string          `db:"client_name"`
	ClientAddress json.RawMessage `db:"client_address"`
	ClientEmail   string          `db:"client_email"`
	ClientTel     string          `db:"client_tel"`
	ClientTaxID   string          `db:"client_tax_id"`
	QuotationID   uuid.UUID       `db:"quotation_id"`
	ValidDate     time.Time       `db:"valid_date"`
	TaxPercentage float64         `db:"tax_percentage"`
	FinalAmount   sql.NullFloat64 `db:"final_amount"`
	Status        string          `db:"status"`
	JobDetails    []JobDetail     `db:"-"`
}

type JobDetail struct {
	Name         string  `db:"name"`
	Unit         string  `db:"unit"`
	Quantity     float64 `db:"quantity"`
	SellingPrice float64 `db:"selling_price"`
	Amount       float64 `db:"amount"`
}
