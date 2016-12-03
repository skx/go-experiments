#
# Compile our driver, and each plugin-subdirectory.
#


#
# Compile the driver.
#
main: libs
	go build .

#
# Compile each plugin; treating each subdirectory as a new one.
#
libs:
	for i in */; do \
		cd $$i ; go build -buildmode=plugin ; cd .. ;\
	done

#
# Format all our code.
#
fmt:
	find . -name '*.go' -print | xargs -n 1 go fmt


clean:
	rm -f */*.so  || true
	rm -f plugins || true
