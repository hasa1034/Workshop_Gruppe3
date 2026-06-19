# CI fuer Go

Eine einfache CI-Pipeline fuer die Go-Implementierung sollte schnell pruefen,
ob der Code kompiliert, formatiert ist und die Tests erfolgreich laufen.

## Empfohlene Checks

```shell
go fmt ./...
go vet ./...
go test ./...
```

Optional kann ein zusaetzlicher Linter verwendet werden:

```shell
golangci-lint run
```

## Reihenfolge

1. Repository auschecken.
2. Go-Version installieren.
3. Dependencies laden.
4. Formatierung pruefen.
5. Statische Pruefung ausfuehren.
6. Tests ausfuehren.

## Minimaler Workflow

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
          go-version: "1.22"
      - run: go fmt ./...
      - run: go vet ./...
      - run: go test ./...
```

## Erwartung

- Formatierungsabweichungen werden sichtbar.
- Offensichtliche Codeprobleme werden erkannt.
- Unit- und Integrationstests laufen automatisiert.
