package requests

import (
	"time"
)

type CreateProjectRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ContractURL string    `json:"contract_url"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type UpdateProjectRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ContractURL string    `json:"contract_url"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}
