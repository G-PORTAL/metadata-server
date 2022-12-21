
build:
	docker build -t gportal/metadata-server:latest .

lint:
	golangci-lint run

proto:
	buf generate buf.build/gportal/gportal-cloud

test:
	go test -v ./...
