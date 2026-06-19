package main

import (
	"log/slog"
	"os"

	"github.com/hasa1034/Workshop_Gruppe3/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		slog.Error("server beendet mit fehler", "err", err)
		os.Exit(1)
	}
}
