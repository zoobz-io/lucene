package lucene

// PipelineAgg is a base for pipeline aggregations.
type PipelineAgg struct {
	agg
	bucketsPath string
	gapPolicy   *string
	format      *string
}

// GapPolicy sets how to handle gaps in data.
func (a *PipelineAgg) GapPolicy(p string) *PipelineAgg { a.gapPolicy = &p; return a }

// Format sets the output format.
func (a *PipelineAgg) Format(f string) *PipelineAgg { a.format = &f; return a }

// BucketsPath returns the buckets_path.
func (a *PipelineAgg) BucketsPath() string { return a.bucketsPath }

// GapPolicyValue returns the gap_policy if set.
func (a *PipelineAgg) GapPolicyValue() *string { return a.gapPolicy }

// FormatValue returns the format if set.
func (a *PipelineAgg) FormatValue() *string { return a.format }

// AvgBucketAgg computes the average of bucket values.
type AvgBucketAgg struct {
	PipelineAgg
}

// AvgBucket creates an avg_bucket pipeline aggregation.
func (b *Builder[T]) AvgBucket(name, bucketsPath string) *AvgBucketAgg {
	return &AvgBucketAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggAvgBucket},
			bucketsPath: bucketsPath,
		},
	}
}

// SumBucketAgg computes the sum of bucket values.
type SumBucketAgg struct {
	PipelineAgg
}

// SumBucket creates a sum_bucket pipeline aggregation.
func (b *Builder[T]) SumBucket(name, bucketsPath string) *SumBucketAgg {
	return &SumBucketAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggSumBucket},
			bucketsPath: bucketsPath,
		},
	}
}

// MaxBucketAgg finds the maximum bucket value.
type MaxBucketAgg struct {
	PipelineAgg
}

// MaxBucket creates a max_bucket pipeline aggregation.
func (b *Builder[T]) MaxBucket(name, bucketsPath string) *MaxBucketAgg {
	return &MaxBucketAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggMaxBucket},
			bucketsPath: bucketsPath,
		},
	}
}

// MinBucketAgg finds the minimum bucket value.
type MinBucketAgg struct {
	PipelineAgg
}

// MinBucket creates a min_bucket pipeline aggregation.
func (b *Builder[T]) MinBucket(name, bucketsPath string) *MinBucketAgg {
	return &MinBucketAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggMinBucket},
			bucketsPath: bucketsPath,
		},
	}
}

// DerivativeAgg computes the derivative of a metric.
type DerivativeAgg struct {
	PipelineAgg
	unit *string
}

// Unit sets the unit for normalization.
func (a *DerivativeAgg) Unit(u string) *DerivativeAgg { a.unit = &u; return a }

// UnitValue returns the unit if set.
func (a *DerivativeAgg) UnitValue() *string { return a.unit }

// Derivative creates a derivative pipeline aggregation.
func (b *Builder[T]) Derivative(name, bucketsPath string) *DerivativeAgg {
	return &DerivativeAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggDerivative},
			bucketsPath: bucketsPath,
		},
	}
}

// CumulativeSumAgg computes the cumulative sum.
type CumulativeSumAgg struct {
	PipelineAgg
}

// CumulativeSum creates a cumulative_sum pipeline aggregation.
func (b *Builder[T]) CumulativeSum(name, bucketsPath string) *CumulativeSumAgg {
	return &CumulativeSumAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggCumulativeSum},
			bucketsPath: bucketsPath,
		},
	}
}

// MovingAvgAgg computes a moving average.
type MovingAvgAgg struct {
	PipelineAgg
	window  *int
	model   *string
	predict *int
}

// Window sets the window size.
func (a *MovingAvgAgg) Window(w int) *MovingAvgAgg { a.window = &w; return a }

// Model sets the smoothing model.
func (a *MovingAvgAgg) Model(m string) *MovingAvgAgg { a.model = &m; return a }

// Predict sets the number of predictions.
func (a *MovingAvgAgg) Predict(p int) *MovingAvgAgg { a.predict = &p; return a }

// WindowValue returns the window if set.
func (a *MovingAvgAgg) WindowValue() *int { return a.window }

// ModelValue returns the model if set.
func (a *MovingAvgAgg) ModelValue() *string { return a.model }

// PredictValue returns the predict if set.
func (a *MovingAvgAgg) PredictValue() *int { return a.predict }

// MovingAvg creates a moving_avg pipeline aggregation.
func (b *Builder[T]) MovingAvg(name, bucketsPath string) *MovingAvgAgg {
	return &MovingAvgAgg{
		PipelineAgg: PipelineAgg{
			agg:         agg{name: name, aggType: AggMovingAvg},
			bucketsPath: bucketsPath,
		},
	}
}
