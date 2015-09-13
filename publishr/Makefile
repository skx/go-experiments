#
#  Utility to make the 2fa example.
#

default: server

clean:
	go clean

fmt:
	go fmt *.go

libs:
	test -d $$GOPATH/github.com/dgryski/dgoogauth || go get github.com/dgryski/dgoogauth
	test -d $$GOPATH/github.com/gorilla/mux || go get github.com/gorilla/mux
	test -d $$GOPATH/github.com/gorilla/handlers || go get github.com/gorilla/handlers
	test -d $$GOPATH/github.com/rakyll/magicmime || go get github.com/rakyll/magicmime
	test -d $$GOPATH/github.com/speps/go-hashids || go get github.com/speps/go-hashids


server: libs main.go
	go build .
