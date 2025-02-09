package requests

import (
	"time"

	"github.com/google/uuid"
)

type CreateContractRequest struct {
	ProjectID           uuid.UUID             `json:"project_id" validate:"required"`
	ProjectDescription  string                `json:"project_description" validate:"required,min=1,max=500"`
	AreaSize            float64               `json:"area_size" validate:"required,min=0"`
	StartDate           time.Time             `json:"start_date" validate:"required,ltefield=EndDate"`
	EndDate             time.Time             `json:"end_date" validate:"required"`
	ForceMajeure        string                `json:"force_majeure" validate:"max=1000"`
	BreachOfContract    string                `json:"breach_of_contract" validate:"max=1000"`
	EndOfContract       string                `json:"end_of_contract" validate:"max=1000"`
	TerminationContract string                `json:"termination_of_contract" validate:"max=1000"`
	Amendment           string                `json:"amendment" validate:"max=1000"`
	GuaranteeWithin     int                   `json:"guarantee_within" validate:"required,min=0,max=365"`
	RetentionMoney      float64               `json:"retention_money" validate:"required,min=0"`
	PayWithin           int                   `json:"pay_within" validate:"required,min=0,max=365"`
	ValidateWithin      int                   `json:"validate_within" validate:"required,min=0,max=365"`
	Format              []string              `json:"format" validate:"required,min=1,dive,required,oneof=pdf doc docx xls xlsx dwg"`
	Periods             []CreatePeriodRequest `json:"periods" validate:"required,min=1,max=100,dive"`
}

type CreatePeriodRequest struct {
	PeriodNumber    int                      `json:"period_number" validate:"required,min=1,max=100"`
	AmountPeriod    float64                  `json:"amount_period" validate:"required,min=0"`
	DeliveredWithin int                      `json:"delivered_within" validate:"required,min=1,max=365"`
	Jobs            []CreateJobPeriodRequest `json:"jobs" validate:"required,max=100,dive"`
}

type CreateJobPeriodRequest struct {
	JobID     uuid.UUID `json:"job_id" validate:"required"`
	JobAmount float64   `json:"job_amount" validate:"required,min=0"`
}

type UpdateContractRequest struct {
	ProjectDescription  string    `json:"project_description" validate:"min=1,max=500"`
	AreaSize            float64   `json:"area_size" validate:"min=0"`
	StartDate           time.Time `json:"start_date" validate:"ltefield=EndDate"`
	EndDate             time.Time `json:"end_date"`
	ForceMajeure        string    `json:"force_majeure" validate:"max=1000"`
	BreachOfContract    string    `json:"breach_of_contract" validate:"max=1000"`
	EndOfContract       string    `json:"end_of_contract" validate:"max=1000"`
	TerminationContract string    `json:"termination_of_contract" validate:"max=1000"`
	Amendment           string    `json:"amendment" validate:"max=1000"`
	GuaranteeWithin     int       `json:"guarantee_within" validate:"min=0,max=365"`
	RetentionMoney      float64   `json:"retention_money" validate:"min=0"`
	PayWithin           int       `json:"pay_within" validate:"min=0,max=365"`
	ValidateWithin      int       `json:"validate_within" validate:"min=0,max=365"`
	Format              []string  `json:"format" validate:"min=1,dive,required,oneof=pdf doc docx xls xlsx dwg"`
}

type UpdatePeriodRequest struct {
	AmountPeriod    float64 `json:"amount_period" validate:"min=0"`
	DeliveredWithin int     `json:"delivered_within" validate:"min=1,max=365"`
}

type UpdateJobPeriodRequest struct {
	JobAmount float64 `json:"job_amount" validate:"min=0"`
}
