# Programmierworkshop am 19.06.2026

## Namen

Sam Haghighi, Ali Arslan, Efe Yueksel

## Link zum Git-Repository

https://github.com/hasa1034/Workshop_Gruppe3

## Kurzbeschreibung

Dieses Repository enthält eine prototypische Go-Implementierung für eine
Kiosk-Verwaltung. Die Anwendung stellt eine REST-Schnittstelle bereit, mit der
Kioske gelesen und neu angelegt werden können. Ein Kiosk besitzt genau einen
Betreiber und mehrere Produkte. Die Daten werden mit GORM in PostgreSQL
gespeichert.

Keycloak ist als optionale lokale Zusatzumgebung vorbereitet, wird im
HTTP-Server aber nicht verpflichtend für die Endpunkte erzwungen.

## KI-Werkzeuge

OpenAI Codex wurde lokal im Repository verwendet. Claude wurde ergänzend von
Ali Arslan für seinen Verantwortungsbereich eingesetzt.

### Agenten

OpenAI Codex (Sam Haghighi)
Claude (Ali Arslan)

### Chat-URLs

Keine separate Chat-URL. Die Arbeit erfolgte lokal mit OpenAI Codex direkt im
Repository.

## Frameworks und Bibliotheken

- Go 1.25
- `net/http` als HTTP-Server
- `github.com/go-chi/chi/v5` als Router
- `github.com/go-playground/validator/v10` für Request-Validierung
- `gorm.io/gorm` mit `gorm.io/driver/postgres` für PostgreSQL
- `github.com/shopspring/decimal` für Produktpreise
- Docker Compose für PostgreSQL und optional Keycloak
- Bruno-Collection unter `extras/bruno/Kiosk-Go` zum manuellen Testen der
  REST-Endpunkte

## Umgesetzte Funktionen

- `GET /kioske` liest alle Kioske.
- `GET /kioske?name=...` filtert nach Name.
- `GET /kioske?email=...` filtert nach E-Mail.
- `GET /kioske/{id}` liest einen einzelnen Kiosk.
- `POST /kioske` legt einen Kiosk inklusive Betreiber und Produkten an.
- Beim erfolgreichen `POST` wird `201 Created` mit `Location: /kioske/{id}`
  zurückgegeben.
- Ungültiges JSON oder ungültige Eingabedaten liefern `400 Bad Request`.
- Nicht vorhandene IDs liefern `404 Not Found`.
- Doppelte E-Mail-Adressen liefern `409 Conflict`.
- Das Datenbankschema `kiosk` und der PostgreSQL-Enum `kiosk.geschlecht`
  werden beim Start automatisch angelegt.

## Projektstruktur

```text
cmd/server                 alternativer Startpunkt für den HTTP-Server
cmd/seed                   Seed-Kommando für Beispieldaten
internal/app               Verdrahtung von Konfiguration, DB, Service und HTTP
internal/config            Umgebungsvariablen und Standardwerte
internal/handler           HTTP-Handler und DTO-zu-Modell-Mapping
internal/http              chi-Router
internal/model             GORM-Modelle Kiosk, Betreiber, Produkt, Geschlecht
internal/repository        PostgreSQL/GORM-Zugriff und Migration
internal/service           Fachlogik, Transaktion, Dublettenprüfung
internal/validation        Create-Request-DTOs und Validator
extras/compose/postgres    lokale PostgreSQL-Compose-Datei
extras/compose/keycloak    optionale Keycloak-Compose-Datei
extras/bruno/Kiosk-Go      Bruno-Requests für manuelle API-Tests
Informationen              ergänzende Aufgaben-, Struktur- und Referenzdoku
```

## Datenmodell

- `Kiosk` enthält Name, E-Mail, Öffnungsstatus, Homepage, Username,
  Zeitstempel und Version.
- `Betreiber` enthält Vorname, Nachname und optionales Geschlecht.
- `Produkt` enthält Name, Preis, Währung und die technische Zuordnung zum
  Kiosk.
- Ein `Kiosk` hat genau einen `Betreiber`.
- Ein `Kiosk` hat mehrere `Produkte`.
- `Betreiber` und `Produkt` enthalten keine fachliche Rücknavigation zum
  `Kiosk`.

## Konfiguration

Ohne Umgebungsvariablen verwendet die Anwendung diese Standardwerte:

```text
SERVER_PORT=8081
PGHOST=localhost
PGPORT=5432
PGUSER=postgres
PGPASSWORD=postgres
PGDATABASE=kiosk
PGSSLMODE=disable
```

