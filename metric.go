package lucene

// AvgAgg computes the average value.
type AvgAgg struct {
	agg
	missing any
}

// Missing sets the value to use for missing fields.
func (a *AvgAgg) Missing(m any) *AvgAgg { a.missing = m; return a }

// MissingValue returns the missing value if set.
func (a *AvgAgg) MissingValue() any { return a.missing }

// Avg creates an avg aggregation.
func (b *Builder[T]) Avg(name, field string) *AvgAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &AvgAgg{agg: agg{name: name, aggType: AggAvg, err: err}}
	}
	return &AvgAgg{agg: agg{name: name, aggType: AggAvg, field: spec.Name}}
}

// SumAgg computes the sum of values.
type SumAgg struct {
	agg
	missing any
}

// Missing sets the value to use for missing fields.
func (a *SumAgg) Missing(m any) *SumAgg { a.missing = m; return a }

// MissingValue returns the missing value if set.
func (a *SumAgg) MissingValue() any { return a.missing }

// Sum creates a sum aggregation.
func (b *Builder[T]) Sum(name, field string) *SumAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &SumAgg{agg: agg{name: name, aggType: AggSum, err: err}}
	}
	return &SumAgg{agg: agg{name: name, aggType: AggSum, field: spec.Name}}
}

// MinAgg computes the minimum value.
type MinAgg struct {
	agg
	missing any
}

// Missing sets the value to use for missing fields.
func (a *MinAgg) Missing(m any) *MinAgg { a.missing = m; return a }

// MissingValue returns the missing value if set.
func (a *MinAgg) MissingValue() any { return a.missing }

// Min creates a min aggregation.
func (b *Builder[T]) Min(name, field string) *MinAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &MinAgg{agg: agg{name: name, aggType: AggMin, err: err}}
	}
	return &MinAgg{agg: agg{name: name, aggType: AggMin, field: spec.Name}}
}

// MaxAgg computes the maximum value.
type MaxAgg struct {
	agg
	missing any
}

// Missing sets the value to use for missing fields.
func (a *MaxAgg) Missing(m any) *MaxAgg { a.missing = m; return a }

// MissingValue returns the missing value if set.
func (a *MaxAgg) MissingValue() any { return a.missing }

// Max creates a max aggregation.
func (b *Builder[T]) Max(name, field string) *MaxAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &MaxAgg{agg: agg{name: name, aggType: AggMax, err: err}}
	}
	return &MaxAgg{agg: agg{name: name, aggType: AggMax, field: spec.Name}}
}

// CountAgg counts values.
type CountAgg struct {
	agg
}

// Count creates a value_count aggregation.
func (b *Builder[T]) Count(name, field string) *CountAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &CountAgg{agg: agg{name: name, aggType: AggCount, err: err}}
	}
	return &CountAgg{agg: agg{name: name, aggType: AggCount, field: spec.Name}}
}

// CardinalityAgg counts distinct values.
type CardinalityAgg struct {
	agg
	precisionThreshold *int
}

// PrecisionThreshold sets the precision threshold.
func (a *CardinalityAgg) PrecisionThreshold(p int) *CardinalityAgg {
	a.precisionThreshold = &p
	return a
}

// PrecisionThresholdValue returns the precision_threshold if set.
func (a *CardinalityAgg) PrecisionThresholdValue() *int { return a.precisionThreshold }

// Cardinality creates a cardinality aggregation.
func (b *Builder[T]) Cardinality(name, field string) *CardinalityAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &CardinalityAgg{agg: agg{name: name, aggType: AggCardinality, err: err}}
	}
	return &CardinalityAgg{agg: agg{name: name, aggType: AggCardinality, field: spec.Name}}
}

// StatsAgg computes basic statistics.
type StatsAgg struct {
	agg
	missing any
}

// Missing sets the value to use for missing fields.
func (a *StatsAgg) Missing(m any) *StatsAgg { a.missing = m; return a }

