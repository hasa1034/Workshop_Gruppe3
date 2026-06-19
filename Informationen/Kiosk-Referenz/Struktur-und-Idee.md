# Struktur und Idee fuer die Go-Implementierung

Diese Dokumentation beschreibt die fachliche Struktur des Kiosk-Beispiels fuer
eine prototypische Implementierung mit Go.

## Zielbild

Die Anwendung verwaltet Kioske. Zu jedem Kiosk wird genau ein Betreiber und
eine Liste von Produkten gespeichert. Die erste Ausbaustufe konzentriert sich auf Lesen und
Neuanlegen ueber REST, Validierung beim Neuanlegen und Speicherung in
PostgreSQL.

## Schichtenstruktur

- `cmd/server`: Startpunkt der Anwendung, Laden der Konfiguration und Start des
  HTTP-Servers.
- `internal/model`: Go-Structs fuer `Kiosk`, `Betreiber`, `Produkt` und
  `Geschlecht`.
- `internal/repository`: Datenbankzugriff mit GORM.
- `internal/service`: Fachlogik, Transaktionen, Dublettenpruefung und
  Fehlerabbildung.
- `internal/http`: Router, Handler, Request-/Response-DTOs und Statuscodes.
- `internal/validation`: zentrale Validierung fuer Create-Requests.

## Entity-Idee

- `Kiosk` ist die zentrale Entity mit Stammdaten, Version, E-Mail und
  Zeitstempeln.
- `Betreiber` beschreibt die Person, die genau einem Kiosk als Betreiber
  zugeordnet wird.
- `Produkt` beschreibt ein Produkt mit Preis und Waehrung.
- Ein Kiosk hat genau einen Betreiber und mehrere Produkte.
- Die fachliche Richtung geht nur vom Kiosk zum Betreiber und zu den Produkten.
- Betreiber und Produkt enthalten kein Kiosk-Objekt als Rueckrichtung.
- Beim Loeschen eines Kiosks sollen die zugehoerigen Produkte ebenfalls
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

- Ungueltiges JSON oder ungueltige Create-Daten: `400 Bad Request`.
- Nicht vorhandene ID beim Lesen: `404 Not Found`.
- Bereits vorhandene E-Mail beim Neuanlegen: `409 Conflict`.
- Erfolgreiches Lesen: `200 OK`.
- Erfolgreiches Neuanlegen: `201 Created`.

## Datenbankzugriff

- PostgreSQL ist die relationale Datenbank.
- GORM bildet die Go-Structs auf Tabellen ab.
- `kiosk.betreiber_id` bildet die 1:1-Beziehung vom Kiosk zum Betreiber ab.
- `produkt.kiosk_id` bildet die 1:N-Beziehung vom Kiosk zu Produkten ab.
- Lesen benoetigt keine explizite Transaktion.
- Neuanlegen wird in einer Transaktion ausgefuehrt, damit Kiosk, Betreiber und
  Produkte gemeinsam gespeichert werden.
- Die E-Mail des Kiosks ist eindeutig und wird vor dem Speichern geprueft.
