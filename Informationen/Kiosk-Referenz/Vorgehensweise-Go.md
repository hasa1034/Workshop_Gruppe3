# Vorgehensweise für Go

## Schritt 1: Projektstruktur

- Go-Modul ist initialisiert.
- Ordner für Startpunkt, Modelle, Repository, Service, HTTP-Handler,
  Validierung und Tests sind angelegt.
- Konfiguration für Server-Port und PostgreSQL-Verbindung liegt unter
  `internal/config`.

## Schritt 2: Modelle

- `Kiosk`, `Betreiber`, `Produkt` und `Geschlecht` sind als Go-Typen
  definiert.
- GORM-Tags für Tabellen, Spalten, Indizes und Beziehungen sind gesetzt.
- Beziehung von `Kiosk` zu `Betreiber` ist als einzelnes Struct abgebildet.
- Beziehung von `Kiosk` zu `Produkt` ist als Slice abgebildet.
- `Betreiber` und `Produkt` besitzen keine fachliche Go-Rücknavigation zu
  `Kiosk`.

## Schritt 3: Datenbank

- PostgreSQL-Treiber für GORM ist konfiguriert.
- Verbindung wird beim Start der Anwendung geöffnet.
- `kiosk.betreiber_id` bildet die 1:1-Beziehung ab.
- `produkt.kiosk_id` bildet die 1:N-Beziehung ab.
- Schema `kiosk`, Enum `kiosk.geschlecht` und Tabellen werden per Migration
  bereitgestellt.

## Schritt 4: REST-Router und Handler

- Router ist mit `net/http` und `github.com/go-chi/chi/v5` umgesetzt.
- Handler für `GET /kioske/{id}`, `GET /kioske` und `POST /kioske` sind
  vorhanden.
- JSON wird gelesen und geschrieben.
- `Location`-Header wird beim Neuanlegen gesetzt.

## Schritt 5: Validierung

- Create-DTOs für Kiosk, Betreiber und Produkt sind definiert.
- Im Create-DTO ist `betreiber` ein einzelnes Objekt und `produkte` ein Array.
- `github.com/go-playground/validator/v10` prüft Pflichtfelder, E-Mail,
  Längen, Geschlecht, URL und Währung.
- Validierung wird nur im POST-Flow ausgeführt.

## Schritt 6: Service und Repository

- Repository kapselt alle GORM-Zugriffe.
- Service prüft fachliche Regeln wie eindeutige E-Mail.
- Service startet die Transaktion für das Neuanlegen.
- Handler wandeln Service-Fehler in HTTP-Statuscodes um.

## Schritt 7: Integrationstest

- Testdatenbank wird über PostgreSQL vorbereitet.
- Tests laufen gegen Service und Repository.
- Erfolgreiches Anlegen eines Kiosks wird geprüft.
- Lesen per ID für den neu angelegten Datensatz wird geprüft.
- Doppelte E-Mail-Adressen, nicht vorhandene IDs und Listenfilter werden
  geprüft.
