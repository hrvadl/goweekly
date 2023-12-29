run:
	cd cmd/goweekly && rm -rf ./goweekly && go build . && ./goweekly

build:
	cd cmd/goweekly && go build .

lint:
	golangci-lint --config ./golangci.yaml run ./...

test:
	go test ./... -v
