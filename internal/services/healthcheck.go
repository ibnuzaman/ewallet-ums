package services

import (
	"context"
	"time"

	"github.com/ibnuzaman/ewallet-ums/database"
	"github.com/ibnuzaman/ewallet-ums/internal/interfaces"
)

const healthCheckTimeout = 3 * time.Second

// Healthcheck service implementation.
type Healthcheck struct {
	HealthcheckRepository interfaces.IHealthcheckRepo
}

// HealthcheckServices performs health check including database.
func (s *Healthcheck) HealthcheckServices() (string, error) {
	// Check database connection
	ctx, cancel := context.WithTimeout(context.Background(), healthCheckTimeout)
	defer cancel()

	if err := database.HealthCheck(ctx); err != nil {
		return "unhealthy - database connection failed", err
	}

	return "healthy", nil
}
