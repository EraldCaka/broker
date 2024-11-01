upgrade:
	@cd tools/upgrade && go run .

tidy:
	@go mod tidy

up:
	@cd build && docker compose up

run:
	@cd cmd/broker && go run .