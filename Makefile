# Define variables
GO := go
BINARY_NAME := gocms

# Phony targets are not actual files
.PHONY: all build run test clean

# Default target
all: build

# Build target
build:
	$(GO) build -o ./cmd/http/$(BINARY_NAME) ./cmd/http/main.go

# Run target
run: build 
	./cmd/http/$(BINARY_NAME)

# Test target
test:
	$(GO) test -v ./...

# Clean target
clean:
	$(GO) clean
	rm -f ./cmd/http/$(BINARY_NAME)