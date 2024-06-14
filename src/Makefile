# Variables
BINARY_NAME=myapp
VERSION=1.0.0
BUILD=`git rev-parse HEAD`

# OS-specific settings
ifeq ($(OS),Windows_NT)
	GOOS=windows
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		GOOS=linux
	endif
	ifeq ($(UNAME_S),Darwin)
		GOOS=darwin
	endif
endif

# Build for current OS
build:
	go build -o ./bin/$(BINARY_NAME)-$(VERSION)-$(GOOS) ./cmd

# Build for Linux
build-linux:
	GOOS=linux go build -o ./bin/$(BINARY_NAME)-$(VERSION)-linux ./cmd

# Build for Windows
build-windows:
	GOOS=windows go build -o ./bin/$(BINARY_NAME)-$(VERSION)-windows.exe ./cmd

# Build for MacOS
build-darwin:
	GOOS=darwin go build -o ./bin/$(BINARY_NAME)-$(VERSION)-darwin ./cmd

# Run
run:
	go run ./cmd

# Clean binaries
clean:
	rm -rf bin/*
