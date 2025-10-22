package services

import (
	"testing"
)

func TestHealthcheck_HealthcheckServices(t *testing.T) {
	// Arrange
	svc := &Healthcheck{}

	// Act
	result, err := svc.HealthcheckServices()

	// Assert
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "service healthy" {
		t.Errorf("Expected 'service healthy', got '%s'", result)
	}
}
