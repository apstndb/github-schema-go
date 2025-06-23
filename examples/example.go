package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/apstndb/github-schema-go/schema"
)

func main() {
	// Create schema instance with embedded introspection data
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
		if desc, ok := typeInfo["description"].(string); ok && len(desc) > 100 {
			fmt.Printf("Description: %.100s...\n", desc)
		}
		if fields, ok := typeInfo["fields"].([]interface{}); ok {
			fmt.Printf("Field count: %d\n", len(fields))
			// Show first 3 fields
			fmt.Println("First 3 fields:")
			for i := 0; i < 3 && i < len(fields); i++ {
				if field, ok := fields[i].(map[string]interface{}); ok {
					fmt.Printf("  - %s: %s\n", field["name"], field["type"])
				}
			}
		}
	}

	// Example 2: Search for types
	fmt.Println("\n=== Search for Review types ===")
	searchResult, err := s.Search("Review")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Found %v types matching 'Review':\n", searchResult["count"])
	if results, ok := searchResult["results"].([]interface{}); ok {
		for _, r := range results {
			if t, ok := r.(map[string]interface{}); ok {
				fmt.Printf("  - %s (%s)\n", t["name"], t["kind"])
			}
		}
	}

	// Example 3: Query a mutation
	fmt.Println("\n=== createIssue Mutation ===")
	mutationResult, err := s.Mutation("createIssue")
	if err != nil {
		log.Fatal(err)
	}
	
	if mutation, ok := mutationResult["mutation"].(map[string]interface{}); ok {
		fmt.Printf("Mutation: %s\n", mutation["name"])
		if inputs, ok := mutation["inputs"].([]interface{}); ok && len(inputs) > 0 {
			if input, ok := inputs[0].(map[string]interface{}); ok {
				fmt.Printf("Input type: %s\n", input["type"])
				fmt.Printf("Required: %v\n", input["required"])
				// Description contains field details
				if desc, ok := input["description"].(string); ok && len(desc) > 200 {
					fmt.Printf("Input fields preview: %.200s...\n", desc)
				}
			}
		}
	}

	// Example 4: Custom jq query - Get all interface types
	fmt.Println("\n=== Custom Query: List all interfaces ===")
	interfaces, err := s.Query(`.data.__schema.types[] | select(.kind == "INTERFACE") | .name`, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	// Results come as individual values, collect them
	if name, ok := interfaces.(string); ok {
		fmt.Printf("First interface: %s\n", name)
	}

	// Example 5: Using custom schema file (commented out)
	/*
	customSchema, err := schema.NewWithFile("./my-schema.json")
	if err != nil {
		log.Fatal(err)
	}
	
	// Use customSchema the same way...
	*/

	// Example 6: Complex query with variables
	fmt.Println("\n=== Query with variables: Find fields of a specific type ===")
	fieldsQuery := `.data.__schema.types[] | select(.name == $typename) | .fields[]? | {name, type: .type.name}`
	vars := map[string]interface{}{"typename": "Repository"}
	
	fields, err := s.Query(fieldsQuery, vars)
	if err != nil {
		log.Fatal(err)
	}
	
	// Pretty print the first result
	if field, ok := fields.(map[string]interface{}); ok {
		jsonBytes, _ := json.MarshalIndent(field, "", "  ")
		fmt.Printf("First field of Repository type:\n%s\n", jsonBytes)
	}
}