# Import .env file
include .env
export

test-sshfs:
	go test -count=1 -v ./pkg/sshfs

test-fusex:
	go test -count=1 -v ./pkg/fusex
