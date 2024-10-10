package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	JobID       uuid.UUID     `db:"job_id"`
	Description string        `db:"description"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	Materials   []JobMaterial `db:"-"` // ไม่ได้มาจาก database โดยตรง
}
