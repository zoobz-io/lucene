package lucene

import (
	"strings"

	"github.com/zoobz-io/sentinel"
)

// FieldKind categorizes field types for validation.
type FieldKind uint8

const (
	// KindUnknown is an unrecognized field type.
	KindUnknown FieldKind = iota
	// KindString is a string field.
	KindString
	// KindInt is an integer field (int, int64, uint, etc.).
	KindInt
	// KindFloat is a floating-point field (float32, float64).
	KindFloat
	// KindBool is a boolean field.
	KindBool
	// KindTime is a time.Time field.
	KindTime
	// KindSlice is a slice field (excluding vector types).
	KindSlice
	// KindVector is a vector embedding field ([]float32, []float64).
	KindVector
)

// FieldSpec describes a single field in the schema.
type FieldSpec struct {
	Name string    // Resolved name (json tag or Go name).
	Type string    // Go type string.
	Kind FieldKind // Categorized type.
}

// Spec holds the extracted schema for a type.
type Spec struct {
	Fields []FieldSpec
}

// buildSpec converts sentinel metadata to a Spec.
func buildSpec(metadata sentinel.Metadata) *Spec {
	fields := make([]FieldSpec, 0, len(metadata.Fields))

	for _, f := range metadata.Fields {
		name := resolveFieldName(f)
		if name == "" || name == "-" {
			continue
		}

		fields = append(fields, FieldSpec{
			Name: name,
			Type: f.Type,
			Kind: resolveFieldKind(f.Type),
		})
	}

	return &Spec{Fields: fields}
}

// resolveFieldName returns the JSON tag name or falls back to the Go field name.
func resolveFieldName(f sentinel.FieldMetadata) string {
	if tag, ok := f.Tags["json"]; ok {
		// Handle json:"name,omitempty" format.
		if idx := strings.Index(tag, ","); idx != -1 {
			return tag[:idx]
		}
		return tag
	}
	return f.Name
}

// resolveFieldKind maps Go type strings to FieldKind.
func resolveFieldKind(typeName string) FieldKind {
	// Handle pointer types.
	typeName = strings.TrimPrefix(typeName, "*")

	// Check for vector types first.
	if typeName == "[]float32" || typeName == "[]float64" {
		return KindVector
	}

	// Check for slice types.
	if strings.HasPrefix(typeName, "[]") {
		return KindSlice
	}

	// Check for time.Time.
	if typeName == "time.Time" || typeName == "Time" {
		return KindTime
	}

	// Check for basic types.
	switch {
	case typeName == "string":
		return KindString
	case typeName == "bool":
		return KindBool
	case strings.HasPrefix(typeName, "int") || strings.HasPrefix(typeName, "uint"):
		return KindInt
	case strings.HasPrefix(typeName, "float"):
		return KindFloat
	default:
		return KindUnknown
	}
}
