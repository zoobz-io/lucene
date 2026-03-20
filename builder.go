package lucene

import (
	"fmt"

	"github.com/zoobz-io/sentinel"
)

// Builder provides schema-validated query building for type T.
type Builder[T any] struct {
	spec   *Spec
	fields map[string]*FieldSpec
}

// New creates a new Builder for type T.
// Returns an error if T is not a struct or cannot be inspected.
func New[T any]() (*Builder[T], error) {
	sentinel.Tag("json")
	sentinel.Tag("lucene")

	metadata, err := sentinel.TryInspect[T]()
	if err != nil {
		return nil, fmt.Errorf("lucene: failed to inspect type: %w", err)
	}

	spec := buildSpec(metadata)
	fields := make(map[string]*FieldSpec, len(spec.Fields))
	for i := range spec.Fields {
		fields[spec.Fields[i].Name] = &spec.Fields[i]
	}

	return &Builder[T]{spec: spec, fields: fields}, nil
}

// Spec returns the extracted schema specification.
func (b *Builder[T]) Spec() *Spec {
	return b.spec
}

// resolveField looks up a field by name and returns its spec.
// Returns an error if the field does not exist.
func (b *Builder[T]) resolveField(name string) (*FieldSpec, error) {
	if spec, ok := b.fields[name]; ok {
		return spec, nil
	}
	return nil, fmt.Errorf("%w: %s", ErrUnknownField, name)
}

// validateField checks if a field exists and returns an error query if not.
// This is used internally by query builders for deferred error handling.
func (b *Builder[T]) validateField(op Op, field string) (*FieldSpec, *query) {
	spec, err := b.resolveField(field)
	if err != nil {
		return nil, errQuery(op, err)
	}
	return spec, nil
}
