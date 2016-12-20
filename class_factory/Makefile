

#
# Compile the driver.
#
main: $(wildcard *.go)
	go build .

#
# Format all our code.
#
fmt:
	find . -name '*.go' -print | xargs -n 1 go fmt


clean:
	rm -f class_factory  || true
