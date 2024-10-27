// usecase/material_usecase.go
package usecase

import (
	"boonkosang/internal/domain/models"
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type MaterialUsecase interface {
	Create(ctx context.Context, req requests.CreateMaterialRequest) (*responses.MaterialResponse, error)
	Update(ctx context.Context, materialID string, req requests.UpdateMaterialRequest) error
	Delete(ctx context.Context, materialID string) error
	GetByID(ctx context.Context, materialID string) (*responses.MaterialResponse, error)
	List(ctx context.Context) (*responses.MaterialListResponse, error)

	GetMaterialPrices(ctx context.Context, projectID uuid.UUID) (*responses.MaterialPriceListResponse, error)
	UpdateEstimatedPrice(ctx context.Context, boqID uuid.UUID, req requests.UpdateMaterialEstimatedPriceRequest) error
}

type materialUsecase struct {
	materialRepo repositories.MaterialRepository
	supplierRepo repositories.SupplierRepository
}

func NewMaterialUsecase(
	materialRepo repositories.MaterialRepository,
	supplierRepo repositories.SupplierRepository,
) MaterialUsecase {
	return &materialUsecase{
		materialRepo: materialRepo,
		supplierRepo: supplierRepo,
	}
}

func (u *materialUsecase) Create(ctx context.Context, req requests.CreateMaterialRequest) (*responses.MaterialResponse, error) {

	material, err := u.materialRepo.Create(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create material: %w", err)
	}

	return u.createMaterialResponse(ctx, material)
}

func (u *materialUsecase) Update(ctx context.Context, materialID string, req requests.UpdateMaterialRequest) error {
	existing, err := u.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("material not found")
	}

	return u.materialRepo.Update(ctx, materialID, req)
}

func (u *materialUsecase) Delete(ctx context.Context, materialID string) error {
	existing, err := u.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("material not found")
	}

	return u.materialRepo.Delete(ctx, materialID)
}

func (u *materialUsecase) GetByID(ctx context.Context, materialID string) (*responses.MaterialResponse, error) {
	material, err := u.materialRepo.GetByID(ctx, materialID)
	if err != nil {
		return nil, err
	}
	if material == nil {
		return nil, errors.New("material not found")
	}

	return u.createMaterialResponse(ctx, material)
}

func (u *materialUsecase) List(ctx context.Context) (*responses.MaterialListResponse, error) {
	materials, err := u.materialRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list materials: %w", err)
	}

	var materialResponses []*responses.MaterialResponse
	for _, material := range materials {
		materialResponse, err := u.createMaterialResponse(ctx, &material)
		if err != nil {
			return nil, fmt.Errorf("failed to create material response: %w", err)
		}
		materialResponses = append(materialResponses, materialResponse)
	}

	materialValues := make([]responses.MaterialResponse, len(materialResponses))
	for i, materialResponse := range materialResponses {
		materialValues[i] = *materialResponse
	}

	return &responses.MaterialListResponse{
		Materials: materialValues,
	}, nil

}

func (u *materialUsecase) createMaterialResponse(ctx context.Context, material *models.Material) (*responses.MaterialResponse, error) {

	return &responses.MaterialResponse{
		MaterialID: material.MaterialID,
		Name:       material.Name,
		Unit:       material.Unit,
	}, nil
}

func (u *materialUsecase) GetMaterialPrices(ctx context.Context, projectID uuid.UUID) (*responses.MaterialPriceListResponse, error) {
	materials, err := u.materialRepo.GetMaterialPricesByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var response []responses.MaterialPriceDetail
	var totalEstimated, totalActual float64

	for _, m := range materials {
		detail := responses.MaterialPriceDetail{
			MaterialID:     m.MaterialID,
			Name:           m.Name,
			TotalQuantity:  m.TotalQuantity,
			Unit:           m.Unit,
			EstimatedPrice: m.EstimatedPrice.Float64,
			AvgActualPrice: m.AvgActualPrice.Float64,
			ActualPrice:    m.ActualPrice.Float64,
			SupplierName:   m.SupplierName.String,
		}

		totalEstimated += m.EstimatedPrice.Float64
		totalActual += m.ActualPrice.Float64

		response = append(response, detail)
	}

	return &responses.MaterialPriceListResponse{
		Materials: response,
		Totals: responses.MaterialPriceTotals{
			TotalEstimatedPrice: totalEstimated,
			TotalActualPrice:    totalActual,
		},
	}, nil
}

func (u *materialUsecase) UpdateEstimatedPrice(ctx context.Context, boqID uuid.UUID, req requests.UpdateMaterialEstimatedPriceRequest) error {
	// Check BOQ status
	status, err := u.materialRepo.GetBOQStatus(ctx, boqID)
	if err != nil {
		return err
	}

	if status != "draft" {
		return errors.New("can only update estimated prices for BOQ in draft status")
	}

	// Validate estimated price
	if req.EstimatedPrice <= 0 {
		return errors.New("estimated price must be greater than 0")
	}

	return u.materialRepo.UpdateEstimatedPrices(ctx, boqID, req.MaterialID, req.EstimatedPrice)
}
