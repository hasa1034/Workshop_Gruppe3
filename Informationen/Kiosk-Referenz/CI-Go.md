# CI für Go

Eine einfache CI-Pipeline für die Go-Implementierung prüft schnell,
ob der Code kompiliert, formatiert ist und die Tests erfolgreich laufen.

## Empfohlene Checks

```shell
go fmt ./...
go vet ./...
go test ./...
```

Optional kann ein zusätzlicher Linter verwendet werden:

```shell
golangci-lint run
```

## Reihenfolge

1. Repository auschecken.
2. Go-Version installieren.
3. Dependencies laden.
4. Formatierung prüfen.
5. Statische Prüfung ausführen.
6. Tests ausführen.

## Minimaler Workflow-Vorschlag

```yaml
name: Go CI

on:
  push:
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.25"
      - run: go fmt ./...
      - run: go vet ./...
      - run: go test ./...
```

## Erwartung

- Formatierungsabweichungen werden sichtbar.
- Offensichtliche Codeprobleme werden erkannt.
- Tests laufen automatisiert, sofern PostgreSQL in der CI bereitgestellt oder
  `DATABASE_URL` auf eine erreichbare Testdatenbank gesetzt wird.

## Aktueller Repository-Stand

Im Repository liegt aktuell kein GitHub-Actions-Workflow. Für die lokale Abgabe
werden die Checks direkt ausgeführt:

```shell
go fmt ./...
go test ./...
```
