# Yuki CLI Foundation Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the first read-only Yuki CLI slice: auth/session handling, domain/admin discovery, and GL account listing.

**Architecture:** Use a small hand-shaped SOAP client in `internal/api`, keyring-backed secret storage in `internal/auth`, non-secret profile config in `internal/config`, Kong command wiring in `internal/cmd`, and table/JSON output helpers in `internal/output`.

**Tech Stack:** Go 1.26, Kong, 99designs/keyring, yaml.v3, httptest fixtures.

---

## File Structure

- Create `go.mod` and `go.sum`: module and dependency lock.
- Create `cmd/yuki/main.go`: binary entrypoint.
- Create `internal/api/soap.go`: SOAP envelope creation, XML escaping, SOAPAction, response helpers.
- Create `internal/api/client.go`: HTTP client, service endpoints, auth/session operations, typed API methods.
- Create `internal/api/types.go`: exported domain/admin/company/GL account structs.
- Create `internal/api/*_test.go`: failing-first tests for envelope generation, parsing, and httptest transport.
- Create `internal/auth/keyring.go`: profile-scoped access key storage and env fallback.
- Create `internal/config/config.go`: config paths, profile defaults, base URL and default administration.
- Create `internal/output/output.go`: JSON and table helpers.
- Create `internal/cmd/*.go`: Kong root, auth, domains, administrations, accounting, version commands.
- Create `Makefile`: build, fmt, fmt-check, lint, test, ci.
- Create `.golangci.yml`: minimal v2 lint config.

## Tasks

### Task 1: API SOAP Core

- [ ] Write failing tests in `internal/api/soap_test.go` for envelope escaping and operation parameter order.
- [ ] Run `go test ./internal/api` and verify missing symbols fail.
- [ ] Implement `Envelope(operation string, params []Param) string` and `SOAPAction(operation string) string`.
- [ ] Run `go test ./internal/api` and verify it passes.

### Task 2: API Response Parsing And Transport

- [ ] Write failing tests in `internal/api/client_test.go` using httptest for `Authenticate`, `Domains`, `Administrations`, `GetCurrentDomain`, and `GetGLAccountScheme`.
- [ ] Run `go test ./internal/api` and verify missing methods fail.
- [ ] Implement `Client`, service endpoint resolution, HTTP POST, SOAP fault parsing, and XML response mappers.
- [ ] Run `go test ./internal/api`.

### Task 3: Auth And Config

- [ ] Write failing tests for access key source precedence: explicit/env/keyring and non-secret config defaults.
- [ ] Run targeted tests and verify failures.
- [ ] Implement keyring store and config read/write using `YUKI_CONFIG_DIR` for tests.
- [ ] Run targeted tests.

### Task 4: Command Wiring

- [ ] Write failing command tests for `version`, no-args help, `auth status --json`, and command parsing for first-slice commands.
- [ ] Run `go test ./internal/cmd` and verify failures.
- [ ] Implement Kong root, global flags, auth commands, domains commands, administrations command, and accounting GL account command.
- [ ] Run `go test ./internal/cmd`.

### Task 5: Project Gates And Docs

- [ ] Add Makefile and lint config.
- [ ] Update README with install/build/auth examples and first-slice commands.
- [ ] Run `make build`.
- [ ] Run `go test ./...`.
- [ ] Run `./bin/yuki --help`.
- [ ] Run `./bin/yuki version`.
- [ ] Run `./bin/yuki auth status --json`.

## Constraints

- Do not commit unless Peter asks.
- Do not run live Yuki calls unless `YUKI_ACCESS_KEY` is intentionally present and the command is read-only.
- Do not add mutating commands until their docs parity rows and dry-run behavior are designed.
