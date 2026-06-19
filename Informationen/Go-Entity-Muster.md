# Go-Entity-Muster fuer Kiosk

Diese Notiz beschreibt das Entity-Muster fuer die Go-Implementierung mit
`Kiosk`, `Betreiber` und `Produkt`.

## Domänenmodell

- Ein `Kiosk` besitzt Stammdaten wie Name, E-Mail, Status, Homepage,
  Username sowie Zeitstempel.
- Ein `Kiosk` hat genau einen `Betreiber`.
- Ein `Kiosk` hat mehrere `Produkte`.
- Die fachliche Richtung geht vom `Kiosk` zum `Betreiber` und vom `Kiosk` zu
  den `Produkten`.
- `Betreiber` und `Produkt` haben kein `Kiosk`-Objekt als Rueckrichtung.
- `Produkt` nutzt nur `KioskID` als technische Datenbankzuordnung.
- Wird ein `Kiosk` geloescht, werden zugehoerige `Produkte` per Cascade
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

Validierung wird nur beim Neuanlegen benoetigt. Dafuer bietet sich ein eigenes
Request-DTO an, damit Datenbankmodell und REST-Eingabe getrennt bleiben.

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

Empfohlene Bibliothek: `github.com/go-playground/validator/v10`.

## REST-Schnittstelle

Empfohlene schlanke Variante: `net/http` mit `github.com/go-chi/chi/v5`.

- `GET /kioske/{id}` liest einen Kiosk mit genau einem `Betreiber` und
  mehreren `Produkten`.
- `GET /kioske` liest Kioske, optional mit Query-Parametern wie `name` oder
  `email`.
- `POST /kioske` legt einen neuen Kiosk inklusive genau einem Betreiber und
  optionalen Produkten an.
- Beim erfolgreichen `POST` wird `201 Created` mit `Location: /kioske/{id}`
  zurueckgegeben.

## Repository- und Service-Verhalten

- Lesen erfolgt ohne explizite Transaktion.
- Neuanlegen erfolgt in einer Transaktion, damit `Kiosk`, `Betreiber` und
  `Produkte` gemeinsam gespeichert werden.
- Vor dem Neuanlegen wird geprueft, ob die E-Mail bereits existiert.
- Fehler werden als passende HTTP-Statuscodes abgebildet:
  - `400 Bad Request` fuer ungueltiges JSON oder Validierungsfehler
  - `404 Not Found` fuer nicht vorhandene IDs
  - `409 Conflict` fuer bereits vorhandene E-Mail
  - `201 Created` fuer erfolgreiches Neuanlegen

## Einfacher Integrationstest

Ein minimaler Integrationstest sollte pruefen:

- `POST /kioske` mit gueltigem Body liefert `201 Created` und `Location`.
- `GET /kioske/{id}` liefert den neu angelegten Kiosk.
- `POST /kioske` mit ungueltiger E-Mail liefert `400 Bad Request`.

Geeignete Go-Werkzeuge:

- `testing`
- `net/http/httptest`
- Testdatenbank oder Testcontainer fuer PostgreSQL
