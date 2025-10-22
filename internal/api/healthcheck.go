package api

import (
	"net/http"

	"github.com/ibnuzaman/ewallet-ums/helpers"
	"github.com/ibnuzaman/ewallet-ums/internal/interfaces"
)

type Healthcheck struct {
	HealthcheckServices interfaces.IHealthcheckServices
}

func (api *Healthcheck) HealthcheckHandlerHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := api.HealthcheckServices.HealthcheckServices()
	if err != nil {
		helpers.SendErrorResponse(w, r, "Health check failed", err, http.StatusInternalServerError)
		return
	}

	helpers.SendResponse(w, r, map[string]string{
		"status": response,
	}, "Health check successful", http.StatusOK)
}
