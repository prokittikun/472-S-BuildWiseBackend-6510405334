package usecase

import (
	"boonkosang/internal/repositories"
	"boonkosang/internal/responses"
	"context"

	"github.com/google/uuid"
)

type BOQUsecase interface {
	GetBoqWithProject(ctx context.Context, project_id uuid.UUID) (*responses.BOQResponse, error)
}

type boqUsecase struct {
	boqRepo     repositories.BOQRepository
	projectRepo repositories.ProjectRepository
}

func NewBOQUsecase(boqRepo repositories.BOQRepository, projectRepo repositories.ProjectRepository) BOQUsecase {
	return &boqUsecase{
		boqRepo:     boqRepo,
		projectRepo: projectRepo,
	}
}
func (u *boqUsecase) GetBoqWithProject(ctx context.Context, project_id uuid.UUID) (*responses.BOQResponse, error) {
	return u.boqRepo.GetBoqWithProject(ctx, project_id)
}
