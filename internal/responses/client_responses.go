package responses

import (
	"boonkosang/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type ClientResponse struct {
	ClientID      uuid.UUID `json:"client_id"`
	CompanyName   string    `json:"company_name"`
	ContactPerson string    `json:"contact_person"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone,omitempty"`
	Address       string    `json:"address,omitempty"`
	TaxID         string    `json:"tax_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func NewClientResponse(client *models.Client) ClientResponse {
	return ClientResponse{
		ClientID:      client.ClientID,
		CompanyName:   client.CompanyName,
		ContactPerson: client.ContactPerson,
		Email:         client.Email,
		Phone:         client.Phone,
		Address:       client.Address,
		TaxID:         client.TaxID,
		CreatedAt:     client.CreatedAt,
		UpdatedAt:     client.UpdatedAt,
	}
}

type CreateClientResponse ClientResponse
type UpdateClientResponse ClientResponse

type ListClientsResponse struct {
	Clients []ClientResponse `json:"clients"`
}

func NewListClientsResponse(clients []*models.Client) ListClientsResponse {
	response := ListClientsResponse{
		Clients: make([]ClientResponse, len(clients)),
	}
	for i, client := range clients {
		response.Clients[i] = NewClientResponse(client)
	}
	return response
}
