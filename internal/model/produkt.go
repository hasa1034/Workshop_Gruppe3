package model

import "github.com/shopspring/decimal"

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
