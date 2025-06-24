# github-schema-go

A Go module and CLI tool for querying GitHub GraphQL schema offline using embedded introspection data.

## Features

- Query GitHub GraphQL schema without API calls
- Embedded schema from GitHub GraphQL API introspection
- Support for custom schema files
- Pure jq queries using [gojq](https://github.com/itchyny/gojq)
- Zero GraphQL client dependencies
- Native compression support using GitHub API gzip
- Consistent YAML/JSON formatting (via [go-yamlformat](https://github.com/apstndb/go-yamlformat))

## Installation

### As a Go module

```bash
go get github.com/apstndb/github-schema-go
```

### CLI tool

```bash
go install github.com/apstndb/github-schema-go/cmd/github-schema@latest
```

## Library Usage

```go
package main

import (
    "fmt"
    "github.com/apstndb/github-schema-go/schema"
)

func main() {
    // Use embedded schema (default)
    s, err := schema.New()
    if err != nil {
        panic(err)
    }

    // Or use a custom schema file
    // s, err := schema.NewWithFile("path/to/schema.json")

    // Query type information
    result, err := s.Type("PullRequest")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", result)

    // Search for types matching a pattern
    results, err := s.Search("review.*thread")
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d matching types\n", results["count"])

    // Get mutation input details
    mutation, err := s.Mutation("createIssue")
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", mutation)

    // Run custom jq queries
    custom, err := s.Query(`.data.__schema.queryType.name`, nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Query type: %v\n", custom)
}
```

## CLI Usage

### Basic Commands

```bash
# Show fields and description for a type
github-schema type PullRequest

# Show input requirements for a mutation
github-schema mutation createIssue

# Search for types matching a pattern
github-schema search ".*Thread"

# Run custom jq query
github-schema query '.data.__schema.types[] | select(.name == "Issue") | .fields[] | .name'

# Output as JSON instead of YAML
github-schema --json type Repository

# Use a custom schema file
github-schema --schema ./my-schema.json type Issue
```

### Downloading Schema

The `download` command fetches the latest schema via GitHub GraphQL API introspection:

```bash
# Download to stdout
github-schema download

# Download to file
github-schema download -o schema.json

# Download with compression (auto-detected by .gz extension)
github-schema download -o schema.json.gz

# Explicitly compress
github-schema download --compress -o my-schema.gz

# Note: Requires GitHub authentication
# Run 'gh auth login' if you haven't already
```

## Development

### Initial Setup

When cloning this repository, the embedded schema file must exist for the build to succeed:

```bash
# Clone the repository
git clone https://github.com/apstndb/github-schema-go.git
cd github-schema-go

# If schema/schema.json.gz is missing (e.g., after make clean)
touch schema/schema.json.gz

# Download the actual schema
make update-schema
```

### Updating the Embedded Schema

The embedded schema should be updated periodically to reflect GitHub's API changes:

```bash
# Update embedded schema (requires gh auth login)
make update-schema

# Or using go generate
go generate ./schema

# Or manually
github-schema download --compress -o schema/schema.json.gz
```

### Building and Testing

```bash
# Build the CLI
make build

# Run tests
make test

# Install locally
make install

# Clean generated files
make clean
```

## Schema Format

The embedded schema uses the standard GraphQL introspection format:

```json
{
  "data": {
    "__schema": {
      "queryType": { "name": "Query" },
      "mutationType": { "name": "Mutation" },
      "types": [
        {
          "kind": "OBJECT",
          "name": "Repository",
          "description": "A repository contains...",
          "fields": [...],
          ...
        }
      ]
    }
  }
}
```

This format is obtained directly from GitHub's GraphQL API using an introspection query, ensuring compatibility with standard GraphQL tooling.

## Performance

- Schema queries are performed using compiled jq expressions for optimal performance
- The embedded schema is compressed with gzip, reducing the binary size by ~92%
- All queries run offline without network calls
- Native GitHub API compression is used when downloading updates

## Requirements

- Go 1.16 or later (for go:embed support)
- `gh` CLI tool (for downloading schema updates)
- GitHub authentication via `gh auth login` (for schema updates only)

## License

This project is licensed under the MIT License. See the LICENSE file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [gojq](https://github.com/itchyny/gojq) - Pure Go implementation of jq
- GitHub for providing the GraphQL API and schema