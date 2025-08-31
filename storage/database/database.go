package database

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	pgx "github.com/jackc/pgx/v5"
	"github.com/zarinit-routers/cloud-connector/storage/repository"
)

const (
	ENV_CONNECTION_STRING = "DATABASE_CONNECTION_STRING"
	ENV_AUTO_MIGRATE      = "DATABASE_AUTO_MIGRATE"
)

var (
	connection *pgx.Conn
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
	connectionString, err := getConnectionString()
	if err != nil {
		return fmt.Errorf("bad connection string: %s", err)
	}
	log.Info("Connecting to database", "connectionString", connectionString)

	if conn, err := pgx.Connect(context.Background(), connectionString); err != nil {
		return err
	} else {
		connection = conn
	}

	repository.Setup(connection)

	if shouldMigrate() {
		if err := migrateDb(); err != nil {
			return fmt.Errorf("migrations failed: %s", err)
		}
	}

	return nil
}
