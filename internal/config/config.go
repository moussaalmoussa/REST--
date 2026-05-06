package config

import (
    "github.com/joho/godotenv"
    "os"
)

type Config struct {
    ServerAddr string
    DB         DBConfig
}

type DBConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    DBName   string
    SSLMode  string
}

func Load() (*Config, error) {
    _ = godotenv.Load()
    return &Config{
        ServerAddr: getEnv("SERVER_ADDR", ":8080"),
        DB: DBConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", "postgres"),
            DBName:   getEnv("DB_NAME", "subscriptions"),
            SSLMode:  getEnv("DB_SSLMODE", "disable"),
        },
    }, nil
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}
