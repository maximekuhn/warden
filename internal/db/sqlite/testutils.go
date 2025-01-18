package sqlite

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createTmpDb() *sql.DB {
	f, err := os.CreateTemp("", "test-db-*.sqlite3")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("sqlite3", f.Name())
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	return db
}
