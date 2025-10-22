package interfaces

import "net/http"

type IHealthcheckServices interface {
	HealthcheckServices() (string, error)
}

type IHealthcheckAPI interface {
	HealthcheckHandlerHTTP(w http.ResponseWriter, r *http.Request)
}

type IHealthcheckRepo interface {
}