// MissingValue returns the missing value if set.
func (a *StatsAgg) MissingValue() any { return a.missing }

// Stats creates a stats aggregation.
func (b *Builder[T]) Stats(name, field string) *StatsAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &StatsAgg{agg: agg{name: name, aggType: AggStats, err: err}}
	}
	return &StatsAgg{agg: agg{name: name, aggType: AggStats, field: spec.Name}}
}

// ExtendedStatsAgg computes extended statistics.
type ExtendedStatsAgg struct {
	agg
	missing any
	sigma   *float64
}

// Missing sets the value to use for missing fields.
func (a *ExtendedStatsAgg) Missing(m any) *ExtendedStatsAgg { a.missing = m; return a }

// Sigma sets the sigma value for bounds.
func (a *ExtendedStatsAgg) Sigma(s float64) *ExtendedStatsAgg { a.sigma = &s; return a }

// MissingValue returns the missing value if set.
func (a *ExtendedStatsAgg) MissingValue() any { return a.missing }

// SigmaValue returns the sigma value if set.
func (a *ExtendedStatsAgg) SigmaValue() *float64 { return a.sigma }

// ExtendedStats creates an extended_stats aggregation.
func (b *Builder[T]) ExtendedStats(name, field string) *ExtendedStatsAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &ExtendedStatsAgg{agg: agg{name: name, aggType: AggExtendedStats, err: err}}
	}
	return &ExtendedStatsAgg{agg: agg{name: name, aggType: AggExtendedStats, field: spec.Name}}
}

// PercentilesAgg computes percentile values.
type PercentilesAgg struct {
	agg
	percents []float64
	missing  any
}

// Percents sets the percentiles to compute.
func (a *PercentilesAgg) Percents(p ...float64) *PercentilesAgg { a.percents = p; return a }

// Missing sets the value to use for missing fields.
func (a *PercentilesAgg) Missing(m any) *PercentilesAgg { a.missing = m; return a }

// PercentsValue returns the percents if set.
func (a *PercentilesAgg) PercentsValue() []float64 { return a.percents }

// MissingValue returns the missing value if set.
func (a *PercentilesAgg) MissingValue() any { return a.missing }

// Percentiles creates a percentiles aggregation.
func (b *Builder[T]) Percentiles(name, field string) *PercentilesAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &PercentilesAgg{agg: agg{name: name, aggType: AggPercentiles, err: err}}
	}
	return &PercentilesAgg{agg: agg{name: name, aggType: AggPercentiles, field: spec.Name}}
}

// TopHitsAgg returns top matching documents.
type TopHitsAgg struct {
	agg
	size   *int
	from   *int
	sort   []SortField
	source []string
}

// Size sets the number of hits to return.
func (a *TopHitsAgg) Size(s int) *TopHitsAgg { a.size = &s; return a }

// From sets the offset.
func (a *TopHitsAgg) From(f int) *TopHitsAgg { a.from = &f; return a }

// Sort adds a sort field.
func (a *TopHitsAgg) Sort(field, order string) *TopHitsAgg {
	a.sort = append(a.sort, SortField{Field: field, Order: order})
	return a
}

// Source sets the fields to return.
func (a *TopHitsAgg) Source(fields ...string) *TopHitsAgg { a.source = fields; return a }

// SizeValue returns the size if set.
func (a *TopHitsAgg) SizeValue() *int { return a.size }

// FromValue returns the from if set.
func (a *TopHitsAgg) FromValue() *int { return a.from }

// SortValue returns the sort fields.
func (a *TopHitsAgg) SortValue() []SortField { return a.sort }

// SourceValue returns the source fields.
func (a *TopHitsAgg) SourceValue() []string { return a.source }

// TopHits creates a top_hits aggregation.
func (b *Builder[T]) TopHits(name string) *TopHitsAgg {
	return &TopHitsAgg{agg: agg{name: name, aggType: AggTopHits}}
}
