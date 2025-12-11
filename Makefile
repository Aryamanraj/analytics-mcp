SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

GO ?= go
GOFMT ?= gofmt
BINARY ?= mcp-server
BIN_DIR ?= bin
PKG := ./...
GOFILES := $(shell find . -type f -name '*.go' -not -path './vendor/*')

.PHONY: help
help:
	@echo "PayRam MCP server automation"
	@echo "Targets:"
	@echo "  make setup          Download module deps (go mod tidy)"
	@echo "  make fmt            Run go fmt on all packages"
	@echo "  make fmt-check      Check formatting (gofmt -l)"
	@echo "  make vet            Run go vet"
	@echo "  make test           Run go test ./..."
	@echo "  make cover          Run tests with coverage report"
	@echo "  make build          Build binary to $(BIN_DIR)/$(BINARY)"
	@echo "  make run            Run server (stdio)"
	@echo "  make run-http       Run server over HTTP (:8080 by default)"
	@echo "  make precommit      fmt-check + vet + test + build"
	@echo "  make commit         Guide an interactive conventional commit"
	@echo "  make clean          Remove build artifacts and coverage files"

.PHONY: setup
setup:
	$(GO) mod tidy

.PHONY: fmt
fmt:
	$(GO) fmt $(PKG)

.PHONY: fmt-check
fmt-check:
	@set -euo pipefail; \
	CHANGED=$$($(GOFMT) -l $(GOFILES)); \
	if [ -n "$$CHANGED" ]; then \
		echo "Files need formatting:"; echo "$$CHANGED"; exit 1; \
	fi

.PHONY: vet
vet:
	$(GO) vet $(PKG)

.PHONY: test
test:
	$(GO) test $(PKG)

.PHONY: cover
cover:
	$(GO) test -coverprofile=coverage.out $(PKG)
	$(GO) tool cover -func=coverage.out

.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(BINARY) .

.PHONY: run
run:
	$(GO) run .

.PHONY: run-http
run-http:
	$(GO) run . --http :8080

.PHONY: precommit
precommit:
	@set -euo pipefail; \
	$(MAKE) fmt-check; \
	$(MAKE) vet; \
	$(MAKE) test; \
	$(MAKE) build

.PHONY: commit
commit: precommit
	@set -euo pipefail; \
	TYPE=${TYPE-}; SCOPE=${SCOPE-}; MSG=${MSG-}; BODY=${BODY-}; BR=${BR-}; BREAKING=${BREAKING-}; FOOTER=${FOOTER-}; ADD=${ADD-}; \
	if ! git rev-parse --git-dir >/dev/null 2>&1; then echo "Not a git repo"; exit 1; fi; \
	TYPES="feat fix chore docs refactor perf test build ci revert"; \
	if [[ -z "$$ADD" ]]; then read -p "Stage all changes (git add -A)? [Y/n]: " ADD || true; fi; \
	if [[ -z "$$ADD" || "$$ADD" =~ ^[Yy] ]]; then git add -A; fi; \
	if git diff --cached --quiet; then echo "No staged changes; aborting commit"; exit 1; fi; \
	if [ -z "$$TYPE" ]; then \
		echo "Select commit type:"; i=1; for t in $$TYPES; do echo "  $$i) $$t"; i=$$((i+1)); done; \
		read -p "Choose number: " N || true; TYPE=$$(echo $$TYPES | awk -v n=$$N '{split($$0,a," "); print a[n]}'); \
	fi; \
	if [ -z "$$TYPE" ]; then echo "Commit type required"; exit 1; fi; \
	if [ -z "$$SCOPE" ]; then read -p "Optional scope (e.g., core/tools): " SCOPE || true; fi; \
	while [ -z "$$MSG" ]; do read -p "Short description (<=72 chars): " MSG || true; done; \
	read -p "Body (optional): " BODY || true; \
	read -p "Breaking change? [y/N]: " BR || true; \
	if [[ $${BR:-N} =~ ^(y|Y)$$ ]]; then read -p "Describe breaking change: " BREAKING || true; else BREAKING=""; fi; \
	read -p "Footer (e.g., Closes #123) (optional): " FOOTER || true; \
	HEADER="$$TYPE"; [ -n "$$SCOPE" ] && HEADER="$$HEADER($$SCOPE)"; [ -n "$$BREAKING" ] && HEADER="$$HEADER!"; HEADER="$$HEADER: $$MSG"; \
	MSGFILE=$$(mktemp); echo "$$HEADER" > $$MSGFILE; \
	[ -n "$$BODY" ] && { echo; echo "$$BODY"; } >> $$MSGFILE; \
	[ -n "$$BREAKING" ] && { echo; echo "BREAKING CHANGE: $$BREAKING"; } >> $$MSGFILE; \
	[ -n "$$FOOTER" ] && { echo; echo "$$FOOTER"; } >> $$MSGFILE; \
	if git commit -F $$MSGFILE; then echo "Commit created"; else echo "Commit failed"; fi; \
	rm -f $$MSGFILE

.PHONY: clean
clean:
	rm -rf $(BIN_DIR) coverage.out
