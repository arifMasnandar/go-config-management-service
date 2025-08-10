package main

import (
	"fmt"
	"log/slog"
	"os"

	"example.com/go-config-management-service/internal/adapter/config"
	"example.com/go-config-management-service/internal/adapter/handler/http"
	"example.com/go-config-management-service/internal/adapter/storage/memory"
	"example.com/go-config-management-service/internal/core/service"

	_ "example.com/go-config-management-service/docs"
)

func main() {
	// Load environment variables
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading environment variables", "error", err)
		os.Exit(1)
	}

	configurationRepo := memory.NewConfigurationRepository()
	configurationService := service.NewConfigurationService(configurationRepo)
	configurationHandler := http.NewConfigurationHandler(configurationService)

	// Init router
	router, err := http.NewRouter(
		config.HTTP,
		*configurationHandler,
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
