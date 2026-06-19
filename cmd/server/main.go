// Command server startet den Kiosk-HTTP-Server. Hier werden Konfiguration,
// Datenbank, Migration, Repository, Service, Handler und Router verdrahtet.
package main

import (
	"context"
	"errors"
	"log/slog"
	stdhttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hasa1034/Workshop_Gruppe3/internal/config"
	"github.com/hasa1034/Workshop_Gruppe3/internal/handler"
	kioskhttp "github.com/hasa1034/Workshop_Gruppe3/internal/http"
	"github.com/hasa1034/Workshop_Gruppe3/internal/repository"
	"github.com/hasa1034/Workshop_Gruppe3/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server beendet mit fehler", "err", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Konfiguration aus Umgebungsvariablen laden.
	cfg := config.Load()

	// 2. Datenbankverbindung oeffnen.
	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}

	// 3. Schema, Enum und Tabellen bereitstellen.
	if err := repository.Migrate(db); err != nil {
		return err
	}

	// 4. Repository, Service und Handler verdrahten.
	repo := repository.NewRepository(db)
	svc := service.NewKioskService(repo)
	kioskHandler := handler.NewKioskHandler(svc)

	// 5. Router bauen.
	router := kioskhttp.NewRouter(kioskHandler)

	srv := &stdhttp.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Server starten und auf Abbruchsignale reagieren (sauberes Beenden).
	errCh := make(chan error, 1)
	go func() {
		slog.Info("server gestartet", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-stop:
		slog.Info("abschaltsignal empfangen, server wird beendet")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(ctx)
}
