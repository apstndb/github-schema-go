package schema

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// Since SchemaURL is a const, we need to test with a mock transport
type mockTransport struct {
	handler http.HandlerFunc
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	m.handler(rec, req)
	return rec.Result(), nil
}

// Test DownloadToWriter with mocked HTTP client
func TestDownloadToWriter(t *testing.T) {
	// Create test schema
	testSchema := `{"data": {"__schema": {"types": []}}}`
	
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() == SchemaURL {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(testSchema))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		},
	}
	
	// Test download to buffer
	var buf bytes.Buffer
	err := DownloadToWriter(&buf)
	if err != nil {
		t.Fatalf("DownloadToWriter failed: %v", err)
	}
	
	if buf.String() != testSchema {
		t.Errorf("Expected %s, got %s", testSchema, buf.String())
	}
}

// Test DownloadAndCompressToWriter
func TestDownloadAndCompressToWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test that requires network access in short mode")
	}
	// Create test schema
	testSchema := `{"data": {"__schema": {"types": []}}}`
	
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport that simulates gzip response from GitHub
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() == SchemaURL {
				// Check if client requested gzip
				if r.Header.Get("Accept-Encoding") == "gzip" {
					// Simulate GitHub's gzip response
					w.Header().Set("Content-Encoding", "gzip")
					w.WriteHeader(http.StatusOK)
					// Write pre-compressed data
					gz := gzip.NewWriter(w)
					gz.Write([]byte(testSchema))
					gz.Close()
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(testSchema))
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		},
	}
	
	// Test download and compress to buffer
	var buf bytes.Buffer
	err := DownloadAndCompressToWriter(&buf)
	if err != nil {
		t.Fatalf("DownloadAndCompressToWriter failed: %v", err)
	}
	
	// Decompress and verify
	reader, err := gzip.NewReader(&buf)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer reader.Close()
	
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	
	if string(decompressed) != testSchema {
		t.Errorf("Expected %s, got %s", testSchema, string(decompressed))
	}
}

// Test DownloadSchema to file
func TestDownloadSchema(t *testing.T) {
	// Create test schema
	testSchema := `{"data": {"__schema": {"types": []}}}`
	
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() == SchemaURL {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(testSchema))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		},
	}
	
	// Create temp file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "schema.json")
	
	// Test download to file
	err := DownloadSchema(outputPath)
	if err != nil {
		t.Fatalf("DownloadSchema failed: %v", err)
	}
	
	// Verify file contents
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	
	if string(content) != testSchema {
		t.Errorf("Expected %s, got %s", testSchema, string(content))
	}
}

// Test DownloadAndCompressSchema to file
func TestDownloadAndCompressSchema(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test that requires network access in short mode")
	}
	// Create test schema
	testSchema := `{"data": {"__schema": {"types": []}}}`
	
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport that simulates gzip response from GitHub
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			if r.URL.String() == SchemaURL {
				// Check if client requested gzip
				if r.Header.Get("Accept-Encoding") == "gzip" {
					// Simulate GitHub's gzip response
					w.Header().Set("Content-Encoding", "gzip")
					w.WriteHeader(http.StatusOK)
					// Write pre-compressed data
					gz := gzip.NewWriter(w)
					gz.Write([]byte(testSchema))
					gz.Close()
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(testSchema))
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		},
	}
	
	// Create temp file
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "schema.json.gz")
	
	// Test download and compress to file
	err := DownloadAndCompressSchema(outputPath)
	if err != nil {
		t.Fatalf("DownloadAndCompressSchema failed: %v", err)
	}
	
	// Verify file contents by decompressing
	file, err := os.Open(outputPath)
	if err != nil {
		t.Fatalf("Failed to open compressed file: %v", err)
	}
	defer file.Close()
	
	reader, err := gzip.NewReader(file)
	if err != nil {
		t.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer reader.Close()
	
	decompressed, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("Failed to decompress: %v", err)
	}
	
	if string(decompressed) != testSchema {
		t.Errorf("Expected %s, got %s", testSchema, string(decompressed))
	}
}

// Test HTTP error handling
func TestDownloadToWriter_HTTPError(t *testing.T) {
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport that returns 404
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		},
	}
	
	var buf bytes.Buffer
	err := DownloadToWriter(&buf)
	if err == nil {
		t.Error("Expected error for HTTP 404")
	}
	if err != nil && !strings.Contains(err.Error(), "HTTP 404") {
		t.Errorf("Expected HTTP 404 error, got: %v", err)
	}
}

// Test network error handling
func TestDownloadToWriter_NetworkError(t *testing.T) {
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set transport that always fails
	http.DefaultTransport = &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return nil, fmt.Errorf("network error")
		},
	}
	
	var buf bytes.Buffer
	err := DownloadToWriter(&buf)
	if err == nil {
		t.Error("Expected error for network failure")
	}
}

// Test file creation error
func TestDownloadSchema_FileCreationError(t *testing.T) {
	// Use an invalid path that cannot be created
	invalidPath := "/root/no-permission/schema.json"
	
	err := DownloadSchema(invalidPath)
	if err == nil {
		t.Error("Expected error for invalid file path")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to create output file") {
		t.Errorf("Expected file creation error, got: %v", err)
	}
}

// Benchmark download performance
func BenchmarkDownloadToWriter(b *testing.B) {
	// Create a large test schema (1MB)
	largeSchema := strings.Repeat(`{"data": {"__schema": {"types": []}}}`, 30000)
	
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(largeSchema))
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		err := DownloadToWriter(&buf)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test concurrent downloads
func TestConcurrentDownloads(t *testing.T) {
	// Create test schema
	testSchema := `{"data": {"__schema": {"types": []}}}`
	
	// Save original default transport
	originalTransport := http.DefaultTransport
	defer func() {
		http.DefaultTransport = originalTransport
	}()
	
	// Set mock transport with delay to test concurrency
	http.DefaultTransport = &mockTransport{
		handler: func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(10 * time.Millisecond) // Small delay
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(testSchema))
		},
	}
	
	// Run multiple downloads concurrently
	const numGoroutines = 10
	errors := make(chan error, numGoroutines)
	
	for i := 0; i < numGoroutines; i++ {
		go func() {
			var buf bytes.Buffer
			err := DownloadToWriter(&buf)
			errors <- err
		}()
	}
	
	// Wait for all downloads
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent download failed: %v", err)
		}
	}
}