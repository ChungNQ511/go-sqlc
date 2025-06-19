package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	db "example.com/m/db/sqlc"
	"example.com/m/docs"
	"example.com/m/modules/examples"
	"example.com/m/pkg/middlewares"
	"example.com/m/pkg/utils/app"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	server *gin.Engine
	ctx    context.Context
)

// @title           Example API
// @version         1.0
// @description     Example API for testing and development.
// @description     Features include:
// @description     - Example of API
// @description     - Example of API
// @description     - Example of API

// @contact.name   Example Developer Team
// @contact.email  example@example.com

// @host            localhost:8080
func main() {
	config, err := app.LoadAppConfig(".")
	if err != nil {
		log.Fatal("Cannot load app config: ", err)
	}
	// swagger
	// Manually update Swagger info with host
	docs.SwaggerInfo.Host = config.SWAGGER_HOST
	fmt.Printf("✅ Swagger configured with host: %s\n", config.SWAGGER_HOST)

	// connect to database
	db, store, pool := app.Connect(config)
	defer pool.Close()

	// connect to redis
	err = app.RSetup(config)
	if err != nil {
		log.Fatal("Cannot connect to redis: ", err)
	}
	defer app.RedisConn.Close()

	// server
	server = gin.Default()

	// middlewares & recovery
	server.Use(
		gin.Recovery(),
		middlewares.CORS(config.CORS_ORIGIN),
	)

	// routes
	router := server.Group("/api")
	versionAPI, err := os.ReadFile(".build-version")
	if err != nil {
		log.Fatal("Cannot read build version: ", err)
	}

	// version
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"version": string(versionAPI),
			"message": "Welcome to the API",
			"service": "API",
			"host":    config.SERVICE_HOST_NAME,
		})
	})

	// healthz
	router.GET("/healthz", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "pong",
		})
	})

	// set app router
	setAppRouter(router, db, store)

	// swagger
	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// no route
	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": fmt.Sprintf("The requested path %s is not found", ctx.Request.URL.Path),
		})
	})

	// run server
	server.Run(":" + config.SERVER_ADDRESS)
}

// setAppRouter setAppRouter
// @Summary Set up application routes
// @Description Configure all routes for the application
// @Tags router
// @Accept json
// @Produce json
// @Success 200 {object} nil
// @Router /api [get]
func setAppRouter(router *gin.RouterGroup, db *db.Queries, store *db.Store) {
	// defined routers
	examplesController := *examples.NewController(ctx, db, store)
	examplesRoutes := examples.NewRouter(examplesController)
	examplesRoutes.RegisterRoutes(router)
}
