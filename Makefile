.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	mkdir -p .build
	CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o .build/webhook-conversion-service .

.PHONY: coverage
coverage:
	go test -v ./... -covermode=count -coverprofile=coverage.out
