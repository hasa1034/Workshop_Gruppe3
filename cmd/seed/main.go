// Command seed befüllt die Datenbank mit Testdaten.
// Aufruf: go run ./cmd/seed
package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/hasa1034/Workshop_Gruppe3/internal/config"
	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"github.com/hasa1034/Workshop_Gruppe3/internal/repository"
	"github.com/shopspring/decimal"
)

const banner = `
 _  __ _           _      ___              _
| |/ /(_) ___  ___| | __ / __| ___ ___  __| |
| ' / | |/ _ \/ __| |/ / \__ \/ -_) -_)/ _' |
|_|\_\|_|\___/|___/_|\_\ |___/\___\___|\__,_|
   Testdaten werden eingespielt · Go + GORM
`

func main() {
	fmt.Print(banner)

	cfg := config.Load()

	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Datenbankverbindung fehlgeschlagen", "err", err)
		os.Exit(1)
	}
	slog.Info("Datenbankverbindung hergestellt")

	if err := repository.Migrate(db); err != nil {
		slog.Error("Migration fehlgeschlagen", "err", err)
		os.Exit(1)
	}
	slog.Info("Schema migriert")

	maennlich := model.Geschlecht("MAENNLICH")
	weiblich := model.Geschlecht("WEIBLICH")
	homepage1 := "https://bahnhof-kiosk.de"
	homepage2 := "https://park-kiosk.de"

	kioske := []model.Kiosk{
		{
			Name:         "Hauptbahnhof Kiosk",
			Email:        "hbf@kiosk.de",
			IstGeoeffnet: true,
			Homepage:     &homepage1,
			Username:     "hbfuser",
			Erzeugt:      time.Now(),
			Aktualisiert: time.Now(),
			Betreiber: model.Betreiber{
				Vorname:    "Thomas",
				Nachname:   "Müller",
				Geschlecht: &maennlich,
			},
			Produkte: []model.Produkt{
				{Name: "Kaffee", Preis: decimal.NewFromFloat(2.50), Waehrung: "EUR"},
				{Name: "Wasser", Preis: decimal.NewFromFloat(1.50), Waehrung: "EUR"},
				{Name: "Brezn", Preis: decimal.NewFromFloat(1.80), Waehrung: "EUR"},
			},
		},
		{
			Name:         "Stadtpark Kiosk",
			Email:        "park@kiosk.de",
			IstGeoeffnet: true,
			Homepage:     &homepage2,
			Username:     "parkuser",
			Erzeugt:      time.Now(),
			Aktualisiert: time.Now(),
			Betreiber: model.Betreiber{
				Vorname:    "Anna",
				Nachname:   "Schmidt",
				Geschlecht: &weiblich,
			},
			Produkte: []model.Produkt{
				{Name: "Eis", Preis: decimal.NewFromFloat(3.00), Waehrung: "EUR"},
				{Name: "Limo", Preis: decimal.NewFromFloat(2.00), Waehrung: "EUR"},
			},
		},
		{
			Name:         "Uni Mensa Kiosk",
			Email:        "mensa@kiosk.de",
			IstGeoeffnet: false,
			Username:     "mensauser",
			Erzeugt:      time.Now(),
			Aktualisiert: time.Now(),
			Betreiber: model.Betreiber{
				Vorname:  "Felix",
				Nachname: "Weber",
			},
			Produkte: []model.Produkt{
				{Name: "Sandwich", Preis: decimal.NewFromFloat(4.50), Waehrung: "EUR"},
				{Name: "Saft", Preis: decimal.NewFromFloat(2.20), Waehrung: "EUR"},
				{Name: "Schokoriegel", Preis: decimal.NewFromFloat(1.20), Waehrung: "EUR"},
			},
		},
	}

	var angelegt, uebersprungen int
	for _, k := range kioske {
		if err := db.Create(&k).Error; err != nil {
			slog.Warn("Kiosk übersprungen (existiert bereits?)", "email", k.Email, "err", err)
			uebersprungen++
			continue
		}
		slog.Info("Kiosk angelegt", "name", k.Name, "id", k.ID, "produkte", len(k.Produkte))
		angelegt++
	}

	fmt.Println()
	slog.Info("Seed abgeschlossen", "angelegt", angelegt, "uebersprungen", uebersprungen)
}
