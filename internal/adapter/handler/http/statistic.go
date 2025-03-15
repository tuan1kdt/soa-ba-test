package http

import (
	"math"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tuan1kdt/soa-ba-test/internal/core/port"
)

// StatisticHandler represents the HTTP handler for Statistic-related requests
type StatisticHandler struct {
	svc port.StatisticService
}

// NewStatisticHandler creates a new StatisticHandler instance
func NewStatisticHandler(svc port.StatisticService) *StatisticHandler {
	return &StatisticHandler{
		svc,
	}
}

type StatisticSupplierProductResponse struct {
	SupplierID   uuid.UUID `json:"supplier_id"`
	SupplierName string    `json:"supplier_name"`
	Percentage   float64   `json:"percentage"`
}

type StatisticCategoryProductResponse struct {
	CategoryID   uuid.UUID `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Percentage   float64   `json:"percentage"`
}

// GetSupplierProduct godoc
//
//	@Summary		Get Statistic of supplier product
//	@Description	Get Statistic of supplier product
//	@Tags			Statistics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	StatisticSupplierProductResponse	"Statistic created"
//	@Failure		400	{object}	errorResponse						"Validation error"
//	@Failure		401	{object}	errorResponse						"Unauthorized error"
//	@Failure		403	{object}	errorResponse						"Forbidden error"
//	@Failure		404	{object}	errorResponse						"Data not found error"
//	@Failure		409	{object}	errorResponse						"Data conflict error"
//	@Failure		500	{object}	errorResponse						"Internal server error"
//	@Router			/statistics/products-per-supplier [get]
//	@Security		BearerAuth
func (ch *StatisticHandler) GetSupplierProduct(ctx *gin.Context) {
	res, err := ch.svc.StatisticSupplierProduct(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	var statistics []StatisticSupplierProductResponse
	for _, item := range res {
		statistics = append(statistics, StatisticSupplierProductResponse{
			SupplierID:   item.SupplierID,
			SupplierName: item.SupplierName,
			Percentage:   roundToTwoDecimalPlaces(item.Percentage),
		})

	}

	handleSuccess(ctx, statistics)
}

// GetCategoryProduct godoc
//
//	@Summary		Get Statistic of category product
//	@Description	Get Statistic of category product
//	@Tags			Statistics
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	StatisticCategoryProductResponse	"Statistic created"
//	@Failure		400	{object}	errorResponse						"Validation error"
//	@Failure		401	{object}	errorResponse						"Unauthorized error"
//	@Failure		403	{object}	errorResponse						"Forbidden error"
//	@Failure		404	{object}	errorResponse						"Data not found error"
//	@Failure		409	{object}	errorResponse						"Data conflict error"
//	@Failure		500	{object}	errorResponse						"Internal server error"
//	@Router			/statistics/products-per-category [get]
//	@Security		BearerAuth
func (ch *StatisticHandler) GetCategoryProduct(ctx *gin.Context) {
	res, err := ch.svc.StatisticCategoryProduct(ctx.Request.Context())
	if err != nil {
		handleError(ctx, err)
		return
	}

	var statistics []StatisticCategoryProductResponse
	for _, item := range res {
		statistics = append(statistics, StatisticCategoryProductResponse{
			CategoryID:   item.CategoryID,
			CategoryName: item.CategoryName,
			Percentage:   roundToTwoDecimalPlaces(item.Percentage),
		})

	}

	handleSuccess(ctx, statistics)
}

// roundToTwoDecimalPlaces round to 2 decimal places
func roundToTwoDecimalPlaces(value float64) float64 {
	return math.Round(value*100) / 100
}
