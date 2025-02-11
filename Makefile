# Variables
BINARY_NAME = url-shortener
PROTO_SRC = api/proto/url_shortener.proto
PROTO_GEN_FLAGS = paths=source_relative:.

.PHONY: all build proto test tidy clean docker-build docker-run run

all: tidy proto build


# Build the Go binary
build:
	go build -o $(BINARY_NAME) ./cmd/url-shortener

# Run tests in the project
test:
	go test ./...

# Tidy up module dependencies
tidy:
	go mod tidy

# Remove binary and generated proto files if needed
clean:
	rm -f $(BINARY_NAME)
	find . -name "*.pb.go" -delete

# Build Docker image
docker-build:
	docker build -t $(BINARY_NAME) .

# Run Docker container
docker-run:
	docker run -p 8080:8080 -p 50051:50051 $(BINARY_NAME)

# Run the binary locally
run: build
	./$(BINARY_NAME)