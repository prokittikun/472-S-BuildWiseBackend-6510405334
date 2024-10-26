package responses

import (
	"time"

	"github.com/google/uuid"
)

type JobResponse struct {
	JobID       uuid.UUID `json:"job_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
}

type JobListResponse struct {
	Jobs []JobResponse `json:"jobs"`
}

type PaginationResponse struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	TotalPages   int `json:"total_pages"`
	TotalRecords int `json:"total_records"`
}

// JobMaterialResponse represents a response after adding/updating materials
type JobMaterialResponse struct {
	JobID        uuid.UUID `json:"job_id"`
	MaterialID   string    `json:"material_id"`
	MaterialName string    `json:"material_name"`
	Unit         string    `json:"unit"`
	Quantity     float64   `json:"quantity"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ApiResponse represents a standard API response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
