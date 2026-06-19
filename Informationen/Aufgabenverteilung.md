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
