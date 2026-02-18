package lucene

// AggType represents an aggregation type.
type AggType uint8

const (
	// AggTerms groups documents by field value.
	AggTerms AggType = iota
	// AggHistogram creates numeric buckets.
	AggHistogram
	// AggDateHistogram creates time-based buckets.
	AggDateHistogram
	// AggRange creates custom range buckets.
	AggRange
	// AggDateRange creates date range buckets.
	AggDateRange
	// AggFilter creates a single filter bucket.
	AggFilter
	// AggFilters creates named filter buckets.
	AggFilters
	// AggNested aggregates nested documents.
	AggNested
	// AggMissing counts documents missing a field.
	AggMissing

	// AggAvg computes the average value.
	AggAvg
	// AggSum computes the sum of values.
	AggSum
	// AggMin computes the minimum value.
	AggMin
	// AggMax computes the maximum value.
	AggMax
	// AggCount counts values.
	AggCount
	// AggCardinality counts distinct values.
	AggCardinality
	// AggStats computes basic statistics.
	AggStats
	// AggExtendedStats computes extended statistics.
	AggExtendedStats
	// AggPercentiles computes percentile values.
	AggPercentiles
	// AggTopHits returns top matching documents.
	AggTopHits

	// AggAvgBucket computes the average of bucket values.
	AggAvgBucket
	// AggSumBucket computes the sum of bucket values.
	AggSumBucket
	// AggMaxBucket finds the maximum bucket value.
	AggMaxBucket
	// AggMinBucket finds the minimum bucket value.
	AggMinBucket
	// AggDerivative computes the derivative of a metric.
	AggDerivative
	// AggCumulativeSum computes the cumulative sum.
	AggCumulativeSum
	// AggMovingAvg computes a moving average.
	AggMovingAvg
)

// agg is the base struct for all aggregations.
type agg struct {
	name    string
	aggType AggType
	field   string
	subAggs []Aggregation
	err     error
}

func (a *agg) Name() string     { return a.name }
func (a *agg) Type() AggType    { return a.aggType }
func (a *agg) Field() string    { return a.field }
func (a *agg) SubAggs() []Aggregation { return a.subAggs }
func (a *agg) Err() error       { return a.err }
func (a *agg) sealed()          {}

// TermsAgg groups documents by field value.
type TermsAgg struct {
	agg
	size        *int
	minDocCount *int
	order       map[string]string
}

// Size sets the maximum number of buckets to return.
func (a *TermsAgg) Size(s int) *TermsAgg { a.size = &s; return a }

// MinDocCount sets the minimum document count for a bucket.
func (a *TermsAgg) MinDocCount(m int) *TermsAgg { a.minDocCount = &m; return a }

// Order sets the bucket sort order.
func (a *TermsAgg) Order(field, dir string) *TermsAgg {
	if a.order == nil {
		a.order = make(map[string]string)
	}
	a.order[field] = dir
	return a
}

// SubAgg adds a sub-aggregation.
func (a *TermsAgg) SubAgg(sub Aggregation) *TermsAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// SizeValue returns the size if set.
func (a *TermsAgg) SizeValue() *int { return a.size }

// MinDocCountValue returns the min_doc_count if set.
func (a *TermsAgg) MinDocCountValue() *int { return a.minDocCount }

// OrderValue returns the order if set.
func (a *TermsAgg) OrderValue() map[string]string { return a.order }

// TermsAgg creates a terms aggregation.
func (b *Builder[T]) TermsAgg(name, field string) *TermsAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &TermsAgg{agg: agg{name: name, aggType: AggTerms, err: err}}
	}
	return &TermsAgg{agg: agg{name: name, aggType: AggTerms, field: spec.Name}}
}

// HistogramAgg creates numeric buckets.
type HistogramAgg struct {
	agg
	interval    *float64
	offset      *float64
	minDocCount *int
}

// Interval sets the bucket interval.
func (a *HistogramAgg) Interval(i float64) *HistogramAgg { a.interval = &i; return a }

// Offset sets the bucket offset.
func (a *HistogramAgg) Offset(o float64) *HistogramAgg { a.offset = &o; return a }

// MinDocCount sets the minimum document count for a bucket.
func (a *HistogramAgg) MinDocCount(m int) *HistogramAgg { a.minDocCount = &m; return a }

// SubAgg adds a sub-aggregation.
func (a *HistogramAgg) SubAgg(sub Aggregation) *HistogramAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// IntervalValue returns the interval if set.
func (a *HistogramAgg) IntervalValue() *float64 { return a.interval }

// OffsetValue returns the offset if set.
func (a *HistogramAgg) OffsetValue() *float64 { return a.offset }

