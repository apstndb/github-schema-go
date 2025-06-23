package schema

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const (
	// GitHubAPIURL is the GitHub GraphQL API endpoint
	GitHubAPIURL = "https://api.github.com/graphql"
	
	// IntrospectionQuery is the GraphQL introspection query
	IntrospectionQuery = `
	{
	  __schema {
	    queryType { name }
	    mutationType { name }
	    subscriptionType { name }
	    types {
	      ...FullType
	    }
	    directives {
	      name
	      description
	      locations
	      args {
	        ...InputValue
	      }
	    }
	  }
	}
	
	fragment FullType on __Type {
	  kind
	  name
	  description
	  fields(includeDeprecated: true) {
	    name
	    description
	    args {
	      ...InputValue
	    }
	    type {
	      ...TypeRef
	    }
	    isDeprecated
	    deprecationReason
	  }
	  inputFields {
	    ...InputValue
	  }
	  interfaces {
	    ...TypeRef
	  }
	  enumValues(includeDeprecated: true) {
	    name
	    description
	    isDeprecated
	    deprecationReason
	  }
	  possibleTypes {
	    ...TypeRef
	  }
	}
	
	fragment InputValue on __InputValue {
	  name
	  description
	  type { ...TypeRef }
	  defaultValue
	}
	
	fragment TypeRef on __Type {
	  kind
	  name
	  ofType {
	    kind
	    name
	    ofType {
	      kind
	      name
	      ofType {
	        kind
	        name
	        ofType {
	          kind
	          name
	          ofType {
	            kind
	            name
	            ofType {
	              kind
	              name
	              ofType {
	                kind
	                name
	              }
	            }
	          }
	        }
	      }
	    }
	  }
	}`
)

// DownloadSchema downloads the schema using GitHub GraphQL API introspection.
// This is an alias for DownloadIntrospectionSchema for backward compatibility.
func DownloadSchema(outputPath string) error {
	return DownloadIntrospectionSchema(outputPath)
}

// DownloadAndCompressSchema downloads the schema with gzip compression.
// When possible, it uses GitHub API's native gzip compression to reduce bandwidth usage.
// The compressed data is saved directly without re-compression.
func DownloadAndCompressSchema(outputPath string) error {
	// Get GitHub token from gh auth
	cmd := exec.Command("gh", "auth", "token")
	tokenBytes, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token (run 'gh auth login'): %w", err)
	}
	token := string(bytes.TrimSpace(tokenBytes))
	
	// Prepare GraphQL request
	requestBody := map[string]string{
		"query": IntrospectionQuery,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", GitHubAPIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	
	// Use custom transport to prevent automatic decompression
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}
	
	// Check if response is compressed
	if resp.Header.Get("Content-Encoding") != "gzip" {
		// Fallback: read uncompressed and compress it
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		
		// Validate it's valid JSON
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("failed to parse response as JSON: %w", err)
		}
		
		// Check for errors in response
		if errors, ok := result["errors"]; ok {
			return fmt.Errorf("GraphQL errors: %v", errors)
		}
		
		// Create output file and compress
		out, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer out.Close()
		
		gz := gzip.NewWriter(out)
		defer gz.Close()
		
		if _, err := gz.Write(body); err != nil {
			return fmt.Errorf("failed to write compressed data: %w", err)
		}
	} else {
		// Response is already compressed, save directly
		out, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer out.Close()
		
		if _, err := io.Copy(out, resp.Body); err != nil {
			return fmt.Errorf("failed to write compressed data: %w", err)
		}
	}
	
	return nil
}

// DownloadToWriter downloads introspection schema and writes to writer
func DownloadToWriter(w io.Writer) error {
	return DownloadIntrospectionToWriter(w)
}

// DownloadAndCompressToWriter downloads introspection schema with native compression and writes to writer
func DownloadAndCompressToWriter(w io.Writer) error {
	// Get GitHub token from gh auth
	cmd := exec.Command("gh", "auth", "token")
	tokenBytes, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token (run 'gh auth login'): %w", err)
	}
	token := string(bytes.TrimSpace(tokenBytes))
	
	// Prepare GraphQL request
	requestBody := map[string]string{
		"query": IntrospectionQuery,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", GitHubAPIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	
	// Use custom transport to prevent automatic decompression
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}
	
	// Check if response is compressed
	if resp.Header.Get("Content-Encoding") != "gzip" {
		// Fallback: compress on the fly
		gz := gzip.NewWriter(w)
		defer gz.Close()
		
		if _, err := io.Copy(gz, resp.Body); err != nil {
			return fmt.Errorf("failed to write compressed response: %w", err)
		}
	} else {
		// Response is already compressed, copy directly
		if _, err := io.Copy(w, resp.Body); err != nil {
			return fmt.Errorf("failed to write compressed response: %w", err)
		}
	}
	
	return nil
}

// DownloadIntrospectionSchema downloads the GitHub GraphQL schema using the standard
// introspection query. The schema is saved in the GraphQL introspection format,
// which includes the data wrapper: {"data": {"__schema": {...}}}.
// Requires GitHub authentication via 'gh auth login'.
func DownloadIntrospectionSchema(outputPath string) error {
	// Get GitHub token from gh auth
	cmd := exec.Command("gh", "auth", "token")
	tokenBytes, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token (run 'gh auth login'): %w", err)
	}
	token := string(bytes.TrimSpace(tokenBytes))
	
	// Prepare GraphQL request
	requestBody := map[string]string{
		"query": IntrospectionQuery,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", GitHubAPIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}
	
	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	
	// Validate it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response as JSON: %w", err)
	}
	
	// Check for errors in response
	if errors, ok := result["errors"]; ok {
		return fmt.Errorf("GraphQL errors: %v", errors)
	}
	
	// Write to file
	if err := os.WriteFile(outputPath, body, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// DownloadIntrospectionToWriter downloads introspection schema and writes to writer
func DownloadIntrospectionToWriter(w io.Writer) error {
	// Get GitHub token from gh auth
	cmd := exec.Command("gh", "auth", "token")
	tokenBytes, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get GitHub token (run 'gh auth login'): %w", err)
	}
	token := string(bytes.TrimSpace(tokenBytes))
	
	// Prepare GraphQL request
	requestBody := map[string]string{
		"query": IntrospectionQuery,
	}
	
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Create HTTP request
	req, err := http.NewRequest("POST", GitHubAPIURL, bytes.NewReader(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API returned HTTP %d", resp.StatusCode)
	}
	
	// Copy response to writer
	if _, err := io.Copy(w, resp.Body); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	
	return nil
}