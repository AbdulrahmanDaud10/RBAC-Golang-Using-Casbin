package repository

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQLDataBaseConnection -> return db instance
func PostgreSQLDataBaseConnection() (*gorm.DB, error) {
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		username = os.Getenv("DB_USER")
		database = os.Getenv("DB_NAME")
		password = os.Getenv("DB_PASSWORD")
	)

	dsn := fmt.Sprintf("host=%s port=%s user=%s database=%s  password=%s sslmode=disable",
		host,
		port,
		username,
		database,
		password,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

}
