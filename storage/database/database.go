package database

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
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

func migrate() error {
	log.Info("Migrating database")
	return fmt.Errorf("migrations not implemented")
}

func Setup() error {
	conn, err := getConnectionString()
	if err != nil {
		return fmt.Errorf("bad connection string: %s", err)
	}
	log.Info("Connecting to database", "connectionString", conn)

	if shouldMigrate() {
		if err := migrate(); err != nil {
			return fmt.Errorf("migrations failed: %s", err)
		}
	}

	return nil
}
