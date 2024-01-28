build_server:
	GOOS=linux GOARCH=arm go build -o bin/server-linux-arm cmd/server/main.go

build_client:
	go build -o bin/client cmd/client/main.go

clean:
	rm -rf bin 
