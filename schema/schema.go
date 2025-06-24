package schema

import (
	"bytes"
	"compress/gzip"
	"context"
	_ "embed"
	"fmt"
	"io"
	"log/slog"
	"os"

	jqyaml "github.com/apstndb/go-jq-yamlformat"
	"github.com/apstndb/go-yamlformat"
)

// Embed the GitHub GraphQL schema in standard introspection format
// This file is obtained via GitHub GraphQL API introspection query
//
//go:embed schema.json.gz
var embeddedSchema []byte

// Schema provides methods to query GitHub GraphQL schema
type Schema struct {
	data interface{} // Parsed JSON schema
}

// New creates a Schema instance using the embedded schema
func New() (*Schema, error) {
	slog.Debug("Creating schema from embedded data", "size", len(embeddedSchema))
	
	reader, err := gzip.NewReader(bytes.NewReader(embeddedSchema))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress schema: %w", err)
	}
	
	slog.Debug("Decompressed schema", "size", len(data))

	return NewWithData(data)
}

// NewWithFile creates a Schema instance from a file
func NewWithFile(path string) (*Schema, error) {
	slog.Debug("Loading schema from file", "path", path)
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	
	slog.Debug("Loaded schema file", "size", len(data))

	return NewWithData(data)
}

// NewWithData creates a Schema instance from raw JSON data
func NewWithData(data []byte) (*Schema, error) {
	var schema interface{}
	// Use consistent unmarshaling with proper number handling
	if err := yamlformat.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return &Schema{data: schema}, nil
}

// Type queries information about a GraphQL type
func (s *Schema) Type(typeName string) (map[string]interface{}, error) {
	query := typeQuery
	return s.runQuery(query, map[string]interface{}{"type": typeName})
}

// Search searches for types matching a pattern
func (s *Schema) Search(pattern string) (map[string]interface{}, error) {
	query := searchQuery
	return s.runQuery(query, map[string]interface{}{"pattern": pattern})
}

// Mutation queries information about a GraphQL mutation
func (s *Schema) Mutation(mutationName string) (map[string]interface{}, error) {
	query := mutationQuery
	return s.runQuery(query, map[string]interface{}{"mutation": mutationName})
}

// Query runs a custom jq query on the schema
func (s *Schema) Query(jqQuery string, variables map[string]interface{}) (interface{}, error) {
	// Create pipeline with the query
	pipeline, err := jqyaml.New(jqyaml.WithQuery(jqQuery))
	if err != nil {
		return nil, fmt.Errorf("failed to create jq pipeline: %w", err)
	}

	// Collect results using a custom callback
	var results []interface{}
	opts := []jqyaml.ExecuteOption{
		jqyaml.WithCallback(func(item interface{}) error {
			results = append(results, item)
			return nil
		}),
	}
	
	// Add variables if provided
	if variables != nil {
		opts = append(opts, jqyaml.WithVariables(variables))
	}

	// Execute the pipeline
	ctx := context.Background()
	if err := pipeline.Execute(ctx, s.data, opts...); err != nil {
		return nil, err
	}

	// Return results based on count
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}

// runQuery is a helper to run predefined queries
func (s *Schema) runQuery(query string, variables map[string]interface{}) (map[string]interface{}, error) {
	slog.Debug("Running predefined query", "variables", variables)
	
	result, err := s.Query(query, variables)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("no results found")
	}

	if m, ok := result.(map[string]interface{}); ok {
		return m, nil
	}

	return nil, fmt.Errorf("unexpected result type: %T", result)
}