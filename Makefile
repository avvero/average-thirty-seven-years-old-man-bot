BINARY=engine
test:
	richgo test ./...

test-full:
	richgo test -v -cover -covermode=atomic ./...

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

.PHONY: clean install unittest build docker run stop vendor lint-prepare lint