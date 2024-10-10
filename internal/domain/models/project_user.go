package models

import (
	"github.com/google/uuid"
)

type ProjectUser struct {
	ProjectID uuid.UUID `db:"project_id"`
	UserID    uuid.UUID `db:"user_id"`
}