Alternativ kann `DATABASE_URL` gesetzt werden. Wenn `DATABASE_URL` vorhanden
ist, wird dieser vollständige PostgreSQL-DSN bevorzugt.

## Start

PostgreSQL starten:

```powershell
docker compose -f extras/compose/postgres/compose.yml up -d
```

Server starten:

```powershell
go run .
```

Alternativ kann der Server über den Startpunkt in `cmd/server` gestartet
werden:

```powershell
go run ./cmd/server
```

Der Server läuft standardmäßig auf:

```text
http://localhost:8081
```

Beispieldaten einspielen:

```powershell
go run ./cmd/seed
```

PostgreSQL wieder stoppen:

```powershell
docker compose -f extras/compose/postgres/compose.yml down
```

PostgreSQL inklusive Volume löschen:

```powershell
docker compose -f extras/compose/postgres/compose.yml down -v
```

## Beispiel-Requests

Alle Kioske lesen:

```powershell
Invoke-RestMethod http://localhost:8081/kioske
```

Kiosk mit ID 1 lesen:

```powershell
Invoke-RestMethod http://localhost:8081/kioske/1
```

Kiosk anlegen:

```powershell
Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8081/kioske `
  -ContentType "application/json" `
  -Body '{
    "name": "Neuer Kiosk",
    "email": "neu@kiosk.de",
    "username": "neuuser",
    "betreiber": {
      "vorname": "Max",
      "nachname": "Mustermann"
    },
    "produkte": [
      { "name": "Kaffee", "preis": "2.50", "waehrung": "EUR" }
    ]
  }'
```

Die gleichen Requests sind auch als Bruno-Collection unter
`extras/bruno/Kiosk-Go` vorbereitet.

## Tests

Die Integrationstests laufen gegen PostgreSQL. Dafür muss die lokale
PostgreSQL-Compose-Umgebung laufen oder `DATABASE_URL` auf eine erreichbare
Testdatenbank zeigen.

```powershell
docker compose -f extras/compose/postgres/compose.yml up -d
go test .
go test ./...
```

Die Tests migrieren das Schema automatisch und räumen die Tabellen vor jedem
Test mit `TRUNCATE ... RESTART IDENTITY CASCADE` auf.

## Optional: Keycloak

Keycloak kann separat im Dev-Modus gestartet werden:

```powershell
docker compose -f extras/compose/keycloak/compose.yml up -d
```

Admin-UI:

```text
http://localhost:8880
Benutzer: admin
Passwort: admin
```

Die Compose-Datei dient als vorbereitete Zusatzumgebung. Die REST-Endpunkte des
Go-Servers sind in dieser Abgabe ohne verpflichtende Token-Prüfung nutzbar.

## Aufgabenverteilung

Sam Haghighi:

- Entity-Modell für `Kiosk`, `Betreiber` und `Produkt`
- Go-Structs und GORM-Mapping
- Dokumentation und README
- KI-gestützte Repository-Arbeit mit OpenAI Codex

Ali Arslan:

- REST-Router und Handler
- Request-/Response-Verarbeitung
- Validierung beim Neuanlegen
- HTTP-Statuscodes und Fehlerantworten
- Keycloak-/Bruno-Vorbereitung in seinem Verantwortungsbereich
- Ergänzende KI-Unterstützung mit Claude

Efe Yueksel:

- PostgreSQL-Anbindung
- Repository- und Service-Logik
- Transaktion für das gemeinsame Speichern
- Integrationstests und Testdatenbank-Setup

## Prompts/Requests an KI-Agenten

- Relevante Kiosk-Referenzdokumente auswerten und in eine Go-Struktur
  übertragen.
- Entity-Muster `Kiosk`, `Betreiber`, `Produkt` für Go, REST, Validierung und
  GORM/PostgreSQL anpassen.
- Beziehung korrigieren: Kiosk zu Betreiber ist 1:1, Kiosk zu Produkt ist 1:N,
  ohne fachliche Rückrichtung.
- PostgreSQL-Compose, Migration, Repository, Service, Handler und Router
  miteinander abgleichen.
- README für die Abgabe fertigstellen und mit dem tatsächlichen Code-Stand
  synchronisieren.
- Ali-bezogene Prompts zu PostgreSQL, Keycloak, Bruno, JWT-Middleware und
  Integrationstests sind zusätzlich in
  `Informationen/Prompts-Kiosk-Go-Workshop.md` dokumentiert.
