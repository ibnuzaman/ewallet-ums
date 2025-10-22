package constants

import "time"

const (
	ServiceHealthyMessage = "service healthy"
	HealthcheckTimeout    = 60 * time.Second
	ShutdownTimeout       = 30 * time.Second
	ReadTimeout           = 15 * time.Second
	WriteTimeout          = 15 * time.Second
	IdleTimeout           = 60 * time.Second
	RequestTimeout        = 60 * time.Second
)
