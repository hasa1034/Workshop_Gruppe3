# Programmierworkshop am 19.6.2026

## Namen
Sam Haghighi, Ali Arslan, Efe Yueksel

## Link zum Git-Repository
https://github.com/hasa1034/Workshop_Gruppe3

## KI-Werkzeuge
OpenAI Codex (Sam Haghighi)

### Agenten
OpenAI Codex (Sam Haghighi)

### Chat-URLs, z.B. https://chatgpt.com
Keine separate Chat-URL; die Arbeit erfolgte lokal mit OpenAI Codex im Repository.

## Frameworks und Bibliotheken

### REST-Schnittstelle (Lesen und Neuanlegen)
Go mit `net/http` und optional `github.com/go-chi/chi/v5` als schlankem Router.
Geplant sind Endpunkte zum Lesen von Kiosken und zum Neuanlegen eines Kiosks
inklusive genau einem Betreiber und mehreren Produkten.

### Validierung (nur Neuanlegen)
Validierung des POST-Request-DTOs mit `github.com/go-playground/validator/v10`.
Validiert werden u.a. Name, E-Mail, Username, Homepage, Betreiberdaten,
Produktname, Preis und Waehrung.

### OR-Mapping (für PostgreSQL)
GORM (`gorm.io/gorm` und `gorm.io/driver/postgres`) fuer PostgreSQL.
Das Entity-Muster aus dem Kiosk-Referenzprojekt wird auf Go-Structs fuer `Kiosk`,
`Betreiber` und `Produkt` uebertragen.

### Optional: OIDC mit Keycloak
Optional. Fuer die prototypische Implementierung steht zuerst REST, Validierung,
PostgreSQL und Integrationstest im Vordergrund.

### Einfacher Integrationstest
Integrationstest mit Go `testing`, `net/http/httptest` und einer PostgreSQL-
Testdatenbank bzw. einem Testcontainer. Geprueft werden mindestens POST zum
Neuanlegen und GET zum Lesen.

## Aufgabenverteilung
Person 1: Entity-Modell, Go-Structs, GORM-Mapping und Dokumentation.
Person 2: REST-Router, Handler, Request-/Response-DTOs und Validierung beim
Neuanlegen.
Person 3: PostgreSQL-Anbindung, Repository-/Service-Logik,
Integrationstests und CI-Checks.

## Prompts/Requests an KI-Agent/en
- Relevante Kiosk-Referenzdokumente auswerten und als Go-Dokumentation in
  `Informationen/` zusammenfassen.
- Das Entity-Muster `Kiosk`, `Betreiber`, `Produkt` fuer Go, REST,
  Validierung und GORM/PostgreSQL anpassen.
- Beziehung korrigieren: Kiosk zu Betreiber ist 1:1, Kiosk zu Produkt ist 1:N,
  ohne fachliche Rueckrichtung.
