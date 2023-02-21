
build:
	docker build -t gportal/metadata-server:latest .

lint:
	golangci-lint run

test:
	go test -v ./...
