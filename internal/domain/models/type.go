package models

import (
	"time"
)

type Type struct {
	TypeName    string    `db:"type_name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
