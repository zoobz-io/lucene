package marshal

// === Bucket aggregation inner types ===

// TermsAggInner is the inner structure for terms aggregations.
type TermsAggInner struct {
	Field       string            `json:"field"`
	Size        *int              `json:"size,omitempty"`
	MinDocCount *int              `json:"min_doc_count,omitempty"`
	Order       map[string]string `json:"order,omitempty"`
}

// HistogramAggInner is the inner structure for histogram aggregations.
type HistogramAggInner struct {
	Field       string   `json:"field"`
	Interval    *float64 `json:"interval,omitempty"`
	Offset      *float64 `json:"offset,omitempty"`
	MinDocCount *int     `json:"min_doc_count,omitempty"`
}

// DateHistogramAggInner is the inner structure for date_histogram aggregations.
type DateHistogramAggInner struct {
	Field            string  `json:"field"`
	CalendarInterval *string `json:"calendar_interval,omitempty"`
	FixedInterval    *string `json:"fixed_interval,omitempty"`
	Format           *string `json:"format,omitempty"`
	TimeZone         *string `json:"time_zone,omitempty"`
	MinDocCount      *int    `json:"min_doc_count,omitempty"`
}

// RangeSpec defines a range bucket.
type RangeSpec struct {
	Key  string `json:"key,omitempty"`
	From any    `json:"from,omitempty"`
	To   any    `json:"to,omitempty"`
}

// RangeAggInner is the inner structure for range aggregations.
type RangeAggInner struct {
	Field  string      `json:"field"`
	Ranges []RangeSpec `json:"ranges"`
	Keyed  *bool       `json:"keyed,omitempty"`
}

// DateRangeAggInner is the inner structure for date_range aggregations.
type DateRangeAggInner struct {
	Field  string      `json:"field"`
	Ranges []RangeSpec `json:"ranges"`
	Format *string     `json:"format,omitempty"`
	Keyed  *bool       `json:"keyed,omitempty"`
}

// NestedAggInner is the inner structure for nested aggregations.
type NestedAggInner struct {
	Path string `json:"path"`
}

// MissingAggInner is the inner structure for missing aggregations.
type MissingAggInner struct {
	Field string `json:"field"`
}

// FilterAggInner wraps a filter query for filter aggregations.
// Note: The filter is the query itself, not wrapped in another object.
type FilterAggInner any

// FiltersAggInner is the inner structure for filters aggregations.
type FiltersAggInner struct {
	Filters map[string]any `json:"filters"`
}

// === Metric aggregation inner types ===

// MetricAggInner is the common structure for simple metric aggregations.
type MetricAggInner struct {
	Field   string `json:"field"`
	Missing any    `json:"missing,omitempty"`
}

// CardinalityAggInner is the inner structure for cardinality aggregations.
type CardinalityAggInner struct {
	Field              string `json:"field"`
	PrecisionThreshold *int   `json:"precision_threshold,omitempty"`
}

// ExtendedStatsAggInner is the inner structure for extended_stats aggregations.
type ExtendedStatsAggInner struct {
	Field   string   `json:"field"`
	Missing any      `json:"missing,omitempty"`
	Sigma   *float64 `json:"sigma,omitempty"`
}

// PercentilesAggInner is the inner structure for percentiles aggregations.
type PercentilesAggInner struct {
	Field    string    `json:"field"`
	Percents []float64 `json:"percents,omitempty"`
	Missing  any       `json:"missing,omitempty"`
}

// SortSpec defines a sort order for top_hits.
type SortSpec struct {
	Field string `json:"-"`
	Order string `json:"-"`
}

// TopHitsAggInner is the inner structure for top_hits aggregations.
type TopHitsAggInner struct {
	Size   *int     `json:"size,omitempty"`
	From   *int     `json:"from,omitempty"`
	Sort   []any    `json:"sort,omitempty"`   // Will be marshaled specially
	Source []string `json:"_source,omitempty"`
}

// === Pipeline aggregation inner types ===

// PipelineAggInner is the common structure for pipeline aggregations.
type PipelineAggInner struct {
	BucketsPath string  `json:"buckets_path"`
	GapPolicy   *string `json:"gap_policy,omitempty"`
	Format      *string `json:"format,omitempty"`
}

// DerivativeAggInner is the inner structure for derivative aggregations.
type DerivativeAggInner struct {
	BucketsPath string  `json:"buckets_path"`
	GapPolicy   *string `json:"gap_policy,omitempty"`
	Format      *string `json:"format,omitempty"`
	Unit        *string `json:"unit,omitempty"`
}

// MovingAvgAggInner is the inner structure for moving_avg aggregations.
type MovingAvgAggInner struct {
	BucketsPath string  `json:"buckets_path"`
	GapPolicy   *string `json:"gap_policy,omitempty"`
	Format      *string `json:"format,omitempty"`
	Window      *int    `json:"window,omitempty"`
	Model       *string `json:"model,omitempty"`
	Predict     *int    `json:"predict,omitempty"`
}
