# Struktur und Idee für die Go-Implementierung

Diese Dokumentation beschreibt die fachliche Struktur des Kiosk-Beispiels für
die prototypische Implementierung mit Go.

## Zielbild

Die Anwendung verwaltet Kioske. Zu jedem Kiosk wird genau ein Betreiber und
eine Liste von Produkten gespeichert. Die umgesetzte Ausbaustufe konzentriert
sich auf Lesen und Neuanlegen über REST, Validierung beim Neuanlegen und
Speicherung in PostgreSQL.

## Schichtenstruktur

- `main.go`: Startpunkt für `go run .`.
- `cmd/server`: alternativer Startpunkt der Anwendung.
- `cmd/seed`: Seed-Kommando für Beispieldaten.
- `internal/model`: Go-Structs für `Kiosk`, `Betreiber`, `Produkt` und
  `Geschlecht`.
- `internal/repository`: Datenbankzugriff mit GORM.
- `internal/service`: Fachlogik, Transaktionen, Dublettenprüfung und
  Fehlerabbildung.
- `internal/http`: Router.
- `internal/handler`: Handler, Request-Verarbeitung und Statuscodes.
- `internal/validation`: zentrale Validierung für Create-Requests.
- `extras/compose`: lokale Compose-Dateien für PostgreSQL und optional
  Keycloak.
- `extras/bruno`: Bruno-Collection für manuelle API-Tests.

## Entity-Idee

- `Kiosk` ist die zentrale Entity mit Stammdaten, Version, E-Mail und
  Zeitstempeln.
- `Betreiber` beschreibt die Person, die genau einem Kiosk als Betreiber
  zugeordnet wird.
- `Produkt` beschreibt ein Produkt mit Preis und Währung.
- Ein Kiosk hat genau einen Betreiber und mehrere Produkte.
- Die fachliche Richtung geht nur vom Kiosk zum Betreiber und zu den Produkten.
- Betreiber und Produkt enthalten kein Kiosk-Objekt als Rückrichtung.
- Beim Löschen eines Kiosks sollen die zugehörigen Produkte ebenfalls
  entfernt werden.

## REST-Fluss

- `GET /kioske/{id}` liest einen Kiosk inklusive genau einem Betreiber und
  mehreren Produkten.
- `GET /kioske` liest eine Liste von Kiosken, optional gefiltert nach einfachen
  Query-Parametern wie `name` oder `email`.
- `POST /kioske` legt einen Kiosk inklusive genau einem Betreiber und Produkten
  an.
- Erfolgreiches Neuanlegen antwortet mit `201 Created` und einem
  `Location`-Header auf den neuen Datensatz.

## Fehlerbehandlung

- Ungültiges JSON oder ungültige Create-Daten: `400 Bad Request`.
- Nicht vorhandene ID beim Lesen: `404 Not Found`.
- Bereits vorhandene E-Mail beim Neuanlegen: `409 Conflict`.
- Erfolgreiches Lesen: `200 OK`.
- Erfolgreiches Neuanlegen: `201 Created`.

## Datenbankzugriff

- PostgreSQL ist die relationale Datenbank.
- GORM bildet die Go-Structs auf Tabellen ab.
- `kiosk.betreiber_id` bildet die 1:1-Beziehung vom Kiosk zum Betreiber ab.
- `produkt.kiosk_id` bildet die 1:N-Beziehung vom Kiosk zu Produkten ab.
- Lesen benötigt keine explizite Transaktion.
- Neuanlegen wird in einer Transaktion ausgeführt, damit Kiosk, Betreiber und
  Produkte gemeinsam gespeichert werden.
- Die E-Mail des Kiosks ist eindeutig und wird vor dem Speichern geprüft.
