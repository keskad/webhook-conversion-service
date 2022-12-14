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

.PHONY: image
image:
	docker build . -t webhook-conversion-service:latest

.PHONY: test_run
test_run:
	./.build/webhook-conversion-service -c ./example-config.yaml --debug
