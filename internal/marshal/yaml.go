package marshal

import (
	"github.com/goccy/go-yaml"
)

// Marshal marshals data to YAML/JSON bytes using consistent options
func Marshal(v interface{}) ([]byte, error) {
	return yaml.MarshalWithOptions(v, yaml.UseJSONMarshaler())
}

// MarshalJSON marshals data to JSON bytes
func MarshalJSON(v interface{}) ([]byte, error) {
	return yaml.MarshalWithOptions(v, yaml.UseJSONMarshaler(), yaml.JSON())
}

// Unmarshal unmarshals YAML/JSON bytes using consistent options
func Unmarshal(data []byte, v interface{}) error {
	return yaml.UnmarshalWithOptions(data, v, yaml.UseJSONUnmarshaler())
}