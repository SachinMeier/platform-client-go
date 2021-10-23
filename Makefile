PHONY: test

.PHONY: test
test:
	@go test -v ./...

.PHONY: lint
lint:
	@golint -set_exit_status