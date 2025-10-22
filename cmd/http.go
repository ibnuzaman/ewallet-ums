package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ibnuzaman/ewallet-ums/helpers"
	"github.com/ibnuzaman/ewallet-ums/internal/api"
	"github.com/ibnuzaman/ewallet-ums/internal/constants"
	"github.com/ibnuzaman/ewallet-ums/internal/interfaces"
	"github.com/ibnuzaman/ewallet-ums/internal/services"
)

// ServerHTTP starts the HTTP server.
func ServerHTTP() {
	dependency := dependencyInject()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(helpers.LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(constants.RequestTimeout))

	// Routes
	r.Get("/healthcheck", dependency.HealthcheckAPI.HealthcheckHandlerHTTP)

	// Server configuration
	port := helpers.GetEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  constants.ReadTimeout,
		WriteTimeout: constants.WriteTimeout,
		IdleTimeout:  constants.IdleTimeout,
	}

	// Start server in a goroutine
	go func() {
		helpers.Logger.Infof("Starting HTTP server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			helpers.Logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	helpers.Logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		helpers.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	helpers.Logger.Info("Server exited properly")
}

// Dependency holds all API dependencies.
type Dependency struct {
	HealthcheckAPI interfaces.IHealthcheckAPI
}

func dependencyInject() Dependency {
	healthcheckSvc := &services.Healthcheck{}
	healthcheckAPI := &api.Healthcheck{
		HealthcheckServices: healthcheckSvc,
	}

	return Dependency{
		HealthcheckAPI: healthcheckAPI,
	}
}
