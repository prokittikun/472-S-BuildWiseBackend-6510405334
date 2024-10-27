package repositories

import (
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type BOQRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*responses.BOQResponse, error)
	Approve(ctx context.Context, boqID uuid.UUID) error
	GetBoqWithProject(ctx context.Context, projectID uuid.UUID) (*responses.BOQResponse, error)
	AddBOQJob(ctx context.Context, boqID uuid.UUID, req requests.BOQJobRequest) error
	DeleteBOQJob(ctx context.Context, boqID uuid.UUID, jobID uuid.UUID) error
}
