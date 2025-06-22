package schema

// This file contains the go:generate directive to update the embedded schema

//go:generate go run ../cmd/github-schema/main.go download --compress -o schema.json.gz