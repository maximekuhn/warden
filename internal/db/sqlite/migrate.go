package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"io/fs"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Migrate accepts a current version number and applies all required migrations.
// Only sql files with a prefix number > currentVerNum will be executed.
func Migrate(db *sql.DB, currentVerNum int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = migrate(ctx, tx, currentVerNum)
	return err
}

func migrate(ctx context.Context, tx *sql.Tx, currentVerNum int) error {
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	for _, mf := range migrationFiles {
		if mf.prefixNumber <= currentVerNum {
			// migration already applied
			continue
		}

		path := filepath.Join("migrations", mf.filename)
		sqlBytes, err := migrations.ReadFile(path)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, string(sqlBytes))
		if err != nil {
			return err
		}
	}
	newVersion := migrationFiles[len(migrationFiles)-1].prefixNumber
	return updateVersionInMetadataTable(ctx, tx, newVersion)
}

func updateVersionInMetadataTable(ctx context.Context, tx *sql.Tx, newVersion int) error {
	query := `
    INSERT INTO migrations_metadata (applied_datetime, current_version)
    VALUES (?, ?)
    `
	_, err := tx.ExecContext(ctx, query, time.Now(), newVersion)
	return err
}

func getMigrationFiles() ([]migrationFile, error) {
	entries, err := fs.ReadDir(migrations, "migrations")
	if err != nil {
		return nil, err
	}

	migrationFiles := make([]migrationFile, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if strings.HasSuffix(entry.Name(), ".sql") {
			mf, err := parseMigrationFile(entry.Name())
			if err != nil {
				// for now, if a file in the migrations dir is not a valid
				// migration file (by name), we return an error.
				// In the future, we might adapt it to log the error and continue
				// with the remaining files.
				return nil, err
			}
			migrationFiles = append(migrationFiles, mf)
		}
	}

	sort.Slice(migrationFiles, func(i, j int) bool {
		return migrationFiles[i].prefixNumber < migrationFiles[j].prefixNumber
	})

	return migrationFiles, nil
}

type migrationFile struct {
	prefixNumber int
	filename     string
}

func parseMigrationFile(filename string) (migrationFile, error) {
	mf := migrationFile{
		prefixNumber: 0,
		filename:     filename,
	}
	versionStr := strings.Split(filename, "_")[0]
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		return mf, err
	}
	mf.prefixNumber = version
	return mf, nil
}
