#
#  Utility to test the hasher.
#

default: hash


clean:
	go clean

libs:
	go get ./...

fmt:
	go fmt *.go

hash: libs main.go
	go build .
