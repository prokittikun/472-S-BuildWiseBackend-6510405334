package requests

type JobMaterialRequest struct {
	MaterialName string  `json:"material_name" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,gt=0"`
}

type CreateJobRequest struct {
	Description string               `json:"description" validate:"required"`
	Materials   []JobMaterialRequest `json:"materials" validate:"required,dive"`
}

type UpdateJobRequest struct {
	Description string               `json:"description" validate:"omitempty"`
	Materials   []JobMaterialRequest `json:"materials" validate:"omitempty,dive"`
}
