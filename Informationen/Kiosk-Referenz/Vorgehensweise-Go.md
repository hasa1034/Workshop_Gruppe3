# Vorgehensweise fuer Go

## Schritt 1: Projektstruktur anlegen

- Go-Modul initialisieren.
- Ordner fuer Startpunkt, Modelle, Repository, Service, HTTP-Handler,
  Validierung und Tests anlegen.
- Konfiguration fuer Server-Port und PostgreSQL-Verbindung vorbereiten.

## Schritt 2: Modelle erstellen

- `Kiosk`, `Betreiber`, `Produkt` und `Geschlecht` als Go-Typen definieren.
- GORM-Tags fuer Tabellen, Spalten, Indizes und Beziehungen setzen.
- Beziehung von `Kiosk` zu `Betreiber` als einzelnes Struct abbilden.
- Beziehung von `Kiosk` zu `Produkt` als Slice abbilden.
- Keine Go-Ruecknavigation von `Betreiber` oder `Produkt` zu `Kiosk`
  modellieren.

## Schritt 3: Datenbank anbinden

- PostgreSQL-Treiber fuer GORM konfigurieren.
- Verbindung beim Start der Anwendung oeffnen.
- `kiosk.betreiber_id` fuer die 1:1-Beziehung und `produkt.kiosk_id` fuer die
  1:N-Beziehung vorsehen.
- Tabellenstruktur entweder per vorhandener SQL-Datei oder per kontrollierter
  Migration bereitstellen.

## Schritt 4: REST-Router und Handler bauen

- Router mit `net/http` und optional `github.com/go-chi/chi/v5` aufbauen.
- Handler fuer `GET /kioske/{id}`, `GET /kioske` und `POST /kioske` erstellen.
- JSON sauber lesen und schreiben.
- `Location`-Header beim Neuanlegen setzen.

## Schritt 5: Validierung beim Neuanlegen

- Create-DTOs fuer Kiosk, Betreiber und Produkt definieren.
- Im Create-DTO ist `betreiber` ein einzelnes Objekt und `produkte` ein Array.
- `github.com/go-playground/validator/v10` fuer Pflichtfelder, E-Mail,
  Laengen und Waehrung verwenden.
- Validierung nur im POST-Flow ausfuehren.

## Schritt 6: Service- und Repository-Logik

- Repository kapselt alle GORM-Zugriffe.
- Service prueft fachliche Regeln wie eindeutige E-Mail.
- Service startet die Transaktion fuer das Neuanlegen.
- Handler wandeln Service-Fehler in HTTP-Statuscodes um.

## Schritt 7: Integrationstest

- Testserver mit `net/http/httptest` starten.
- Testdatenbank vorbereiten.
- Erfolgreiches `POST /kioske` pruefen.
- Danach `GET /kioske/{id}` fuer den neu angelegten Datensatz pruefen.
- Ungueltigen Create-Request pruefen und `400 Bad Request` erwarten.
