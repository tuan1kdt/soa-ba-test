package service

import (
	"context"

	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
	"github.com/tuan1kdt/soa-ba-test/internal/core/port"
)

/**
 * StatisticService implements port.StatisticService interface
 * and provides an access to the statistic repository
 */
type StatisticService struct {
	repo port.StatisticService
}

// NewStatisticService creates a new category service instance
func NewStatisticService(repo port.StatisticRepository) *StatisticService {
	return &StatisticService{
		repo,
	}
}

func (ss *StatisticService) StatisticSupplierProduct(ctx context.Context) ([]*domain.StatisticSupplierProduct, error) {
	return ss.repo.StatisticSupplierProduct(ctx)
}

func (ss *StatisticService) StatisticCategoryProduct(ctx context.Context) ([]*domain.StatisticCategoryProduct, error) {
	return ss.repo.StatisticCategoryProduct(ctx)
}
