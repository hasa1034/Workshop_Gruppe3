package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"gorm.io/gorm"
)

// KioskFilter bietet optionale Filterkriterien fuer das Lesen mehrerer Kioske.
// Leere Felder werden ignoriert.
type KioskFilter struct {
	Name  string
	Email string
}

// Repository kapselt alle GORM-Zugriffe auf Kiosk-Daten. Der Service arbeitet
// ausschliesslich gegen dieses Interface.
type Repository interface {
	// FindByID liest einen Kiosk inklusive Betreiber und Produkten.
	FindByID(ctx context.Context, id uint) (*model.Kiosk, error)
	// FindAll liest mehrere Kioske, optional gefiltert.
	FindAll(ctx context.Context, filter KioskFilter) ([]model.Kiosk, error)
	// ExistsByEmail prueft, ob bereits ein Kiosk mit dieser E-Mail existiert.
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	// Create speichert einen Kiosk samt Betreiber und Produkten.
	Create(ctx context.Context, kiosk *model.Kiosk) error
	// Transaction fuehrt fn in einer Datenbanktransaktion aus. Innerhalb von fn
	// arbeitet das uebergebene Repository auf derselben Transaktion.
	Transaction(ctx context.Context, fn func(repo Repository) error) error
}

// ErrNotFound wird zurueckgegeben, wenn kein passender Datensatz existiert.
var ErrNotFound = errors.New("datensatz nicht gefunden")

type gormRepository struct {
	db *gorm.DB
}

// NewRepository erstellt ein GORM-basiertes Repository.
func NewRepository(db *gorm.DB) Repository {
	return &gormRepository{db: db}
}

func (r *gormRepository) FindByID(ctx context.Context, id uint) (*model.Kiosk, error) {
	var kiosk model.Kiosk
	err := r.db.WithContext(ctx).
		Preload("Betreiber").
		Preload("Produkte").
		First(&kiosk, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("kiosk lesen fehlgeschlagen: %w", err)
	}
	return &kiosk, nil
}

func (r *gormRepository) FindAll(ctx context.Context, filter KioskFilter) ([]model.Kiosk, error) {
	query := r.db.WithContext(ctx).
		Preload("Betreiber").
		Preload("Produkte")

	if filter.Name != "" {
		query = query.Where("name = ?", filter.Name)
	}
	if filter.Email != "" {
		query = query.Where("email = ?", filter.Email)
	}

	var kioske []model.Kiosk
	if err := query.Order("id").Find(&kioske).Error; err != nil {
		return nil, fmt.Errorf("kioske lesen fehlgeschlagen: %w", err)
	}
	return kioske, nil
}

func (r *gormRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Kiosk{}).
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("e-mail-pruefung fehlgeschlagen: %w", err)
	}
	return count > 0, nil
}

func (r *gormRepository) Create(ctx context.Context, kiosk *model.Kiosk) error {
	if err := r.db.WithContext(ctx).Create(kiosk).Error; err != nil {
		return fmt.Errorf("kiosk speichern fehlgeschlagen: %w", err)
	}
	return nil
}

func (r *gormRepository) Transaction(ctx context.Context, fn func(repo Repository) error) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(&gormRepository{db: tx})
	})
}
