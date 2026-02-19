// Package marshal provides typed JSON marshaling for Elasticsearch/OpenSearch queries.
package marshal

import "encoding/json"

// FieldQuery wraps a query type with a dynamic field name as the JSON key.
// Used for queries like term, match, range where the field name is a key.
// Example: {"term": {"status": {"value": "active"}}}.
type FieldQuery[T any] struct {
	QueryType string // e.g., "term", "match", "range"
	Field     string // The field name (becomes JSON key)
	Inner     T      // The inner query structure
}

// MarshalJSON implements json.Marshaler.
func (f FieldQuery[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]map[string]T{
		f.QueryType: {f.Field: f.Inner},
	})
}

// SimpleQuery wraps a query type without a dynamic field name.
// Used for queries like bool, match_all, multi_match.
// Example: {"bool": {"must": [...]}}.
type SimpleQuery[T any] struct {
	QueryType string // e.g., "bool", "match_all"
	Inner     T      // The inner query structure
}

// MarshalJSON implements json.Marshaler.
func (s SimpleQuery[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]T{
		s.QueryType: s.Inner,
	})
}

// Agg wraps an aggregation with its type key and optional sub-aggregations.
// Example: {"terms": {"field": "status"}, "aggs": {...}}.
type Agg[T any] struct {
	AggType string         // e.g., "terms", "avg"
	Inner   T              // The aggregation configuration
	SubAggs map[string]any // Optional sub-aggregations
}

// MarshalJSON implements json.Marshaler.
func (a Agg[T]) MarshalJSON() ([]byte, error) {
	result := map[string]any{a.AggType: a.Inner}
	if len(a.SubAggs) > 0 {
		result["aggs"] = a.SubAggs
	}
	return json.Marshal(result)
}
