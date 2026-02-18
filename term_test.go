package lucene

import "testing"

type termTestDoc struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Price  float64  `json:"price"`
	Tags   []string `json:"tags"`
}

func TestTermQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("valid field", func(t *testing.T) {
		q := b.Term("status", "active")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Op() != OpTerm {
			t.Errorf("Op() = %v, want %v", q.Op(), OpTerm)
		}
		if q.Field() != "status" {
			t.Errorf("Field() = %v, want status", q.Field())
		}
		if q.Value() != "active" {
			t.Errorf("Value() = %v, want active", q.Value())
		}
	})

	t.Run("with boost", func(t *testing.T) {
		q := b.Term("status", "active").Boost(1.5)
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Term("invalid", "value")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestTermsQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("valid field", func(t *testing.T) {
		q := b.Terms("status", "active", "pending", "complete")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if len(q.Values()) != 3 {
			t.Errorf("len(Values()) = %d, want 3", len(q.Values()))
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Terms("invalid", "a", "b")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestRangeQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("valid field with bounds", func(t *testing.T) {
		q := b.Range("price").Gte(10).Lt(100)
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.GteValue() != 10 {
			t.Errorf("GteValue() = %v, want 10", q.GteValue())
		}
		if q.LtValue() != 100 {
			t.Errorf("LtValue() = %v, want 100", q.LtValue())
		}
	})

	t.Run("with format", func(t *testing.T) {
		q := b.Range("price").Format("epoch_millis")
		if q.FormatValue() == nil || *q.FormatValue() != "epoch_millis" {
			t.Errorf("FormatValue() = %v, want epoch_millis", q.FormatValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Range("invalid").Gte(10)
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestExistsQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("valid field", func(t *testing.T) {
		q := b.Exists("status")
		if q.Err() != nil {
			t.Errorf("Err() = %v, want nil", q.Err())
		}
		if q.Field() != "status" {
			t.Errorf("Field() = %v, want status", q.Field())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Exists("invalid")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestIDsQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	q := b.IDs("doc1", "doc2", "doc3")
	if q.Err() != nil {
		t.Errorf("Err() = %v, want nil", q.Err())
	}
	if len(q.IDValues()) != 3 {
		t.Errorf("len(IDValues()) = %d, want 3", len(q.IDValues()))
	}
}

func TestPrefixQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic prefix", func(t *testing.T) {
		q := b.Prefix("status", "act")
		if q.Op() != OpPrefix {
			t.Errorf("Op() = %v, want OpPrefix", q.Op())
		}
		if q.Field() != "status" {
			t.Errorf("Field() = %v, want status", q.Field())
		}
		if q.Value() != "act" {
			t.Errorf("Value() = %v, want act", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.Prefix("status", "act").
			Rewrite("constant_score").
			CaseInsensitive(true).
			Boost(1.5)

		if q.RewriteValue() == nil || *q.RewriteValue() != "constant_score" {
			t.Errorf("RewriteValue() = %v, want constant_score", q.RewriteValue())
		}
		if q.CaseInsensitiveValue() == nil || *q.CaseInsensitiveValue() != true {
			t.Errorf("CaseInsensitiveValue() = %v, want true", q.CaseInsensitiveValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Prefix("invalid", "test")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestWildcardQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic wildcard", func(t *testing.T) {
		q := b.Wildcard("status", "act*")
		if q.Op() != OpWildcard {
			t.Errorf("Op() = %v, want OpWildcard", q.Op())
		}
		if q.Value() != "act*" {
			t.Errorf("Value() = %v, want act*", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.Wildcard("status", "act?ve").
			Rewrite("constant_score").
			CaseInsensitive(true).
			Boost(1.5)

		if q.RewriteValue() == nil || *q.RewriteValue() != "constant_score" {
			t.Errorf("RewriteValue() = %v, want constant_score", q.RewriteValue())
		}
		if q.CaseInsensitiveValue() == nil || *q.CaseInsensitiveValue() != true {
			t.Errorf("CaseInsensitiveValue() = %v, want true", q.CaseInsensitiveValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Wildcard("invalid", "test*")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestRegexpQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic regexp", func(t *testing.T) {
		q := b.Regexp("status", "[a-z]+")
		if q.Op() != OpRegexp {
			t.Errorf("Op() = %v, want OpRegexp", q.Op())
		}
		if q.Value() != "[a-z]+" {
			t.Errorf("Value() = %v, want [a-z]+", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.Regexp("status", "[A-Z]{3}[0-9]+").
			Flags("ALL").
			Rewrite("constant_score").
			CaseInsensitive(true).
			Boost(1.5)

		if q.FlagsValue() == nil || *q.FlagsValue() != "ALL" {
			t.Errorf("FlagsValue() = %v, want ALL", q.FlagsValue())
		}
		if q.RewriteValue() == nil || *q.RewriteValue() != "constant_score" {
			t.Errorf("RewriteValue() = %v, want constant_score", q.RewriteValue())
		}
		if q.CaseInsensitiveValue() == nil || *q.CaseInsensitiveValue() != true {
			t.Errorf("CaseInsensitiveValue() = %v, want true", q.CaseInsensitiveValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Regexp("invalid", ".*")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestFuzzyQuery(t *testing.T) {
	b, err := New[termTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic fuzzy", func(t *testing.T) {
		q := b.Fuzzy("status", "actve")
		if q.Op() != OpFuzzy {
			t.Errorf("Op() = %v, want OpFuzzy", q.Op())
		}
		if q.Value() != "actve" {
			t.Errorf("Value() = %v, want actve", q.Value())
		}
	})

	t.Run("with options", func(t *testing.T) {
		q := b.Fuzzy("status", "actve").
			Fuzziness("AUTO").
			PrefixLength(2).
			MaxExpansions(50).
			Transpositions(true).
			Rewrite("constant_score").
			Boost(1.5)

		if q.FuzzinessValue() == nil || *q.FuzzinessValue() != "AUTO" {
			t.Errorf("FuzzinessValue() = %v, want AUTO", q.FuzzinessValue())
		}
		if q.PrefixLengthValue() == nil || *q.PrefixLengthValue() != 2 {
			t.Errorf("PrefixLengthValue() = %v, want 2", q.PrefixLengthValue())
		}
		if q.MaxExpansionsValue() == nil || *q.MaxExpansionsValue() != 50 {
			t.Errorf("MaxExpansionsValue() = %v, want 50", q.MaxExpansionsValue())
		}
		if q.TranspositionsValue() == nil || *q.TranspositionsValue() != true {
			t.Errorf("TranspositionsValue() = %v, want true", q.TranspositionsValue())
		}
		if q.RewriteValue() == nil || *q.RewriteValue() != "constant_score" {
			t.Errorf("RewriteValue() = %v, want constant_score", q.RewriteValue())
		}
		if q.BoostValue() == nil || *q.BoostValue() != 1.5 {
			t.Errorf("BoostValue() = %v, want 1.5", q.BoostValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		q := b.Fuzzy("invalid", "test")
		if q.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}
