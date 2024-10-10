package responses

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type NullableUUID struct {
	UUID  uuid.UUID
	Valid bool
}

func (nu NullableUUID) MarshalJSON() ([]byte, error) {
	if !nu.Valid || nu.UUID == uuid.Nil {
		return []byte("null"), nil
	}
	return json.Marshal(nu.UUID.String())
}

func (nu *NullableUUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		if string(data) == "null" {
			nu.Valid = false
			return nil
		}
		return err
	}
	if s == "" {
		nu.Valid = false
		return nil
	}
	u, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	nu.UUID = u
	nu.Valid = true
	return nil
}

type CreateProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ContractURL string    `json:"contract_url"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type UpdateProjectResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ContractURL string    `json:"contract_url"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type ProjectResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	ContractURL string         `json:"contract_url"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     time.Time      `json:"end_date"`
	QuotationID NullableUUID   `json:"quotation_id"`
	ContractID  NullableUUID   `json:"contract_id"`
	InvoiceID   NullableUUID   `json:"invoice_id"`
	BID         NullableUUID   `json:"b_id"`
	ClientID    NullableUUID   `json:"client_id"`
	Client      ClientResponse `json:"client"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
