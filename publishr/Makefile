#
#  Utility to make the 2fa example.
#

default: server

clean:
	go clean

fmt:
	go fmt *.go

libs:
	go get ./...


server: libs main.go
	go build .
