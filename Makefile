.PHONY: deps
deps:
	go mod tidy;
	go mod vendor;

.PHONY: test-check
test-check:
	go test -race -v -count 1 -timeout 5s ./...