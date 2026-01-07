SERVER_BIN := specter-server
CLIENT_BIN := specter-client
build:
	go build -o $(SERVER_BIN) server/main.go
	go build -o $(CLIENT_BIN) client/main.go
