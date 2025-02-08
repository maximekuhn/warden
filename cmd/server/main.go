package main

import (
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/maximekuhn/warden/internal/apps/web/server"
	"github.com/maximekuhn/warden/internal/infra/db/sqlite"
)

//go:embed banner.txt
var banner string

func main() {
	fmt.Println(banner)

	logFile, err := os.OpenFile("warden-logs.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	defer logFile.Close()
	logsOutput := io.MultiWriter(os.Stdout, logFile)
	l := slog.New(slog.NewJSONHandler(logsOutput, nil))

	db := setupDB()
	defer db.Close()

	config := parseConfig()

	s := server.NewServer(l, db, config)
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

func parseConfig() *server.Config {
	if len(os.Args) < 3 {
		log.Fatal("Configuration file path is required as the second argument.")
	}
	configFilePath := os.Args[2]
	conf, err := server.ParseConfigFromFile(configFilePath)
	if err != nil {
		log.Fatalf("Failed to parse configuration file: %s", err)
	}
	return conf
}
