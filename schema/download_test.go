package schema

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestDownloadIntrospectionSchema tests the introspection download functionality
// Note: This test requires GitHub authentication via 'gh auth login'
func TestDownloadIntrospectionSchema(t *testing.T) {
	// Skip if in short mode or CI environment
	if testing.Short() {
		t.Skip("Skipping introspection download test in short mode")
	}
	
	// Create temp file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "test_schema.json")
	
	// Try to download
	err := DownloadIntrospectionSchema(outputPath)
	if err != nil {
		// If auth fails, skip the test
		if strings.Contains(err.Error(), "gh auth login") {
			t.Skip("Skipping test: GitHub authentication not available")
		}
		t.Fatalf("Failed to download schema: %v", err)
	}
	
	// Verify file exists
	info, err := os.Stat(outputPath)
	if err != nil {
		t.Fatalf("Failed to stat downloaded file: %v", err)
	}
	
	// Verify it's not empty
	if info.Size() == 0 {
		t.Error("Downloaded file is empty")
	}
	
	// Try to load it
	s, err := NewWithFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to load downloaded schema: %v", err)
	}
	
	// Try a simple query
	result, err := s.Type("Query")
	if err != nil {
		t.Fatalf("Failed to query type: %v", err)
	}
	
	if result == nil {
		t.Error("Query returned nil result")
	}
}