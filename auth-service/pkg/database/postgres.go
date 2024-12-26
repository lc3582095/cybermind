package database

import (
    "fmt"
    "cybermind/auth-service/configs"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func NewPostgresDB(config *configs.DatabaseConfig) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
        config.Host,
        config.User,
        config.Password,
        config.DBName,
        config.Port,
        config.SSLMode,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return db, nil
} 