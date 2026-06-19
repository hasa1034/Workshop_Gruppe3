package model

type Betreiber struct {
	ID         uint        `gorm:"primaryKey;column:id"`
	Vorname    string      `gorm:"not null;index;column:vorname"`
	Nachname   string      `gorm:"not null;index;column:nachname"`
	Geschlecht *Geschlecht `gorm:"type:kiosk.geschlecht;column:geschlecht"`
}

func (Betreiber) TableName() string {
	return "kiosk.betreiber"
}
