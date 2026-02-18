package lucene

import "errors"

// BoolQuery combines queries with boolean logic.
type BoolQuery struct {
	query
	must               []Query
	should             []Query
	mustNot            []Query
	filter             []Query
	minimumShouldMatch *int
	boost              *float64
}

// Must adds queries that must match.
func (q *BoolQuery) Must(queries ...Query) *BoolQuery {
	q.must = append(q.must, queries...)
	return q
}

// Should adds queries that should match.
func (q *BoolQuery) Should(queries ...Query) *BoolQuery {
	q.should = append(q.should, queries...)
	return q
}

// MustNot adds queries that must not match.
func (q *BoolQuery) MustNot(queries ...Query) *BoolQuery {
	q.mustNot = append(q.mustNot, queries...)
	return q
}

// Filter adds queries that must match but don't affect scoring.
func (q *BoolQuery) Filter(queries ...Query) *BoolQuery {
	q.filter = append(q.filter, queries...)
	return q
}

// MinimumShouldMatch sets the minimum number of should clauses that must match.
func (q *BoolQuery) MinimumShouldMatch(n int) *BoolQuery {
	q.minimumShouldMatch = &n
	return q
}

// Boost sets the relevance score multiplier.
func (q *BoolQuery) Boost(b float64) *BoolQuery {
	q.boost = &b
	return q
}

// MustClauses returns the must clauses.
func (q *BoolQuery) MustClauses() []Query { return q.must }

// ShouldClauses returns the should clauses.
func (q *BoolQuery) ShouldClauses() []Query { return q.should }

// MustNotClauses returns the must_not clauses.
func (q *BoolQuery) MustNotClauses() []Query { return q.mustNot }

// FilterClauses returns the filter clauses.
func (q *BoolQuery) FilterClauses() []Query { return q.filter }

// MinimumShouldMatchValue returns the minimum_should_match value if set.
func (q *BoolQuery) MinimumShouldMatchValue() *int { return q.minimumShouldMatch }

// BoostValue returns the boost value if set.
func (q *BoolQuery) BoostValue() *float64 { return q.boost }

// Err returns the first error found in this query or any of its children.
func (q *BoolQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	for _, child := range q.must {
		if err := child.Err(); err != nil {
			return err
		}
	}
	for _, child := range q.should {
		if err := child.Err(); err != nil {
			return err
		}
	}
	for _, child := range q.mustNot {
		if err := child.Err(); err != nil {
			return err
		}
	}
	for _, child := range q.filter {
		if err := child.Err(); err != nil {
			return err
		}
	}
	return nil
}

// Bool creates a bool query for combining queries with boolean logic.
func (b *Builder[T]) Bool() *BoolQuery {
	return &BoolQuery{
		query: query{op: OpBool},
	}
}

// MatchAllQuery matches all documents.
type MatchAllQuery struct {
	query
	boost *float64
}

// Boost sets the relevance score multiplier.
func (q *MatchAllQuery) Boost(b float64) *MatchAllQuery {
	q.boost = &b
	return q
}

// BoostValue returns the boost value if set.
func (q *MatchAllQuery) BoostValue() *float64 { return q.boost }

// MatchAll creates a query that matches all documents.
func (b *Builder[T]) MatchAll() *MatchAllQuery {
	return &MatchAllQuery{
		query: query{op: OpMatchAll},
	}
}

// MatchNoneQuery matches no documents.
type MatchNoneQuery struct {
	query
}

// MatchNone creates a query that matches no documents.
func (b *Builder[T]) MatchNone() *MatchNoneQuery {
	return &MatchNoneQuery{
		query: query{op: OpMatchNone},
	}
}

// And is a convenience method that creates a bool query with must clauses.
func (b *Builder[T]) And(queries ...Query) *BoolQuery {
	return b.Bool().Must(queries...)
}

// Or is a convenience method that creates a bool query with should clauses
// and minimum_should_match set to 1.
func (b *Builder[T]) Or(queries ...Query) *BoolQuery {
	return b.Bool().Should(queries...).MinimumShouldMatch(1)
}