// MinDocCountValue returns the min_doc_count if set.
func (a *HistogramAgg) MinDocCountValue() *int { return a.minDocCount }

// Histogram creates a histogram aggregation.
func (b *Builder[T]) Histogram(name, field string) *HistogramAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &HistogramAgg{agg: agg{name: name, aggType: AggHistogram, err: err}}
	}
	return &HistogramAgg{agg: agg{name: name, aggType: AggHistogram, field: spec.Name}}
}

// DateHistogramAgg creates time-based buckets.
type DateHistogramAgg struct {
	agg
	calendarInterval *string
	fixedInterval    *string
	format           *string
	timeZone         *string
	minDocCount      *int
}

// CalendarInterval sets the calendar-aware interval (e.g., "month", "week").
func (a *DateHistogramAgg) CalendarInterval(i string) *DateHistogramAgg {
	a.calendarInterval = &i
	return a
}

// FixedInterval sets the fixed interval (e.g., "1d", "12h").
func (a *DateHistogramAgg) FixedInterval(i string) *DateHistogramAgg {
	a.fixedInterval = &i
	return a
}

// Format sets the date format for keys.
func (a *DateHistogramAgg) Format(f string) *DateHistogramAgg { a.format = &f; return a }

// TimeZone sets the time zone.
func (a *DateHistogramAgg) TimeZone(tz string) *DateHistogramAgg { a.timeZone = &tz; return a }

// MinDocCount sets the minimum document count for a bucket.
func (a *DateHistogramAgg) MinDocCount(m int) *DateHistogramAgg { a.minDocCount = &m; return a }

// SubAgg adds a sub-aggregation.
func (a *DateHistogramAgg) SubAgg(sub Aggregation) *DateHistogramAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// CalendarIntervalValue returns the calendar_interval if set.
func (a *DateHistogramAgg) CalendarIntervalValue() *string { return a.calendarInterval }

// FixedIntervalValue returns the fixed_interval if set.
func (a *DateHistogramAgg) FixedIntervalValue() *string { return a.fixedInterval }

// FormatValue returns the format if set.
func (a *DateHistogramAgg) FormatValue() *string { return a.format }

// TimeZoneValue returns the time_zone if set.
func (a *DateHistogramAgg) TimeZoneValue() *string { return a.timeZone }

// MinDocCountValue returns the min_doc_count if set.
func (a *DateHistogramAgg) MinDocCountValue() *int { return a.minDocCount }

// DateHistogram creates a date histogram aggregation.
func (b *Builder[T]) DateHistogram(name, field string) *DateHistogramAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &DateHistogramAgg{agg: agg{name: name, aggType: AggDateHistogram, err: err}}
	}
	return &DateHistogramAgg{agg: agg{name: name, aggType: AggDateHistogram, field: spec.Name}}
}

// RangeSpec defines a range bucket.
type RangeSpec struct {
	Key  string
	From any
	To   any
}

// RangeAgg creates custom range buckets.
type RangeAgg struct {
	agg
	ranges []RangeSpec
	keyed  *bool
}

// AddRange adds a range bucket.
func (a *RangeAgg) AddRange(from, to any) *RangeAgg {
	a.ranges = append(a.ranges, RangeSpec{From: from, To: to})
	return a
}

// AddKeyedRange adds a named range bucket.
func (a *RangeAgg) AddKeyedRange(key string, from, to any) *RangeAgg {
	a.ranges = append(a.ranges, RangeSpec{Key: key, From: from, To: to})
	return a
}

// Keyed sets whether to return buckets as a map.
func (a *RangeAgg) Keyed(k bool) *RangeAgg { a.keyed = &k; return a }

// SubAgg adds a sub-aggregation.
func (a *RangeAgg) SubAgg(sub Aggregation) *RangeAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// Ranges returns the range specifications.
func (a *RangeAgg) Ranges() []RangeSpec { return a.ranges }

// KeyedValue returns the keyed value if set.
func (a *RangeAgg) KeyedValue() *bool { return a.keyed }

// RangeAgg creates a range aggregation.
func (b *Builder[T]) RangeAgg(name, field string) *RangeAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &RangeAgg{agg: agg{name: name, aggType: AggRange, err: err}}
	}
	return &RangeAgg{agg: agg{name: name, aggType: AggRange, field: spec.Name}}
}

// FilterAgg creates a single filter bucket.
type FilterAgg struct {
	agg
	filter Query
}

// SubAgg adds a sub-aggregation.
func (a *FilterAgg) SubAgg(sub Aggregation) *FilterAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// FilterQuery returns the filter query.
func (a *FilterAgg) FilterQuery() Query { return a.filter }

