package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
	"github.com/tuan1kdt/soa-ba-test/internal/core/port"
)

// CategoryHandler represents the HTTP handler for category-related requests
type CategoryHandler struct {
	svc port.CategoryService
}

// NewCategoryHandler creates a new CategoryHandler instance
func NewCategoryHandler(svc port.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		svc,
	}
}

// createCategoryRequest represents a request body for creating a new category
type createCategoryRequest struct {
	Name string `json:"name" binding:"required" example:"Foods"`
}

// CreateCategory godoc
//
//	@Summary		Create a new category
//	@Description	create a new category with name
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			createCategoryRequest	body		createCategoryRequest	true	"Create category request"
//	@Success		200						{object}	categoryResponse		"Category created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/categories [post]
//	@Security		BearerAuth
func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	var req createCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	category := domain.Category{
		Name: req.Name,
	}

	_, err := ch.svc.CreateCategory(ctx, &category)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newCategoryResponse(&category)

	handleSuccess(ctx, rsp)
}

// getCategoryRequest represents a request body for retrieving a category
type getCategoryRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetCategory godoc
//
//	@Summary		Get a category
//	@Description	get a category by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string				true	"Category ID"
//	@Success		200	{object}	categoryResponse	"Category retrieved"
//	@Failure		400	{object}	errorResponse		"Validation error"
//	@Failure		404	{object}	errorResponse		"Data not found error"
//	@Failure		500	{object}	errorResponse		"Internal server error"
//	@Router			/categories/{id} [get]
//	@Security		BearerAuth
func (ch *CategoryHandler) GetCategory(ctx *gin.Context) {
	var req getCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	id, _ := uuid.Parse(req.ID)

	category, err := ch.svc.GetCategory(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newCategoryResponse(category)

	handleSuccess(ctx, rsp)
}

// listCategoriesRequest represents a request body for listing categories
type listCategoriesRequest struct {
	Skip  uint64 `form:"skip"`
	Limit uint64 `form:"limit"`
}

// ListCategories godoc
//
//	@Summary		List categories
//	@Description	List categories with pagination
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		uint64			false	"Skip"
//	@Param			limit	query		uint64			false	"Limit"
//	@Success		200		{object}	meta			"Categories displayed"
//	@Failure		400		{object}	errorResponse	"Validation error"
//	@Failure		500		{object}	errorResponse	"Internal server error"
//	@Router			/categories [get]
//	@Security		BearerAuth
func (ch *CategoryHandler) ListCategories(ctx *gin.Context) {
	var req listCategoriesRequest
	var categoriesList []categoryResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if req.Limit == 0 || req.Limit > 1000 {
		req.Limit = 10
	}

	categories, err := ch.svc.ListCategories(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, category := range categories {
		categoriesList = append(categoriesList, newCategoryResponse(&category))
	}

	total := uint64(len(categoriesList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, categoriesList, "categories")

	handleSuccess(ctx, rsp)
}

// updateCategoryRequest represents a request body for updating a category
type updateCategoryRequest struct {
	ID   string `uri:"id" binding:"required,uuid"`
	Name string `json:"name"`
}

// UpdateCategory godoc
//
//	@Summary		Update a category
//	@Description	update a category's name by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id						path		string					true	"Category ID"
//	@Param			updateCategoryRequest	body		updateCategoryRequest	true	"Update category request"
//	@Success		200						{object}	categoryResponse		"Category updated"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/categories/{id} [patch]
//	@Security		BearerAuth
func (ch *CategoryHandler) UpdateCategory(ctx *gin.Context) {
	var req updateCategoryRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	id, _ := uuid.Parse(req.ID)

	category := domain.Category{
		ID:   id,
		Name: req.Name,
	}

	_, err := ch.svc.UpdateCategory(ctx, &category)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newCategoryResponse(&category)

	handleSuccess(ctx, rsp)
}

// deleteCategoryRequest represents a request body for deleting a category
type deleteCategoryRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// DeleteCategory godoc
//
//	@Summary		Delete a category
//	@Description	Delete a category by id
//	@Tags			Categories
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string			true	"Category ID"
//	@Success		200	{object}	response		"Category deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/categories/{id} [delete]
//	@Security		BearerAuth
func (ch *CategoryHandler) DeleteCategory(ctx *gin.Context) {
	var req deleteCategoryRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	id, _ := uuid.Parse(req.ID)

	err := ch.svc.DeleteCategory(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
