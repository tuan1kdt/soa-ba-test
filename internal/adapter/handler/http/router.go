package http

import (
	"log/slog"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/config"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.HTTP,
	categoryHandler CategoryHandler,
	productHandler ProductHandler,
	statisticHandler StatisticHandler,
) (*Router, error) {
	// Disable debug mode in production
	if config.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// CORS
	ginConfig := cors.DefaultConfig()
	allowedOrigins := config.AllowedOrigins
	originsList := strings.Split(allowedOrigins, ",")
	ginConfig.AllowOrigins = originsList

	router := gin.New()
	router.Use(sloggin.New(slog.Default()), gin.Recovery(), cors.New(ginConfig))

	// Custom validators
	//v, ok := binding.Validator.Engine().(*validator.Validate)
	//if ok {
	//	if err := v.RegisterValidation("user_role", userRoleValidator); err != nil {
	//		return nil, err
	//	}
	//
	//}

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/v1")
	{
		category := v1.Group("/categories")
		{
			category.GET("/", categoryHandler.ListCategories)
			category.GET("/:id", categoryHandler.GetCategory)

			admin := category
			{
				admin.POST("/", categoryHandler.CreateCategory)
				admin.PATCH("/:id", categoryHandler.UpdateCategory)
				admin.DELETE("/:id", categoryHandler.DeleteCategory)
			}
		}
		product := v1.Group("/products")
		{
			product.GET("/", productHandler.ListProducts)
			product.GET("/export", productHandler.ExportProducts)
			product.GET("/:id", productHandler.GetProduct)
			product.GET("/:id/distance", productHandler.GetProductDistance)

			admin := product
			{
				admin.POST("/", productHandler.CreateProduct)
				admin.PATCH("/:id", productHandler.UpdateProduct)
				admin.DELETE("/:id", productHandler.DeleteProduct)
			}
		}
		statistic := v1.Group("/statistics")
		{
			statistic.GET("/products-per-category", statisticHandler.GetCategoryProduct)
			statistic.GET("/products-per-supplier", statisticHandler.GetSupplierProduct)
		}
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
