package main_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/hasa1034/Workshop_Gruppe3/internal/config"
	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"github.com/hasa1034/Workshop_Gruppe3/internal/repository"
	"github.com/hasa1034/Workshop_Gruppe3/internal/service"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// setupDB stellt eine Verbindung zur Test-Datenbank her und migriert das Schema.
// Ist DATABASE_URL nicht gesetzt, nutzt der Test die Standardwerte aus der
// Anwendungskonfiguration. PostgreSQL muss dafuer lokal per Compose laufen.
func setupDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = config.Load().DatabaseURL
	}

	db, err := repository.Connect(dsn)
	if err != nil {
		t.Fatalf("verbindung fehlgeschlagen: %v", err)
	}
	if err := repository.Migrate(db); err != nil {
		t.Fatalf("migration fehlgeschlagen: %v", err)
	}

	// Vor jedem Test eine saubere Ausgangslage schaffen.
	if err := db.Exec(`TRUNCATE kiosk.kiosk, kiosk.betreiber, kiosk.produkt RESTART IDENTITY CASCADE`).Error; err != nil {
		t.Fatalf("aufraeumen fehlgeschlagen: %v", err)
	}

	return db
}

func sampleKiosk(email string) *model.Kiosk {
	geschlecht := model.GeschlechtDivers
	homepage := "https://kiosk.example"
	return &model.Kiosk{
		Name:         "Kiosk am Park",
		Email:        email,
		IstGeoeffnet: true,
		Homepage:     &homepage,
		Username:     "kiosk_park",
		Betreiber: model.Betreiber{
			Vorname:    "Max",
			Nachname:   "Mustermann",
			Geschlecht: &geschlecht,
		},
		Produkte: []model.Produkt{
			{Name: "Kaffee", Preis: decimal.RequireFromString("2.50"), Waehrung: "EUR"},
			{Name: "Brezel", Preis: decimal.RequireFromString("1.20"), Waehrung: "EUR"},
		},
	}
}

func TestCreateAndGet(t *testing.T) {
	db := setupDB(t)
	svc := service.NewKioskService(repository.NewRepository(db))
	ctx := context.Background()

	kiosk := sampleKiosk("park@example.com")
	if err := svc.Create(ctx, kiosk); err != nil {
		t.Fatalf("create fehlgeschlagen: %v", err)
	}
	if kiosk.ID == 0 {
		t.Fatal("erwartete eine generierte Kiosk-ID")
	}
	if kiosk.BetreiberID == 0 {
		t.Fatal("erwartete eine generierte Betreiber-ID")
	}

	got, err := svc.GetByID(ctx, kiosk.ID)
	if err != nil {
		t.Fatalf("getByID fehlgeschlagen: %v", err)
	}
	if got.Email != "park@example.com" {
		t.Errorf("unerwartete e-mail: %q", got.Email)
	}
	if got.Betreiber.Nachname != "Mustermann" {
		t.Errorf("betreiber nicht geladen: %+v", got.Betreiber)
	}
	if len(got.Produkte) != 2 {
		t.Errorf("erwartete 2 produkte, erhielt %d", len(got.Produkte))
	}
}

func TestCreateDuplicateEmail(t *testing.T) {
	db := setupDB(t)
	svc := service.NewKioskService(repository.NewRepository(db))
	ctx := context.Background()

	if err := svc.Create(ctx, sampleKiosk("dup@example.com")); err != nil {
		t.Fatalf("erster create fehlgeschlagen: %v", err)
	}

	err := svc.Create(ctx, sampleKiosk("dup@example.com"))
	if !errors.Is(err, service.ErrEmailExists) {
		t.Fatalf("erwartete ErrEmailExists, erhielt: %v", err)
	}

	// Sicherstellen, dass durch den fehlgeschlagenen zweiten Versuch kein
	// verwaister Betreiber-Datensatz uebrig blieb (Transaktion griff).
	var betreiberCount int64
	if err := db.Table("kiosk.betreiber").Count(&betreiberCount).Error; err != nil {
		t.Fatalf("betreiber zaehlen fehlgeschlagen: %v", err)
	}
	if betreiberCount != 1 {
		t.Errorf("erwartete genau 1 betreiber, erhielt %d", betreiberCount)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	db := setupDB(t)
	svc := service.NewKioskService(repository.NewRepository(db))

	_, err := svc.GetByID(context.Background(), 99999)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("erwartete ErrNotFound, erhielt: %v", err)
	}
}

func TestListFilter(t *testing.T) {
	db := setupDB(t)
	svc := service.NewKioskService(repository.NewRepository(db))
	ctx := context.Background()

	if err := svc.Create(ctx, sampleKiosk("a@example.com")); err != nil {
		t.Fatalf("create fehlgeschlagen: %v", err)
	}
	if err := svc.Create(ctx, sampleKiosk("b@example.com")); err != nil {
		t.Fatalf("create fehlgeschlagen: %v", err)
	}

	all, err := svc.List(ctx, repository.KioskFilter{})
	if err != nil {
		t.Fatalf("list fehlgeschlagen: %v", err)
	}
	if len(all) != 2 {
		t.Errorf("erwartete 2 kioske, erhielt %d", len(all))
	}

	filtered, err := svc.List(ctx, repository.KioskFilter{Email: "b@example.com"})
	if err != nil {
		t.Fatalf("gefilterte list fehlgeschlagen: %v", err)
	}
	if len(filtered) != 1 || filtered[0].Email != "b@example.com" {
		t.Errorf("filter nach e-mail lieferte unerwartetes ergebnis: %+v", filtered)
	}
}
