package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	govalidator "github.com/go-playground/validator/v10"
	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"github.com/hasa1034/Workshop_Gruppe3/internal/repository"
	"github.com/hasa1034/Workshop_Gruppe3/internal/service"
	"github.com/hasa1034/Workshop_Gruppe3/internal/validation"
	"github.com/shopspring/decimal"
)

// KioskService definiert die Methoden, die der Handler vom Service benoetigt.
// Entspricht der tatsaechlichen Signatur von service.KioskService (Efe).
type KioskService interface {
	GetByID(ctx context.Context, id uint) (*model.Kiosk, error)
	List(ctx context.Context, filter repository.KioskFilter) ([]model.Kiosk, error)
	Create(ctx context.Context, kiosk *model.Kiosk) error
}

// KioskHandler haelt den Service und beantwortet HTTP-Anfragen.
type KioskHandler struct {
	svc KioskService
}

func NewKioskHandler(svc KioskService) *KioskHandler {
	return &KioskHandler{svc: svc}
}

// GET /kioske
// Optionale Query-Parameter: ?name=... und ?email=...
func (h *KioskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filter := repository.KioskFilter{
		Name:  r.URL.Query().Get("name"),
		Email: r.URL.Query().Get("email"),
	}

	kioske, err := h.svc.List(r.Context(), filter)
	if err != nil {
		slog.Error("GetAll fehlgeschlagen", "err", err)
		writeError(w, http.StatusInternalServerError, "interner Fehler")
		return
	}

	writeJSON(w, http.StatusOK, kioske)
}

// GET /kioske/{id}
func (h *KioskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ungueltige ID")
		return
	}

	kiosk, err := h.svc.GetByID(r.Context(), uint(id))
	if errors.Is(err, service.ErrNotFound) {
		writeError(w, http.StatusNotFound, "Kiosk nicht gefunden")
		return
	}
	if err != nil {
		slog.Error("GetByID fehlgeschlagen", "err", err)
		writeError(w, http.StatusInternalServerError, "interner Fehler")
		return
	}

	writeJSON(w, http.StatusOK, kiosk)
}

// POST /kioske
func (h *KioskHandler) Create(w http.ResponseWriter, r *http.Request) {
	// 1. JSON einlesen
	var req validation.KioskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungueltiges JSON: "+err.Error())
		return
	}

	// 2. Validierung
	if err := validation.ValidateKioskCreate(&req); err != nil {
		var ve govalidator.ValidationErrors
		if errors.As(err, &ve) {
			msgs := make([]string, len(ve))
			for i, e := range ve {
				msgs[i] = e.Field() + ": " + e.Tag()
			}
			writeJSON(w, http.StatusBadRequest, map[string]any{"errors": msgs})
			return
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 3. DTO → model.Kiosk mappen
	kiosk, err := mapToModel(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Service aufrufen
	if err := h.svc.Create(r.Context(), kiosk); err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			writeError(w, http.StatusConflict, "E-Mail bereits vergeben")
			return
		}
		slog.Error("Create fehlgeschlagen", "err", err)
		writeError(w, http.StatusInternalServerError, "interner Fehler")
		return
	}

	// 5. 201 + Location-Header
	w.Header().Set("Location", fmt.Sprintf("/kioske/%d", kiosk.ID))
	writeJSON(w, http.StatusCreated, kiosk)
}

// mapToModel wandelt das Create-DTO in ein model.Kiosk um.
func mapToModel(req *validation.KioskCreateRequest) (*model.Kiosk, error) {
	istGeoeffnet := true
	if req.IstGeoeffnet != nil {
		istGeoeffnet = *req.IstGeoeffnet
	}

	kiosk := &model.Kiosk{
		Name:         req.Name,
		Email:        req.Email,
		IstGeoeffnet: istGeoeffnet,
		Homepage:     req.Homepage,
		Username:     req.Username,
		Erzeugt:      time.Now(),
		Aktualisiert: time.Now(),
		Betreiber: model.Betreiber{
			Vorname:    req.Betreiber.Vorname,
			Nachname:   req.Betreiber.Nachname,
			Geschlecht: req.Betreiber.Geschlecht,
		},
	}

	for _, p := range req.Produkte {
		preis, err := decimal.NewFromString(p.Preis)
		if err != nil {
			return nil, fmt.Errorf("ungültiger Preis '%s': %w", p.Preis, err)
		}
		kiosk.Produkte = append(kiosk.Produkte, model.Produkt{
			Name:     p.Name,
			Preis:    preis,
			Waehrung: p.Waehrung,
		})
	}

	return kiosk, nil
}

// --- Hilfsfunktionen ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
