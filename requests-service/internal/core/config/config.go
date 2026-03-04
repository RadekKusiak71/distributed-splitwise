package config

import (
	"os"
	"strconv"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type APIConfig struct {
	Port int
}

type JWTConfig struct {
	SecretKey string
}

type AWSConfig struct {
	AccessKey  string
	SecretKey  string
	Region     string
	BucketName string
	S3BaseURL  string
}

type Config struct {
	API APIConfig
	DB  DBConfig
	JWT JWTConfig
	AWS AWSConfig
}

func Load() *Config {
	return &Config{
		API: APIConfig{
			Port: GetEnvAsInt("GO_PORT", 80),
		},
		DB: DBConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnvAsInt("DB_PORT", 5432),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "postgres"),
			Name:     GetEnv("DB_NAME", "database"),
		},
		JWT: JWTConfig{
			SecretKey: GetEnv("JWT_SECRET_KEY", "very-secret"),
		},
		AWS: AWSConfig{
			AccessKey:  GetEnv("AWS_ACCESS_KEY", "test"),
			SecretKey:  GetEnv("AWS_SECRET_KEY", "test"),
			Region:     GetEnv("AWS_S3_REGION", "eu-central-1"),
			BucketName: GetEnv("AWS_S3_BUCKET_NAME", "test"),
			S3BaseURL:  GetEnv("AWS_S3_BASE_URL", "https://test.s3.eu-central-1.amazonaws.com/"),
		},
	}
}

func GetEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func GetEnvAsInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}
