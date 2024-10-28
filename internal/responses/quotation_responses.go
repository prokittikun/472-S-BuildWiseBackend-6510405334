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
	Summary     QuotationSummary     `json:"summary"`
}

type QuotationJobDetail struct {
	Name         string  `json:"name"`
	Unit         string  `json:"unit"`
	Quantity     float64 `json:"quantity"`
	LaborCost    float64 `json:"labor_cost"`
	MaterialCost float64 `json:"material_cost"`
	TotalCost    float64 `json:"total_cost"`
	SellingPrice float64 `json:"selling_price,omitempty"`
}

type GeneralCostDetail struct {
	TypeName      string  `json:"type_name"`
	EstimatedCost float64 `json:"estimated_cost"`
}

type QuotationSummary struct {
	TotalLaborCost    float64 `json:"total_labor_cost"`
	TotalMaterialCost float64 `json:"total_material_cost"`
	TotalGeneralCost  float64 `json:"total_general_cost"`
	SubTotal          float64 `json:"subtotal"`
	Tax               float64 `json:"tax"`
	Total             float64 `json:"total"`
}
