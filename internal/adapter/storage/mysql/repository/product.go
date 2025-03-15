package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/storage/mysql"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
)

/**
 * ProductRepository implements port.ProductRepository interface
 * and provides an access to the mysql database
 */
type ProductRepository struct {
	db *mysql.DB
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *mysql.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

// CreateProduct creates a new product record in the database
func (pr *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	tx := pr.db.WithContext(ctx)

	if err := tx.Create(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product record from the database by id
func (pr *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product

	if err := pr.db.Take(&product, id).Error; err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}
	return &product, nil
}

// ListProducts retrieves a list of products from the database
func (pr *ProductRepository) ListProducts(ctx context.Context, search string, categoryIds []uuid.UUID, skip, limit uint64) ([]domain.Product, error) {
	//var product domain.Product
	var products []domain.Product
	////
	////query := pr.db.QueryBuilder.Select("*").
	////	From("products").
	////	OrderBy("id").
	////	Limit(limit).
	////	Offset((skip - 1) * limit)
	////
	////if categoryId != 0 {
	////	query = query.Where(sq.Eq{"category_id": categoryId})
	////}
	////
	////if search != "" {
	////	query = query.Where(sq.ILike{"name": "%" + search + "%"})
	////}
	////
	////sql, args, err := query.ToSql()
	////if err != nil {
	////	return nil, err
	////}
	////
	////rows, err := pr.db.Query(ctx, sql, args...)
	////if err != nil {
	////	return nil, err
	////}
	////
	////for rows.Next() {
	////	err := rows.Scan(
	////		&product.ID,
	////		&product.CategoryID,
	////		&product.SKU,
	////		&product.Name,
	////		&product.Stock,
	////		&product.Price,
	////		&product.Image,
	////		&product.CreatedAt,
	////		&product.UpdatedAt,
	////	)
	////	if err != nil {
	////		return nil, err
	////	}
	////
	////	products = append(products, product)
	//}

	return products, nil
}

// UpdateProduct updates a product record in the database
func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product, updatedFields ...string) (*domain.Product, error) {
	tx := pr.db.WithContext(ctx).Model(product)

	if len(updatedFields) > 0 {
		tx = tx.Select(updatedFields)
	}

	if err := tx.Updates(product).Error; err != nil {
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product record from the database by id
func (pr *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {

	return nil
}
