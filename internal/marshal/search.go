package marshal

import "encoding/json"

// SearchRequest is the typed structure for search request JSON.
type SearchRequest struct {
	Query          any            `json:"query,omitempty"`
	Aggs           map[string]any `json:"aggs,omitempty"`
	Size           *int           `json:"size,omitempty"`
	From           *int           `json:"from,omitempty"`
	Sort           []SortEntry    `json:"sort,omitempty"`
	Source         any            `json:"_source,omitempty"` // []string or SourceFilter
	Highlight      *Highlight     `json:"highlight,omitempty"`
	TrackTotalHits any            `json:"track_total_hits,omitempty"` // bool or int
	MinScore       *float64       `json:"min_score,omitempty"`
	Timeout        *string        `json:"timeout,omitempty"`
}

// SourceFilter specifies include/exclude patterns for _source filtering.
type SourceFilter struct {
	Includes []string `json:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
}

// SortEntry represents a sort specification.
type SortEntry struct {
	Field string
	Order string // "" for default, "asc", or "desc"
}

// MarshalJSON implements json.Marshaler for SortEntry.
func (s SortEntry) MarshalJSON() ([]byte, error) {
	if s.Order == "" {
		return json.Marshal(s.Field)
	}
	return json.Marshal(map[string]map[string]string{
		s.Field: {"order": s.Order},
	})
}

// Highlight is the typed structure for highlight configuration.
type Highlight struct {
	PreTags      []string                  `json:"pre_tags,omitempty"`
	PostTags     []string                  `json:"post_tags,omitempty"`
	Encoder      *string                   `json:"encoder,omitempty"`
	FragmentSize *int                      `json:"fragment_size,omitempty"`
	NumFragments *int                      `json:"number_of_fragments,omitempty"`
	Order        *string                   `json:"order,omitempty"`
	Type         *string                   `json:"type,omitempty"`
	Fields       map[string]HighlightField `json:"fields,omitempty"`
}

// HighlightField is the typed structure for per-field highlight configuration.
type HighlightField struct {
	FragmentSize      *int     `json:"fragment_size,omitempty"`
	NumFragments      *int     `json:"number_of_fragments,omitempty"`
	PreTags           []string `json:"pre_tags,omitempty"`
	PostTags          []string `json:"post_tags,omitempty"`
	MatchedFields     []string `json:"matched_fields,omitempty"`
	FragmentOffset    *int     `json:"fragment_offset,omitempty"`
	NoMatchSize       *int     `json:"no_match_size,omitempty"`
	RequireFieldMatch *bool    `json:"require_field_match,omitempty"`
	HighlightQuery    any      `json:"highlight_query,omitempty"`
}
