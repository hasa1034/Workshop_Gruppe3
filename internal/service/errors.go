// Package service enthaelt die Fachlogik fuer Kioske: Dublettenpruefung,
// Transaktionen und die Abbildung fachlicher Fehler.
package service

import "errors"

// ErrNotFound signalisiert, dass ein angeforderter Kiosk nicht existiert.
// Der HTTP-Handler bildet dies auf 404 Not Found ab.
var ErrNotFound = errors.New("kiosk nicht gefunden")

// ErrEmailExists signalisiert, dass beim Neuanlegen bereits ein Kiosk mit
// derselben E-Mail existiert. Der HTTP-Handler bildet dies auf 409 Conflict ab.
var ErrEmailExists = errors.New("kiosk mit dieser e-mail existiert bereits")
