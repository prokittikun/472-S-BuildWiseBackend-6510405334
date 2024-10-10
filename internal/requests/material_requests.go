package requests

type CreateMaterialRequest struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	UnitOfMeasure string `json:"unit_of_measure"`
}

type UpdateMaterialRequest struct {
	Type          string `json:"type"`
	UnitOfMeasure string `json:"unit_of_measure"`
}
