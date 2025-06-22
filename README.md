# github-schema-go

Go module for querying GitHub GraphQL schema using gojq, with embedded schema support.

## Features

- Query GitHub GraphQL schema without API calls
- Embedded public schema from [github/docs](https://github.com/github/docs/tree/main/src/graphql/data/fpt)
- Support for custom schema files
- Pure jq queries using [gojq](https://github.com/itchyny/gojq)
- Zero GraphQL dependencies

## Installation

```bash
go get github.com/apstndb/github-schema-go
```

## Usage

```go
import "github.com/apstndb/github-schema-go/schema"

// Use embedded schema
s := schema.New()

// Use custom schema file
s := schema.NewWithFile("path/to/schema.json")

// Query type information
result, err := s.Type("PullRequest")

// Search types
results, err := s.Search("review.*thread")

// Get mutation details
result, err := s.Mutation("createIssue")
```

## CLI Tool

```bash
# Install CLI
go install github.com/apstndb/github-schema-go/cmd/github-schema@latest

# Query embedded schema
github-schema type PullRequest
github-schema mutation createIssue
github-schema search ".*Thread"

# Use custom schema
github-schema --schema ./my-schema.json type Issue

# Download latest schema
github-schema download                           # Download to stdout
github-schema download -o schema.json            # Download to file
github-schema download -o schema.json.gz         # Auto-compress based on .gz extension
github-schema download --compress                # Download compressed to stdout
github-schema download -c -o schema.json.gz      # Explicitly compress to file
```

## Updating Schema

```bash
# Update embedded schema (for development)
make update-schema

# Or using go generate
go generate ./schema

# Or using the CLI directly
github-schema download --compress -o schema/schema.json.gz
```

## Schema Source

The embedded schema is sourced from:
https://github.com/github/docs/tree/main/src/graphql/data/fpt

This is the official public schema published by GitHub.