.PHONY: build install clean test

BINARY_NAME=converttomd-jira
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME)

install: build
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

test:
	go test -v ./...

run: build
	./$(BINARY_NAME)
