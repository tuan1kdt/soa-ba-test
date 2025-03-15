package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/storage/postgres"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
)

/**
 * ProductRepository implements port.ProductRepository interface
 * and provides an access to the postgres database
 */
type ProductRepository struct {
	db *postgres.DB
}

// NewProductRepository creates a new product repository instance
func NewProductRepository(db *postgres.DB) *ProductRepository {
	return &ProductRepository{
		db,
	}
}

// CreateProduct creates a new product record in the database
func (pr *ProductRepository) CreateProduct(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := pr.db.QueryBuilder.Insert("products").
		Columns("id", "reference", "name", "added_date", "status", "category_id", "price", "stock_city", "supplier_id", "quantity").
		Values(
			product.ID,
			product.Reference,
			product.Name,
			product.AddedDate,
			product.Status,
			product.CategoryID,
			product.Price,
			product.StockCity,
			product.SupplierID,
			product.Quantity,
		).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.Reference,
		&product.Name,
		&product.AddedDate,
		&product.Status,
		&product.CategoryID,
		&product.Price,
		&product.StockCity,
		&product.SupplierID,
		&product.Quantity,
	)
	if err != nil {
		if errCode := pr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return product, nil
}

// GetProductByID retrieves a product record from the database by id
func (pr *ProductRepository) GetProductByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	var product domain.Product

	query := pr.db.QueryBuilder.Select("*").
		From("products").
		Where(sq.Eq{"id": id}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.Reference,
		&product.Name,
		&product.AddedDate,
		&product.Status,
		&product.CategoryID,
		&product.Price,
		&product.StockCity,
		&product.SupplierID,
		&product.Quantity,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &product, nil
}

// ListProducts retrieves a list of products from the database
func (pr *ProductRepository) ListProducts(ctx context.Context, search string, categoryIds []uuid.UUID, skip, limit uint64) ([]domain.Product, error) {
	var product domain.Product
	var products []domain.Product

	query := pr.db.QueryBuilder.Select("*").
		From("products").
		OrderBy("id").
		Limit(limit).
		Offset((skip) * limit)

	if len(categoryIds) != 0 {
		query = query.Where(sq.Eq{"category_id": categoryIds})
	}

	if search != "" {
		query = query.Where(sq.ILike{"name": "%" + search + "%"})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := pr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&product.ID,
			&product.Reference,
			&product.Name,
			&product.AddedDate,
			&product.Status,
			&product.CategoryID,
			&product.Price,
			&product.StockCity,
			&product.SupplierID,
			&product.Quantity,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// UpdateProduct updates a product record in the database
func (pr *ProductRepository) UpdateProduct(ctx context.Context, product *domain.Product, updatedFields ...string) (*domain.Product, error) {
	query := pr.db.QueryBuilder.Update("products")

	for _, field := range updatedFields {
		switch field {
		case "reference":
			query = query.Set(field, product.Reference)
		case "name":
			query = query.Set(field, product.Name)
		case "added_date":
			query = query.Set(field, product.AddedDate)
		case "status":
			query = query.Set(field, product.Status)
		case "category_id":
			query = query.Set(field, product.CategoryID)
		case "price":
			query = query.Set(field, nullFloat64(product.Price))
		case "stock_city":
			query = query.Set(field, product.StockCity)
		case "supplier_id":
			query = query.Set(field, product.SupplierID)
		case "quantity":
			query = query.Set(field, product.Quantity)
		}
	}

	query = query.Where(sq.Eq{"id": product.ID}).Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = pr.db.QueryRow(ctx, sql, args...).Scan(
		&product.ID,
		&product.Reference,
		&product.Name,
		&product.AddedDate,
		&product.Status,
		&product.CategoryID,
		&product.Price,
		&product.StockCity,
		&product.SupplierID,
		&product.Quantity,
	)
	if err != nil {
		if errCode := pr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return product, nil
}

// DeleteProduct deletes a product record from the database by id
func (pr *ProductRepository) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	query := pr.db.QueryBuilder.Delete("products").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = pr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pr *ProductRepository) StatisticSupplierProduct(ctx context.Context) ([]*domain.StatisticSupplierProduct, error) {
	query := `
		SELECT 
			p.supplier_id, 
			s.name, 
			COUNT(*) * 100.0 / SUM(COUNT(*)) OVER () AS percentage
		FROM products AS p
		INNER JOIN public.suppliers s ON p.supplier_id = s.id
		GROUP BY p.supplier_id, s.name;
	`
	rows, err := pr.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*domain.StatisticSupplierProduct, 0)

	for rows.Next() {
		var stat domain.StatisticSupplierProduct
		err := rows.Scan(&stat.SupplierID, &stat.SupplierName, &stat.Percentage)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}

func (pr *ProductRepository) StatisticCategoryProduct(ctx context.Context) ([]*domain.StatisticCategoryProduct, error) {
	query := `
		SELECT 
			p.category_id, 
			c.name, 
			COUNT(*) * 100.0 / SUM(COUNT(*)) OVER () AS percentage
		FROM products AS p
		INNER JOIN public.categories c ON p.category_id = c.id
		GROUP BY p.category_id, c.name;
	`
	rows, err := pr.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*domain.StatisticCategoryProduct, 0)

	for rows.Next() {
		var stat domain.StatisticCategoryProduct
		err := rows.Scan(&stat.CategoryID, &stat.CategoryName, &stat.Percentage)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &stat)
	}

	return stats, nil
}
