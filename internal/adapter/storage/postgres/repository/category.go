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
 * CategoryRepository implements port.CategoryRepository interface
 * and provides an access to the postgres database
 */
type CategoryRepository struct {
	db *postgres.DB
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(db *postgres.DB) *CategoryRepository {
	return &CategoryRepository{
		db,
	}
}

// CreateCategory creates a new category record in the database
func (cr *CategoryRepository) CreateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := cr.db.QueryBuilder.Insert("categories").
		Columns("name").
		Values(category.Name).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
	)
	if err != nil {
		if errCode := cr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return category, nil
}

// GetCategoryByID retrieves a category record from the database by id
func (cr *CategoryRepository) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	var category domain.Category

	query := cr.db.QueryBuilder.Select("*").
		From("categories").
		Where(sq.Eq{"id": id.String()}).
		Limit(1)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDataNotFound
		}
		return nil, err
	}

	return &category, nil
}

// ListCategories retrieves a list of categories from the database
func (cr *CategoryRepository) ListCategories(ctx context.Context, skip, limit uint64) ([]domain.Category, error) {
	var category domain.Category
	var categories []domain.Category

	query := cr.db.QueryBuilder.Select("*").
		From("categories").
		OrderBy("id").
		Limit(limit).
		Offset((skip) * limit)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := cr.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err := rows.Scan(
			&category.ID,
			&category.Name,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategory updates a category record in the database
func (cr *CategoryRepository) UpdateCategory(ctx context.Context, category *domain.Category) (*domain.Category, error) {
	query := cr.db.QueryBuilder.Update("categories").
		Set("name", category.Name).
		Where(sq.Eq{"id": category.ID}).
		Suffix("RETURNING *")

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = cr.db.QueryRow(ctx, sql, args...).Scan(
		&category.ID,
		&category.Name,
	)
	if err != nil {
		if errCode := cr.db.ErrorCode(err); errCode == "23505" {
			return nil, domain.ErrConflictingData
		}
		return nil, err
	}

	return category, nil
}

// DeleteCategory deletes a category record from the database by id
func (cr *CategoryRepository) DeleteCategory(ctx context.Context, id uuid.UUID) error {
	query := cr.db.QueryBuilder.Delete("categories").
		Where(sq.Eq{"id": id})

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = cr.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
