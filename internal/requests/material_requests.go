package requests

type CreateMaterialRequest struct {
	Name string `json:"name" validate:"required"`
	Unit string `json:"unit" validate:"required"`
}

type UpdateMaterialRequest struct {
	Name string `json:"name" validate:"required"`
	Unit string `json:"unit" validate:"required"`
}
