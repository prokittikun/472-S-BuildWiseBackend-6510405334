package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobUsecase interface {
	CreateJob(ctx context.Context, req requests.CreateJobRequest) (*models.Job, error)
	ListJobs(ctx context.Context) ([]*models.Job, error)
	GetJob(ctx context.Context, id uuid.UUID) (*models.Job, error)
	UpdateJob(ctx context.Context, id uuid.UUID, req requests.UpdateJobRequest) (*models.Job, error)
	DeleteJob(ctx context.Context, id uuid.UUID) error
}

type jobUsecase struct {
	jobRepo      repositories.JobRepository
	materialRepo repositories.MaterialRepository
}

func NewJobUsecase(jobRepo repositories.JobRepository, materialRepo repositories.MaterialRepository) JobUsecase {
	return &jobUsecase{
		jobRepo:      jobRepo,
		materialRepo: materialRepo,
	}
}

func (ju *jobUsecase) CreateJob(ctx context.Context, req requests.CreateJobRequest) (*models.Job, error) {
	// ตรวจสอบว่า material ทั้งหมดมีอยู่จริง
	for _, m := range req.Materials {
		exists, err := ju.materialRepo.MaterialExists(ctx, m.MaterialName)
		if err != nil {
			return nil, fmt.Errorf("error checking material existence: %w", err)
		}
		if !exists {
			return nil, fmt.Errorf("material '%s' does not exist", m.MaterialName)
		}
	}

	job := &models.Job{
		JobID:       uuid.New(),
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	materials := make([]models.JobMaterial, len(req.Materials))
	for i, m := range req.Materials {
		materials[i] = models.JobMaterial{
			JobID:        job.JobID,
			MaterialName: m.MaterialName,
			Quantity:     m.Quantity,
		}
	}

	err := ju.jobRepo.CreateJob(ctx, job, materials)
	if err != nil {
		return nil, fmt.Errorf("error creating job: %w", err)
	}

	job.Materials = materials
	return job, nil
}

func (ju *jobUsecase) ListJobs(ctx context.Context) ([]*models.Job, error) {
	return ju.jobRepo.ListJobs(ctx)
}

func (ju *jobUsecase) GetJob(ctx context.Context, id uuid.UUID) (*models.Job, error) {
	return ju.jobRepo.GetJobByID(ctx, id)
}

func (ju *jobUsecase) UpdateJob(ctx context.Context, id uuid.UUID, req requests.UpdateJobRequest) (*models.Job, error) {
	job, err := ju.jobRepo.GetJobByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Description != "" {
		job.Description = req.Description
	}
	job.UpdatedAt = time.Now()

	materials := make([]models.JobMaterial, len(req.Materials))
	for i, m := range req.Materials {
		materials[i] = models.JobMaterial{
			JobID:        job.JobID,
			MaterialName: m.MaterialName,
			Quantity:     m.Quantity,
		}
	}

	err = ju.jobRepo.UpdateJob(ctx, job, materials)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (ju *jobUsecase) DeleteJob(ctx context.Context, id uuid.UUID) error {
	return ju.jobRepo.DeleteJob(ctx, id)
}
