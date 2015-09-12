#
#  Utility to test the hasher.
#

default: hash


clean:
	go clean

libs:
	test -d $$GOPATH/github.com/speps/go-hashids || go get github.com/speps/go-hashids


hash: libs main.go
	go build .
