package requests

import "github.com/google/uuid"

type JobMaterialItem struct {
	MaterialID string  `json:"material_id" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required,gt=0"`
}

type CreateJobRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Unit        string            `json:"unit" validate:"required"`
	Materials   []JobMaterialItem `json:"materials" validate:"dive"`
}

type UpdateJobRequest struct {
	Name        string            `json:"name" validate:"required"`
	Description string            `json:"description"`
	Unit        string            `json:"unit" validate:"required"`
	Materials   []JobMaterialItem `json:"materials" validate:"dive"`
}
type BOQJobRequest struct {
	JobID        uuid.UUID `json:"job_id" validate:"required"`
	Quantity     int       `json:"quantity" validate:"required,gt=0"`
	LaborCost    float64   `json:"labor_cost" validate:"required,gte=0"`
	SellingPrice float64   `json:"selling_price" validate:"required,gte=0"`
}
