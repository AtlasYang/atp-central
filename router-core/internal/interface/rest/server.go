package rest

import (
	"context"
	"log"

	api_client_delivery "aigendrug.com/router-core/internal/api_client/delivery"
	api_docs_delivery "aigendrug.com/router-core/internal/api_docs/delivery"
	client_service "aigendrug.com/router-core/internal/client/application/service"
	client_delivery "aigendrug.com/router-core/internal/client/delivery"
	client_persistence "aigendrug.com/router-core/internal/client/infrastructure/persistence"
	"aigendrug.com/router-core/internal/config"
	"aigendrug.com/router-core/internal/shared/database/postgres"
	lambda_wrapper "aigendrug.com/router-core/internal/shared/lambda-wrapper"
	"aigendrug.com/router-core/internal/shared/selector"
	tool_service "aigendrug.com/router-core/internal/tool/application/service"
	tool_delivery "aigendrug.com/router-core/internal/tool/delivery"
	tool_persistence "aigendrug.com/router-core/internal/tool/infrastructure/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if config.Server.RunMode == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	port := config.Server.Port
	if port == "" {
		port = "8080"
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://aigendrug.com",
			"https://dev.aigendrug.com",
			"http://localhost:12352",
		},
		AllowCredentials: true,
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
			"X-API-Key",
		},
	}))

	router.Handle("GET", "/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	ctx := context.Background()

	success, err := postgres.AutoMigrateFromConnectionString(ctx, config.Database.Postgres.Connection)
	if err != nil {
		log.Printf("Failed to auto migrate: %v", err)
	}
	if !success {
		log.Fatalf("Exiting...")
	}

	pgPool, err := postgres.NewPostgresPoolFromConnectionString(ctx, config.Database.Postgres.Connection)
	if err != nil {
		log.Fatalf("Failed to create postgres pool: %v", err)
	}

	lambdaClient := lambda_wrapper.NewLambdaWrapperClient(config)

	selectorService := selector.NewSelectorService(config)

	clientRepo := client_persistence.NewPgClientRepository(pgPool)
	toolRepo := tool_persistence.NewPgToolRepository(pgPool)

	clientService := client_service.NewClientService(pgPool, clientRepo)
	toolService := tool_service.NewToolService(pgPool, toolRepo, selectorService, lambdaClient)

	apiDocsHandler := api_docs_delivery.NewAPIDocsHandler(config)
	apiClientHandler := api_client_delivery.NewAPIClientHandler(config)
	clientHandler := client_delivery.NewClientHandler(clientService)
	toolHandler := tool_delivery.NewToolHandler(toolService)

	api_docs_delivery.SetupAPIDocsRoutes(router, apiDocsHandler)
	api_client_delivery.SetupAPIClientRoutes(router, apiClientHandler)
	client_delivery.SetupClientRoutes(router, pgPool, clientHandler)
	tool_delivery.SetupToolRoutes(router, pgPool, toolHandler)

	router.Run(":" + port)
}
