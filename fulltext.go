package lucene

// MatchQuery performs analyzed text search.
type MatchQuery struct {
	query
	fuzziness *string
	operator  *string
	analyzer  *string
	boost     *float64
}

// Fuzziness sets the edit distance for fuzzy matching ("AUTO", "0", "1", "2").
func (q *MatchQuery) Fuzziness(f string) *MatchQuery { q.fuzziness = &f; return q }

// Operator sets the boolean operator for terms ("and", "or").
func (q *MatchQuery) Operator(o string) *MatchQuery { q.operator = &o; return q }

// Analyzer sets the analyzer to use for the query.
func (q *MatchQuery) Analyzer(a string) *MatchQuery { q.analyzer = &a; return q }

// Boost sets the relevance score multiplier.
func (q *MatchQuery) Boost(b float64) *MatchQuery { q.boost = &b; return q }

// FuzzinessValue returns the fuzziness value if set.
func (q *MatchQuery) FuzzinessValue() *string { return q.fuzziness }

// OperatorValue returns the operator value if set.
func (q *MatchQuery) OperatorValue() *string { return q.operator }

// AnalyzerValue returns the analyzer value if set.
func (q *MatchQuery) AnalyzerValue() *string { return q.analyzer }

// BoostValue returns the boost value if set.
func (q *MatchQuery) BoostValue() *float64 { return q.boost }

// Match creates a match query for analyzed text search.
func (b *Builder[T]) Match(field string, text string) *MatchQuery {
	spec, errQ := b.validateField(OpMatch, field)
	if errQ != nil {
		return &MatchQuery{query: *errQ}
	}
	return &MatchQuery{
		query: query{op: OpMatch, field: spec.Name, value: text},
	}
}

// MatchPhraseQuery matches an exact phrase.
type MatchPhraseQuery struct {
	query
	slop     *int
	analyzer *string
	boost    *float64
}

// Slop sets the number of positions allowed between terms.
func (q *MatchPhraseQuery) Slop(s int) *MatchPhraseQuery { q.slop = &s; return q }

// Analyzer sets the analyzer to use for the query.
func (q *MatchPhraseQuery) Analyzer(a string) *MatchPhraseQuery { q.analyzer = &a; return q }

// Boost sets the relevance score multiplier.
func (q *MatchPhraseQuery) Boost(b float64) *MatchPhraseQuery { q.boost = &b; return q }

// SlopValue returns the slop value if set.
func (q *MatchPhraseQuery) SlopValue() *int { return q.slop }

// AnalyzerValue returns the analyzer value if set.
func (q *MatchPhraseQuery) AnalyzerValue() *string { return q.analyzer }

// BoostValue returns the boost value if set.
func (q *MatchPhraseQuery) BoostValue() *float64 { return q.boost }

// MatchPhrase creates a match phrase query for exact phrase matching.
func (b *Builder[T]) MatchPhrase(field string, phrase string) *MatchPhraseQuery {
	spec, errQ := b.validateField(OpMatchPhrase, field)
	if errQ != nil {
		return &MatchPhraseQuery{query: *errQ}
	}
	return &MatchPhraseQuery{
		query: query{op: OpMatchPhrase, field: spec.Name, value: phrase},
	}
}

// MatchPhrasePrefixQuery matches a phrase prefix for autocomplete.
type MatchPhrasePrefixQuery struct {
	query
	slop          *int
	maxExpansions *int
	analyzer      *string
	boost         *float64
}

// Slop sets the number of positions allowed between terms.
func (q *MatchPhrasePrefixQuery) Slop(s int) *MatchPhrasePrefixQuery { q.slop = &s; return q }

// MaxExpansions sets the maximum number of terms to match.
func (q *MatchPhrasePrefixQuery) MaxExpansions(m int) *MatchPhrasePrefixQuery {
	q.maxExpansions = &m
	return q
}

// Analyzer sets the analyzer to use for the query.
func (q *MatchPhrasePrefixQuery) Analyzer(a string) *MatchPhrasePrefixQuery {
	q.analyzer = &a
	return q
}

// Boost sets the relevance score multiplier.
func (q *MatchPhrasePrefixQuery) Boost(b float64) *MatchPhrasePrefixQuery { q.boost = &b; return q }

