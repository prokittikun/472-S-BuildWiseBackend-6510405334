package usecase

import (
	"boonkosang/internal/repositories"
	"boonkosang/internal/requests"
	"boonkosang/internal/responses"
	"context"
	"errors"

	"github.com/google/uuid"
)

type SupplierUsecase interface {
	Create(ctx context.Context, req requests.CreateSupplierRequest) (*responses.SupplierResponse, error)
	Update(ctx context.Context, id uuid.UUID, req requests.UpdateSupplierRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*responses.SupplierResponse, error)
	List(ctx context.Context, page, pageSize int) (*responses.SupplierListResponse, error)
}

type supplierUsecase struct {
	supplierRepo repositories.SupplierRepository
}

func NewSupplierUsecase(supplierRepo repositories.SupplierRepository) SupplierUsecase {
	return &supplierUsecase{
		supplierRepo: supplierRepo,
	}
}

func (u *supplierUsecase) Create(ctx context.Context, req requests.CreateSupplierRequest) (*responses.SupplierResponse, error) {
	existing, err := u.supplierRepo.GetByEmail(ctx, req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("supplier with this email already exists")
	}

	supplier, err := u.supplierRepo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return &responses.SupplierResponse{
		ID:      supplier.SupplierID,
		Name:    supplier.Name,
		Email:   supplier.Email,
		Tel:     supplier.Tel,
		Address: supplier.Address,
	}, nil
}

func (u *supplierUsecase) Update(ctx context.Context, id uuid.UUID, req requests.UpdateSupplierRequest) error {
	existing, err := u.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if existing.Email != req.Email {
		supplier, err := u.supplierRepo.GetByEmail(ctx, req.Email)
		if err == nil && supplier != nil {
			return errors.New("supplier with this email already exists")
		}
	}

	return u.supplierRepo.Update(ctx, id, req)
}

func (u *supplierUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return u.supplierRepo.Delete(ctx, id)
}

func (u *supplierUsecase) GetByID(ctx context.Context, id uuid.UUID) (*responses.SupplierResponse, error) {
	supplier, err := u.supplierRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &responses.SupplierResponse{
		ID:      supplier.SupplierID,
		Name:    supplier.Name,
		Email:   supplier.Email,
		Tel:     supplier.Tel,
		Address: supplier.Address,
	}, nil
}

func (u *supplierUsecase) List(ctx context.Context, page, pageSize int) (*responses.SupplierListResponse, error) {

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	suppliers, total, err := u.supplierRepo.List(ctx, pageSize, offset)
	if err != nil {
		return nil, err
	}

	supplierResponses := make([]responses.SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		supplierResponses[i] = responses.SupplierResponse{
			ID:      supplier.SupplierID,
			Name:    supplier.Name,
			Email:   supplier.Email,
			Tel:     supplier.Tel,
			Address: supplier.Address,
		}
	}

	return &responses.SupplierListResponse{
		Suppliers: supplierResponses,
		Total:     total,
	}, nil
}
