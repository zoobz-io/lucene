package lucene

import (
	"errors"
	"testing"
)

type searchTestDoc struct {
	Title    string  `json:"title"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func TestNewSearch(t *testing.T) {
	s := NewSearch()
	if s == nil {
		t.Fatal("NewSearch() returned nil")
	}
}

func TestSearch_Query(t *testing.T) {
	b, err := New[searchTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	q := b.Match("title", "test")
	s := NewSearch().Query(q)

	if s.QueryValue() != q {
		t.Error("QueryValue() should return the query")
	}
}

func TestSearch_Aggs(t *testing.T) {
	b, err := New[searchTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	agg1 := b.TermsAgg("by_category", "category")
	agg2 := b.Avg("avg_price", "price")
	s := NewSearch().Aggs(agg1, agg2)

	if len(s.AggsValue()) != 2 {
		t.Errorf("AggsValue() len = %d, want 2", len(s.AggsValue()))
	}
}

func TestSearch_SizeFrom(t *testing.T) {
	s := NewSearch().Size(20).From(10)

	if s.SizeValue() == nil || *s.SizeValue() != 20 {
		t.Errorf("SizeValue() = %v, want 20", s.SizeValue())
	}
	if s.FromValue() == nil || *s.FromValue() != 10 {
		t.Errorf("FromValue() = %v, want 10", s.FromValue())
	}
}

func TestSearch_Sort(t *testing.T) {
	s := NewSearch().Sort(
		SortField{Field: "price", Order: "asc"},
		SortField{Field: "title", Order: "desc"},
	)

	sorts := s.SortValue()
	if len(sorts) != 2 {
		t.Fatalf("SortValue() len = %d, want 2", len(sorts))
	}
	if sorts[0].Field != "price" || sorts[0].Order != "asc" {
		t.Errorf("SortValue()[0] = %+v, want {price asc}", sorts[0])
	}
}

func TestSearch_Source(t *testing.T) {
	s := NewSearch().Source("title", "price")

	includes := s.SourceIncludesValue()
	if len(includes) != 2 {
		t.Fatalf("SourceIncludesValue() len = %d, want 2", len(includes))
	}
	if includes[0] != "title" || includes[1] != "price" {
		t.Errorf("SourceIncludesValue() = %v, want [title price]", includes)
	}
}

func TestSearch_SourceIncludesExcludes(t *testing.T) {
	s := NewSearch().
		SourceIncludes("title", "price").
		SourceExcludes("internal_field")

	if len(s.SourceIncludesValue()) != 2 {
		t.Errorf("SourceIncludesValue() len = %d, want 2", len(s.SourceIncludesValue()))
	}
	if len(s.SourceExcludesValue()) != 1 {
		t.Errorf("SourceExcludesValue() len = %d, want 1", len(s.SourceExcludesValue()))
	}
}

func TestSearch_TrackTotalHits(t *testing.T) {
	tests := []struct {
		name  string
		value any
	}{
		{"bool true", true},
		{"bool false", false},
		{"int threshold", 10000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSearch().TrackTotalHits(tt.value)
			if s.TrackTotalHitsValue() != tt.value {
				t.Errorf("TrackTotalHitsValue() = %v, want %v", s.TrackTotalHitsValue(), tt.value)
			}
		})
	}
}

func TestSearch_MinScore(t *testing.T) {
	s := NewSearch().MinScore(0.5)

	if s.MinScoreValue() == nil || *s.MinScoreValue() != 0.5 {
		t.Errorf("MinScoreValue() = %v, want 0.5", s.MinScoreValue())
	}
}

func TestSearch_Timeout(t *testing.T) {
	s := NewSearch().Timeout("10s")

	if s.TimeoutValue() == nil || *s.TimeoutValue() != "10s" {
		t.Errorf("TimeoutValue() = %v, want 10s", s.TimeoutValue())
	}
}

func TestSearch_Highlight(t *testing.T) {
	h := NewHighlight().Fields("title", "description")
	s := NewSearch().Highlight(h)

	if s.HighlightValue() != h {
		t.Error("HighlightValue() should return the highlight")
	}
}

func TestSearch_Err(t *testing.T) {
	b, err := New[searchTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Search with invalid query field
	q := b.Match("invalid_field", "test")
	s := NewSearch().Query(q)

	if s.Err() == nil {
		t.Error("Err() should not be nil for invalid query")
	}
}

func TestSearch_Err_InvalidAgg(t *testing.T) {
	b, err := New[searchTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	agg := b.TermsAgg("by_invalid", "invalid_field")
	s := NewSearch().Aggs(agg)

	if s.Err() == nil {
		t.Error("Err() should not be nil for invalid aggregation")
	}
}

func TestSearch_ChainedBuilder(t *testing.T) {
	b, err := New[searchTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	h := NewHighlight().
		Fields("title").
		PreTags("<em>").
		PostTags("</em>")

	s := NewSearch().
		Query(b.Match("title", "search term")).
		Aggs(b.TermsAgg("by_category", "category")).
		Size(20).
		From(0).
		Sort(SortField{Field: "price", Order: "asc"}).
		Source("title", "price").
		Highlight(h).
		TrackTotalHits(true).
		MinScore(0.1).
		Timeout("5s")

	if s.Err() != nil {
		t.Errorf("Err() = %v, want nil", s.Err())
	}
	if s.QueryValue() == nil {
		t.Error("QueryValue() should not be nil")
	}
	if len(s.AggsValue()) != 1 {
		t.Error("AggsValue() should have 1 aggregation")
	}
	if s.SizeValue() == nil || *s.SizeValue() != 20 {
		t.Error("SizeValue() should be 20")
	}
	if s.FromValue() == nil || *s.FromValue() != 0 {
		t.Error("FromValue() should be 0")
	}
	if len(s.SortValue()) != 1 {
		t.Error("SortValue() should have 1 sort")
	}
	if len(s.SourceIncludesValue()) != 2 {
		t.Error("SourceIncludesValue() should have 2 fields")
	}
	if s.HighlightValue() == nil {
		t.Error("HighlightValue() should not be nil")
	}
	if s.TrackTotalHitsValue() != true {
		t.Error("TrackTotalHitsValue() should be true")
	}
	if s.MinScoreValue() == nil || *s.MinScoreValue() != 0.1 {
		t.Error("MinScoreValue() should be 0.1")
	}
	if s.TimeoutValue() == nil || *s.TimeoutValue() != "5s" {
		t.Error("TimeoutValue() should be 5s")
	}
}

func TestSearch_Err_Highlight(t *testing.T) {
	b, err := New[searchTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	// Create a highlight with invalid query
	invalidQuery := b.Match("invalid_field", "test")
	h := NewHighlight().Field(
		NewHighlightField("title").HighlightQuery(invalidQuery).Build(),
	)
	s := NewSearch().Highlight(h)

	if !errors.Is(s.Err(), ErrUnknownField) {
		t.Error("Err() should return error for highlight with invalid query")
	}
}
