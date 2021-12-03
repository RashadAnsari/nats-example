all: format lint

format:
	@find . -type f -name "*.go" -not -path "./vendor/*" | xargs -n 1 -I R goimports -w R
	@find . -type f -name "*.go" -not -path "./vendor/*" | xargs -n 1 -I R gofmt -s -w R

lint:
	@golangci-lint run --deadline 10m ./...

up:
	@docker-compose up -d

down:
	@docker-compose down
