.PHONY: test-deps
test-deps:
	cd linter/testdata/src/a && go mod vendor

.PHONY: test
test: test-deps
	cd linter && go test -v -covermode=atomic -coverprofile=coverage.txt -coverpkg .

.PHONY: build
build:
	golangci-lint custom -v