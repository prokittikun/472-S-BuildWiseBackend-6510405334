package responses

import (
	"time"

	"github.com/google/uuid"
)

type QuotationResponse struct {
	QuotationID uuid.UUID            `json:"quotation_id"`
	Status      string               `json:"status"`
	ValidDate   time.Time            `json:"valid_date"`
	Jobs        []QuotationJobDetail `json:"jobs"`
	Costs       []GeneralCostDetail  `json:"general_costs"`
}

type QuotationJobDetail struct {
	Name               string  `json:"name"`
	Unit               string  `json:"unit"`
	Quantity           float64 `json:"quantity"`
	LaborCost          float64 `json:"labor_cost"`
	SellingPrice       float64 `json:"selling_price"`
	TotalMaterialPrice float64 `json:"total_material_price"`
	Total              float64 `json:"total"`
	OverallCost        float64 `json:"overall_cost"`
	TotalSellingPrice  float64 `json:"total_selling_price"`
}

type GeneralCostDetail struct {
	TypeName      string  `json:"type_name"`
	EstimatedCost float64 `json:"estimated_cost"`
}
