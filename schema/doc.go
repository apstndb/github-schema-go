// Package schema provides methods to query GitHub GraphQL schema offline.
//
// The embedded schema is obtained via GitHub GraphQL API introspection and
// stored in the standard GraphQL introspection format. This allows querying
// GitHub's GraphQL type system without making API calls.
//
// Basic usage:
//
//	// Use embedded schema
//	s, err := schema.New()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Query type information
//	result, err := s.Type("Repository")
//
// The schema can be updated using go:generate or the CLI tool:
//
//	go generate ./schema
//	# or
//	github-schema download --compress -o schema/schema.json.gz
//
// Custom schemas can be loaded from files:
//
//	s, err := schema.NewWithFile("custom-schema.json")
//
// The schema file must be in GraphQL introspection format with the standard
// structure: {"data": {"__schema": {...}}}
package schema