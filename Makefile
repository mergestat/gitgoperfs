GIT_GO_PERFS_TARGET ?= $(PWD)
TAGS = "static,system_libgit2"

bench:
	go test -tags=$(TAGS) . -bench .

vet:
	go vet -v -tags=$(TAGS) ./...

lint:
	golangci-lint run --build-tags $(TAGS)
	