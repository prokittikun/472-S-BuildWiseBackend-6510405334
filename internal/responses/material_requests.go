package responses

import (
	"time"
)

type MaterialResponse struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateMaterialResponse MaterialResponse
type UpdateMaterialResponse MaterialResponse

type MaterialPriceLogResponse struct {
	MaterialName   string    `json:"material_name"`
	ActualPrice    float64   `json:"actual_price"`
	SalePrice      float64   `json:"sale_price"`
	EstimatedPrice float64   `json:"estimated_price"`
	CreatedAt      time.Time `json:"created_at"`
}
