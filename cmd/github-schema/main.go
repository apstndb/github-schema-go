package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/apstndb/github-schema-go/schema"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	schemaFile string
	outputJSON bool
	debug      bool
)

var rootCmd = &cobra.Command{
	Use:   "github-schema",
	Short: "Query GitHub GraphQL schema offline",
	Long: `Query GitHub GraphQL schema using embedded data or custom schema files.
The embedded schema is obtained via GitHub GraphQL API introspection.`,
}

var typeCmd = &cobra.Command{
	Use:   "type <TypeName>",
	Short: "Show fields and descriptions for a type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := getSchema()
		if err != nil {
			return err
		}

		result, err := s.Type(args[0])
		if err != nil {
			return fmt.Errorf("failed to query type: %w", err)
		}

		return outputResult(result)
	},
}

var mutationCmd = &cobra.Command{
	Use:   "mutation <MutationName>",
	Short: "Show mutation input requirements",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := getSchema()
		if err != nil {
			return err
		}

		result, err := s.Mutation(args[0])
		if err != nil {
			return fmt.Errorf("failed to query mutation: %w", err)
		}

		return outputResult(result)
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <pattern>",
	Short: "Search schema for matching types/fields",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := getSchema()
		if err != nil {
			return err
		}

		result, err := s.Search(args[0])
		if err != nil {
			return fmt.Errorf("failed to search schema: %w", err)
		}

		return outputResult(result)
	},
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download latest schema via GraphQL introspection",
	Long: `Download the latest GitHub GraphQL schema using introspection query.
Requires 'gh auth login' to be configured.

Examples:
  github-schema download                           # Download to stdout
  github-schema download -o schema.json            # Download to file
  github-schema download -o schema.json.gz         # Auto-compress (detected by .gz extension)
  github-schema download --compress                # Download compressed to stdout
  github-schema download -c -o schema.json.gz      # Explicitly compress to file`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		compressFlag, _ := cmd.Flags().GetBool("compress")
		outputFile, _ := cmd.Flags().GetString("output")
		
		// If no output file specified, write to stdout
		toStdout := outputFile == ""
		
		// Determine if we should compress
		// Priority: --compress flag > .gz extension > default (no compression)
		compress := compressFlag
		if !toStdout && !compress && strings.HasSuffix(outputFile, ".gz") {
			compress = true
		}
		
		if toStdout {
			// Write to stdout
			if compress {
				return schema.DownloadAndCompressToWriter(os.Stdout)
			} else {
				return schema.DownloadToWriter(os.Stdout)
			}
		}
		
		// Write to file
		slog.Info("Downloading schema via introspection", 
			"endpoint", schema.GitHubAPIURL,
			"output", outputFile,
			"compress", compress)
		
		var err error
		if compress {
			err = schema.DownloadAndCompressSchema(outputFile)
		} else {
			err = schema.DownloadSchema(outputFile)
		}
		
		if err != nil {
			return err
		}
		
		// Get file info
		info, err := os.Stat(outputFile)
		if err != nil {
			return err
		}
		
		logAttrs := []any{
			"file", outputFile,
			"size_kb", fmt.Sprintf("%.2f", float64(info.Size())/1024),
		}
		
		if compress && !compressFlag {
			logAttrs = append(logAttrs, "auto_compressed", true)
		}
		
		slog.Info("Schema downloaded successfully", logAttrs...)
		
		return nil
	},
}

var queryCmd = &cobra.Command{
	Use:   "query <jq-expression>",
	Short: "Run custom jq query on schema",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s, err := getSchema()
		if err != nil {
			return err
		}

		result, err := s.Query(args[0], nil)
		if err != nil {
			return fmt.Errorf("failed to run query: %w", err)
		}

		return outputResult(result)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&schemaFile, "schema", "s", "", "Path to custom schema file")
	rootCmd.PersistentFlags().BoolVarP(&outputJSON, "json", "j", false, "Output as JSON instead of YAML")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug logging")

	downloadCmd.Flags().BoolP("compress", "c", false, "Compress downloaded schema with gzip")
	downloadCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")

	rootCmd.AddCommand(typeCmd, mutationCmd, searchCmd, downloadCmd, queryCmd)
}

func main() {
	// Parse flags early to get debug setting
	rootCmd.ParseFlags(os.Args[1:])
	
	// Configure slog to write to stderr with text handler
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}
	
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)
	
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command failed", "error", err)
		os.Exit(1)
	}
}

func getSchema() (*schema.Schema, error) {
	if schemaFile != "" {
		return schema.NewWithFile(schemaFile)
	}
	return schema.New()
}

func outputResult(result interface{}) error {
	if outputJSON {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(result)
	}

	// Default to YAML
	encoder := yaml.NewEncoder(os.Stdout)
	return encoder.Encode(result)
}