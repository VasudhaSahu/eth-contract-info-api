package main

import (
	"log"

	app "github.com/VasudhaSahu/eth-contract-info-api.git/internal/application/contract"
	"github.com/VasudhaSahu/eth-contract-info-api.git/internal/config"
	infura "github.com/VasudhaSahu/eth-contract-info-api.git/internal/infrastructure/infura"
	httpHandler "github.com/VasudhaSahu/eth-contract-info-api.git/internal/interfaces/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/VasudhaSahu/eth-contract-info-api.git/docs"
)

// @title 			eth-contract-info-api Contract API
// @version 		1.0
// @description 	Demo API to fetch Ethereum smart-contract info and ERC20 metadata using Infura
// @host 			localhost:8080
// @BasePath 		/

func main() {
	cfg, err := config.Load()
	if err != nil {
		// In a real service we’d integrate structured logging here.
		log.Fatalf("failed to load config: %v", err)
	}

	// The dependency graph: RPC client -> use case -> HTTP handlers
	infuraClient := infura.NewClient(cfg.InfuraNetwork, cfg.InfuraAPIKey)
	useCase := app.NewUseCase(infuraClient, cfg.InfuraNetwork)
	handler := httpHandler.NewContractHandler(useCase)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// two read-only endpoints
	api.GET("/contracts/info", handler.GetContractInfo)
	api.GET("/contracts/erc20", handler.GetERC20Metadata)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed to start on port %s: %v", cfg.Port, err)
	}
}
