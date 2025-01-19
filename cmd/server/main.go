package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/maximekuhn/warden/internal/db/sqlite"
	"github.com/maximekuhn/warden/internal/server"
)

//go:embed banner.txt
var banner string

func main() {
	fmt.Println(banner)
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	db := setupDB()
	defer db.Close()
	s := server.NewServer(l, db)
	log.Fatal(s.Start())
}

func setupDB() *sql.DB {
	if len(os.Args) < 2 {
		log.Fatal("Database file path is required as the first argument.")
	}
	dbFile := os.Args[1]
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Failed to open SQLite3 database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := sqlite.Migrate(db); err != nil {
		db.Close()
		log.Fatalf("Failed to apply migrations: %v", err)
	}
	return db
}
