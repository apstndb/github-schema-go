.PHONY: update-schema test build install

# Update embedded schema using the CLI tool
update-schema:
	@echo "Updating embedded schema..."
	go run ./cmd/github-schema download --compress -o schema/schema.json.gz
	@echo "Schema updated successfully"

# Run tests
test:
	go test -short ./cmd/... ./schema/... ./examples/...

# Build CLI
build:
	go build -o bin/github-schema ./cmd/github-schema

# Install CLI
install:
	go install ./cmd/github-schema

# Clean generated files
clean:
	rm -f schema/schema.json
	rm -f bin/github-schema

# Check if schema needs update
check-schema:
	@echo "Current embedded schema:"
	@ls -lh schema/schema.json.gz
	@echo "\nTo update: make update-schema"