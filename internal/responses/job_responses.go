package responses

import (
	"github.com/google/uuid"
)

type JobResponse struct {
	JobID       uuid.UUID `json:"job_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Unit        string    `json:"unit"`
}

type JobMaterialResponse struct {
	JobID       uuid.UUID         `json:"job_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Unit        string            `json:"unit"`
	Materials   []JobMaterialItem `json:"materials"`
}

type JobMaterialItem struct {
	MaterialID string  `json:"material_id" db:"material_id"`
	Name       string  `json:"name" db:"name"`
	Unit       string  `json:"unit" db:"unit"`
	Quantity   float64 `json:"quantity" db:"quantity"`
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

// ApiResponse represents a standard API response structure
type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