// Not is a convenience method that creates a bool query with a must_not clause.
func (b *Builder[T]) Not(q Query) *BoolQuery {
	if q == nil {
		return &BoolQuery{
			query: query{op: OpBool, err: errors.New("cannot negate nil query")},
		}
	}
	return b.Bool().MustNot(q)
}

// BoostingQuery demotes documents matching a negative query.
type BoostingQuery struct {
	query
	positive      Query
	negative      Query
	negativeBoost *float64
}

// Positive sets the query that must match.
func (q *BoostingQuery) Positive(p Query) *BoostingQuery { q.positive = p; return q }

// Negative sets the query to demote matching documents.
func (q *BoostingQuery) Negative(n Query) *BoostingQuery { q.negative = n; return q }

// NegativeBoost sets the score multiplier for negative matches (0-1).
func (q *BoostingQuery) NegativeBoost(b float64) *BoostingQuery { q.negativeBoost = &b; return q }

// PositiveQuery returns the positive query.
func (q *BoostingQuery) PositiveQuery() Query { return q.positive }

// NegativeQuery returns the negative query.
func (q *BoostingQuery) NegativeQuery() Query { return q.negative }

// NegativeBoostValue returns the negative_boost value if set.
func (q *BoostingQuery) NegativeBoostValue() *float64 { return q.negativeBoost }

// Err returns any error in this query or its children.
func (q *BoostingQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	if q.positive != nil {
		if err := q.positive.Err(); err != nil {
			return err
		}
	}
	if q.negative != nil {
		if err := q.negative.Err(); err != nil {
			return err
		}
	}
	return nil
}

// Boosting creates a boosting query.
func (b *Builder[T]) Boosting() *BoostingQuery {
	return &BoostingQuery{
		query: query{op: OpBoosting},
	}
}

// DisMaxQuery returns the best match from multiple queries.
type DisMaxQuery struct {
	query
	queries    []Query
	tieBreaker *float64
	boost      *float64
}

// TieBreaker sets the tie breaker multiplier (0-1).
func (q *DisMaxQuery) TieBreaker(t float64) *DisMaxQuery { q.tieBreaker = &t; return q }

// Boost sets the relevance score multiplier.
func (q *DisMaxQuery) Boost(b float64) *DisMaxQuery { q.boost = &b; return q }

// Queries returns the dis_max queries.
func (q *DisMaxQuery) Queries() []Query { return q.queries }

// TieBreakerValue returns the tie_breaker value if set.
func (q *DisMaxQuery) TieBreakerValue() *float64 { return q.tieBreaker }

// BoostValue returns the boost value if set.
func (q *DisMaxQuery) BoostValue() *float64 { return q.boost }

// Err returns any error in this query or its children.
func (q *DisMaxQuery) Err() error {
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

// DisMax creates a dis_max query.
func (b *Builder[T]) DisMax(queries ...Query) *DisMaxQuery {
	return &DisMaxQuery{
		query:   query{op: OpDisMax},
		queries: queries,
	}
}

// ConstantScoreQuery wraps a filter with a constant score.
type ConstantScoreQuery struct {
	query
	filter Query
	boost  *float64
}

// Boost sets the constant score value.
func (q *ConstantScoreQuery) Boost(b float64) *ConstantScoreQuery { q.boost = &b; return q }

// FilterQuery returns the wrapped filter query.
func (q *ConstantScoreQuery) FilterQuery() Query { return q.filter }

// BoostValue returns the boost value if set.
func (q *ConstantScoreQuery) BoostValue() *float64 { return q.boost }

// Err returns any error in this query or its filter.
func (q *ConstantScoreQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	if q.filter != nil {
		return q.filter.Err()
	}
	return nil
}

// ConstantScore creates a constant_score query.
func (b *Builder[T]) ConstantScore(filter Query) *ConstantScoreQuery {
	return &ConstantScoreQuery{
		query:  query{op: OpConstantScore},
		filter: filter,
	}
}
