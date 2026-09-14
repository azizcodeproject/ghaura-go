package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"time"

	"github.com/azizcodeproject/ghaura-go/apps/api/migrations"
)

func RunMigrations(database *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if _, err := database.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(64) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	upFiles, err := listUpMigrations()
	if err != nil {
		return err
	}

	for _, fileName := range upFiles {
		version := strings.TrimSuffix(fileName, ".up.sql")
		applied, err := isMigrationApplied(ctx, database, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		script, err := fs.ReadFile(migrations.FS, fileName)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", fileName, err)
		}

		if err := applyMigration(ctx, database, version, string(script)); err != nil {
			return err
		}
	}

	return nil
}

func listUpMigrations() ([]string, error) {
	entries, err := fs.ReadDir(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("list migrations: %w", err)
	}

	fileNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".up.sql") {
			fileNames = append(fileNames, entry.Name())
		}
	}
	sort.Strings(fileNames)
	return fileNames, nil
}

func isMigrationApplied(ctx context.Context, database *sql.DB, version string) (bool, error) {
	var exists int
	err := database.QueryRowContext(ctx, `SELECT 1 FROM schema_migrations WHERE version = $1`, version).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check migration %s: %w", version, err)
	}
	return true, nil
}

func applyMigration(ctx context.Context, database *sql.DB, version, script string) error {
	tx, err := database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration %s: %w", version, err)
	}

	if _, err := tx.ExecContext(ctx, script); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("apply migration %s: %w", version, err)
	}

	if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("record migration %s: %w", version, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration %s: %w", version, err)
	}

	return nil
}
