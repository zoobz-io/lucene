// Package lucene provides a type-safe query builder for OpenSearch and Elasticsearch.
package lucene

import "errors"

// ErrUnknownField is returned when a field is not found in the schema.
var ErrUnknownField = errors.New("unknown field")

// Op represents a query operator type.
type Op uint8

const (
	// OpMatch is an analyzed text search query.
	OpMatch Op = iota
	// OpMatchPhrase matches an exact phrase.
	OpMatchPhrase
	// OpMatchPhrasePrefix matches a phrase prefix for autocomplete.
	OpMatchPhrasePrefix
	// OpMultiMatch searches across multiple fields.
	OpMultiMatch
	// OpQueryString parses Lucene query syntax.
	OpQueryString
	// OpSimpleQueryString parses user-friendly query syntax.
	OpSimpleQueryString

	// OpTerm matches an exact value.
	OpTerm
	// OpTerms matches multiple exact values.
	OpTerms
	// OpRange matches a numeric or date range.
	OpRange
	// OpPrefix matches a field prefix.
	OpPrefix
	// OpWildcard matches a wildcard pattern.
	OpWildcard
	// OpRegexp matches a regular expression.
	OpRegexp
	// OpFuzzy matches with edit distance tolerance.
	OpFuzzy
	// OpExists matches documents where the field exists.
	OpExists
	// OpIDs matches specific document IDs.
	OpIDs

	// OpBool combines queries with boolean logic.
	OpBool
	// OpBoosting demotes documents matching a negative query.
	OpBoosting
	// OpConstantScore wraps a query with a fixed score.
	OpConstantScore
	// OpDisMax returns the best match from multiple queries.
	OpDisMax

	// OpMatchAll matches all documents.
	OpMatchAll
	// OpMatchNone matches no documents.
	OpMatchNone
	// OpNested queries nested object fields.
	OpNested
	// OpHasChild matches parents by child documents.
	OpHasChild
	// OpHasParent matches children by parent documents.
	OpHasParent

	// OpKnn performs k-nearest neighbors vector search.
	OpKnn
	// OpHybrid combines multiple queries with server-side score normalization.
	OpHybrid

	// OpGeoDistance matches documents within a radius.
	OpGeoDistance
	// OpGeoBoundingBox matches documents within a bounding box.
	OpGeoBoundingBox
)

// Query is the interface implemented by all query types.
type Query interface {
	// Op returns the operator type for this query.
	Op() Op

	// Err returns any error associated with this query.
	// Errors are deferred until explicitly checked.
	Err() error

	// sealed prevents external implementations.
	sealed()
}

// query is the base struct embedded by all query types.
type query struct {
	op    Op
	field string
	value any
	err   error
}

func (q *query) Op() Op       { return q.op }
func (q *query) Err() error   { return q.err }
func (q *query) Field() string { return q.field }
func (q *query) Value() any   { return q.value }
func (q *query) sealed()      {}

// errQuery creates a query that carries an error.
func errQuery(op Op, err error) *query {
	return &query{op: op, err: err}
}
