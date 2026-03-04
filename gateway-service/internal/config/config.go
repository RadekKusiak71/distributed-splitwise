package config

import (
	"os"
)

type Config struct {
	JWTSecret          string
	IdentityServiceURL string
	RequestsServiceURL string
}

func Load() *Config {
	return &Config{
		JWTSecret:          GetEnv("JWT_SECRET", "mysecretkey"),
		IdentityServiceURL: GetEnv("IDENTITY_SERVICE_URL", "http://identity-service:80"),
		RequestsServiceURL: GetEnv("REQUESTS_SERVICE_URL", "http://requests-service:80"),
	}
}

func GetEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
