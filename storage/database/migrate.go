package database

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	"github.com/zarinit-routers/cloud-connector/sql/migrations"
)

func migrateDb() error {
	log.Info("Migrating database...")
	files, err := migrations.MigrationsFS.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed read embed directory .: %s", err)
	}
	for _, migrationFile := range files {
		if migrationFile.IsDir() {
			continue
		}
		log.Info("Migration started", "name", migrationFile.Name())
		content, err := migrations.MigrationsFS.ReadFile(migrationFile.Name())
		if err != nil {
			return fmt.Errorf("failed read embed file %q: %s", migrationFile.Name(), err)
		}

		tag, err := connection.Exec(context.TODO(), string(content))
		if err != nil {
			log.Error("Failed execute migration",
				"error", err,
				"migrationFile", migrationFile.Name(),
				"migrationFileContent", content,
				"tag/status", tag.String(),
			)
			return fmt.Errorf("failed execute migration %q: %s", migrationFile.Name(), err)
		}
	}
	log.Info("Database migrated")
	return nil

}
