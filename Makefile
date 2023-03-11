.PHONY: deps
deps:
	go mod tidy && go mod vendor;

.PHONY: test
test:
	export DEBUG_ENABLED=false && go test -race -count 1 -timeout 10s ./...;

.PHONY: lint
lint: fmt
	docker run --rm -v `pwd`:/app -w /app golangci/golangci-lint:v1.58 golangci-lint run -v;

.PHONY: fmt
fmt:
	find . -type f -name '*.go' -not -path "./vendor/*" -exec gofmt -s -w {} \;


docker-prom:
	docker compose -f tests/prom/docker-compose.yaml up -d

push-wip:
	git add -A; git commit -m WIP; git push;

utilities:
	go get github.com/mailru/easyjson && go install github.com/mailru/easyjson/...@latest