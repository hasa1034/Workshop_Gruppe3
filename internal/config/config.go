// Package config laedt die Laufzeitkonfiguration aus Umgebungsvariablen.
package config

import (
	"fmt"
	"os"
)

// Config buendelt die Einstellungen fuer den HTTP-Server und die Datenbank.
type Config struct {
	// ServerPort ist der Port, auf dem der HTTP-Server lauscht.
	ServerPort string
	// DatabaseURL ist der vollstaendige PostgreSQL-DSN fuer GORM.
	DatabaseURL string
}

// Load liest die Konfiguration aus Umgebungsvariablen und faellt auf sinnvolle
// Standardwerte zurueck. Ist DATABASE_URL gesetzt, wird dieser DSN bevorzugt,
// andernfalls wird er aus den einzelnen PG*-Variablen zusammengesetzt.
func Load() Config {
	return Config{
		ServerPort:  getEnv("SERVER_PORT", "8081"),
		DatabaseURL: databaseURL(),
	}
}

func databaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("PGHOST", "localhost"),
		getEnv("PGPORT", "5432"),
		getEnv("PGUSER", "postgres"),
		getEnv("PGPASSWORD", "postgres"),
		getEnv("PGDATABASE", "kiosk"),
		getEnv("PGSSLMODE", "disable"),
	)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
