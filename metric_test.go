package lucene

import "testing"

type metricTestDoc struct {
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func TestAvgAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic avg", func(t *testing.T) {
		a := b.Avg("avg_price", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Name() != "avg_price" {
			t.Errorf("Name() = %v, want avg_price", a.Name())
		}
		if a.Type() != AggAvg {
			t.Errorf("Type() = %v, want AggAvg", a.Type())
		}
		if a.Field() != "price" {
			t.Errorf("Field() = %v, want price", a.Field())
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.Avg("avg_price", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Avg("avg_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestSumAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic sum", func(t *testing.T) {
		a := b.Sum("total_price", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggSum {
			t.Errorf("Type() = %v, want AggSum", a.Type())
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.Sum("total_price", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Sum("total_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestMinAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic min", func(t *testing.T) {
		a := b.Min("min_price", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggMin {
			t.Errorf("Type() = %v, want AggMin", a.Type())
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.Min("min_price", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Min("min_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestMaxAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic max", func(t *testing.T) {
		a := b.Max("max_price", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggMax {
			t.Errorf("Type() = %v, want AggMax", a.Type())
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.Max("max_price", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Max("max_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestCountAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic count", func(t *testing.T) {
		a := b.Count("count_price", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggCount {
			t.Errorf("Type() = %v, want AggCount", a.Type())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Count("count_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestCardinalityAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic cardinality", func(t *testing.T) {
		a := b.Cardinality("unique_categories", "category")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggCardinality {
			t.Errorf("Type() = %v, want AggCardinality", a.Type())
		}
	})

	t.Run("with precision threshold", func(t *testing.T) {
		a := b.Cardinality("unique_categories", "category").PrecisionThreshold(100)
		if a.PrecisionThresholdValue() == nil || *a.PrecisionThresholdValue() != 100 {
			t.Errorf("PrecisionThresholdValue() = %v, want 100", a.PrecisionThresholdValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Cardinality("unique_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestStatsAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic stats", func(t *testing.T) {
		a := b.Stats("price_stats", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggStats {
			t.Errorf("Type() = %v, want AggStats", a.Type())
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.Stats("price_stats", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Stats("stats_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestExtendedStatsAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic extended_stats", func(t *testing.T) {
		a := b.ExtendedStats("price_extended", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggExtendedStats {
			t.Errorf("Type() = %v, want AggExtendedStats", a.Type())
		}
	})

	t.Run("with sigma", func(t *testing.T) {
		a := b.ExtendedStats("price_extended", "price").Sigma(2.0)
		if a.SigmaValue() == nil || *a.SigmaValue() != 2.0 {
			t.Errorf("SigmaValue() = %v, want 2.0", a.SigmaValue())
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.ExtendedStats("price_extended", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.ExtendedStats("extended_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestPercentilesAgg(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic percentiles", func(t *testing.T) {
		a := b.Percentiles("price_pct", "price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggPercentiles {
			t.Errorf("Type() = %v, want AggPercentiles", a.Type())
		}
	})

	t.Run("with percents", func(t *testing.T) {
		a := b.Percentiles("price_pct", "price").Percents(25, 50, 75, 99)
		if len(a.PercentsValue()) != 4 {
			t.Errorf("len(PercentsValue()) = %d, want 4", len(a.PercentsValue()))
		}
	})

	t.Run("with missing", func(t *testing.T) {
		a := b.Percentiles("price_pct", "price").Missing(0)
		if a.MissingValue() != 0 {
			t.Errorf("MissingValue() = %v, want 0", a.MissingValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.Percentiles("pct_invalid", "invalid")
		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestTopHitsAgg_Options(t *testing.T) {
	b, err := New[metricTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic top_hits", func(t *testing.T) {
		a := b.TopHits("top_products")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Name() != "top_products" {
			t.Errorf("Name() = %v, want top_products", a.Name())
		}
		if a.Type() != AggTopHits {
			t.Errorf("Type() = %v, want AggTopHits", a.Type())
		}
	})

	t.Run("with size and from", func(t *testing.T) {
		a := b.TopHits("top_products").Size(5).From(10)
		if a.SizeValue() == nil || *a.SizeValue() != 5 {
			t.Errorf("SizeValue() = %v, want 5", a.SizeValue())
		}
		if a.FromValue() == nil || *a.FromValue() != 10 {
			t.Errorf("FromValue() = %v, want 10", a.FromValue())
		}
	})

	t.Run("with sort", func(t *testing.T) {
		a := b.TopHits("top_products").Sort("price", "desc").Sort("category", "asc")
		if len(a.SortValue()) != 2 {
			t.Errorf("len(SortValue()) = %d, want 2", len(a.SortValue()))
		}
		if a.SortValue()[0].Field != "price" || a.SortValue()[0].Order != "desc" {
			t.Errorf("SortValue()[0] = %+v, want {price, desc}", a.SortValue()[0])
		}
	})

	t.Run("with source", func(t *testing.T) {
		a := b.TopHits("top_products").Source("category", "price")
		if len(a.SourceValue()) != 2 {
			t.Errorf("len(SourceValue()) = %d, want 2", len(a.SourceValue()))
		}
	})
}
