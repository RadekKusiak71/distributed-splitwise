package config

import (
	"os"
)

type Config struct {
	JWTSecret          string
	IdentityServiceURL string
}

func Load() *Config {
	return &Config{
		JWTSecret:          GetEnv("JWT_SECRET", "mysecretkey"),
		IdentityServiceURL: GetEnv("IDENTITY_SERVICE_URL", "http://localhost:8081"),
	}
}

func GetEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}
