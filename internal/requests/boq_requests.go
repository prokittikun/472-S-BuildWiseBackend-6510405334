package requests

import (
	"github.com/google/uuid"
)

type CreateBOQRequest struct {
	ProjectID uuid.UUID `json:"project_id" validate:"required"`
}

type UpdateBOQRequest struct {
	Status             string  `json:"status" validate:"required,oneof=draft approved"`
	SellingGeneralCost float64 `json:"selling_general_cost" validate:"required"`
}
