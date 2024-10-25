package responses

import "github.com/google/uuid"

type JobMaterialResponse struct {
	MaterialID string   `json:"material_id"`
	Name       string   `json:"name"`
	Unit       string   `json:"unit"`
	Quantity   float64  `json:"quantity"`
	LastPrice  *float64 `json:"last_price,omitempty"`
}

type JobResponse struct {
	ID          uuid.UUID             `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Unit        string                `json:"unit"`
	Materials   []JobMaterialResponse `json:"materials"`
}

type JobListResponse struct {
	Jobs  []JobResponse `json:"jobs"`
	Total int64         `json:"total"`
}

type BOQJobResponse struct {
	JobID        uuid.UUID `json:"job_id"`
	Name         string    `json:"name"`
	Unit         string    `json:"unit"`
	Quantity     int       `json:"quantity"`
	LaborCost    float64   `json:"labor_cost"`
	SellingPrice float64   `json:"selling_price"`
}
