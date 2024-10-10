package models

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ClientID      uuid.UUID `db:"client_id"`
	CompanyName   string    `db:"company_name"`
	ContactPerson string    `db:"contact_person"`
	Email         string    `db:"email"`
	Phone         string    `db:"phone"`
	Address       string    `db:"address"`
	TaxID         string    `db:"tax_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
