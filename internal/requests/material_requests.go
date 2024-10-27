package requests

type CreateMaterialRequest struct {
	Name string `json:"name" validate:"required"`
	Unit string `json:"unit" validate:"required"`
}

type UpdateMaterialRequest struct {
	Name string `json:"name" validate:"required"`
	Unit string `json:"unit" validate:"required"`
}

type UpdateMaterialEstimatedPriceRequest struct {
	MaterialID     string  `json:"material_id" validate:"required"`
	EstimatedPrice float64 `json:"estimated_price" validate:"required,gt=0"`
}
