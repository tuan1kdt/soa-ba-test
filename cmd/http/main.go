package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/tuan1kdt/soa-ba-test/docs"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/config"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/geohelper"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/handler/http"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/logger"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/storage/postgres"
	"github.com/tuan1kdt/soa-ba-test/internal/adapter/storage/postgres/repository"
	"github.com/tuan1kdt/soa-ba-test/internal/core/service"
)

// @title						Go SOA Test (Source of Asia) API
// @version					1.0
// @description				This is a simple RESTful Product Backend Service API written in Go using Gin web framework, MySQL database
//
// @contact.name				tuanla
// @contact.url				https://github.com/tuanla/soa-be-test
// @contact.email				leanhtuan1998hl@gmail.com
//
// @host						localhost
// @BasePath					/v1
// @schemes					http https
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the access token.
func main() {
	// Load environment variables
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	// Set logger
	logger.Set(config.App)

	slog.Info("Starting the application", "app", config.App.Name, "env", config.App.Env)

	// Init database
	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing database connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	slog.Info("Successfully connected to the database", "db", config.DB.Connection)

	// Category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryService := service.NewCategoryService(categoryRepo, nil)
	categoryHandler := http.NewCategoryHandler(categoryService)

	// Product
	productRepo := repository.NewProductRepository(db)
	geoClient := geohelper.New(config.GEO)
	productService := service.NewProductService(productRepo, categoryRepo, nil, geoClient)
	productHandler := http.NewProductHandler(productService)

	// Statistic
	statisticService := service.NewStatisticService(productRepo)
	statisticHandler := http.NewStatisticHandler(statisticService)

	// Init router
	router, err := http.NewRouter(
		config.HTTP,
		*categoryHandler,
		*productHandler,
		*statisticHandler,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	// Start server
	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
