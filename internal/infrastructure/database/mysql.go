// Package database
package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/fzndps/eventcheck/config"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.GetDSN()

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return db, nil
}