// SlopValue returns the slop value if set.
func (q *MatchPhrasePrefixQuery) SlopValue() *int { return q.slop }

// MaxExpansionsValue returns the max_expansions value if set.
func (q *MatchPhrasePrefixQuery) MaxExpansionsValue() *int { return q.maxExpansions }

// AnalyzerValue returns the analyzer value if set.
func (q *MatchPhrasePrefixQuery) AnalyzerValue() *string { return q.analyzer }

// BoostValue returns the boost value if set.
func (q *MatchPhrasePrefixQuery) BoostValue() *float64 { return q.boost }

// MatchPhrasePrefix creates a match phrase prefix query for autocomplete.
func (b *Builder[T]) MatchPhrasePrefix(field string, phrase string) *MatchPhrasePrefixQuery {
	spec, errQ := b.validateField(OpMatchPhrasePrefix, field)
	if errQ != nil {
		return &MatchPhrasePrefixQuery{query: *errQ}
	}
	return &MatchPhrasePrefixQuery{
		query: query{op: OpMatchPhrasePrefix, field: spec.Name, value: phrase},
	}
}

// MultiMatchQuery searches across multiple fields.
type MultiMatchQuery struct {
	query
	fields     []string
	queryType  *string
	tieBreaker *float64
	fuzziness  *string
	operator   *string
	analyzer   *string
	boost      *float64
}

// Type sets the multi-match type ("best_fields", "most_fields", "cross_fields", "phrase").
func (q *MultiMatchQuery) Type(t string) *MultiMatchQuery { q.queryType = &t; return q }

// TieBreaker sets the tie breaker for best_fields and most_fields types.
func (q *MultiMatchQuery) TieBreaker(t float64) *MultiMatchQuery { q.tieBreaker = &t; return q }

// Fuzziness sets the edit distance for fuzzy matching.
func (q *MultiMatchQuery) Fuzziness(f string) *MultiMatchQuery { q.fuzziness = &f; return q }

// Operator sets the boolean operator for terms.
func (q *MultiMatchQuery) Operator(o string) *MultiMatchQuery { q.operator = &o; return q }

// Analyzer sets the analyzer to use for the query.
func (q *MultiMatchQuery) Analyzer(a string) *MultiMatchQuery { q.analyzer = &a; return q }

// Boost sets the relevance score multiplier.
func (q *MultiMatchQuery) Boost(b float64) *MultiMatchQuery { q.boost = &b; return q }

// Fields returns the fields to search.
func (q *MultiMatchQuery) Fields() []string { return q.fields }

// TypeValue returns the multi-match type if set.
func (q *MultiMatchQuery) TypeValue() *string { return q.queryType }

// TieBreakerValue returns the tie breaker value if set.
func (q *MultiMatchQuery) TieBreakerValue() *float64 { return q.tieBreaker }

// FuzzinessValue returns the fuzziness value if set.
func (q *MultiMatchQuery) FuzzinessValue() *string { return q.fuzziness }

// OperatorValue returns the operator value if set.
func (q *MultiMatchQuery) OperatorValue() *string { return q.operator }

// AnalyzerValue returns the analyzer value if set.
func (q *MultiMatchQuery) AnalyzerValue() *string { return q.analyzer }

// BoostValue returns the boost value if set.
func (q *MultiMatchQuery) BoostValue() *float64 { return q.boost }

// QueryStringQuery parses Lucene query syntax.
type QueryStringQuery struct {
	query
	defaultField    *string
	defaultOperator *string
	analyzer        *string
	allowWildcard   *bool
	fuzziness       *string
	boost           *float64
}

// DefaultField sets the default field for terms without a field prefix.
func (q *QueryStringQuery) DefaultField(f string) *QueryStringQuery {
	q.defaultField = &f
	return q
}

// DefaultOperator sets the default operator ("AND" or "OR").
func (q *QueryStringQuery) DefaultOperator(o string) *QueryStringQuery {
	q.defaultOperator = &o
	return q
}

// Analyzer sets the analyzer to use.
func (q *QueryStringQuery) Analyzer(a string) *QueryStringQuery { q.analyzer = &a; return q }

// AllowLeadingWildcard enables or disables leading wildcards.
func (q *QueryStringQuery) AllowLeadingWildcard(b bool) *QueryStringQuery {
	q.allowWildcard = &b
	return q
}

