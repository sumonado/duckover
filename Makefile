GOCMD:=$(shell which go)

test:
	@$(GOCMD) test -v ./...

cover:
	@$(GOCMD) test -v ./... -coverprofile=coverage.txt -covermode=atomic

deps:
	@$(GOCMD) mod download

build:
	@$(GOCMD) build -o bin/duckover cmd/duckover.go

build-win:
	@GOOS=windows GOARCH=amd64 $(GOCMD) build -o bin/duckover.exe cmd/duckover.go

run:
	@$(GOCMD) run cmd/duckover.go