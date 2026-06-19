# Aufgabenverteilung

Diese Datei dokumentiert die interne Aufteilung der Arbeit. Die finale
Abgabezusammenfassung steht in der Root-README `ReadMe.md`.

## Sam Haghighi

- Entity-Modell für `Kiosk`, `Betreiber` und `Produkt` gepflegt.
- Go-Structs und GORM-Mapping für die Beziehungen definiert.
- Codebasis für Entitys unter `internal/model` angelegt und mit der
  Referenzstruktur abgeglichen.
- Dokumentation und Root-README aktualisiert.
- Darauf geachtet, dass `Kiosk` genau einen `Betreiber` und mehrere `Produkte`
  referenziert.
- OpenAI Codex lokal im Repository genutzt.

## Ali Arslan

- REST-Router und Handler für Lesen und Neuanlegen umgesetzt.
- Request-DTOs für den Kiosk-Create-Flow definiert.
- Validierung beim Neuanlegen mit `validator` angebunden.
- Fehlerfälle in passende HTTP-Statuscodes übersetzt.
- Handler an Service-Signaturen mit `context.Context` und
  `repository.KioskFilter` angepasst.
- Keycloak- und Bruno-Vorbereitung dokumentiert.
- JWT-Middleware wieder aus dem Router entfernt bzw. nicht weiter benötigt.
- Claude ergänzend für seinen Verantwortungsbereich verwendet.

## Efe Yueksel

- PostgreSQL-Anbindung und GORM-Konfiguration umgesetzt.
- Repository- und Service-Logik für Lesen und Neuanlegen gebaut.
- Transaktion für das gemeinsame Speichern von Kiosk, Betreiber und Produkten
  abgesichert.
- Integrationstests gegen PostgreSQL erstellt.
- PostgreSQL-Compose-Datei bereitgestellt.

## Gemeinsamer Abschluss

- Projektstruktur angelegt und bereinigt.
- `go.mod` mit den benötigten Bibliotheken gepflegt.
- PostgreSQL- und Keycloak-Compose-Dateien unter `extras/compose` abgelegt.
- Bruno-Collection für manuelle REST-Tests unter `extras/bruno/Kiosk-Go`
  ergänzt.
- Prompt-Dokumentation unter `Informationen/Prompts-Kiosk-Go-Workshop.md`
  ergänzt.
- Root-README mit dem tatsächlichen Code-Stand synchronisiert.
