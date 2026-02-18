package lucene

import "testing"

type aggTestDoc struct {
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

func TestTermsAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.TermsAgg("by_category", "category").Size(10).MinDocCount(1)
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.Name() != "by_category" {
		t.Errorf("Name() = %v, want by_category", a.Name())
	}
	if a.Field() != "category" {
		t.Errorf("Field() = %v, want category", a.Field())
	}
	if a.SizeValue() == nil || *a.SizeValue() != 10 {
		t.Errorf("SizeValue() = %v, want 10", a.SizeValue())
	}
}

func TestTermsAgg_InvalidField(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.TermsAgg("by_invalid", "invalid")
	if a.Err() == nil {
		t.Error("Err() should not be nil for invalid field")
	}
}

func TestHistogramAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.Histogram("price_hist", "price").Interval(10)
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.IntervalValue() == nil || *a.IntervalValue() != 10 {
		t.Errorf("IntervalValue() = %v, want 10", a.IntervalValue())
	}
}

func TestDateHistogramAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.DateHistogram("by_month", "timestamp").CalendarInterval("month")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.CalendarIntervalValue() == nil || *a.CalendarIntervalValue() != "month" {
		t.Errorf("CalendarIntervalValue() = %v, want month", a.CalendarIntervalValue())
	}
}

func TestRangeAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.RangeAgg("price_ranges", "price").
		AddRange(0, 50).
		AddRange(50, 100).
		AddRange(100, nil)

	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if len(a.Ranges()) != 3 {
		t.Errorf("len(Ranges()) = %d, want 3", len(a.Ranges()))
	}
}

func TestNestedAggs(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.TermsAgg("by_category", "category").
		SubAgg(b.Avg("avg_price", "price")).
		SubAgg(b.Max("max_price", "price"))

	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if len(a.SubAggs()) != 2 {
		t.Errorf("len(SubAggs()) = %d, want 2", len(a.SubAggs()))
	}
}

func TestMetricAggs(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name string
		agg  Aggregation
	}{
		{"avg", b.Avg("avg_price", "price")},
		{"sum", b.Sum("sum_price", "price")},
		{"min", b.Min("min_price", "price")},
		{"max", b.Max("max_price", "price")},
		{"count", b.Count("count_price", "price")},
		{"cardinality", b.Cardinality("unique_categories", "category")},
		{"stats", b.Stats("price_stats", "price")},
		{"extended_stats", b.ExtendedStats("price_extended", "price")},
		{"percentiles", b.Percentiles("price_pct", "price")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agg.Err() != nil {
				t.Errorf("Err() = %v, want nil", tt.agg.Err())
			}
		})
	}
}

func TestTopHitsAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.TopHits("top_products").Size(3).Sort("price", "desc")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.SizeValue() == nil || *a.SizeValue() != 3 {
		t.Errorf("SizeValue() = %v, want 3", a.SizeValue())
	}
	if len(a.SortValue()) != 1 {
		t.Errorf("len(SortValue()) = %d, want 1", len(a.SortValue()))
	}
}

func TestPipelineAggs(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name string
		agg  Aggregation
	}{
		{"avg_bucket", b.AvgBucket("avg_monthly", "by_month>avg_price")},
		{"sum_bucket", b.SumBucket("sum_monthly", "by_month>sum_price")},
		{"max_bucket", b.MaxBucket("max_monthly", "by_month>max_price")},
		{"min_bucket", b.MinBucket("min_monthly", "by_month>min_price")},
		{"derivative", b.Derivative("price_change", "by_month>avg_price")},
		{"cumulative_sum", b.CumulativeSum("running_total", "by_month>sum_price")},
		{"moving_avg", b.MovingAvg("smoothed", "by_month>avg_price").Window(5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agg.Err() != nil {
				t.Errorf("Err() = %v, want nil", tt.agg.Err())
			}
		})
	}
}

func TestFilterAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic filter", func(t *testing.T) {
		filter := b.Term("category", "electronics")
		a := b.FilterAgg("electronics_only", filter)

		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Name() != "electronics_only" {
			t.Errorf("Name() = %v, want electronics_only", a.Name())
		}
		if a.Type() != AggFilter {
			t.Errorf("Type() = %v, want AggFilter", a.Type())
		}
		if a.FilterQuery() != filter {
			t.Error("FilterQuery() should return the filter query")
		}
	})

	t.Run("with sub-agg", func(t *testing.T) {
		a := b.FilterAgg("filtered", b.Term("category", "electronics")).
			SubAgg(b.Avg("avg_price", "price"))

		if len(a.SubAggs()) != 1 {
			t.Errorf("len(SubAggs()) = %d, want 1", len(a.SubAggs()))
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		a := b.FilterAgg("filtered", b.Term("invalid", "value"))

		if a.Err() == nil {
			t.Error("Err() should propagate filter query error")
		}
	})
}

func TestDateRangeAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic date_range", func(t *testing.T) {
		a := b.DateRangeAgg("time_ranges", "timestamp").
			AddRange("now-1M/M", "now/M").
			AddRange("now/M", nil)

		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Name() != "time_ranges" {
			t.Errorf("Name() = %v, want time_ranges", a.Name())
		}
		if a.Type() != AggDateRange {
			t.Errorf("Type() = %v, want AggDateRange", a.Type())
		}
		if len(a.Ranges()) != 2 {
			t.Errorf("len(Ranges()) = %d, want 2", len(a.Ranges()))
		}
	})

	t.Run("with keyed ranges", func(t *testing.T) {
		a := b.DateRangeAgg("time_ranges", "timestamp").
			AddKeyedRange("last_month", "now-1M/M", "now/M").
			AddKeyedRange("this_month", "now/M", nil).
			Keyed(true)

		if len(a.Ranges()) != 2 {
			t.Errorf("len(Ranges()) = %d, want 2", len(a.Ranges()))
		}
		if a.Ranges()[0].Key != "last_month" {
			t.Errorf("Ranges()[0].Key = %v, want last_month", a.Ranges()[0].Key)
		}
		if a.KeyedValue() == nil || *a.KeyedValue() != true {
			t.Errorf("KeyedValue() = %v, want true", a.KeyedValue())
		}
	})

	t.Run("with format", func(t *testing.T) {
		a := b.DateRangeAgg("time_ranges", "timestamp").
			Format("yyyy-MM-dd")

		if a.FormatValue() == nil || *a.FormatValue() != "yyyy-MM-dd" {
			t.Errorf("FormatValue() = %v, want yyyy-MM-dd", a.FormatValue())
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		a := b.DateRangeAgg("time_ranges", "invalid")

		if a.Err() == nil {
			t.Error("Err() should not be nil for invalid field")
		}
	})
}

func TestFiltersAgg(t *testing.T) {
	b, err := New[aggTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic filters", func(t *testing.T) {
		a := b.FiltersAgg("status_filters").
			Filter("cheap", b.Range("price").Lt(50)).
			Filter("mid", b.Range("price").Gte(50).Lt(100)).
			Filter("expensive", b.Range("price").Gte(100))

		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Name() != "status_filters" {
			t.Errorf("Name() = %v, want status_filters", a.Name())
		}
		if a.Type() != AggFilters {
			t.Errorf("Type() = %v, want AggFilters", a.Type())
		}
		if len(a.Filters()) != 3 {
			t.Errorf("len(Filters()) = %d, want 3", len(a.Filters()))
		}
	})

	t.Run("with sub-agg", func(t *testing.T) {
		a := b.FiltersAgg("by_status").
			Filter("active", b.Term("category", "electronics")).
			SubAgg(b.Avg("avg_price", "price"))

		if len(a.SubAggs()) != 1 {
			t.Errorf("len(SubAggs()) = %d, want 1", len(a.SubAggs()))
		}
	})

	t.Run("error propagation", func(t *testing.T) {
		a := b.FiltersAgg("by_status").
			Filter("valid", b.Term("category", "electronics")).
			Filter("invalid", b.Term("invalid", "value"))

		if a.Err() == nil {
			t.Error("Err() should propagate filter query error")
		}
	})
}
