package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"context"
	"time"
)

type MaterialUsecase interface {
	CreateMaterial(ctx context.Context, req requests.CreateMaterialRequest) (*models.Material, error)
	ListMaterials(ctx context.Context) ([]*models.Material, error)
	GetMaterial(ctx context.Context, name string) (*models.Material, error)
	UpdateMaterial(ctx context.Context, name string, req requests.UpdateMaterialRequest) (*models.Material, error)
	DeleteMaterial(ctx context.Context, name string) error
	GetMaterialPriceHistory(ctx context.Context, name string) ([]*models.MaterialPriceLog, error)
}

type materialUsecase struct {
	materialRepo repositories.MaterialRepository
}

func NewMaterialUsecase(materialRepo repositories.MaterialRepository) MaterialUsecase {
	return &materialUsecase{
		materialRepo: materialRepo,
	}
}

func (mu *materialUsecase) CreateMaterial(ctx context.Context, req requests.CreateMaterialRequest) (*models.Material, error) {
	material := &models.Material{
		Name:          req.Name,
		Type:          req.Type,
		UnitOfMeasure: req.UnitOfMeasure,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := mu.materialRepo.CreateMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return material, nil
}

func (mu *materialUsecase) ListMaterials(ctx context.Context) ([]*models.Material, error) {
	return mu.materialRepo.ListMaterials(ctx)
}

func (mu *materialUsecase) GetMaterial(ctx context.Context, name string) (*models.Material, error) {
	return mu.materialRepo.GetMaterialByName(ctx, name)
}

func (mu *materialUsecase) UpdateMaterial(ctx context.Context, name string, req requests.UpdateMaterialRequest) (*models.Material, error) {
	material, err := mu.materialRepo.GetMaterialByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Type != "" {
		material.Type = req.Type
	}
	if req.UnitOfMeasure != "" {
		material.UnitOfMeasure = req.UnitOfMeasure
	}
	material.UpdatedAt = time.Now()

	err = mu.materialRepo.UpdateMaterial(ctx, material)
	if err != nil {
		return nil, err
	}

	return material, nil
}

func (mu *materialUsecase) DeleteMaterial(ctx context.Context, name string) error {
	return mu.materialRepo.DeleteMaterial(ctx, name)
}

func (mu *materialUsecase) GetMaterialPriceHistory(ctx context.Context, name string) ([]*models.MaterialPriceLog, error) {
	return mu.materialRepo.GetMaterialPriceHistory(ctx, name)
}
