package models

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectStatus string

const (
	ProjectStatusPlanning   ProjectStatus = "planning"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusCancelled  ProjectStatus = "cancelled"
)

type Project struct {
	ProjectID   uuid.UUID       `db:"project_id"`
	Name        string          `db:"name"`
	Description string          `db:"description"`
	Address     json.RawMessage `db:"address"`
	Status      ProjectStatus   `db:"status"`
	ClientID    uuid.UUID       `db:"client_id"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
}
