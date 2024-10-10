package models

import (
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	SupplierID    uuid.UUID `db:"supplier_id"`
	Name          string    `db:"name"`
	ContactPerson string    `db:"contact_person"`
	Email         string    `db:"email"`
	Phone         string    `db:"phone"`
	Address       string    `db:"address"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
