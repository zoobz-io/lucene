package lucene

import "testing"

type pipelineTestDoc struct {
	Category  string  `json:"category"`
	Price     float64 `json:"price"`
	Timestamp string  `json:"timestamp"`
}

func TestPipelineAggBase(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("buckets_path", func(t *testing.T) {
		a := b.AvgBucket("avg_monthly", "by_month>avg_price")
		if a.BucketsPath() != "by_month>avg_price" {
			t.Errorf("BucketsPath() = %v, want by_month>avg_price", a.BucketsPath())
		}
	})

	t.Run("gap_policy", func(t *testing.T) {
		a := b.AvgBucket("avg_monthly", "by_month>avg_price").GapPolicy("skip")
		if a.GapPolicyValue() == nil || *a.GapPolicyValue() != "skip" {
			t.Errorf("GapPolicyValue() = %v, want skip", a.GapPolicyValue())
		}
	})

	t.Run("format", func(t *testing.T) {
		a := b.AvgBucket("avg_monthly", "by_month>avg_price").Format("0.00")
		if a.FormatValue() == nil || *a.FormatValue() != "0.00" {
			t.Errorf("FormatValue() = %v, want 0.00", a.FormatValue())
		}
	})
}

func TestAvgBucketAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.AvgBucket("avg_monthly", "by_month>avg_price")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.Name() != "avg_monthly" {
		t.Errorf("Name() = %v, want avg_monthly", a.Name())
	}
	if a.Type() != AggAvgBucket {
		t.Errorf("Type() = %v, want AggAvgBucket", a.Type())
	}
}

func TestSumBucketAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.SumBucket("sum_monthly", "by_month>sum_price")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.Type() != AggSumBucket {
		t.Errorf("Type() = %v, want AggSumBucket", a.Type())
	}
}

func TestMaxBucketAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.MaxBucket("max_monthly", "by_month>max_price")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.Type() != AggMaxBucket {
		t.Errorf("Type() = %v, want AggMaxBucket", a.Type())
	}
}

func TestMinBucketAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.MinBucket("min_monthly", "by_month>min_price")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.Type() != AggMinBucket {
		t.Errorf("Type() = %v, want AggMinBucket", a.Type())
	}
}

func TestDerivativeAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic derivative", func(t *testing.T) {
		a := b.Derivative("price_change", "by_month>avg_price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggDerivative {
			t.Errorf("Type() = %v, want AggDerivative", a.Type())
		}
	})

	t.Run("with unit", func(t *testing.T) {
		a := b.Derivative("price_change", "by_month>avg_price").Unit("day")
		if a.UnitValue() == nil || *a.UnitValue() != "day" {
			t.Errorf("UnitValue() = %v, want day", a.UnitValue())
		}
	})
}

func TestCumulativeSumAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	a := b.CumulativeSum("running_total", "by_month>sum_price")
	if a.Err() != nil {
		t.Errorf("Err() = %v, want nil", a.Err())
	}
	if a.Type() != AggCumulativeSum {
		t.Errorf("Type() = %v, want AggCumulativeSum", a.Type())
	}
}

func TestMovingAvgAgg(t *testing.T) {
	b, err := New[pipelineTestDoc]()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	t.Run("basic moving_avg", func(t *testing.T) {
		a := b.MovingAvg("smoothed", "by_month>avg_price")
		if a.Err() != nil {
			t.Errorf("Err() = %v, want nil", a.Err())
		}
		if a.Type() != AggMovingAvg {
			t.Errorf("Type() = %v, want AggMovingAvg", a.Type())
		}
	})

	t.Run("with window", func(t *testing.T) {
		a := b.MovingAvg("smoothed", "by_month>avg_price").Window(5)
		if a.WindowValue() == nil || *a.WindowValue() != 5 {
			t.Errorf("WindowValue() = %v, want 5", a.WindowValue())
		}
	})

	t.Run("with model", func(t *testing.T) {
		a := b.MovingAvg("smoothed", "by_month>avg_price").Model("ewma")
		if a.ModelValue() == nil || *a.ModelValue() != "ewma" {
			t.Errorf("ModelValue() = %v, want ewma", a.ModelValue())
		}
	})

	t.Run("with predict", func(t *testing.T) {
		a := b.MovingAvg("smoothed", "by_month>avg_price").Predict(10)
		if a.PredictValue() == nil || *a.PredictValue() != 10 {
			t.Errorf("PredictValue() = %v, want 10", a.PredictValue())
		}
	})

	t.Run("with all options", func(t *testing.T) {
		a := b.MovingAvg("smoothed", "by_month>avg_price").
			Window(5).
			Model("holt").
			Predict(3)

		// These are available via PipelineAgg embedding
		a.GapPolicy("insert_zeros")
		a.Format("0.00")

		if a.WindowValue() == nil || *a.WindowValue() != 5 {
			t.Errorf("WindowValue() = %v, want 5", a.WindowValue())
		}
		if a.ModelValue() == nil || *a.ModelValue() != "holt" {
			t.Errorf("ModelValue() = %v, want holt", a.ModelValue())
		}
		if a.PredictValue() == nil || *a.PredictValue() != 3 {
			t.Errorf("PredictValue() = %v, want 3", a.PredictValue())
		}
		if a.GapPolicyValue() == nil || *a.GapPolicyValue() != "insert_zeros" {
			t.Errorf("GapPolicyValue() = %v, want insert_zeros", a.GapPolicyValue())
		}
		if a.FormatValue() == nil || *a.FormatValue() != "0.00" {
			t.Errorf("FormatValue() = %v, want 0.00", a.FormatValue())
		}
	})
}
