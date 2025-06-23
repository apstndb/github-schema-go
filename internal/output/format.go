package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

// Format represents the output format for structured data
type Format string

const (
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
)

// IsValid checks if the format is supported
func (f Format) IsValid() bool {
	return f == FormatYAML || f == FormatJSON
}

// ParseFormat parses a string into a Format
func ParseFormat(s string) (Format, error) {
	format := Format(strings.ToLower(s))
	if !format.IsValid() {
		return "", fmt.Errorf("invalid format: %s (valid: yaml, json)", s)
	}
	return format, nil
}

// NewEncoder creates a new encoder for the specified format using goccy/go-yaml
func NewEncoder(w io.Writer, format Format) *yaml.Encoder {
	switch format {
	case FormatJSON:
		return yaml.NewEncoder(w, 
			yaml.JSON(),
			yaml.UseJSONMarshaler(),
		)
	case FormatYAML:
		return yaml.NewEncoder(w,
			yaml.UseJSONMarshaler(),
			yaml.UseLiteralStyleIfMultiline(true),
		)
	default:
		return yaml.NewEncoder(w, yaml.UseJSONMarshaler()) // Default to YAML
	}
}

// ResolveFormat resolves the output format from command flags
// Handles --json flag, defaults to YAML
func ResolveFormat(cmd *cobra.Command) Format {
	if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
		return FormatJSON
	}
	return FormatYAML
}