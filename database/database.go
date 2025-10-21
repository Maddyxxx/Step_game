package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// go:embed migrations/001_init_userstate.sql
var initialSchema string

// InitDBMust - Must-функция для инициализации БД с автоматической миграцией
func InitDBMust(path string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", path)
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %s", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("failed to ping database: %s", err))
	}

	// Применяем начальную миграцию
	_, err = db.Exec(initialSchema)
	if err != nil {
		panic(fmt.Sprintf("failed to create tables: %s", err))
	}

	log.Printf("database %s initialized successfully", path)
	return db
}
