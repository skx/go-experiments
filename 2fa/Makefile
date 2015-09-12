#
#  Utility to make the 2fa example.
#

default: 2fa

clean:
	go clean

libs:
	test -d $$GOPATH/github.com/dgryski/dgoogauth || go get github.com/dgryski/dgoogauth


2fa: libs main.go
	go build .
