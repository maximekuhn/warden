package main

import (
	_ "embed"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/maximekuhn/warden/internal/server"
)

//go:embed banner.txt
var banner string

func main() {
	fmt.Println(banner)
	l := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	s := server.NewServer(l, nil) // FIXME: nil db
	log.Fatal(s.Start())
}
