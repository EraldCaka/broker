upgrade:
	@cd tools/upgrade && go run .

tidy:
	@go mod tidy

up:	# start kafka and zookeeper docker images
	@cd build && docker compose up

run: # start broker
	@cd cmd/broker && go run .
