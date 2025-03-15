package port

import (
	"context"

	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
)

//go:generate mockgen -source=statistic.go -destination=mock/statistic.go -package=mock

// StatisticRepository is an interface for interacting with statistic-related data
type StatisticRepository interface {
	StatisticSupplierProduct(ctx context.Context) ([]*domain.StatisticSupplierProduct, error)
	StatisticCategoryProduct(ctx context.Context) ([]*domain.StatisticCategoryProduct, error)
}

// StatisticService is an interface for interacting with statistic-related business logic
type StatisticService interface {
	StatisticSupplierProduct(ctx context.Context) ([]*domain.StatisticSupplierProduct, error)
	StatisticCategoryProduct(ctx context.Context) ([]*domain.StatisticCategoryProduct, error)
}
