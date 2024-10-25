package models

import (
	"database/sql"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID      `db:"user_id"`
	Username  string         `db:"username"`
	Password  string         `db:"password"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	Email     sql.NullString `db:"email"`
	Tel       sql.NullString `db:"tel"`
	CompanyID *uuid.UUID     `db:"company_id"`
}
