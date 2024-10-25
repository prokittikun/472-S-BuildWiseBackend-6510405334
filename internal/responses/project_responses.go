package responses

import (
	"boonkosang/internal/domain/models"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ProjectResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Address     json.RawMessage      `json:"address"`
	Status      models.ProjectStatus `json:"status"`
	ClientID    uuid.UUID            `json:"client_id"`
	Client      *ClientResponse      `json:"client,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type ProjectListResponse struct {
	Projects []ProjectResponse `json:"projects"`
	Total    int64             `json:"total"`
}
