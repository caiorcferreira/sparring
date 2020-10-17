
dev:
	SPARRING_CONFIG=samples/config.yml go run *.go

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/sparring -ldflags="-s -w -extldflags -static" *.go
