package database

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

const (
	ENV_CONNECTION_STRING = "DATABASE_CONNECTION_STRING"
)

func getConnectionString() (string, error) {
	conn := os.Getenv(ENV_CONNECTION_STRING)
	if conn == "" {
		return "", fmt.Errorf("environment variable %q is not set or empty", ENV_CONNECTION_STRING)
	}

	return conn, nil
}
func Setup() error {
	conn, err := getConnectionString()
	if err != nil {
		return fmt.Errorf("bad connection string: %s", err)
	}
	log.Info("Connecting to database", "connectionString", conn)
	return nil
}
