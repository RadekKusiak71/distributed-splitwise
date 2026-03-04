run-dev:
	@air -c .air.toml

build:
	@go build -o ./tmp/gateway ./cmd/gateway

run-prod: build
	@./tmp/gateway