package responses

import (
	"boonkosang/internal/domain/models"

	"github.com/google/uuid"
)

type BOQResponse struct {
	ID                 uuid.UUID        `json:"id"`
	ProjectID          uuid.UUID        `json:"project_id"`
	Status             models.BOQStatus `json:"status"`
	SellingGeneralCost float64          `json:"selling_general_cost"`
	Project            ProjectResponse  `json:"project,omitempty"`
}

type BOQListResponse struct {
	BOQs  []BOQResponse `json:"boqs"`
	Total int64         `json:"total"`
}
