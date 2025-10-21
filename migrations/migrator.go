package migrations

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

//go:embed *.sql
var migrationFiles embed.FS

// RunMigrations - применяет все миграции из папки migrations
func RunMigrations(db *sqlx.DB) {
	migrations, err := migrationFiles.ReadDir(".")
	if err != nil {
		panic(fmt.Sprintf("failed to read migrations: %s", err))
	}

	for _, migration := range migrations {
		if migration.IsDir() {
			continue
		}

		content, err := migrationFiles.ReadFile(migration.Name())
		if err != nil {
			panic(fmt.Sprintf("failed to read migration file %s: %s", migration.Name(), err))
		}

		_, err = db.Exec(string(content))
		if err != nil {
			panic(fmt.Sprintf("failed to execute migration %s: %s", migration.Name(), err))
		}

		log.Printf("migration applied: %s", migration.Name())
	}

	log.Printf("all migrations applied successfully")
}
