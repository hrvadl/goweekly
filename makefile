run:
	cd cmd/goweekly && rm -rf ./goweekly && go build . && ./goweekly

install:
	go mod tidy

build:
	cd cmd/goweekly && go build .

lint:
	golangci-lint --config ./golangci.yaml run ./...

test:
	go test ./... -v
