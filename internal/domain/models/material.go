package models

type Material struct {
	MaterialID string `db:"material_id"`
	Name       string `db:"name"`
	Unit       string `db:"unit"`
}
