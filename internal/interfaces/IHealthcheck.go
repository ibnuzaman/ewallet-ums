package interfaces

import "net/http"

// IHealthcheckServices defines the interface for healthcheck service.
type IHealthcheckServices interface {
	HealthcheckServices() (string, error)
}

// IHealthcheckAPI defines the interface for healthcheck API handler.
type IHealthcheckAPI interface {
	HealthcheckHandlerHTTP(w http.ResponseWriter, r *http.Request)
}

// IHealthcheckRepo defines the interface for healthcheck repository.
type IHealthcheckRepo interface{}
