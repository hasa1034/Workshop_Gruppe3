package model

import "time"

type Kiosk struct {
	ID           uint      `gorm:"primaryKey;column:id"`
	Version      int       `gorm:"not null;default:0;column:version"`
	Name         string    `gorm:"not null;index;column:name"`
	Email        string    `gorm:"not null;unique;column:email"`
	IstGeoeffnet bool      `gorm:"not null;default:true;column:ist_geoeffnet"`
	Homepage     *string   `gorm:"column:homepage"`
	Username     string    `gorm:"not null;column:username"`
	Erzeugt      time.Time `gorm:"not null;column:erzeugt"`
	Aktualisiert time.Time `gorm:"not null;column:aktualisiert"`
	BetreiberID  uint      `gorm:"not null;uniqueIndex;column:betreiber_id"`
	Betreiber    Betreiber `gorm:"foreignKey:BetreiberID;references:ID;constraint:OnUpdate:NO ACTION,OnDelete:RESTRICT;"`
	Produkte     []Produkt `gorm:"foreignKey:KioskID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE;"`
}

func (Kiosk) TableName() string {
	return "kiosk.kiosk"
}
