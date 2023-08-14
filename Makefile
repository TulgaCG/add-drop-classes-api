.PHONY: build
build: deps lint compile test install

.PHONY: deps
deps:
	@ go mod tidy --compat=1.20

.PHONY: compile
compile:
	@ go build -o bin/catgo cmd/add-drop-classes-api/main.go

.PHONY: install
install:
	@ go install ./cmd/add-drop-classes-api

.PHONY: test
test:
	@ go test ./...

.PHONY: lint
lint:
	@ golangci-lint run --config=.golangci.yaml --fix

.PHONY: generate
generate: generate-sqlc-code

generate-sqlc-code:	
	@ sqlc generate --file=sqlc.yaml