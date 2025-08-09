package http

import (
	"log/slog"
	"strings"

	"example.com/go-config-management-service/internal/adapter/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router is a wrapper for HTTP router
type Router struct {
	*gin.Engine
}

// NewRouter creates a new HTTP router
func NewRouter(
	config *config.HTTP,
	configurationHandler ConfigurationHandler,
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

	// Swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	configuration := router.Group("/configs")
	{
		configuration.GET("/", configurationHandler.ListConfigurations)
		configuration.PUT("/:name", configurationHandler.PutConfiguration)
		configuration.GET("/:name", configurationHandler.GetConfiguration)
		configuration.GET("/:name/versions", configurationHandler.ListConfigurationVersions)
		configuration.GET("/:name/versions/:version", configurationHandler.GetConfigurationVersion)
		configuration.POST("/:name/versions/:version/rollback", configurationHandler.RollbackConfigurationVersion)
	}

	return &Router{
		router,
	}, nil
}

// Serve starts the HTTP server
func (r *Router) Serve(listenAddr string) error {
	return r.Run(listenAddr)
}
