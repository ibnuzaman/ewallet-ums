package helpers

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	once sync.Once
	env  map[string]string
	mu   sync.RWMutex
)

func SetupConfig() error {
	var err error
	once.Do(func() {
		// Load .env file if exists
		if err = godotenv.Load(".env"); err != nil {
			// .env file is optional, so just log warning
			if Logger != nil {
				Logger.Warn("No .env file found, using system environment variables")
			}
		}

		// Load all environment variables into map
		mu.Lock()
		env = make(map[string]string)
		for _, e := range os.Environ() {
			// Parse environment variables
			var key, value string
			for i := 0; i < len(e); i++ {
				if e[i] == '=' {
					key = e[:i]
					value = e[i+1:]
					break
				}
			}
			if key != "" {
				env[key] = value
			}
		}
		mu.Unlock()

		if Logger != nil {
			Logger.Info("Configuration loaded successfully")
		}
	})

	return err
}

func GetEnv(key string, defaultVal string) string {
	mu.RLock()
	defer mu.RUnlock()

	// Try from map first
	if val, exists := env[key]; exists && val != "" {
		return val
	}

	// Fallback to os.Getenv
	if val := os.Getenv(key); val != "" {
		return val
	}

	return defaultVal
}

func GetRequiredEnv(key string) (string, error) {
	val := GetEnv(key, "")
	if val == "" {
		if Logger != nil {
			Logger.Errorf("Required environment variable %s is not set", key)
		}
		return "", os.ErrNotExist
	}
	return val, nil
}
