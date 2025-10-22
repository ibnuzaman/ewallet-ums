package helpers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func SetupLogger() {
	log := logrus.New()

	// Get environment
	env := GetEnv("ENVIRONMENT", "development")

	// Set formatter based on environment
	if env == "production" {
		log.SetFormatter(&logrus.JSONFormatter{})
		log.SetLevel(logrus.InfoLevel)
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			ForceColors:   true,
		})
		log.SetLevel(logrus.DebugLevel)
	}

	Logger = log
}

// LoggerMiddleware is a middleware that logs HTTP requests.
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			Logger.WithFields(logrus.Fields{
				"method":     r.Method,
				"path":       r.URL.Path,
				"status":     ww.Status(),
				"duration":   time.Since(start).String(),
				"request_id": middleware.GetReqID(r.Context()),
				"remote":     r.RemoteAddr,
			}).Info("HTTP request")
		}()

		next.ServeHTTP(ww, r)
	})
}
