package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/tuan1kdt/soa-ba-test/internal/core/domain"
)

// response represents a response body format
type response struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data,omitempty"`
}

// newResponse is a helper function to create a response body
func newResponse(success bool, message string, data any) response {
	return response{
		Success: success,
		Message: message,
		Data:    data,
	}
}

// meta represents metadata for a paginated response
type meta struct {
	Total uint64 `json:"total" example:"100"`
	Limit uint64 `json:"limit" example:"10"`
	Skip  uint64 `json:"skip" example:"0"`
}

// newMeta is a helper function to create metadata for a paginated response
func newMeta(total, limit, skip uint64) meta {
	return meta{
		Total: total,
		Limit: limit,
		Skip:  skip,
	}
}

// authResponse represents an authentication response body
type authResponse struct {
	AccessToken string `json:"token" example:"v2.local.Gdh5kiOTyyaQ3_bNykYDeYHO21Jg2..."`
}

// newAuthResponse is a helper function to create a response body for handling authentication data
func newAuthResponse(token string) authResponse {
	return authResponse{
		AccessToken: token,
	}
}

// userResponse represents a user response body
type userResponse struct {
	ID        uint64    `json:"id" example:"1"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"test@example.com"`
	CreatedAt time.Time `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// categoryResponse represents a category response body
type categoryResponse struct {
	ID   uuid.UUID `json:"id,omitempty" example:"1"`
	Name string    `json:"name,omitempty" example:"Foods"`
}

// newCategoryResponse is a helper function to create a response body for handling category data
func newCategoryResponse(category *domain.Category) categoryResponse {
	if category == nil {
		return categoryResponse{}
	}
	return categoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}

// productResponse represents a product response body
type productResponse struct {
	ID         uuid.UUID `json:"id" example:"1"`
	Reference  string
	Name       string
	AddedDate  time.Time
	Status     string
	CategoryID *uuid.UUID
	Price      float64
	StockCity  string
	SupplierID *uuid.UUID
	Quantity   int
	Category   categoryResponse `json:"category,omitempty"`
}

// newProductResponse is a helper function to create a response body for handling product data
func newProductResponse(product *domain.Product) productResponse {
	return productResponse{
		ID:         product.ID,
		Reference:  product.Reference,
		Name:       product.Name,
		AddedDate:  product.AddedDate,
		Status:     product.Status.String(),
		CategoryID: product.CategoryID,
		Price:      product.Price,
		StockCity:  product.StockCity,
		SupplierID: product.SupplierID,
		Quantity:   product.Quantity,
		Category:   newCategoryResponse(product.Category),
	}
}

type ProductDistancesResponse struct {
	DistanceKM float64 `json:"distance_km" example:"1.5"`
}

func newProductDistancesResponse(distances float64) ProductDistancesResponse {

	return ProductDistancesResponse{DistanceKM: distances}
}

// orderProductResponse represents an order product response body
type orderProductResponse struct {
	ID               uint64          `json:"id" example:"1"`
	OrderID          uint64          `json:"order_id" example:"1"`
	ProductID        uint64          `json:"product_id" example:"1"`
	Quantity         int64           `json:"qty" example:"1"`
	Price            float64         `json:"price" example:"100000"`
	TotalNormalPrice float64         `json:"total_normal_price" example:"100000"`
	TotalFinalPrice  float64         `json:"total_final_price" example:"100000"`
	Product          productResponse `json:"product"`
	CreatedAt        time.Time       `json:"created_at" example:"1970-01-01T00:00:00Z"`
	UpdatedAt        time.Time       `json:"updated_at" example:"1970-01-01T00:00:00Z"`
}

// errorStatusMap is a map of defined error messages and their corresponding http status codes
var errorStatusMap = map[error]int{
	domain.ErrInternal:                   http.StatusInternalServerError,
	domain.ErrDataNotFound:               http.StatusNotFound,
	domain.ErrConflictingData:            http.StatusConflict,
	domain.ErrInvalidCredentials:         http.StatusUnauthorized,
	domain.ErrUnauthorized:               http.StatusUnauthorized,
	domain.ErrEmptyAuthorizationHeader:   http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationHeader: http.StatusUnauthorized,
	domain.ErrInvalidAuthorizationType:   http.StatusUnauthorized,
	domain.ErrInvalidToken:               http.StatusUnauthorized,
	domain.ErrExpiredToken:               http.StatusUnauthorized,
	domain.ErrForbidden:                  http.StatusForbidden,
	domain.ErrNoUpdatedData:              http.StatusBadRequest,
	domain.ErrInsufficientStock:          http.StatusBadRequest,
	domain.ErrInsufficientPayment:        http.StatusBadRequest,
	domain.ErrInvalidStatus:              http.StatusBadRequest,
}

// validationError sends an error response for some specific request validation error
func validationError(ctx *gin.Context, err error) {
	errMsgs := parseError(err)
	errRsp := newErrorResponse(errMsgs)
	ctx.JSON(http.StatusBadRequest, errRsp)
}

// handleError determines the status code of an error and returns a JSON response with the error message and status code
func handleError(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.JSON(statusCode, errRsp)
}

// handleAbort sends an error response and aborts the request with the specified status code and error message
func handleAbort(ctx *gin.Context, err error) {
	statusCode, ok := errorStatusMap[err]
	if !ok {
		statusCode = http.StatusInternalServerError
	}

	errMsg := parseError(err)
	errRsp := newErrorResponse(errMsg)
	ctx.AbortWithStatusJSON(statusCode, errRsp)
}

// parseError parses error messages from the error object and returns a slice of error messages
func parseError(err error) []string {
	var errMsgs []string

	if errors.As(err, &validator.ValidationErrors{}) {
		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Error())
		}
	} else {
		errMsgs = append(errMsgs, err.Error())
	}

	return errMsgs
}

// errorResponse represents an error response body format
type errorResponse struct {
	Success  bool     `json:"success" example:"false"`
	Messages []string `json:"messages" example:"Error message 1, Error message 2"`
}

// newErrorResponse is a helper function to create an error response body
func newErrorResponse(errMsgs []string) errorResponse {
	return errorResponse{
		Success:  false,
		Messages: errMsgs,
	}
}

// handleSuccess sends a success response with the specified status code and optional data
func handleSuccess(ctx *gin.Context, data any) {
	rsp := newResponse(true, "Success", data)
	ctx.JSON(http.StatusOK, rsp)
}
