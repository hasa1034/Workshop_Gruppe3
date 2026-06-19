// Package app verdrahtet Konfiguration, Datenbank, Service, Handler und Router.
package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
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

// Run startet den Kiosk-HTTP-Server.
func Run() error {
	cfg := config.Load()

	db, err := repository.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}

	if err := repository.Migrate(db); err != nil {
		return err
	}

	repo := repository.NewRepository(db)
	svc := service.NewKioskService(repo)
	kioskHandler := handler.NewKioskHandler(svc)
	router := kioskhttp.NewRouter(kioskHandler)

	srv := &stdhttp.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server gestartet", "addr", listener.Addr().String())
		if err := srv.Serve(listener); err != nil && !errors.Is(err, stdhttp.ErrServerClosed) {
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
