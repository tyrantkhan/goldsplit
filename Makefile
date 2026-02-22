.PHONY: lint frontend-check test dev build ci

lint: frontend-build
	golangci-lint run

frontend-build:
	cd frontend && npm install && npm run build

frontend-check:
	cd frontend && npm install && npm run check

test:
	go test ./internal/...

dev:
	wails dev

build:
	wails build

ci: lint frontend-check test
