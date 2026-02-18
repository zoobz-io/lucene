package lucene

// Renderer converts queries and search requests to JSON.
type Renderer interface {
	// Render converts a complete search request to JSON.
	Render(search *Search) ([]byte, error)

	// RenderQuery converts a single query to JSON.
	RenderQuery(query Query) ([]byte, error)

	// RenderAggs converts aggregations to JSON.
	RenderAggs(aggs []Aggregation) ([]byte, error)
}

// Search represents a complete search request.
type Search struct {
	query          Query
	aggs           []Aggregation
	size           *int
	from           *int
	sort           []SortField
	sourceIncludes []string
	sourceExcludes []string
	highlight      *Highlight
	trackTotalHits any // bool or int
	minScore       *float64
	timeout        *string
}

// NewSearch creates a new search request builder.
func NewSearch() *Search {
	return &Search{}
}

// Query sets the query for the search.
func (s *Search) Query(q Query) *Search {
	s.query = q
	return s
}

// Aggs adds aggregations to the search.
func (s *Search) Aggs(aggs ...Aggregation) *Search {
	s.aggs = append(s.aggs, aggs...)
	return s
}

// Size sets the number of hits to return.
func (s *Search) Size(n int) *Search {
	s.size = &n
	return s
}

// From sets the starting offset for results.
func (s *Search) From(n int) *Search {
	s.from = &n
	return s
}

// Sort adds sort fields to the search.
func (s *Search) Sort(fields ...SortField) *Search {
	s.sort = append(s.sort, fields...)
	return s
}

// Source sets the fields to include in _source.
func (s *Search) Source(fields ...string) *Search {
	s.sourceIncludes = fields
	return s
}

// SourceIncludes sets fields to include in _source.
func (s *Search) SourceIncludes(fields ...string) *Search {
	s.sourceIncludes = fields
	return s
}

// SourceExcludes sets fields to exclude from _source.
func (s *Search) SourceExcludes(fields ...string) *Search {
	s.sourceExcludes = fields
	return s
}

// Highlight sets the highlight configuration.
func (s *Search) Highlight(h *Highlight) *Search {
	s.highlight = h
	return s
}

// TrackTotalHits sets whether to track the total number of hits.
// Pass true for accurate count, false for bounded count, or an int for a threshold.
func (s *Search) TrackTotalHits(v any) *Search {
	s.trackTotalHits = v
	return s
}

// MinScore sets the minimum score threshold.
func (s *Search) MinScore(score float64) *Search {
	s.minScore = &score
	return s
}

// Timeout sets the search timeout.
func (s *Search) Timeout(t string) *Search {
	s.timeout = &t
	return s
}

// QueryValue returns the query.
func (s *Search) QueryValue() Query { return s.query }

// AggsValue returns the aggregations.
func (s *Search) AggsValue() []Aggregation { return s.aggs }

// SizeValue returns the size if set.
func (s *Search) SizeValue() *int { return s.size }

// FromValue returns the from offset if set.
func (s *Search) FromValue() *int { return s.from }

// SortValue returns the sort fields.
func (s *Search) SortValue() []SortField { return s.sort }

// SourceIncludesValue returns the source includes.
func (s *Search) SourceIncludesValue() []string { return s.sourceIncludes }

// SourceExcludesValue returns the source excludes.
func (s *Search) SourceExcludesValue() []string { return s.sourceExcludes }

// HighlightValue returns the highlight configuration.
func (s *Search) HighlightValue() *Highlight { return s.highlight }

// TrackTotalHitsValue returns the track_total_hits value.
func (s *Search) TrackTotalHitsValue() any { return s.trackTotalHits }

// MinScoreValue returns the minimum score if set.
func (s *Search) MinScoreValue() *float64 { return s.minScore }

// TimeoutValue returns the timeout if set.
func (s *Search) TimeoutValue() *string { return s.timeout }

// Err returns any error in the search request.
func (s *Search) Err() error {
	if s.query != nil {
		if err := s.query.Err(); err != nil {
			return err
		}
	}
	for _, agg := range s.aggs {
		if err := agg.Err(); err != nil {
			return err
		}
	}
	if s.highlight != nil {
		if err := s.highlight.Err(); err != nil {
			return err
		}
	}
	return nil
}

// Aggregation is the interface for aggregation types.
type Aggregation interface {
	// Name returns the aggregation name.
	Name() string

	// Type returns the aggregation type.
	Type() AggType

	// Field returns the field name.
	Field() string

	// SubAggs returns sub-aggregations.
	SubAggs() []Aggregation

	// Err returns any error.
	Err() error

	// sealed prevents external implementations.
	sealed()
}

// SortField represents a sort specification.
type SortField struct {
	Field string
	Order string // "asc" or "desc".
}
