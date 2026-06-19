// Command server startet den Kiosk-HTTP-Server.
package main

import (
	"log/slog"
	"os"
	"fmt"
	"github.com/hasa1034/Workshop_Gruppe3/internal/app"
)

func main() {
		banner := `
 _  __ _           _
| |/ /(_) ___  ___| | __
| ' / | |/ _ \/ __| |/ /
| . \ | | (_) \__ \   <
|_|\_\|_|\___/|___/_|\_\
   Kiosk-Server · Go + GORM
`
	fmt.Println(banner)
	if err := app.Run(); err != nil {
		slog.Error("server beendet mit fehler", "err", err)
		os.Exit(1)
	}
}
