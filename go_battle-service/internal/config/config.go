package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    Environment string
}

func LoadConfig() (*Config, error) {
    if err := godotenv.Load(); err != nil {
        log.Println("Warning: .env file not found, using environment variables")
    }

    cfg := &Config{
        Port:        getEnv("PORT", "8080"),
        Environment: getEnv("ENVIRONMENT", "development"),
    }
    return cfg, nil
}

func getEnv(key, defaultVal string) string {
    if val, exists := os.LookupEnv(key); exists {
        return val
    }
    return defaultVal
}
