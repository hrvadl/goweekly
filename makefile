run:
	cd cmd/goweekly && go run .

lint:
	golangci-lint --config ./golangci.yaml run ./...

test:
	go test ./...
