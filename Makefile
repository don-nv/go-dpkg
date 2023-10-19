.PHONY: deps
deps:
	go mod tidy;
	go mod vendor;

.PHONY: test-check
test-check:
	go test -race -v -count 1 -timeout 5s ./...;

.PHONY: fmt
fmt:
	find . -type f -name '*.go' -not -path "./vendor/*" -exec gofmt -s -w {} \;