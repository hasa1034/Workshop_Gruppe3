// Package repository kapselt den Datenbankzugriff mit GORM.
package repository

import (
	"fmt"

	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect oeffnet eine GORM-Verbindung zu PostgreSQL anhand des uebergebenen DSN.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("verbindung zu postgres fehlgeschlagen: %w", err)
	}
	return db, nil
}

// Migrate stellt das Datenbankschema bereit: das Schema "kiosk", den Enum-Typ
// fuer das Geschlecht und die Tabellen fuer Kiosk, Betreiber und Produkt.
func Migrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE SCHEMA IF NOT EXISTS kiosk`).Error; err != nil {
		return fmt.Errorf("schema kiosk anlegen fehlgeschlagen: %w", err)
	}

	if err := ensureGeschlechtEnum(db); err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.Betreiber{}, &model.Kiosk{}, &model.Produkt{}); err != nil {
		return fmt.Errorf("automigrate fehlgeschlagen: %w", err)
	}
	return nil
}

// ensureGeschlechtEnum legt den PostgreSQL-Enum-Typ kiosk.geschlecht an, falls er
// noch nicht existiert. AutoMigrate erzeugt benutzerdefinierte Typen nicht selbst.
func ensureGeschlechtEnum(db *gorm.DB) error {
	const stmt = `
DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1 FROM pg_type t
		JOIN pg_namespace n ON n.oid = t.typnamespace
		WHERE t.typname = 'geschlecht' AND n.nspname = 'kiosk'
	) THEN
		CREATE TYPE kiosk.geschlecht AS ENUM ('MAENNLICH', 'WEIBLICH', 'DIVERS');
	END IF;
END
$$;`

	if err := db.Exec(stmt).Error; err != nil {
		return fmt.Errorf("enum kiosk.geschlecht anlegen fehlgeschlagen: %w", err)
	}
	return nil
}
