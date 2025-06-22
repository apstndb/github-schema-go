package schema

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	// SchemaURL is the official GitHub GraphQL schema location
	SchemaURL = "https://raw.githubusercontent.com/github/docs/main/src/graphql/data/fpt/schema.json"
)

// DownloadSchema downloads the latest schema from github/docs
func DownloadSchema(outputPath string) error {
	// Request uncompressed content
	req, err := http.NewRequest("GET", SchemaURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download schema: HTTP %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy data
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write schema: %w", err)
	}

	return nil
}

// DownloadAndCompressSchema downloads the schema already compressed from GitHub
func DownloadAndCompressSchema(outputPath string) error {
	// Request gzip-compressed content
	req, err := http.NewRequest("GET", SchemaURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	
	// Use custom transport to prevent automatic decompression
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download schema: HTTP %d", resp.StatusCode)
	}

	// Create output file
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	// Copy already-compressed data directly
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write compressed schema: %w", err)
	}

	return nil
}

// DownloadToWriter downloads the schema and writes to the provided writer
func DownloadToWriter(w io.Writer) error {
	resp, err := http.Get(SchemaURL)
	if err != nil {
		return fmt.Errorf("failed to download schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download schema: HTTP %d", resp.StatusCode)
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write schema: %w", err)
	}

	return nil
}

// DownloadAndCompressToWriter downloads and writes already compressed data from GitHub
func DownloadAndCompressToWriter(w io.Writer) error {
	// Request gzip-compressed content
	req, err := http.NewRequest("GET", SchemaURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept-Encoding", "gzip")
	
	// Use custom transport to prevent automatic decompression
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download schema: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download schema: HTTP %d", resp.StatusCode)
	}

	// Verify we got compressed content
	if resp.Header.Get("Content-Encoding") != "gzip" {
		// Fallback: compress it ourselves if server didn't
		gz := gzip.NewWriter(w)
		defer gz.Close()
		_, err = io.Copy(gz, resp.Body)
	} else {
		// Copy already-compressed data directly
		_, err = io.Copy(w, resp.Body)
	}
	
	if err != nil {
		return fmt.Errorf("failed to write compressed schema: %w", err)
	}

	return nil
}