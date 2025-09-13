package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ENV_CONNECTION_STRING = "DATABASE_CONNECTION_STRING"
	ENV_AUTO_MIGRATE      = "DATABASE_AUTO_MIGRATE"
)

func getConnectionString() (string, error) {
	conn := os.Getenv(ENV_CONNECTION_STRING)
	if conn == "" {
		return "", fmt.Errorf("environment variable %q is not set or empty", ENV_CONNECTION_STRING)
	}

	return conn, nil
}

func shouldMigrate() bool {
	envVar := strings.ToLower(os.Getenv(ENV_AUTO_MIGRATE))
	if envVar == "1" || envVar == "yes" || envVar == "true" {
		return true
	}
	return false
}

func Setup() error {

	if !shouldMigrate() {
		return nil
	}
	db, err := Connect()
	if err != nil {
		log.Error("Error connecting to database", "error", err)
		return err
	}
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	if err := migrateDb(sqlDb); err != nil {
		return fmt.Errorf("migrations failed: %s", err)
	}

	return nil
}

func Connect() (*gorm.DB, error) {
	connectionString, err := getConnectionString()
	if err != nil {
		return nil, fmt.Errorf("bad connection string: %s", err)
	}
	log.Info("Connecting to database", "connectionString", connectionString)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
