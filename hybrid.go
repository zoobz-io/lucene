package lucene

// HybridQuery combines multiple queries with server-side score normalization.
// This is an OpenSearch-specific query type used with the neural search plugin.
type HybridQuery struct {
	query
	queries []Query
	boost   *float64
}

// Boost sets the relevance score multiplier.
func (q *HybridQuery) Boost(b float64) *HybridQuery { q.boost = &b; return q }

// Queries returns the sub-queries to combine.
func (q *HybridQuery) Queries() []Query { return q.queries }

// BoostValue returns the boost value if set.
func (q *HybridQuery) BoostValue() *float64 { return q.boost }

// Err returns any error in this query or its children.
func (q *HybridQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	for _, child := range q.queries {
		if err := child.Err(); err != nil {
			return err
		}
	}
	return nil
}

// Hybrid creates a hybrid query for combining queries with server-side score normalization.
func (b *Builder[T]) Hybrid(queries ...Query) *HybridQuery {
	return &HybridQuery{
		query:   query{op: OpHybrid},
		queries: queries,
	}
}
