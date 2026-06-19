# Go-Entity-Muster für Kiosk

Diese Notiz beschreibt das umgesetzte Entity-Muster der Go-Implementierung mit
`Kiosk`, `Betreiber` und `Produkt`.

Die Codebasis liegt unter `internal/model`.

## Domänenmodell

- Ein `Kiosk` besitzt Stammdaten wie Name, E-Mail, Status, Homepage,
  Username sowie Zeitstempel.
- Ein `Kiosk` hat genau einen `Betreiber`.
- Ein `Kiosk` hat mehrere `Produkte`.
- Die fachliche Richtung geht vom `Kiosk` zum `Betreiber` und vom `Kiosk` zu
  den `Produkten`.
- `Betreiber` und `Produkt` haben kein `Kiosk`-Objekt als Rückrichtung.
- `Produkt` nutzt nur `KioskID` als technische Datenbankzuordnung.
- Wird ein `Kiosk` gelöscht, werden zugehörige `Produkte` per Cascade
  ebenfalls entfernt.

## Go-Structs mit GORM

```go
package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Geschlecht string

const (
	GeschlechtMaennlich Geschlecht = "MAENNLICH"
	GeschlechtWeiblich  Geschlecht = "WEIBLICH"
	GeschlechtDivers    Geschlecht = "DIVERS"
)

type Kiosk struct {
	ID           uint        `gorm:"primaryKey;column:id"`
	Version      int         `gorm:"not null;default:0;column:version"`
	Name         string      `gorm:"not null;index;column:name"`
	Email        string      `gorm:"not null;unique;column:email"`
	IstGeoeffnet bool        `gorm:"not null;default:true;column:ist_geoeffnet"`
	Homepage     *string     `gorm:"column:homepage"`
	Username     string      `gorm:"not null;column:username"`
	Erzeugt      time.Time   `gorm:"not null;column:erzeugt"`
	Aktualisiert time.Time   `gorm:"not null;column:aktualisiert"`
	BetreiberID  uint        `gorm:"not null;uniqueIndex;column:betreiber_id"`
	Betreiber    Betreiber   `gorm:"foreignKey:BetreiberID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:RESTRICT;"`
	Produkte     []Produkt   `gorm:"foreignKey:KioskID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
}

func (Kiosk) TableName() string {
	return "kiosk.kiosk"
}

type Betreiber struct {
	ID         uint        `gorm:"primaryKey;column:id"`
	Vorname    string      `gorm:"not null;index;column:vorname"`
	Nachname   string      `gorm:"not null;index;column:nachname"`
	Geschlecht *Geschlecht `gorm:"type:kiosk.geschlecht;column:geschlecht"`
}

func (Betreiber) TableName() string {
	return "kiosk.betreiber"
}

type Produkt struct {
	ID       uint            `gorm:"primaryKey;column:id"`
	Name     string          `gorm:"not null;column:name"`
	Preis    decimal.Decimal `gorm:"type:numeric(10,2);not null;column:preis"`
	Waehrung string          `gorm:"not null;size:3;column:waehrung"`
	KioskID  uint            `gorm:"not null;index;column:kiosk_id"`
}

func (Produkt) TableName() string {
	return "kiosk.produkt"
}
```

## DTO und Validierung beim Neuanlegen

Validierung wird nur beim Neuanlegen ausgeführt. Dafür nutzt die Anwendung ein
eigenes Request-DTO, damit Datenbankmodell und REST-Eingabe getrennt bleiben.

```go
type BetreiberCreateRequest struct {
	Vorname    string      `json:"vorname" validate:"required,max=40"`
	Nachname   string      `json:"nachname" validate:"required,max=40"`
	Geschlecht *Geschlecht `json:"geschlecht,omitempty" validate:"omitempty,oneof=MAENNLICH WEIBLICH DIVERS"`
}

type ProduktCreateRequest struct {
	Name     string `json:"name" validate:"required,max=40"`
	Preis    string `json:"preis" validate:"required"`
	Waehrung string `json:"waehrung" validate:"required,len=3,uppercase"`
}

type KioskCreateRequest struct {
	Name         string                   `json:"name" validate:"required,max=40"`
	Email        string                   `json:"email" validate:"required,email"`
	IstGeoeffnet *bool                    `json:"istGeoeffnet,omitempty"`
	Homepage     *string                  `json:"homepage,omitempty" validate:"omitempty,url"`
	Username     string                   `json:"username" validate:"required,max=20"`
	Betreiber    BetreiberCreateRequest   `json:"betreiber" validate:"required"`
	Produkte     []ProduktCreateRequest   `json:"produkte,omitempty" validate:"dive"`
}
```

Umgesetzte Bibliothek: `github.com/go-playground/validator/v10`.

## REST-Schnittstelle

Umgesetzt ist `net/http` mit `github.com/go-chi/chi/v5`.

- `GET /kioske/{id}` liest einen Kiosk mit genau einem `Betreiber` und
  mehreren `Produkten`.
- `GET /kioske` liest Kioske, optional mit Query-Parametern wie `name` oder
  `email`.
- `POST /kioske` legt einen neuen Kiosk inklusive genau einem Betreiber und
  optionalen Produkten an.
- Beim erfolgreichen `POST` wird `201 Created` mit `Location: /kioske/{id}`
  zurückgegeben.

## Repository- und Service-Verhalten

- Lesen erfolgt ohne explizite Transaktion.
- Neuanlegen erfolgt in einer Transaktion, damit `Kiosk`, `Betreiber` und
  `Produkte` gemeinsam gespeichert werden.
- Vor dem Neuanlegen wird geprüft, ob die E-Mail bereits existiert.
- Fehler werden als passende HTTP-Statuscodes abgebildet:
  - `400 Bad Request` für ungültiges JSON oder Validierungsfehler
  - `404 Not Found` für nicht vorhandene IDs
  - `409 Conflict` für bereits vorhandene E-Mail
  - `201 Created` für erfolgreiches Neuanlegen

## Integrationstest

Der vorhandene Integrationstest `kiosk_service_integration_test.go` prüft
Service und Repository gegen eine echte PostgreSQL-Datenbank:

- Kiosk inklusive Betreiber und Produkten anlegen.
- Neu angelegten Kiosk per ID lesen.
- Doppelte E-Mail-Adressen als `service.ErrEmailExists` erkennen.
- Nicht vorhandene IDs als `service.ErrNotFound` erkennen.
- Liste aller Kioske und Filter nach E-Mail prüfen.

Geeignete Go-Werkzeuge:

- `testing`
- PostgreSQL-Testdatenbank, z.B. über `extras/compose/postgres/compose.yml`
