package services

import (
	"testing"
)

func TestHealthcheck_HealthcheckServices(t *testing.T) {
	t.Parallel()

	t.Run("without database connection", func(t *testing.T) {
		t.Parallel()

		// Arrange
		svc := &Healthcheck{}

		// Act
		result, err := svc.HealthcheckServices()

		// Assert - Since DB is not initialized in test, it should return error
		if err == nil {
			t.Error("Expected error when database is not initialized, got nil")
		}

		if result != "unhealthy - database connection failed" {
			t.Errorf("Expected 'unhealthy - database connection failed', got '%s'", result)
		}
	})
}
