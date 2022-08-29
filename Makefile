GIT_GO_PERFS_TARGET ?= $(PWD)

bench:
	
	go test . -bench .
