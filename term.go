package lucene

// TermQuery matches documents with an exact value.
type TermQuery struct {
	query
	boost *float64
}

// Boost sets the relevance score multiplier.
func (q *TermQuery) Boost(b float64) *TermQuery {
	q.boost = &b
	return q
}

// BoostValue returns the boost value if set.
func (q *TermQuery) BoostValue() *float64 { return q.boost }

// Term creates a term query for exact value matching.
func (b *Builder[T]) Term(field string, value any) *TermQuery {
	spec, errQ := b.validateField(OpTerm, field)
	if errQ != nil {
		return &TermQuery{query: *errQ}
	}
	return &TermQuery{
		query: query{op: OpTerm, field: spec.Name, value: value},
	}
}

// TermsQuery matches documents with any of the specified values.
type TermsQuery struct {
	query
	values []any
	boost  *float64
}

// Boost sets the relevance score multiplier.
func (q *TermsQuery) Boost(b float64) *TermsQuery {
	q.boost = &b
	return q
}

// Values returns the term values.
func (q *TermsQuery) Values() []any { return q.values }

// BoostValue returns the boost value if set.
func (q *TermsQuery) BoostValue() *float64 { return q.boost }

// Terms creates a terms query for matching any of multiple values.
func (b *Builder[T]) Terms(field string, values ...any) *TermsQuery {
	spec, errQ := b.validateField(OpTerms, field)
	if errQ != nil {
		return &TermsQuery{query: *errQ}
	}
	return &TermsQuery{
		query:  query{op: OpTerms, field: spec.Name},
		values: values,
	}
}

// RangeQuery matches documents within a range.
type RangeQuery struct {
	query
	gt, gte, lt, lte any
	format           *string
	boost            *float64
}

// Gt sets the exclusive lower bound.
func (q *RangeQuery) Gt(v any) *RangeQuery { q.gt = v; return q }

// Gte sets the inclusive lower bound.
func (q *RangeQuery) Gte(v any) *RangeQuery { q.gte = v; return q }

// Lt sets the exclusive upper bound.
func (q *RangeQuery) Lt(v any) *RangeQuery { q.lt = v; return q }

// Lte sets the inclusive upper bound.
func (q *RangeQuery) Lte(v any) *RangeQuery { q.lte = v; return q }

// Format sets the date format for parsing string values.
func (q *RangeQuery) Format(f string) *RangeQuery { q.format = &f; return q }

// Boost sets the relevance score multiplier.
func (q *RangeQuery) Boost(b float64) *RangeQuery { q.boost = &b; return q }

// GtValue returns the gt value if set.
func (q *RangeQuery) GtValue() any { return q.gt }

// GteValue returns the gte value if set.
func (q *RangeQuery) GteValue() any { return q.gte }

// LtValue returns the lt value if set.
func (q *RangeQuery) LtValue() any { return q.lt }

// LteValue returns the lte value if set.
func (q *RangeQuery) LteValue() any { return q.lte }

// FormatValue returns the format value if set.
func (q *RangeQuery) FormatValue() *string { return q.format }

// BoostValue returns the boost value if set.
func (q *RangeQuery) BoostValue() *float64 { return q.boost }

// Range creates a range query builder.
func (b *Builder[T]) Range(field string) *RangeQuery {
	spec, errQ := b.validateField(OpRange, field)
	if errQ != nil {
		return &RangeQuery{query: *errQ}
	}
	return &RangeQuery{
		query: query{op: OpRange, field: spec.Name},
	}
}

// ExistsQuery matches documents where the field exists.
type ExistsQuery struct {
	query
}

// Exists creates an exists query.
func (b *Builder[T]) Exists(field string) *ExistsQuery {
	spec, errQ := b.validateField(OpExists, field)
	if errQ != nil {
		return &ExistsQuery{query: *errQ}
	}
	return &ExistsQuery{
		query: query{op: OpExists, field: spec.Name},
	}
}

// IDsQuery matches documents by their IDs.
type IDsQuery struct {
	query
	ids []string
}

// IDs creates an IDs query.
func (b *Builder[T]) IDs(ids ...string) *IDsQuery {
	return &IDsQuery{
		query: query{op: OpIDs},
		ids:   ids,
	}
}

// IDValues returns the document IDs.
func (q *IDsQuery) IDValues() []string { return q.ids }

// PrefixQuery matches documents with a field prefix.
type PrefixQuery struct {
	query
	rewrite         *string
	caseInsensitive *bool
	boost           *float64
}

// Rewrite sets the rewrite method.
func (q *PrefixQuery) Rewrite(r string) *PrefixQuery { q.rewrite = &r; return q }

// CaseInsensitive enables case-insensitive matching.
func (q *PrefixQuery) CaseInsensitive(b bool) *PrefixQuery { q.caseInsensitive = &b; return q }

// Boost sets the relevance score multiplier.
func (q *PrefixQuery) Boost(b float64) *PrefixQuery { q.boost = &b; return q }

// RewriteValue returns the rewrite value if set.
func (q *PrefixQuery) RewriteValue() *string { return q.rewrite }

// CaseInsensitiveValue returns the case_insensitive value if set.
func (q *PrefixQuery) CaseInsensitiveValue() *bool { return q.caseInsensitive }

// BoostValue returns the boost value if set.
func (q *PrefixQuery) BoostValue() *float64 { return q.boost }

// Prefix creates a prefix query.
func (b *Builder[T]) Prefix(field string, prefix string) *PrefixQuery {
	spec, errQ := b.validateField(OpPrefix, field)
	if errQ != nil {
		return &PrefixQuery{query: *errQ}
	}
	return &PrefixQuery{
		query: query{op: OpPrefix, field: spec.Name, value: prefix},
	}
}

