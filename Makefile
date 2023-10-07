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


#go build -gcflags="-m=3" ./dlog/v1/*.go |& grep -i escapes
#go test  -bench=Wrap  ./dlog/v1/  -test.benchmem -test.memprofile dlog.zap.wrapped.mem.pprof -test.outputdir ./pprof -test.count 5
#go test  -bench=Orig  ./dlog/v1/  -test.benchmem -test.memprofile dlog.zap.orig.mem.pprof -test.outputdir ./pprof -test.count 5

#go test  -bench=Wrap  ./dlog/v1/  -test.cpuprofile dlog.zap.wrapped.mem.pprof -test.outputdir ./pprof -test.count 5