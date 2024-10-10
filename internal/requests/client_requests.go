package requests

type CreateClientRequest struct {
	CompanyName   string `json:"company_name" validate:"required"`
	ContactPerson string `json:"contact_person" validate:"required"`
	Email         string `json:"email" validate:"required,email"`
	Phone         string `json:"phone" validate:"omitempty"`
	Address       string `json:"address" validate:"omitempty"`
	TaxID         string `json:"tax_id" validate:"omitempty"`
}

type UpdateClientRequest struct {
	CompanyName   string `json:"company_name" validate:"omitempty"`
	ContactPerson string `json:"contact_person" validate:"omitempty"`
	Email         string `json:"email" validate:"omitempty,email"`
	Phone         string `json:"phone" validate:"omitempty"`
	Address       string `json:"address" validate:"omitempty"`
	TaxID         string `json:"tax_id" validate:"omitempty"`
}
