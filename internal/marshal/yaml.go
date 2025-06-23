package marshal

import (
	"io"

	"github.com/goccy/go-yaml"
)

// Common options for consistent behavior
var (
	// MarshalOptions are the default options for marshaling
	MarshalOptions = []yaml.EncodeOption{
		yaml.UseJSONMarshaler(),
		yaml.AutoInt(),
	}
	
	// UnmarshalOptions are the default options for unmarshaling
	UnmarshalOptions = []yaml.DecodeOption{
		yaml.UseJSONUnmarshaler(),
	}
)

// Marshal marshals data to YAML bytes using consistent options
func Marshal(v interface{}) ([]byte, error) {
	return yaml.MarshalWithOptions(v, MarshalOptions...)
}

// MarshalJSON marshals data to JSON bytes
func MarshalJSON(v interface{}) ([]byte, error) {
	opts := append([]yaml.EncodeOption{}, MarshalOptions...)
	opts = append(opts, yaml.JSON())
	return yaml.MarshalWithOptions(v, opts...)
}

// Unmarshal unmarshals YAML/JSON bytes using consistent options
func Unmarshal(data []byte, v interface{}) error {
	return yaml.UnmarshalWithOptions(data, v, UnmarshalOptions...)
}

// NewEncoder creates a new YAML encoder with consistent options
func NewEncoder(w io.Writer, opts ...yaml.EncodeOption) *yaml.Encoder {
	allOpts := append([]yaml.EncodeOption{}, MarshalOptions...)
	allOpts = append(allOpts, opts...)
	return yaml.NewEncoder(w, allOpts...)
}

// NewJSONEncoder creates a new JSON encoder with consistent options
func NewJSONEncoder(w io.Writer, opts ...yaml.EncodeOption) *yaml.Encoder {
	allOpts := append([]yaml.EncodeOption{}, MarshalOptions...)
	allOpts = append(allOpts, yaml.JSON())
	allOpts = append(allOpts, opts...)
	return yaml.NewEncoder(w, allOpts...)
}