package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/hasa1034/Workshop_Gruppe3/internal/model"
	"github.com/hasa1034/Workshop_Gruppe3/internal/validation"
	"github.com/go-chi/chi/v5"
	govalidator "github.com/go-playground/validator/v10"
)

// KioskService definiert die Methoden, die der Handler vom Service benoetigt.
// Person 3 (Efe) implementiert dieses Interface im Service-Package.
type KioskService interface {
	GetByID(id uint) (*model.Kiosk, error)
	GetAll(name, email string) ([]model.Kiosk, error)
	Create(req *validation.KioskCreateRequest) (*model.Kiosk, error)
}

// Sentinel-Fehler – werden vom Service zurueckgegeben.
var (
	ErrNotFound      = errors.New("kiosk nicht gefunden")
	ErrEmailConflict = errors.New("e-mail bereits vergeben")
)

// KioskHandler haelt den Service und beantwortet HTTP-Anfragen.
type KioskHandler struct {
	svc KioskService
}

func NewKioskHandler(svc KioskService) *KioskHandler {
	return &KioskHandler{svc: svc}
}

// GET /kioske
func (h *KioskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")

	kioske, err := h.svc.GetAll(name, email)
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

	kiosk, err := h.svc.GetByID(uint(id))
	if errors.Is(err, ErrNotFound) {
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
	var req validation.KioskCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "ungueltiges JSON: "+err.Error())
		return
	}

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

	kiosk, err := h.svc.Create(&req)
	if errors.Is(err, ErrEmailConflict) {
		writeError(w, http.StatusConflict, "E-Mail bereits vergeben")
		return
	}
	if err != nil {
		slog.Error("Create fehlgeschlagen", "err", err)
		writeError(w, http.StatusInternalServerError, "interner Fehler")
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/kioske/%d", kiosk.ID))
	writeJSON(w, http.StatusCreated, kiosk)
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
