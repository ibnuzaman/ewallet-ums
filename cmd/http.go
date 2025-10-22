package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ibnuzaman/ewallet-ums/helpers"
	"github.com/ibnuzaman/ewallet-ums/internal/api"
	"github.com/ibnuzaman/ewallet-ums/internal/interfaces"
	"github.com/ibnuzaman/ewallet-ums/internal/services"
)

func ServerHttp() {
	dependency := dependencyInject()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(helpers.LoggerMiddleware)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Routes
	r.Get("/healthcheck", dependency.HealthcheckAPI.HealthcheckHandlerHTTP)

	// Server configuration
	port := helpers.GetEnv("PORT", "8080")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		helpers.Logger.Fatalf("Server forced to shutdown: %v", err)
	}

	helpers.Logger.Info("Server exited properly")
}

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
