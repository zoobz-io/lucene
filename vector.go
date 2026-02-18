package lucene

// KnnQuery performs k-nearest neighbors vector search.
type KnnQuery struct {
	query
	vector        []float32
	k             *int
	numCandidates *int
	filter        Query
	boost         *float64
}

// K sets the number of nearest neighbors to return.
func (q *KnnQuery) K(k int) *KnnQuery { q.k = &k; return q }

// NumCandidates sets the number of candidates to consider.
func (q *KnnQuery) NumCandidates(n int) *KnnQuery { q.numCandidates = &n; return q }

// Filter sets a filter to apply to candidates.
func (q *KnnQuery) Filter(f Query) *KnnQuery { q.filter = f; return q }

// Boost sets the relevance score multiplier.
func (q *KnnQuery) Boost(b float64) *KnnQuery { q.boost = &b; return q }

// Vector returns the query vector.
func (q *KnnQuery) Vector() []float32 { return q.vector }

// KValue returns the k value if set.
func (q *KnnQuery) KValue() *int { return q.k }

// NumCandidatesValue returns the num_candidates value if set.
func (q *KnnQuery) NumCandidatesValue() *int { return q.numCandidates }

// FilterQuery returns the filter query if set.
func (q *KnnQuery) FilterQuery() Query { return q.filter }

// BoostValue returns the boost value if set.
func (q *KnnQuery) BoostValue() *float64 { return q.boost }

// Err returns any error in this query or its filter.
func (q *KnnQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	if q.filter != nil {
		return q.filter.Err()
	}
	return nil
}

// Knn creates a kNN query.
func (b *Builder[T]) Knn(field string, vector []float32) *KnnQuery {
	spec, errQ := b.validateField(OpKnn, field)
	if errQ != nil {
		return &KnnQuery{query: *errQ}
	}
	return &KnnQuery{
		query:  query{op: OpKnn, field: spec.Name},
		vector: vector,
	}
}
