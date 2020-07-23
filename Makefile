check-formatter:
	which goimports || GO111MODULE=off go get -u golang.org/x/tools/cmd/goimports

format: check-formatter
	find . -type f -name "*.go" -not -path "./vendor/*" | xargs -n 1 -I R goimports -w R
	find . -type f -name "*.go" -not -path "./vendor/*" | xargs -n 1 -I R gofmt -s -w R

check-linter:
	which golangci-lint || (GO111MODULE=off go get -u -v github.com/golangci/golangci-lint/cmd/golangci-lint)

lint: check-linter
	golangci-lint run --deadline 10m ./...

up:
	docker-compose up -d

down:
	docker-compose down
