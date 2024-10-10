package models

import (
	"time"

	"github.com/google/uuid"
)

type PurchaseOrder struct {
	POID        uuid.UUID `db:"po_id"`
	ProjectID   uuid.UUID `db:"project_id"`
	SupplierID  uuid.UUID `db:"supplier_id"`
	OrderDate   time.Time `db:"order_date"`
	TotalAmount float64   `db:"total_amount"`
	Status      string    `db:"status"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
