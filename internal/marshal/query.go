package marshal

// === Term-level query inner types ===

// TermInner is the inner structure for term queries.
type TermInner struct {
	Value any      `json:"value"`
	Boost *float64 `json:"boost,omitempty"`
}

// TermsInner is the inner structure for terms queries.
// Note: The field name maps to the values slice, boost is at the same level.
type TermsInner struct {
	Values []any    `json:"-"` // Handled specially
	Boost  *float64 `json:"boost,omitempty"`
}

// RangeInner is the inner structure for range queries.
type RangeInner struct {
	Gt     any      `json:"gt,omitempty"`
	Gte    any      `json:"gte,omitempty"`
	Lt     any      `json:"lt,omitempty"`
	Lte    any      `json:"lte,omitempty"`
	Format *string  `json:"format,omitempty"`
	Boost  *float64 `json:"boost,omitempty"`
}

// ExistsInner is the inner structure for exists queries.
type ExistsInner struct {
	Field string `json:"field"`
}

// IDsInner is the inner structure for ids queries.
type IDsInner struct {
	Values []string `json:"values"`
}

// PrefixInner is the inner structure for prefix queries.
type PrefixInner struct {
	Value           string   `json:"value"`
	Rewrite         *string  `json:"rewrite,omitempty"`
	CaseInsensitive *bool    `json:"case_insensitive,omitempty"`
	Boost           *float64 `json:"boost,omitempty"`
}

// WildcardInner is the inner structure for wildcard queries.
type WildcardInner struct {
	Value           string   `json:"value"`
	Rewrite         *string  `json:"rewrite,omitempty"`
	CaseInsensitive *bool    `json:"case_insensitive,omitempty"`
	Boost           *float64 `json:"boost,omitempty"`
}

// RegexpInner is the inner structure for regexp queries.
type RegexpInner struct {
	Value           string   `json:"value"`
	Flags           *string  `json:"flags,omitempty"`
	Rewrite         *string  `json:"rewrite,omitempty"`
	CaseInsensitive *bool    `json:"case_insensitive,omitempty"`
	Boost           *float64 `json:"boost,omitempty"`
}

// FuzzyInner is the inner structure for fuzzy queries.
type FuzzyInner struct {
	Value          string   `json:"value"`
	Fuzziness      *string  `json:"fuzziness,omitempty"`
	PrefixLength   *int     `json:"prefix_length,omitempty"`
	MaxExpansions  *int     `json:"max_expansions,omitempty"`
	Transpositions *bool    `json:"transpositions,omitempty"`
	Rewrite        *string  `json:"rewrite,omitempty"`
	Boost          *float64 `json:"boost,omitempty"`
}

// === Full-text query inner types ===

// MatchInner is the inner structure for match queries.
type MatchInner struct {
	Query     string   `json:"query"`
	Fuzziness *string  `json:"fuzziness,omitempty"`
	Operator  *string  `json:"operator,omitempty"`
	Analyzer  *string  `json:"analyzer,omitempty"`
	Boost     *float64 `json:"boost,omitempty"`
}

// MatchPhraseInner is the inner structure for match_phrase queries.
type MatchPhraseInner struct {
	Query    string   `json:"query"`
	Slop     *int     `json:"slop,omitempty"`
	Analyzer *string  `json:"analyzer,omitempty"`
	Boost    *float64 `json:"boost,omitempty"`
}

// MatchPhrasePrefixInner is the inner structure for match_phrase_prefix queries.
type MatchPhrasePrefixInner struct {
	Query         string   `json:"query"`
	Slop          *int     `json:"slop,omitempty"`
	MaxExpansions *int     `json:"max_expansions,omitempty"`
	Analyzer      *string  `json:"analyzer,omitempty"`
	Boost         *float64 `json:"boost,omitempty"`
}

// MultiMatchInner is the inner structure for multi_match queries.
type MultiMatchInner struct {
	Query      string   `json:"query"`
	Fields     []string `json:"fields"`
	Type       *string  `json:"type,omitempty"`
	TieBreaker *float64 `json:"tie_breaker,omitempty"`
	Fuzziness  *string  `json:"fuzziness,omitempty"`
	Operator   *string  `json:"operator,omitempty"`
	Analyzer   *string  `json:"analyzer,omitempty"`
	Boost      *float64 `json:"boost,omitempty"`
}

// QueryStringInner is the inner structure for query_string queries.
type QueryStringInner struct {
	Query                string   `json:"query"`
	DefaultField         *string  `json:"default_field,omitempty"`
	DefaultOperator      *string  `json:"default_operator,omitempty"`
	Analyzer             *string  `json:"analyzer,omitempty"`
	AllowLeadingWildcard *bool    `json:"allow_leading_wildcard,omitempty"`
	Fuzziness            *string  `json:"fuzziness,omitempty"`
	Boost                *float64 `json:"boost,omitempty"`
}

// SimpleQueryStringInner is the inner structure for simple_query_string queries.
type SimpleQueryStringInner struct {
	Query           string   `json:"query"`
	Fields          []string `json:"fields,omitempty"`
	DefaultOperator *string  `json:"default_operator,omitempty"`
	Analyzer        *string  `json:"analyzer,omitempty"`
	Flags           *string  `json:"flags,omitempty"`
	Boost           *float64 `json:"boost,omitempty"`
}