// Err returns any error in this aggregation or its filter.
func (a *FilterAgg) Err() error {
	if a.err != nil {
		return a.err
	}
	if a.filter != nil {
		return a.filter.Err()
	}
	return nil
}

// FilterAgg creates a filter aggregation.
func (b *Builder[T]) FilterAgg(name string, filter Query) *FilterAgg {
	return &FilterAgg{
		agg:    agg{name: name, aggType: AggFilter},
		filter: filter,
	}
}

// NestedAgg aggregates nested documents.
type NestedAgg struct {
	agg
	path string
}

// SubAgg adds a sub-aggregation.
func (a *NestedAgg) SubAgg(sub Aggregation) *NestedAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// Path returns the nested path.
func (a *NestedAgg) Path() string { return a.path }

// NestedAgg creates a nested aggregation.
func (b *Builder[T]) NestedAgg(name, path string) *NestedAgg {
	return &NestedAgg{
		agg:  agg{name: name, aggType: AggNested},
		path: path,
	}
}

// MissingAgg counts documents missing a field.
type MissingAgg struct {
	agg
}

// SubAgg adds a sub-aggregation.
func (a *MissingAgg) SubAgg(sub Aggregation) *MissingAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// MissingAgg creates a missing aggregation.
func (b *Builder[T]) MissingAgg(name, field string) *MissingAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &MissingAgg{agg: agg{name: name, aggType: AggMissing, err: err}}
	}
	return &MissingAgg{agg: agg{name: name, aggType: AggMissing, field: spec.Name}}
}

// DateRangeSpec defines a date range bucket.
type DateRangeSpec struct {
	Key  string
	From any
	To   any
}

// DateRangeAgg creates date range buckets.
type DateRangeAgg struct {
	agg
	ranges []DateRangeSpec
	format *string
	keyed  *bool
}

// AddRange adds a date range bucket.
func (a *DateRangeAgg) AddRange(from, to any) *DateRangeAgg {
	a.ranges = append(a.ranges, DateRangeSpec{From: from, To: to})
	return a
}

// AddKeyedRange adds a named date range bucket.
func (a *DateRangeAgg) AddKeyedRange(key string, from, to any) *DateRangeAgg {
	a.ranges = append(a.ranges, DateRangeSpec{Key: key, From: from, To: to})
	return a
}

// Format sets the date format for parsing string values.
func (a *DateRangeAgg) Format(f string) *DateRangeAgg { a.format = &f; return a }

// Keyed sets whether to return buckets as a map.
func (a *DateRangeAgg) Keyed(k bool) *DateRangeAgg { a.keyed = &k; return a }

// SubAgg adds a sub-aggregation.
func (a *DateRangeAgg) SubAgg(sub Aggregation) *DateRangeAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// Ranges returns the range specifications.
func (a *DateRangeAgg) Ranges() []DateRangeSpec { return a.ranges }

// FormatValue returns the format if set.
func (a *DateRangeAgg) FormatValue() *string { return a.format }

// KeyedValue returns the keyed value if set.
func (a *DateRangeAgg) KeyedValue() *bool { return a.keyed }

// DateRangeAgg creates a date_range aggregation.
func (b *Builder[T]) DateRangeAgg(name, field string) *DateRangeAgg {
	spec, err := b.resolveField(field)
	if err != nil {
		return &DateRangeAgg{agg: agg{name: name, aggType: AggDateRange, err: err}}
	}
	return &DateRangeAgg{agg: agg{name: name, aggType: AggDateRange, field: spec.Name}}
}

// FiltersAgg creates named filter buckets.
type FiltersAgg struct {
	agg
	filters map[string]Query
}

// Filter adds a named filter bucket.
func (a *FiltersAgg) Filter(name string, q Query) *FiltersAgg {
	if a.filters == nil {
		a.filters = make(map[string]Query)
	}
	a.filters[name] = q
	return a
}

// SubAgg adds a sub-aggregation.
func (a *FiltersAgg) SubAgg(sub Aggregation) *FiltersAgg {
	a.subAggs = append(a.subAggs, sub)
	return a
}

// Filters returns the named filters.
func (a *FiltersAgg) Filters() map[string]Query { return a.filters }

// Err returns any error in this aggregation or its filters.
func (a *FiltersAgg) Err() error {
	if a.err != nil {
		return a.err
	}
	for _, q := range a.filters {
		if err := q.Err(); err != nil {
			return err
		}
	}
	return nil
}

// FiltersAgg creates a filters aggregation with named buckets.
func (b *Builder[T]) FiltersAgg(name string) *FiltersAgg {
	return &FiltersAgg{
		agg: agg{name: name, aggType: AggFilters},
	}
}
