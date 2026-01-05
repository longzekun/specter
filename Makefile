BIN := specter-server

build:
	go build -o $(BIN) server/main.go
