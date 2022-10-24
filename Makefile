GIT_GO_PERFS_TARGET ?= $(PWD)

bench:
	
	go test . -bench .

vet:
    go vet -v ./...
	
lint:
    golangci-lint run