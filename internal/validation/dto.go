package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
)

var validate = validator.New()

// BetreiberCreateRequest – Eingabe fuer einen neuen Betreiber.
type BetreiberCreateRequest struct {
	Vorname    string            `json:"vorname"              validate:"required,max=40"`
	Nachname   string            `json:"nachname"             validate:"required,max=40"`
	Geschlecht *model.Geschlecht `json:"geschlecht,omitempty" validate:"omitempty,oneof=MAENNLICH WEIBLICH DIVERS"`
}

// ProduktCreateRequest – Eingabe fuer ein neues Produkt.
type ProduktCreateRequest struct {
	Name     string `json:"name"     validate:"required,max=40"`
	Preis    string `json:"preis"    validate:"required"`
	Waehrung string `json:"waehrung" validate:"required,len=3,uppercase"`
}

// KioskCreateRequest – Eingabe fuer einen neuen Kiosk (POST /kioske).
type KioskCreateRequest struct {
	Name         string                 `json:"name"                   validate:"required,max=40"`
	Email        string                 `json:"email"                  validate:"required,email"`
	IstGeoeffnet *bool                  `json:"istGeoeffnet,omitempty"`
	Homepage     *string                `json:"homepage,omitempty"     validate:"omitempty,url"`
	Username     string                 `json:"username"               validate:"required,max=20"`
	Betreiber    BetreiberCreateRequest `json:"betreiber"              validate:"required"`
	Produkte     []ProduktCreateRequest `json:"produkte,omitempty"     validate:"dive"`
}

// ValidateKioskCreate prueft den Request und gibt Validierungsfehler zurueck.
func ValidateKioskCreate(req *KioskCreateRequest) error {
	return validate.Struct(req)
}
