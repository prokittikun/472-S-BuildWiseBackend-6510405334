package responses

type MaterialResponse struct {
	MaterialID string `json:"material_id"`
	Name       string `json:"name"`
	Unit       string `json:"unit"`
}

type MaterialListResponse struct {
	Materials []MaterialResponse `json:"materials"`
}
