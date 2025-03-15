package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
	"github.com/tuan1kdt/soa-ba-test/internal/core/port"
)

/**
 * CategoryService implements port.CategoryService interface
 * and provides an access to the category repository
 * and cache service
 */
type CategoryService struct {
	repo  port.CategoryRepository
	cache port.CacheRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(repo port.CategoryRepository, cache port.CacheRepository) *CategoryService {
	return &CategoryService{
		repo,
		cache,
	}
}

// CreateCategory creates a new category
func (cs *CategoryService) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	category.ID = uuid.New()
	category, err := cs.repo.CreateCategory(ctx, category)
	if err != nil {
		if errors.Is(err, domain.ErrConflictingData) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}
	return category, nil
}

// GetCategory retrieves a category by id
func (cs *CategoryService) GetCategory(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	category, err := cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return category, nil
}

// ListCategories retrieves a list of categories
func (cs *CategoryService) ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error) {
	categories, err := cs.repo.ListCategories(ctx, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	return categories, nil
}

// UpdateCategory updates a category
func (cs *CategoryService) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	if category.Name == "" {
		return nil, domain.ErrNoUpdatedData
	}
	_, err := cs.repo.GetCategoryByID(ctx, category.ID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	_, err = cs.repo.UpdateCategory(ctx, category)
	if err != nil {
		if err == domain.ErrConflictingData {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return category, nil
}

// DeleteCategory deletes a category
func (cs *CategoryService) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	_, err := cs.repo.GetCategoryByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}
	return cs.repo.DeleteCategory(ctx, id)
}