// Fuzziness sets the default fuzziness.
func (q *QueryStringQuery) Fuzziness(f string) *QueryStringQuery { q.fuzziness = &f; return q }

// Boost sets the relevance score multiplier.
func (q *QueryStringQuery) Boost(b float64) *QueryStringQuery { q.boost = &b; return q }

// DefaultFieldValue returns the default_field value if set.
func (q *QueryStringQuery) DefaultFieldValue() *string { return q.defaultField }

// DefaultOperatorValue returns the default_operator value if set.
func (q *QueryStringQuery) DefaultOperatorValue() *string { return q.defaultOperator }

// AnalyzerValue returns the analyzer value if set.
func (q *QueryStringQuery) AnalyzerValue() *string { return q.analyzer }

// AllowLeadingWildcardValue returns the allow_leading_wildcard value if set.
func (q *QueryStringQuery) AllowLeadingWildcardValue() *bool { return q.allowWildcard }

// FuzzinessValue returns the fuzziness value if set.
func (q *QueryStringQuery) FuzzinessValue() *string { return q.fuzziness }

// BoostValue returns the boost value if set.
func (q *QueryStringQuery) BoostValue() *float64 { return q.boost }

// QueryString creates a query_string query.
func (b *Builder[T]) QueryString(queryStr string) *QueryStringQuery {
	return &QueryStringQuery{
		query: query{op: OpQueryString, value: queryStr},
	}
}

// SimpleQueryStringQuery parses user-friendly query syntax.
type SimpleQueryStringQuery struct {
	query
	fields          []string
	defaultOperator *string
	analyzer        *string
	flags           *string
	boost           *float64
}

// Fields sets the fields to search.
func (q *SimpleQueryStringQuery) Fields(fields ...string) *SimpleQueryStringQuery {
	q.fields = fields
	return q
}

// DefaultOperator sets the default operator ("AND" or "OR").
func (q *SimpleQueryStringQuery) DefaultOperator(o string) *SimpleQueryStringQuery {
	q.defaultOperator = &o
	return q
}

// Analyzer sets the analyzer to use.
func (q *SimpleQueryStringQuery) Analyzer(a string) *SimpleQueryStringQuery {
	q.analyzer = &a
	return q
}

// Flags sets the enabled query features.
func (q *SimpleQueryStringQuery) Flags(f string) *SimpleQueryStringQuery { q.flags = &f; return q }

// Boost sets the relevance score multiplier.
func (q *SimpleQueryStringQuery) Boost(b float64) *SimpleQueryStringQuery { q.boost = &b; return q }

// FieldsValue returns the fields value.
func (q *SimpleQueryStringQuery) FieldsValue() []string { return q.fields }

// DefaultOperatorValue returns the default_operator value if set.
func (q *SimpleQueryStringQuery) DefaultOperatorValue() *string { return q.defaultOperator }

// AnalyzerValue returns the analyzer value if set.
func (q *SimpleQueryStringQuery) AnalyzerValue() *string { return q.analyzer }

// FlagsValue returns the flags value if set.
func (q *SimpleQueryStringQuery) FlagsValue() *string { return q.flags }

// BoostValue returns the boost value if set.
func (q *SimpleQueryStringQuery) BoostValue() *float64 { return q.boost }

// SimpleQueryString creates a simple_query_string query.
func (b *Builder[T]) SimpleQueryString(queryStr string) *SimpleQueryStringQuery {
	return &SimpleQueryStringQuery{
		query: query{op: OpSimpleQueryString, value: queryStr},
	}
}

// MultiMatch creates a multi-match query for searching across multiple fields.
// Fields are validated; if any field is invalid, the query carries an error.
func (b *Builder[T]) MultiMatch(text string, fields ...string) *MultiMatchQuery {
	validFields := make([]string, 0, len(fields))
	for _, field := range fields {
		spec, errQ := b.validateField(OpMultiMatch, field)
		if errQ != nil {
			return &MultiMatchQuery{query: *errQ}
		}
		validFields = append(validFields, spec.Name)
	}
	return &MultiMatchQuery{
		query:  query{op: OpMultiMatch, value: text},
		fields: validFields,
	}
}
