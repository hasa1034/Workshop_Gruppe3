package http

import (
	"net/http"

	"github.com/hasa1034/Workshop_Gruppe3/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter baut den chi-Router mit allen Kiosk-Endpunkten.
//
//	GET  /kioske       – Liste aller Kioske
//	GET  /kioske/{id}  – Einzelnen Kiosk lesen
//	POST /kioske       – Neuen Kiosk anlegen
func NewRouter(kh *handler.KioskHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/kioske", kh.GetAll)
	r.Get("/kioske/{id}", kh.GetByID)
	r.Post("/kioske", kh.Create)

	return r
}
