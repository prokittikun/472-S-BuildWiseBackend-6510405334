package repositories

import (
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type BOQRepository interface {
	GetBoqWithProject(ctx context.Context, projectID uuid.UUID) (*responses.BOQResponse, error)
}