// === Compound query inner types ===

// BoolInner is the inner structure for bool queries.
type BoolInner struct {
	Must               []any    `json:"must,omitempty"`
	Should             []any    `json:"should,omitempty"`
	MustNot            []any    `json:"must_not,omitempty"`
	Filter             []any    `json:"filter,omitempty"`
	MinimumShouldMatch *int     `json:"minimum_should_match,omitempty"`
	Boost              *float64 `json:"boost,omitempty"`
}

// MatchAllInner is the inner structure for match_all queries.
type MatchAllInner struct {
	Boost *float64 `json:"boost,omitempty"`
}

// MatchNoneInner is the inner structure for match_none queries.
type MatchNoneInner struct{}

// BoostingInner is the inner structure for boosting queries.
type BoostingInner struct {
	Positive      any      `json:"positive,omitempty"`
	Negative      any      `json:"negative,omitempty"`
	NegativeBoost *float64 `json:"negative_boost,omitempty"`
}

// DisMaxInner is the inner structure for dis_max queries.
type DisMaxInner struct {
	Queries    []any    `json:"queries,omitempty"`
	TieBreaker *float64 `json:"tie_breaker,omitempty"`
	Boost      *float64 `json:"boost,omitempty"`
}

// ConstantScoreInner is the inner structure for constant_score queries.
type ConstantScoreInner struct {
	Filter any      `json:"filter,omitempty"`
	Boost  *float64 `json:"boost,omitempty"`
}

// === Nested/Join query inner types ===

// NestedInner is the inner structure for nested queries.
type NestedInner struct {
	Path           string   `json:"path"`
	Query          any      `json:"query,omitempty"`
	ScoreMode      *string  `json:"score_mode,omitempty"`
	IgnoreUnmapped *bool    `json:"ignore_unmapped,omitempty"`
	Boost          *float64 `json:"boost,omitempty"`
}

// HasChildInner is the inner structure for has_child queries.
type HasChildInner struct {
	Type           string   `json:"type"`
	Query          any      `json:"query,omitempty"`
	ScoreMode      *string  `json:"score_mode,omitempty"`
	MinChildren    *int     `json:"min_children,omitempty"`
	MaxChildren    *int     `json:"max_children,omitempty"`
	IgnoreUnmapped *bool    `json:"ignore_unmapped,omitempty"`
	Boost          *float64 `json:"boost,omitempty"`
}

// HasParentInner is the inner structure for has_parent queries.
type HasParentInner struct {
	ParentType     string   `json:"parent_type"`
	Query          any      `json:"query,omitempty"`
	Score          *bool    `json:"score,omitempty"`
	IgnoreUnmapped *bool    `json:"ignore_unmapped,omitempty"`
	Boost          *float64 `json:"boost,omitempty"`
}

// === Geo query inner types ===

// GeoPoint represents a lat/lon coordinate.
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// GeoDistanceInner is the inner structure for geo_distance queries.
// Note: Field name maps to GeoPoint, other fields are at the same level.
type GeoDistanceInner struct {
	Field        string   `json:"-"` // Handled specially
	Point        GeoPoint `json:"-"` // Handled specially
	Distance     *string  `json:"distance,omitempty"`
	DistanceType *string  `json:"distance_type,omitempty"`
	Boost        *float64 `json:"boost,omitempty"`
}

// GeoBoundingBoxCorners holds the corner points for a bounding box.
type GeoBoundingBoxCorners struct {
	TopLeft     *GeoPoint `json:"top_left,omitempty"`
	BottomRight *GeoPoint `json:"bottom_right,omitempty"`
}

// GeoBoundingBoxInner is the inner structure for geo_bounding_box queries.
type GeoBoundingBoxInner struct {
	Field   string                `json:"-"` // Handled specially
	Corners GeoBoundingBoxCorners `json:"-"` // Handled specially
	Boost   *float64              `json:"boost,omitempty"`
}

// === Vector query inner types ===

// KnnInnerES is the inner structure for knn queries (Elasticsearch format).
type KnnInnerES struct {
	Field         string    `json:"field"`
	Vector        []float32 `json:"vector"`
	K             *int      `json:"k,omitempty"`
	NumCandidates *int      `json:"num_candidates,omitempty"`
	Filter        any       `json:"filter,omitempty"`
	Boost         *float64  `json:"boost,omitempty"`
}

// KnnFieldInnerOS is the field-specific structure for OpenSearch knn.
type KnnFieldInnerOS struct {
	Vector        []float32 `json:"vector"`
	K             *int      `json:"k,omitempty"`
	NumCandidates *int      `json:"num_candidates,omitempty"`
	Filter        any       `json:"filter,omitempty"`
	Boost         *float64  `json:"boost,omitempty"`
}

// KnnInnerOS is the inner structure for knn queries (OpenSearch format).
// Note: Field name maps to KnnFieldInnerOS.
type KnnInnerOS struct {
	Field string          `json:"-"` // Handled specially
	Inner KnnFieldInnerOS `json:"-"` // Handled specially
}
