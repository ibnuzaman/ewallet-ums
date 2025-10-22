package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock service for testing
type mockHealthcheckService struct {
	shouldError bool
}

func (m *mockHealthcheckService) HealthcheckServices() (string, error) {
	if m.shouldError {
		return "", errors.New("service error")
	}
	return "service healthy", nil
}

func TestHealthcheck_HealthcheckHandlerHTTP_Success(t *testing.T) {
	// Arrange
	mockSvc := &mockHealthcheckService{shouldError: false}
	handler := &Healthcheck{
		HealthcheckServices: mockSvc,
	}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	// Act
	handler.HealthcheckHandlerHTTP(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHealthcheck_HealthcheckHandlerHTTP_Error(t *testing.T) {
	// Arrange
	mockSvc := &mockHealthcheckService{shouldError: true}
	handler := &Healthcheck{
		HealthcheckServices: mockSvc,
	}

	req := httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
	w := httptest.NewRecorder()

	// Act
	handler.HealthcheckHandlerHTTP(w, req)

	// Assert
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
