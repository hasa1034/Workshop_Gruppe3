# Prompts - Kiosk-Go Workshop

Verwendete Prompts während der Entwicklung, geordnet nach Thema. Die Prompts
beziehen sich auf die prototypische Go-Implementierung mit PostgreSQL,
optionalem Keycloak und Bruno-Requests.

## KI-Werkzeuge

- OpenAI Codex wurde lokal im Repository eingesetzt.
- Claude wurde ergänzend von Ali Arslan für seinen Verantwortungsbereich
  verwendet.

## PostgreSQL

- Muss PostgreSQL für dieses Projekt separat eingerichtet werden, oder kann ich
  die bereits vorhandene Instanz aus dem TypeScript-Projekt wiederverwenden?
- Lass uns zunächst PostgreSQL zum Laufen bringen.
- Gibt es ein Seed-Skript zum Befüllen der Datenbank, analog zu Prisma im
  TypeScript-Projekt? Verwende die gleiche Struktur wie im Kiosk-TS-Projekt.
- Ist die PostgreSQL-Einrichtung abgeschlossen?
- Wie kann ich überprüfen, ob PostgreSQL läuft und das Datenbankschema korrekt
  angelegt wurde?

## Keycloak

- Lass uns jetzt Keycloak einrichten - ohne TLS und unnötige Komplexität.
- Kann der bereits für das TypeScript-Projekt konfigurierte Keycloak-Container
  wiederverwendet werden, anstatt einen neuen aufzusetzen?
- Schreib bitte eine README für die Keycloak-Einrichtung, die zur
  Projektstruktur passt und alle nötigen Schritte dokumentiert.

## Bruno / API-Tests

- Lass uns die API-Endpunkte mit Bruno testen.
- Wie kann ich den Speicherort der Bruno-Kollektion auf `extras/bruno` ändern?
- PostgreSQL, Keycloak und der Server laufen. Wie integriere ich die
  Keycloak-Token-Anfrage in Bruno?
- Kannst du die Bruno-Request-Dateien für die Kollektion erstellen?
- Die übrigen Requests sollen ohne Authentifizierung bleiben - lediglich ein
  separater Request soll den Token von Keycloak holen.

## JWT-Middleware / Authentifizierung

- Sollten wir eine JWT-Middleware zur Absicherung der Routen einbinden?
- Bitte die JWT-Middleware wieder aus dem Router entfernen - die Routen sollen
  vorerst ungeschützt bleiben.
- Wird die Datei `internal/middleware/auth.go` noch benötigt, oder kann sie
  gelöscht werden?

Ergebnis: Die REST-Routen sind in dieser Abgabe ungeschützt. Keycloak bleibt als
optionale lokale Zusatzumgebung vorhanden. Eine JWT-Middleware ist nicht mehr
eingebunden; `internal/middleware/auth.go` wird nicht benötigt und ist nicht im
Repository enthalten.

## Allgemein

- Alle 4 Integrationstests sollen gegen die echte Datenbank laufen und vor jedem
  Test die Tabellen zurücksetzen.
- Bitte nur meinen Verantwortungsbereich bearbeiten: `internal/handler/`,
  `internal/http/`, `internal/validation/` und `internal/middleware/`.
