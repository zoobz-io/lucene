package lucene

import "testing"

type fulltextTestDoc struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Content     string `json:"content"`
}

func TestMatchQuery(t *testing.T) {
	b, err := New[fulltextTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic match", func(t *testing.T) {
		q := b.Match("title", "search term")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpMatch {
			t.Errorf("Op() = %v, want %v", q.Op(), OpMatch)
		}
		if q.Field() != "title" {
			t.Errorf("Field() = %v, want title", q.Field())
		}
		if q.Value() != "search term" {
			t.Errorf("Value() = %v, want search term", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.Match("title", "search").
			Fuzziness("AUTO").
			Operator("and").
			Analyzer("standard").
			Boost(2.0)

		if q.FuzzinessValue() == nil || *q.FuzzinessValue() != "AUTO" {
			t.Errorf("FuzzinessValue() = %v, want AUTO", q.FuzzinessValue())
		}
		if q.OperatorValue() == nil || *q.OperatorValue() != "and" {
			t.Errorf("OperatorValue() = %v, want and", q.OperatorValue())
		}
		if q.AnalyzerValue() == nil || *q.AnalyzerValue() != "standard" {
			t.Errorf("AnalyzerValue() = %v, want standard", q.AnalyzerValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 2.0 {
			t.Errorf("BoostValue() = %v, want 2.0", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Match("invalid", "search")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestMatchPhraseQuery(t *testing.T) {
	b, err := New[fulltextTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic match phrase", func(t *testing.T) {
		q := b.MatchPhrase("title", "exact phrase")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpMatchPhrase {
			t.Errorf("Op() = %v, want %v", q.Op(), OpMatchPhrase)
		}
	})

	t.Run("with slop", func(t *testing.T) {
		q := b.MatchPhrase("title", "phrase").Slop(2)
		if q.SlopValue() == nil || *q.SlopValue() != 2 {
			t.Errorf("SlopValue() = %v, want 2", q.SlopValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.MatchPhrase("invalid", "phrase")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestMultiMatchQuery(t *testing.T) {
	b, err := New[fulltextTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic multi match", func(t *testing.T) {
		q := b.MultiMatch("search term", "title", "description")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpMultiMatch {
			t.Errorf("Op() = %v, want %v", q.Op(), OpMultiMatch)
		}
		if len(q.Fields()) != 2 {
			t.Errorf("len(Fields()) = %d, want 2", len(q.Fields()))
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.MultiMatch("search", "title", "content").
			Type("best_fields").
			TieBreaker(0.3)

		if q.TypeValue() == nil || *q.TypeValue() != "best_fields" {
			t.Errorf("TypeValue() = %v, want best_fields", q.TypeValue())
		}
		if q.TieBreakerValue() == nil || *q.TieBreakerValue() != 0.3 {
			t.Errorf("TieBreakerValue() = %v, want 0.3", q.TieBreakerValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.MultiMatch("search", "title", "invalid")
		if q.Err() == nil {
			t.Error("Err() should not be nil when one field is invalid")
		}
	})
}

func TestMatchPhrasePrefixQuery(t *testing.T) {
	b, err := New[fulltextTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic match_phrase_prefix", func(t *testing.T) {
		q := b.MatchPhrasePrefix("title", "quick bro")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpMatchPhrasePrefix {
			t.Errorf("Op() = %v, want OpMatchPhrasePrefix", q.Op())
		}
		if q.Field() != "title" {
			t.Errorf("Field() = %v, want title", q.Field())
		}
		if q.Value() != "quick bro" {
			t.Errorf("Value() = %v, want quick bro", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.MatchPhrasePrefix("title", "quick bro").
			Slop(2).
			MaxExpansions(50).
			Analyzer("standard").
			Boost(1.5)

		if q.SlopValue() == nil || *q.SlopValue() != 2 {
			t.Errorf("SlopValue() = %v, want 2", q.SlopValue())
		}
		if q.MaxExpansionsValue() == nil || *q.MaxExpansionsValue() != 50 {
			t.Errorf("MaxExpansionsValue() = %v, want 50", q.MaxExpansionsValue())
		}
		if q.AnalyzerValue() == nil || *q.AnalyzerValue() != "standard" {
			t.Errorf("AnalyzerValue() = %v, want standard", q.AnalyzerValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.MatchPhrasePrefix("invalid", "test")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestQueryStringQuery(t *testing.T) {
	b, err := New[fulltextTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic query_string", func(t *testing.T) {
		q := b.QueryString("title:test AND content:example")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpQueryString {
			t.Errorf("Op() = %v, want OpQueryString", q.Op())
		}
		if q.Value() != "title:test AND content:example" {
			t.Errorf("Value() = %v, want title:test AND content:example", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.QueryString("test query").
			DefaultField("title").
			DefaultOperator("AND").
			Analyzer("standard").
			AllowLeadingWildcard(false).
			Fuzziness("AUTO").
			Boost(1.5)

		if q.DefaultFieldValue() == nil || *q.DefaultFieldValue() != "title" {
			t.Errorf("DefaultFieldValue() = %v, want title", q.DefaultFieldValue())
		}
		if q.DefaultOperatorValue() == nil || *q.DefaultOperatorValue() != "AND" {
			t.Errorf("DefaultOperatorValue() = %v, want AND", q.DefaultOperatorValue())
		}
		if q.AnalyzerValue() == nil || *q.AnalyzerValue() != "standard" {
			t.Errorf("AnalyzerValue() = %v, want standard", q.AnalyzerValue())
		}
		if q.AllowLeadingWildcardValue() == nil || *q.AllowLeadingWildcardValue() != false {
			t.Errorf("AllowLeadingWildcardValue() = %v, want false", q.AllowLeadingWildcardValue())
		}
		if q.FuzzinessValue() == nil || *q.FuzzinessValue() != "AUTO" {
			t.Errorf("FuzzinessValue() = %v, want AUTO", q.FuzzinessValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})
}

func TestSimpleQueryStringQuery(t *testing.T) {
	b, err := New[fulltextTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic simple_query_string", func(t *testing.T) {
		q := b.SimpleQueryString("test + example -exclude")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpSimpleQueryString {
			t.Errorf("Op() = %v, want OpSimpleQueryString", q.Op())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.SimpleQueryString("test query").
			Fields("title", "description").
			DefaultOperator("AND").
			Analyzer("standard").
			Flags("ALL").
			Boost(1.5)

		if len(q.FieldsValue()) != 2 {
			t.Errorf("len(FieldsValue()) = %d, want 2", len(q.FieldsValue()))
		}
		if q.DefaultOperatorValue() == nil || *q.DefaultOperatorValue() != "AND" {
			t.Errorf("DefaultOperatorValue() = %v, want AND", q.DefaultOperatorValue())
		}
		if q.AnalyzerValue() == nil || *q.AnalyzerValue() != "standard" {
			t.Errorf("AnalyzerValue() = %v, want standard", q.AnalyzerValue())
		}
		if q.FlagsValue() == nil || *q.FlagsValue() != "ALL" {
			t.Errorf("FlagsValue() = %v, want ALL", q.FlagsValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})
}
