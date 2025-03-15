package http

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
	"github.com/tuan1kdt/soa-ba-test/internal/core/port"
)

// ProductHandler represents the HTTP handler for product-related requests
type ProductHandler struct {
	svc port.ProductService
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(svc port.ProductService) *ProductHandler {
	return &ProductHandler{
		svc,
	}
}

// createProductRequest represents a request body for creating a new product
type createProductRequest struct {
	ID         string
	Reference  string
	Name       string
	Status     string
	CategoryID string
	Price      float64
	StockCity  string
	SupplierID string
	Quantity   int
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	create a new product with name, image, price, and stock
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			createProductRequest	body		createProductRequest	true	"Create product request"
//	@Success		200						{object}	productResponse			"Product created"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/products [post]
//	@Security		BearerAuth
func (ph *ProductHandler) CreateProduct(ctx *gin.Context) {
	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	var categoryID, supplierID *uuid.UUID
	if req.CategoryID != "" {
		tempID, err := uuid.Parse(req.CategoryID)
		if err != nil {
			validationError(ctx, err)
			return
		}
		categoryID = &tempID
	}

	if req.SupplierID != "" {
		tempID, err := uuid.Parse(req.SupplierID)
		if err != nil {
			validationError(ctx, err)
			return
		}
		supplierID = &tempID
	}

	product := domain.Product{
		ID:         uuid.UUID{},
		Reference:  req.Reference,
		Name:       req.Name,
		Status:     domain.ParseProductStatus(req.Status),
		CategoryID: categoryID,
		Price:      req.Price,
		StockCity:  req.StockCity,
		SupplierID: supplierID,
		Quantity:   req.Quantity,
	}

	_, err := ph.svc.CreateProduct(ctx, &product)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newProductResponse(&product)

	handleSuccess(ctx, rsp)
}

// getProductRequest represents a request body for retrieving a product
type getProductRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// GetProduct godoc
//
//	@Summary		Get a product
//	@Description	get a product by id with its category
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Product ID"
//	@Success		200	{object}	productResponse	"Product retrieved"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/products/{id} [get]
//	@Security		BearerAuth
func (ph *ProductHandler) GetProduct(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	product, err := ph.svc.GetProduct(ctx, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newProductResponse(product)

	handleSuccess(ctx, rsp)
}

// GetProductDistance godoc
//
//	@Summary		Get a product
//	@Description	get a product by id with its category
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Product ID"
//	@Success		200	{object}	productResponse	"Product retrieved"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/products/{id}/distance [get]
//	@Security		BearerAuth
func (ph *ProductHandler) GetProductDistance(ctx *gin.Context) {
	var req getProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	id, err := uuid.Parse(req.ID)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	sourceIP := ctx.ClientIP()

	distanceKM, err := ph.svc.GetProductDistance(ctx, sourceIP, id)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newProductDistancesResponse(distanceKM)
	handleSuccess(ctx, rsp)
}

// listProductsRequest represents a request body for listing products
type listProductsRequest struct {
	CategoryIDs []string `form:"category_ids"`
	Query       string   `form:"q"`
	Skip        uint64   `form:"skip"`
	Limit       uint64   `form:"limit"`

	paging
}

// ListProducts godoc
//
//	@Summary		List products
//	@Description	List products with pagination
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			category_id	query		string			false	"Category ID"
//	@Param			q			query		string			false	"Query"
//	@Param			skip		query		uint64			true	"Skip"
//	@Param			limit		query		uint64			true	"Limit"
//	@Success		200			{object}	meta			"Products retrieved"
//	@Failure		400			{object}	errorResponse	"Validation error"
//	@Failure		500			{object}	errorResponse	"Internal server error"
//	@Router			/products [get]
//	@Security		BearerAuth
func (ph *ProductHandler) ListProducts(ctx *gin.Context) {
	var req listProductsRequest
	var productsList []productResponse

	if err := ctx.ShouldBind(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	categories := make([]uuid.UUID, len(req.CategoryIDs))
	for i, id := range req.CategoryIDs {
		categoryID, err := uuid.Parse(id)
		if err != nil {
			validationError(ctx, err)
			return
		}
		categories[i] = categoryID
	}

	products, err := ph.svc.ListProducts(ctx, req.Query, categories, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, product := range products {
		productsList = append(productsList, newProductResponse(&product))
	}

	total := uint64(len(productsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, productsList, "products")

	handleSuccess(ctx, rsp)
}

// ExportProducts godoc
//
//	@Summary		Export products
//	@Description	Export a list of products as a PDF file
//	@Tags			Products
//	@Accept			json
//	@Produce		application/pdf
//	@Param			category_id	query		string			false	"Category ID"
//	@Param			q			query		string			false	"Query"
//	@Param			skip		query		uint64			true	"Skip"
//	@Param			limit		query		uint64			true	"Limit"
//	@Success		200			{file}		application/pdf	"PDF file generated"
//	@Failure		400			{object}	errorResponse	"Validation error"
//	@Failure		500			{object}	errorResponse	"Internal server error"
//	@Router			/products/export [get]
//	@Security		BearerAuth
func (ph *ProductHandler) ExportProducts(ctx *gin.Context) {
	var req listProductsRequest

	if err := ctx.ShouldBind(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if req.Limit == 0 {
		req.Limit = 10
	}

	categories := make([]uuid.UUID, len(req.CategoryIDs))
	for i, id := range req.CategoryIDs {
		categoryID, err := uuid.Parse(id)
		if err != nil {
			validationError(ctx, err)
			return
		}
		categories[i] = categoryID
	}

	products, err := ph.svc.ListProducts(ctx, req.Query, categories, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Product List")

	pdf.SetFont("Arial", "", 12)
	for _, product := range products {
		pdf.Ln(10)
		pdf.Cell(35, 5, product.Reference)
		pdf.Cell(35, 5, product.Name)
		pdf.Cell(35, 5, product.AddedDate.Format("2006/01/02"))
		pdf.Cell(35, 5, product.Status.String())
		pdf.Cell(35, 5, product.Category.Name)
		pdf.Cell(35, 5, fmt.Sprintf("%.2f", product.Price))
	}

	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=products.pdf")
	ctx.Header("Content-Transfer-Encoding", "binary")
	err = pdf.Output(ctx.Writer)
	if err != nil {
		handleError(ctx, err)
		return
	}
}

// updateProductRequest represents a request body for updating a product
type updateProductRequest struct {
	ID         string  `uri:"id" binding:"required,uuid"`
	CategoryID string  `json:"category_id" binding:"omitempty,required,uuid"`
	Name       string  `json:"name" binding:"omitempty"`
	Price      float64 `json:"price" binding:"omitempty,min=0" example:"2000"`
	Stock      int64   `json:"stock" binding:"omitempty,min=0" example:"200"`
	Status     string  `json:"status" binding:"omitempty,oneof=Available On Order Out of Stock"`
}

// UpdateProduct godoc
//
//	@Summary		Update a product
//	@Description	update a product's name, image, price, or stock by id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id						path		uint64					true	"Product ID"
//	@Param			updateProductRequest	body		updateProductRequest	true	"Update product request"
//	@Success		200						{object}	productResponse			"Product updated"
//	@Failure		400						{object}	errorResponse			"Validation error"
//	@Failure		401						{object}	errorResponse			"Unauthorized error"
//	@Failure		403						{object}	errorResponse			"Forbidden error"
//	@Failure		404						{object}	errorResponse			"Data not found error"
//	@Failure		409						{object}	errorResponse			"Data conflict error"
//	@Failure		500						{object}	errorResponse			"Internal server error"
//	@Router			/products/{id} [patch]
//	@Security		BearerAuth
func (ph *ProductHandler) UpdateProduct(ctx *gin.Context) {
	var req updateProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		validationError(ctx, err)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		validationError(ctx, err)
		return
	}

	var categoryID *uuid.UUID
	if req.CategoryID != "" {
		tempId, err := uuid.Parse(req.CategoryID)
		if err != nil {
			validationError(ctx, err)
			return
		}
		categoryID = &tempId
	}

	product := domain.Product{
		ID:         id,
		CategoryID: categoryID,
		Name:       req.Name,
		Price:      req.Price,
		Status:     domain.ParseProductStatus(req.Status),
	}

	_, err = ph.svc.UpdateProduct(ctx, &product)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newProductResponse(&product)

	handleSuccess(ctx, rsp)
}

// deleteProductRequest represents a request body for deleting a product
type deleteProductRequest struct {
	ID uuid.UUID `uri:"id" binding:"required,min=1" example:"1"`
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete a product by id
//	@Tags			Products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		uint64			true	"Product ID"
//	@Success		200	{object}	response		"Product deleted"
//	@Failure		400	{object}	errorResponse	"Validation error"
//	@Failure		401	{object}	errorResponse	"Unauthorized error"
//	@Failure		403	{object}	errorResponse	"Forbidden error"
//	@Failure		404	{object}	errorResponse	"Data not found error"
//	@Failure		500	{object}	errorResponse	"Internal server error"
//	@Router			/products/{id} [delete]
//	@Security		BearerAuth
func (ph *ProductHandler) DeleteProduct(ctx *gin.Context) {
	var req deleteProductRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	err := ph.svc.DeleteProduct(ctx, req.ID)
	if err != nil {
		handleError(ctx, err)
		return
	}

	handleSuccess(ctx, nil)
}
