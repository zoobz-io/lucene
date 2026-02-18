package lucene

// NestedQuery queries nested object fields.
type NestedQuery struct {
	query
	path          string
	innerQuery    Query
	scoreMode     *string
	ignoreUnmapped *bool
}

// ScoreMode sets how nested scores are combined ("avg", "sum", "min", "max", "none").
func (q *NestedQuery) ScoreMode(m string) *NestedQuery { q.scoreMode = &m; return q }

// IgnoreUnmapped sets whether to ignore unmapped paths.
func (q *NestedQuery) IgnoreUnmapped(b bool) *NestedQuery { q.ignoreUnmapped = &b; return q }

// Path returns the nested path.
func (q *NestedQuery) Path() string { return q.path }

// InnerQuery returns the inner query.
func (q *NestedQuery) InnerQuery() Query { return q.innerQuery }

// ScoreModeValue returns the score_mode value if set.
func (q *NestedQuery) ScoreModeValue() *string { return q.scoreMode }

// IgnoreUnmappedValue returns the ignore_unmapped value if set.
func (q *NestedQuery) IgnoreUnmappedValue() *bool { return q.ignoreUnmapped }

// Err returns any error in this query or its inner query.
func (q *NestedQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	if q.innerQuery != nil {
		return q.innerQuery.Err()
	}
	return nil
}

// Nested creates a nested query.
func (b *Builder[T]) Nested(path string, inner Query) *NestedQuery {
	return &NestedQuery{
		query:      query{op: OpNested},
		path:       path,
		innerQuery: inner,
	}
}

// HasChildQuery matches parents by child documents.
type HasChildQuery struct {
	query
	childType     string
	innerQuery    Query
	scoreMode     *string
	minChildren   *int
	maxChildren   *int
	ignoreUnmapped *bool
}

// ScoreMode sets how child scores affect parent score.
func (q *HasChildQuery) ScoreMode(m string) *HasChildQuery { q.scoreMode = &m; return q }

// MinChildren sets the minimum number of children that must match.
func (q *HasChildQuery) MinChildren(n int) *HasChildQuery { q.minChildren = &n; return q }

// MaxChildren sets the maximum number of children that must match.
func (q *HasChildQuery) MaxChildren(n int) *HasChildQuery { q.maxChildren = &n; return q }

// IgnoreUnmapped sets whether to ignore unmapped types.
func (q *HasChildQuery) IgnoreUnmapped(b bool) *HasChildQuery { q.ignoreUnmapped = &b; return q }

// ChildType returns the child type.
func (q *HasChildQuery) ChildType() string { return q.childType }

// InnerQuery returns the inner query.
func (q *HasChildQuery) InnerQuery() Query { return q.innerQuery }

// ScoreModeValue returns the score_mode value if set.
func (q *HasChildQuery) ScoreModeValue() *string { return q.scoreMode }

// MinChildrenValue returns the min_children value if set.
func (q *HasChildQuery) MinChildrenValue() *int { return q.minChildren }

// MaxChildrenValue returns the max_children value if set.
func (q *HasChildQuery) MaxChildrenValue() *int { return q.maxChildren }

// IgnoreUnmappedValue returns the ignore_unmapped value if set.
func (q *HasChildQuery) IgnoreUnmappedValue() *bool { return q.ignoreUnmapped }

// Err returns any error in this query or its inner query.
func (q *HasChildQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	if q.innerQuery != nil {
		return q.innerQuery.Err()
	}
	return nil
}

// HasChild creates a has_child query.
func (b *Builder[T]) HasChild(childType string, inner Query) *HasChildQuery {
	return &HasChildQuery{
		query:      query{op: OpHasChild},
		childType:  childType,
		innerQuery: inner,
	}
}

// HasParentQuery matches children by parent documents.
type HasParentQuery struct {
	query
	parentType    string
	innerQuery    Query
	score         *bool
	ignoreUnmapped *bool
}

// Score sets whether to include the parent score.
func (q *HasParentQuery) Score(b bool) *HasParentQuery { q.score = &b; return q }

// IgnoreUnmapped sets whether to ignore unmapped types.
func (q *HasParentQuery) IgnoreUnmapped(b bool) *HasParentQuery { q.ignoreUnmapped = &b; return q }

// ParentType returns the parent type.
func (q *HasParentQuery) ParentType() string { return q.parentType }

// InnerQuery returns the inner query.
func (q *HasParentQuery) InnerQuery() Query { return q.innerQuery }

// ScoreValue returns the score value if set.
func (q *HasParentQuery) ScoreValue() *bool { return q.score }

// IgnoreUnmappedValue returns the ignore_unmapped value if set.
func (q *HasParentQuery) IgnoreUnmappedValue() *bool { return q.ignoreUnmapped }

// Err returns any error in this query or its inner query.
func (q *HasParentQuery) Err() error {
	if q.err != nil {
		return q.err
	}
	if q.innerQuery != nil {
		return q.innerQuery.Err()
	}
	return nil
}

// HasParent creates a has_parent query.
func (b *Builder[T]) HasParent(parentType string, inner Query) *HasParentQuery {
	return &HasParentQuery{
		query:      query{op: OpHasParent},
		parentType: parentType,
		innerQuery: inner,
	}
}
