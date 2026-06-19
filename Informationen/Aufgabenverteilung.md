# Aufgabenverteilung

## Sam Haghighi

- Entity-Modell fuer `Kiosk`, `Betreiber` und `Produkt` pflegen.
- Go-Structs und GORM-Mapping fuer die Beziehungen definieren.
- Codebasis fuer Entitys unter `internal/model` anlegen und aktuell halten.
- Informationsdokumentation und README aktuell halten.
- Darauf achten, dass `Kiosk` genau einen `Betreiber` und mehrere `Produkte`
  referenziert.

## Ali Arslan

- REST-Router und Handler fuer Lesen und Neuanlegen umsetzen.
- Request- und Response-DTOs fuer den Kiosk-Flow definieren.
- Validierung beim Neuanlegen mit `validator` anbinden.
- Fehlerfaelle in passende HTTP-Statuscodes uebersetzen.

## Efe Yueksel

- PostgreSQL-Anbindung und GORM-Konfiguration umsetzen.
- Repository- und Service-Logik fuer Lesen und Neuanlegen bauen.
- Transaktion fuer das gemeinsame Speichern von Kiosk, Betreiber und Produkten
  absichern.
- Integrationstests und CI-Checks betreuen.

## Phase 2: Offene Aufgaben

### Sam Haghighi

- `go.mod` um `gorm.io/gorm` und `gorm.io/driver/postgres` ergaenzen.
- `debug.log` entfernen und `.gitignore` anlegen.
- Danach pruefen:
  - `go mod tidy`
  - `go test ./...`
  - `git status --short`

### Ali Arslan

- Handler an Efe-Service anpassen.
- Handler-Interface mit `context.Context` und `repository.KioskFilter`
  kompatibel machen.
- Create-Request in `model.Kiosk` umwandeln.
- Eigene Handler-Fehler entfernen und `service.ErrNotFound` sowie
  `service.ErrEmailExists` verwenden.

### Efe Yueksel

- PostgreSQL-Compose-Datei fuellen.
- `cmd/server/main.go` erstellen.
- Config, DB, Migration, Repository, Service, Handler und Router verdrahten.
- Integrationstest mit echter DB pruefen.

### Gemeinsamer Abschluss

- `gofmt`
- `go mod tidy`
- `go test ./...`
- README aktualisieren, falls Ports oder Befehle abweichen.
