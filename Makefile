
build:
	go build cmd/server.go

lint:
	golangci-lint run

proto:
	buf generate buf.build/gportal/gportal-cloud

test:
	go test ./...
