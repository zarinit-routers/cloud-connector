package database

import (
	"database/sql"

	"github.com/charmbracelet/log"

	migrate "github.com/rubenv/sql-migrate"
)

func migrateDb(db *sql.DB) error {
	source := migrate.FileMigrationSource{
		Dir: "./sql/migrations",
	}
	log.Info("Migrating database...")
	count, err := migrate.Exec(db, "postgres", source, migrate.Up)
	if err != nil {
		return err
	}
	log.Info("Database migrated", "migrations", count)
	return nil
}
