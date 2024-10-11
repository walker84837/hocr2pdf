BINARY_NAME := hocr2pdf
SRC_DIR := ./src
#
# Default target to build the project
all: build

# Target to build the Go project
build:
	go build -ldflags "-linkmode external -s -w" -trimpath -buildmode=pie -mod=readonly -modcacherw -o $(BINARY_NAME) src/main.go src/hocr_converter.go

# Target to run the Go project
run: build
	./$(BINARY_NAME)

# Target to clean the build artifacts
clean:
	rm -f $(BINARY_NAME)

# Target to format the Go code
fmt:
	go fmt $(SRC_DIR)

# Target to check for dependencies
deps:
	go mod tidy

# Target to display help
help:
	@echo "Usage:"
	@echo "  make       - Build the project"
	@echo "  make run   - Build and run the project"
	@echo "  make clean - Remove build artifacts"
	@echo "  make fmt   - Format the Go code"
	@echo "  make deps  - Tidy up Go module dependencies"

.PHONY: all build run clean fmt deps help
