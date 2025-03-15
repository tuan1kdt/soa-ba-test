package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
	"github.com/tuan1kdt/soa-ba-test/internal/core/port"
	"github.com/tuan1kdt/soa-ba-test/internal/core/util"
)

/**
 * ProductService implements port.ProductService and port.CategoryService
 * interfaces and provides an access to the product and category repositories
 * and cache service
 */
type ProductService struct {
	productRepo  port.ProductRepository
	categoryRepo port.CategoryRepository
	cache        port.CacheRepository
	geoClient    port.GeoClient
}

func (ps *ProductService) ListProducts2(ctx context.Context, search string, categoryIDs []uuid.UUID, cursor *string, perPage uint64) ([]domain.Product, error) {
	//TODO implement me
	panic("implement me")
}

// NewProductService creates a new product service instance
func NewProductService(productRepo port.ProductRepository, categoryRepo port.CategoryRepository, cache port.CacheRepository, geoClient port.GeoClient) *ProductService {
	return &ProductService{
		productRepo,
		categoryRepo,
		cache,
		geoClient,
	}
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	if product.CategoryID != nil {
		_, err := ps.categoryRepo.GetCategoryByID(ctx, *product.CategoryID)
		if err != nil {
			if errors.Is(err, domain.ErrDataNotFound) {
				return nil, err
			}
			return nil, domain.ErrInternal
		}
	}

	id := uuid.New()
	product.ID = id

	product.AddedDate = time.Now()

	product, err := ps.productRepo.CreateProduct(ctx, product)
	if err != nil {
		if errors.Is(err, domain.ErrConflictingData) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	//cacheKey := util.GenerateCacheKey("product", product.ID)
	//productSerialized, err := util.Serialize(product)
	//if err != nil {
	//	return nil, domain.ErrInternal
	//}
	//
	//err = ps.cache.Set(ctx, cacheKey, productSerialized, 0)
	//if err != nil {
	//	return nil, domain.ErrInternal
	//}
	//
	//err = ps.cache.DeleteByPrefix(ctx, "products:*")
	//if err != nil {
	//	return nil, domain.ErrInternal
	//}

	return product, nil
}

// GetProduct retrieves a product by id
func (ps *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product *domain.Product

	product, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	category, err := ps.categoryRepo.GetCategoryByID(ctx, *product.CategoryID)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	product.Category = category

	return product, nil
}

// ListProducts retrieves a list of products
func (ps *ProductService) ListProducts(ctx context.Context, search string, categoryIds []uuid.UUID, skip, limit uint64) ([]domain.Product, error) {
	products, err := ps.productRepo.ListProducts(ctx, search, categoryIds, skip, limit)
	if err != nil {
		return nil, domain.ErrInternal
	}

	categories := make(map[uuid.UUID]domain.Category)
	for i, product := range products {
		category := &domain.Category{}
		if _, ok := categories[*product.CategoryID]; !ok && product.CategoryID != nil {
			category, err = ps.categoryRepo.GetCategoryByID(ctx, *product.CategoryID)
			if err != nil {
				slog.Error("Error getting category by id", "error", err)
			}
		}

		products[i].Category = category
	}

	//productsSerialized, err := util.Serialize(products)
	//if err != nil {
	//	return nil, domain.ErrInternal
	//}

	//err = ps.cache.Set(ctx, cacheKey, productsSerialized, 0)
	//if err != nil {
	//	return nil, domain.ErrInternal
	//}

	return products, nil
}

func (ps *ProductService) GetProductDistance(ctx context.Context, ip string, id uuid.UUID) (float64, error) {
	product, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return 0, err
		}
		return 0, domain.ErrInternal
	}

	if product.StockCity == "" {
		return 0, domain.ErrDataNotFound
	}

	dstLat, dstLon, err := ps.geoClient.GetCityLocation(product.StockCity)
	if err != nil {
		return 0, domain.ErrInternal
	}

	srcLat, srcLon, err := ps.geoClient.GetIPLocation(ip)
	if err != nil {
		return 0, domain.ErrInternal
	}

	distance := ps.geoClient.GetDistance(srcLat, srcLon, dstLat, dstLon)

	return distance, nil
}

// UpdateProduct updates a product
func (ps *ProductService) UpdateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	_, err := ps.productRepo.GetProductByID(ctx, product.ID)
	if err != nil {
		if errors.Is(err, domain.ErrDataNotFound) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	var updatedFields []string
	if product.Reference != "" {
		updatedFields = append(updatedFields, "reference")
	}
	if product.Name != "" {
		updatedFields = append(updatedFields, "name")
	}
	if product.Price != 0 {
		updatedFields = append(updatedFields, "price")
	}
	if product.Quantity != 0 {
		updatedFields = append(updatedFields, "quantity")
	}
	if product.StockCity != "" {
		updatedFields = append(updatedFields, "stock_city")
	}
	if product.CategoryID != nil {
		_, err := ps.categoryRepo.GetCategoryByID(ctx, *product.CategoryID)
		if err != nil {
			if errors.Is(err, domain.ErrDataNotFound) {
				return nil, err
			}
			return nil, domain.ErrInternal
		}
		updatedFields = append(updatedFields, "category_id")
	}

	if product.Status != domain.StatusUnknown {
		updatedFields = append(updatedFields, "status")
	}

	if len(updatedFields) == 0 {
		return nil, domain.ErrNoUpdatedData
	}

	_, err = ps.productRepo.UpdateProduct(ctx, product, updatedFields...)
	if err != nil {
		if errors.Is(err, domain.ErrConflictingData) {
			return nil, err
		}
		return nil, domain.ErrInternal
	}

	return product, nil
}

// DeleteProduct deletes a product
func (ps *ProductService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	_, err := ps.productRepo.GetProductByID(ctx, id)
	if err != nil {
		if err == domain.ErrDataNotFound {
			return err
		}
		return domain.ErrInternal
	}

	cacheKey := util.GenerateCacheKey("product", id)

	err = ps.cache.Delete(ctx, cacheKey)
	if err != nil {
		return domain.ErrInternal
	}

	err = ps.cache.DeleteByPrefix(ctx, "products:*")
	if err != nil {
		return domain.ErrInternal
	}

	return ps.productRepo.DeleteProduct(ctx, id)
}
