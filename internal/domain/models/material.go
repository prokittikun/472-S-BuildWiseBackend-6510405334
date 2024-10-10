package models

import (
	"time"
)

type Material struct {
	Name          string    `db:"name"`
	Type          string    `db:"type"`
	UnitOfMeasure string    `db:"unit_of_measure"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