// WildcardQuery matches documents using wildcard patterns.
type WildcardQuery struct {
	query
	rewrite         *string
	caseInsensitive *bool
	boost           *float64
}

// Rewrite sets the rewrite method.
func (q *WildcardQuery) Rewrite(r string) *WildcardQuery { q.rewrite = &r; return q }

// CaseInsensitive enables case-insensitive matching.
func (q *WildcardQuery) CaseInsensitive(b bool) *WildcardQuery { q.caseInsensitive = &b; return q }

// Boost sets the relevance score multiplier.
func (q *WildcardQuery) Boost(b float64) *WildcardQuery { q.boost = &b; return q }

// RewriteValue returns the rewrite value if set.
func (q *WildcardQuery) RewriteValue() *string { return q.rewrite }

// CaseInsensitiveValue returns the case_insensitive value if set.
func (q *WildcardQuery) CaseInsensitiveValue() *bool { return q.caseInsensitive }

// BoostValue returns the boost value if set.
func (q *WildcardQuery) BoostValue() *float64 { return q.boost }

// Wildcard creates a wildcard query.
func (b *Builder[T]) Wildcard(field string, pattern string) *WildcardQuery {
	spec, errQ := b.validateField(OpWildcard, field)
	if errQ != nil {
		return &WildcardQuery{query: *errQ}
	}
	return &WildcardQuery{
		query: query{op: OpWildcard, field: spec.Name, value: pattern},
	}
}

// RegexpQuery matches documents using regular expressions.
type RegexpQuery struct {
	query
	flags           *string
	rewrite         *string
	caseInsensitive *bool
	boost           *float64
}

// Flags sets the regex flags.
func (q *RegexpQuery) Flags(f string) *RegexpQuery { q.flags = &f; return q }

// Rewrite sets the rewrite method.
func (q *RegexpQuery) Rewrite(r string) *RegexpQuery { q.rewrite = &r; return q }

// CaseInsensitive enables case-insensitive matching.
func (q *RegexpQuery) CaseInsensitive(b bool) *RegexpQuery { q.caseInsensitive = &b; return q }

// Boost sets the relevance score multiplier.
func (q *RegexpQuery) Boost(b float64) *RegexpQuery { q.boost = &b; return q }

// FlagsValue returns the flags value if set.
func (q *RegexpQuery) FlagsValue() *string { return q.flags }

// RewriteValue returns the rewrite value if set.
func (q *RegexpQuery) RewriteValue() *string { return q.rewrite }

// CaseInsensitiveValue returns the case_insensitive value if set.
func (q *RegexpQuery) CaseInsensitiveValue() *bool { return q.caseInsensitive }

// BoostValue returns the boost value if set.
func (q *RegexpQuery) BoostValue() *float64 { return q.boost }

// Regexp creates a regexp query.
func (b *Builder[T]) Regexp(field string, pattern string) *RegexpQuery {
	spec, errQ := b.validateField(OpRegexp, field)
	if errQ != nil {
		return &RegexpQuery{query: *errQ}
	}
	return &RegexpQuery{
		query: query{op: OpRegexp, field: spec.Name, value: pattern},
	}
}

// FuzzyQuery matches documents using edit distance.
type FuzzyQuery struct {
	query
	fuzziness    *string
	prefixLength *int
	maxExpansions *int
	transpositions *bool
	rewrite      *string
	boost        *float64
}

// Fuzziness sets the maximum edit distance.
func (q *FuzzyQuery) Fuzziness(f string) *FuzzyQuery { q.fuzziness = &f; return q }

// PrefixLength sets the number of initial characters that must match exactly.
func (q *FuzzyQuery) PrefixLength(p int) *FuzzyQuery { q.prefixLength = &p; return q }

// MaxExpansions sets the maximum number of terms to match.
func (q *FuzzyQuery) MaxExpansions(m int) *FuzzyQuery { q.maxExpansions = &m; return q }

// Transpositions enables or disables transpositions.
func (q *FuzzyQuery) Transpositions(t bool) *FuzzyQuery { q.transpositions = &t; return q }

// Rewrite sets the rewrite method.
func (q *FuzzyQuery) Rewrite(r string) *FuzzyQuery { q.rewrite = &r; return q }

// Boost sets the relevance score multiplier.
func (q *FuzzyQuery) Boost(b float64) *FuzzyQuery { q.boost = &b; return q }

// FuzzinessValue returns the fuzziness value if set.
func (q *FuzzyQuery) FuzzinessValue() *string { return q.fuzziness }

// PrefixLengthValue returns the prefix_length value if set.
func (q *FuzzyQuery) PrefixLengthValue() *int { return q.prefixLength }

// MaxExpansionsValue returns the max_expansions value if set.
func (q *FuzzyQuery) MaxExpansionsValue() *int { return q.maxExpansions }

// TranspositionsValue returns the transpositions value if set.
func (q *FuzzyQuery) TranspositionsValue() *bool { return q.transpositions }

// RewriteValue returns the rewrite value if set.
func (q *FuzzyQuery) RewriteValue() *string { return q.rewrite }

// BoostValue returns the boost value if set.
func (q *FuzzyQuery) BoostValue() *float64 { return q.boost }

// Fuzzy creates a fuzzy query.
func (b *Builder[T]) Fuzzy(field string, value string) *FuzzyQuery {
	spec, errQ := b.validateField(OpFuzzy, field)
	if errQ != nil {
		return &FuzzyQuery{query: *errQ}
	}
	return &FuzzyQuery{
		query: query{op: OpFuzzy, field: spec.Name, value: value},
	}
}
