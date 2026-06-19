package service

import (
	"context"
	"errors"
	"time"

	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"github.com/hasa1034/Workshop_Gruppe3/internal/repository"
)

// KioskService buendelt die Fachlogik rund um Kioske.
type KioskService struct {
	repo repository.Repository
	now  func() time.Time
}

// NewKioskService erstellt einen Service auf Basis des uebergebenen Repositories.
func NewKioskService(repo repository.Repository) *KioskService {
	return &KioskService{
		repo: repo,
		now:  time.Now,
	}
}

// GetByID liest einen einzelnen Kiosk. Existiert er nicht, wird ErrNotFound
// zurueckgegeben. Lesen benoetigt keine Transaktion.
func (s *KioskService) GetByID(ctx context.Context, id uint) (*model.Kiosk, error) {
	kiosk, err := s.repo.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return kiosk, nil
}

// List liest mehrere Kioske, optional gefiltert nach Name oder E-Mail.
func (s *KioskService) List(ctx context.Context, filter repository.KioskFilter) ([]model.Kiosk, error) {
	return s.repo.FindAll(ctx, filter)
}

// Create legt einen Kiosk inklusive genau einem Betreiber und mehreren Produkten
// an. Dublettenpruefung und Speichern laufen gemeinsam in einer Transaktion, damit
// Kiosk, Betreiber und Produkte nur zusammen gespeichert werden. Existiert die
// E-Mail bereits, wird ErrEmailExists zurueckgegeben.
func (s *KioskService) Create(ctx context.Context, kiosk *model.Kiosk) error {
	now := s.now()
	kiosk.Erzeugt = now
	kiosk.Aktualisiert = now

	return s.repo.Transaction(ctx, func(repo repository.Repository) error {
		exists, err := repo.ExistsByEmail(ctx, kiosk.Email)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailExists
		}
		return repo.Create(ctx, kiosk)
	})
}
