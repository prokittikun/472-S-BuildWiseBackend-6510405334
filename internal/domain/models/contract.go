package models

import (
	"time"

	"github.com/google/uuid"
)

type Contract struct {
	ContractID uuid.UUID `db:"contract_id"`
	ProjectID  uuid.UUID `db:"project_id"`
	FileURL    string    `db:"file_url"`
	SignedDate time.Time `db:"signed_date"`
	Status     string    `db:"status"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
