package services

import (
	"github.com/ibnuzaman/ewallet-ums/internal/constants"
	"github.com/ibnuzaman/ewallet-ums/internal/interfaces"
)

type Healthcheck struct {
	HealthcheckRepository interfaces.IHealthcheckRepo
}

func (s *Healthcheck) HealthcheckServices() (string, error) {
	return constants.ServiceHealthyMessage, nil
}
