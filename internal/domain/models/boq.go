package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type BOQStatus string

const (
	BOQStatusDraft    BOQStatus = "draft"
	BOQStatusApproved BOQStatus = "approved"
)

type BOQ struct {
	BOQID              uuid.UUID       `db:"boq_id"`
	ProjectID          uuid.UUID       `db:"project_id"`
	Status             BOQStatus       `db:"status"`
	SellingGeneralCost sql.NullFloat64 `db:"selling_general_cost"`
}
