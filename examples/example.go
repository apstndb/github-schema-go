package main

import (
	"fmt"
	"log"

	"github.com/apstndb/github-schema-go/schema"
)

func main() {
	// Create schema instance with embedded data
	s, err := schema.New()
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Query a type
	fmt.Println("=== PullRequest Type ===")
	result, err := s.Type("PullRequest")
	if err != nil {
		log.Fatal(err)
	}
	
	if typeInfo, ok := result["type"].(map[string]interface{}); ok {
		fmt.Printf("Type: %s\n", typeInfo["name"])
		fmt.Printf("Kind: %s\n", typeInfo["kind"])
		if fields, ok := typeInfo["fields"].([]interface{}); ok {
			fmt.Printf("Field count: %d\n\n", len(fields))
		}
	}

	// Example 2: Search for types
	fmt.Println("=== Search for Thread types ===")
	searchResult, err := s.Search("Thread")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Found %v types\n\n", searchResult["count"])

	// Example 3: Query a mutation
	fmt.Println("=== createIssue Mutation ===")
	mutationResult, err := s.Mutation("createIssue")
	if err != nil {
		log.Fatal(err)
	}
	
	if mutation, ok := mutationResult["mutation"].(map[string]interface{}); ok {
		fmt.Printf("Mutation: %s\n", mutation["name"])
		fmt.Printf("Description: %s\n", mutation["description"])
	}

	// Example 4: Custom jq query
	fmt.Println("\n=== Custom Query: List all mutations ===")
	mutations, err := s.Query(`.data.__schema.types[] | select(.name == "Mutation") | .fields[].name`, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	// Results come as separate values from jq
	fmt.Println("First few mutations:", mutations)
}