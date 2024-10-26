package requests

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateProjectRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Address     json.RawMessage `json:"address" validate:"required"`
	ClientID    uuid.UUID       `json:"client_id" validate:"required"`
}

type UpdateProjectRequest struct {
	Name        string          `json:"name" validate:"required"`
	Description string          `json:"description" validate:"required"`
	Address     json.RawMessage `json:"address" validate:"required"`
	ClientID    uuid.UUID       `json:"client_id" validate:"required"`
}
