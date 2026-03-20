// Package elasticsearch provides an Elasticsearch-specific query renderer.
package elasticsearch

import (
	"encoding/json"
	"fmt"

	"github.com/zoobz-io/lucene"
	"github.com/zoobz-io/lucene/internal/marshal"
)

// Version represents an Elasticsearch version.
type Version int

const (
	// V7 targets Elasticsearch 7.x.
	V7 Version = 7
	// V8 targets Elasticsearch 8.x.
	V8 Version = 8
)

// Renderer converts lucene queries to Elasticsearch JSON.
type Renderer struct {
	version Version
}

// NewRenderer creates a new Elasticsearch renderer for the specified version.
func NewRenderer(v Version) *Renderer {
	return &Renderer{version: v}
}

// Render converts a complete search request to JSON.
func (r *Renderer) Render(s *lucene.Search) ([]byte, error) {
	if s == nil {
		return nil, fmt.Errorf("search request is nil")
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("search request has error: %w", err)
	}

	req := marshal.SearchRequest{}

	// Query
	if q := s.QueryValue(); q != nil {
		rendered, err := r.renderQuery(q)
		if err != nil {
			return nil, err
		}
		req.Query = rendered
	}

	// Aggregations
	if aggs := s.AggsValue(); len(aggs) > 0 {
		aggsResult := make(map[string]any)
		for _, agg := range aggs {
			rendered, err := r.renderAgg(agg)
			if err != nil {
				return nil, err
			}
			aggsResult[agg.Name()] = rendered
		}
		req.Aggs = aggsResult
	}

	// Size
	req.Size = s.SizeValue()

	// From
	req.From = s.FromValue()

	// Sort
	if sorts := s.SortValue(); len(sorts) > 0 {
		sortEntries := make([]marshal.SortEntry, 0, len(sorts))
		for _, sort := range sorts {
			sortEntries = append(sortEntries, marshal.SortEntry{
				Field: sort.Field,
				Order: sort.Order,
			})
		}
		req.Sort = sortEntries
	}

	// Source filtering
	if includes := s.SourceIncludesValue(); len(includes) > 0 {
		if excludes := s.SourceExcludesValue(); len(excludes) > 0 {
			req.Source = marshal.SourceFilter{
				Includes: includes,
				Excludes: excludes,
			}
		} else {
			req.Source = includes
		}
	} else if excludes := s.SourceExcludesValue(); len(excludes) > 0 {
		req.Source = marshal.SourceFilter{
			Excludes: excludes,
		}
	}

	// Highlight
	if h := s.HighlightValue(); h != nil {
		highlight, err := r.renderHighlight(h)
		if err != nil {
			return nil, err
		}
		req.Highlight = highlight
	}

	// Track total hits
	req.TrackTotalHits = s.TrackTotalHitsValue()

	// Min score
	req.MinScore = s.MinScoreValue()

	// Timeout
	req.Timeout = s.TimeoutValue()

	return json.Marshal(req)
}

func (r *Renderer) renderHighlight(h *lucene.Highlight) (*marshal.Highlight, error) {
	result := &marshal.Highlight{
		PreTags:      h.PreTagsValue(),
		PostTags:     h.PostTagsValue(),
		Encoder:      h.EncoderValue(),
		FragmentSize: h.FragmentSizeValue(),
		NumFragments: h.NumFragmentsValue(),
		Order:        h.OrderValue(),
		Type:         h.HighlighterValue(),
	}

	// Fields
	if fields := h.FieldsValue(); len(fields) > 0 {
		fieldsMap := make(map[string]marshal.HighlightField)
		for _, f := range fields {
			hf := marshal.HighlightField{
				FragmentSize:      f.FragmentSize,
				NumFragments:      f.NumFragments,
				PreTags:           f.PreTags,
				PostTags:          f.PostTags,
				MatchedFields:     f.MatchedFields,
				FragmentOffset:    f.FragmentOffset,
				NoMatchSize:       f.NoMatchSize,
				RequireFieldMatch: f.RequireFieldMatch,
			}
			if q := f.HighlightQuery; q != nil {
				rendered, err := r.renderQuery(q)
				if err != nil {
					return nil, err
				}
				hf.HighlightQuery = rendered
			}
			fieldsMap[f.Name] = hf
		}
		result.Fields = fieldsMap
	}

	return result, nil
}

// RenderQuery converts a single query to JSON.
func (r *Renderer) RenderQuery(q lucene.Query) ([]byte, error) {
	if err := q.Err(); err != nil {
		return nil, fmt.Errorf("query has error: %w", err)
	}

	obj, err := r.renderQuery(q)
	if err != nil {
		return nil, err
	}

	return json.Marshal(obj)
}

// RenderAggs converts aggregations to JSON.
func (r *Renderer) RenderAggs(aggs []lucene.Aggregation) ([]byte, error) {
	if len(aggs) == 0 {
		return []byte("{}"), nil
	}

	result := make(map[string]any)
	for _, agg := range aggs {
		if err := agg.Err(); err != nil {
			return nil, fmt.Errorf("aggregation %s has error: %w", agg.Name(), err)
		}
		rendered, err := r.renderAgg(agg)
		if err != nil {
			return nil, err
		}
		result[agg.Name()] = rendered
	}

	return json.Marshal(result)
}

func (r *Renderer) renderAgg(a lucene.Aggregation) (any, error) {
	var inner any
	var aggKey string

	switch v := a.(type) {
	case *lucene.TermsAgg:
		aggKey = "terms"
		inner = marshal.TermsAggInner{
			Field:       v.Field(),
			Size:        v.SizeValue(),
			MinDocCount: v.MinDocCountValue(),
			Order:       v.OrderValue(),
		}
	case *lucene.HistogramAgg:
		aggKey = "histogram"
		inner = marshal.HistogramAggInner{
			Field:       v.Field(),
			Interval:    v.IntervalValue(),
			Offset:      v.OffsetValue(),
			MinDocCount: v.MinDocCountValue(),
		}
	case *lucene.DateHistogramAgg:
		aggKey = "date_histogram"
		inner = marshal.DateHistogramAggInner{
			Field:            v.Field(),
			CalendarInterval: v.CalendarIntervalValue(),
			FixedInterval:    v.FixedIntervalValue(),
			Format:           v.FormatValue(),
			TimeZone:         v.TimeZoneValue(),
			MinDocCount:      v.MinDocCountValue(),
		}
	case *lucene.RangeAgg:
		aggKey = "range"
		ranges := make([]marshal.RangeSpec, 0, len(v.Ranges()))
		for _, rng := range v.Ranges() {
			ranges = append(ranges, marshal.RangeSpec{
				Key:  rng.Key,
				From: rng.From,
				To:   rng.To,
			})
		}
		inner = marshal.RangeAggInner{
			Field:  v.Field(),
			Ranges: ranges,
			Keyed:  v.KeyedValue(),
		}
	case *lucene.DateRangeAgg:
		aggKey = "date_range"
		ranges := make([]marshal.RangeSpec, 0, len(v.Ranges()))
		for _, rng := range v.Ranges() {
			ranges = append(ranges, marshal.RangeSpec{
				Key:  rng.Key,
				From: rng.From,
				To:   rng.To,
			})
		}
		inner = marshal.DateRangeAggInner{
			Field:  v.Field(),
			Ranges: ranges,
			Format: v.FormatValue(),
			Keyed:  v.KeyedValue(),
		}
	case *lucene.FilterAgg:
		return r.renderFilterAgg(v)
	case *lucene.FiltersAgg:
		return r.renderFiltersAgg(v)
	case *lucene.NestedAgg:
		aggKey = "nested"
		inner = marshal.NestedAggInner{Path: v.Path()}
	case *lucene.MissingAgg:
		aggKey = "missing"
		inner = marshal.MissingAggInner{Field: v.Field()}
	case *lucene.AvgAgg:
		aggKey = "avg"
		inner = marshal.MetricAggInner{Field: v.Field(), Missing: v.MissingValue()}
	case *lucene.SumAgg:
		aggKey = "sum"
		inner = marshal.MetricAggInner{Field: v.Field(), Missing: v.MissingValue()}
	case *lucene.MinAgg:
		aggKey = "min"
		inner = marshal.MetricAggInner{Field: v.Field(), Missing: v.MissingValue()}
	case *lucene.MaxAgg:
		aggKey = "max"
		inner = marshal.MetricAggInner{Field: v.Field(), Missing: v.MissingValue()}
	case *lucene.CountAgg:
		aggKey = "value_count"
		inner = marshal.MetricAggInner{Field: v.Field()}
	case *lucene.CardinalityAgg:
		aggKey = "cardinality"
		inner = marshal.CardinalityAggInner{
			Field:              v.Field(),
			PrecisionThreshold: v.PrecisionThresholdValue(),
		}
	case *lucene.StatsAgg:
		aggKey = "stats"
		inner = marshal.MetricAggInner{Field: v.Field(), Missing: v.MissingValue()}
	case *lucene.ExtendedStatsAgg:
		aggKey = "extended_stats"
		inner = marshal.ExtendedStatsAggInner{
			Field:   v.Field(),
			Missing: v.MissingValue(),
			Sigma:   v.SigmaValue(),
		}
	case *lucene.PercentilesAgg:
		aggKey = "percentiles"
		inner = marshal.PercentilesAggInner{
			Field:    v.Field(),
			Percents: v.PercentsValue(),
			Missing:  v.MissingValue(),
		}
	case *lucene.TopHitsAgg:
		aggKey = "top_hits"
		inner = r.renderTopHitsAgg(v)
	case *lucene.AvgBucketAgg:
		aggKey = "avg_bucket"
		inner = marshal.PipelineAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
		}
	case *lucene.SumBucketAgg:
		aggKey = "sum_bucket"
		inner = marshal.PipelineAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
		}
	case *lucene.MaxBucketAgg:
		aggKey = "max_bucket"
		inner = marshal.PipelineAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
		}
	case *lucene.MinBucketAgg:
		aggKey = "min_bucket"
		inner = marshal.PipelineAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
		}
	case *lucene.DerivativeAgg:
		aggKey = "derivative"
		inner = marshal.DerivativeAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
			Unit:        v.UnitValue(),
		}
	case *lucene.CumulativeSumAgg:
		aggKey = "cumulative_sum"
		inner = marshal.PipelineAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
		}
	case *lucene.MovingAvgAgg:
		aggKey = "moving_avg"
		inner = marshal.MovingAvgAggInner{
			BucketsPath: v.BucketsPath(),
			GapPolicy:   v.GapPolicyValue(),
			Format:      v.FormatValue(),
			Window:      v.WindowValue(),
			Model:       v.ModelValue(),
			Predict:     v.PredictValue(),
		}
	default:
		return nil, fmt.Errorf("unsupported aggregation type: %T", a)
	}

	// Build result with sub-aggregations
	var subAggs map[string]any
	if subs := a.SubAggs(); len(subs) > 0 {
		subAggs = make(map[string]any)
		for _, sub := range subs {
			rendered, err := r.renderAgg(sub)
			if err != nil {
				return nil, err
			}
			subAggs[sub.Name()] = rendered
		}
	}

	return marshal.Agg[any]{
		AggType: aggKey,
		Inner:   inner,
		SubAggs: subAggs,
	}, nil
}

func (r *Renderer) renderFilterAgg(a *lucene.FilterAgg) (any, error) {
	filter, err := r.renderQuery(a.FilterQuery())
	if err != nil {
		return nil, err
	}

	var subAggs map[string]any
	if subs := a.SubAggs(); len(subs) > 0 {
		subAggs = make(map[string]any)
		for _, sub := range subs {
			rendered, err := r.renderAgg(sub)
			if err != nil {
				return nil, err
			}
			subAggs[sub.Name()] = rendered
		}
	}

	return marshal.Agg[any]{
		AggType: "filter",
		Inner:   filter,
		SubAggs: subAggs,
	}, nil
}

func (r *Renderer) renderFiltersAgg(a *lucene.FiltersAgg) (any, error) {
	filters := make(map[string]any)
	for name, q := range a.Filters() {
		rendered, err := r.renderQuery(q)
		if err != nil {
			return nil, err
		}
		filters[name] = rendered
	}

	var subAggs map[string]any
	if subs := a.SubAggs(); len(subs) > 0 {
		subAggs = make(map[string]any)
		for _, sub := range subs {
			rendered, err := r.renderAgg(sub)
			if err != nil {
				return nil, err
			}
			subAggs[sub.Name()] = rendered
		}
	}

	return marshal.Agg[marshal.FiltersAggInner]{
		AggType: "filters",
		Inner:   marshal.FiltersAggInner{Filters: filters},
		SubAggs: subAggs,
	}, nil
}

func (r *Renderer) renderTopHitsAgg(a *lucene.TopHitsAgg) marshal.TopHitsAggInner {
	result := marshal.TopHitsAggInner{
		Size:   a.SizeValue(),
		From:   a.FromValue(),
		Source: a.SourceValue(),
	}

	if sorts := a.SortValue(); len(sorts) > 0 {
		sortList := make([]any, 0, len(sorts))
		for _, s := range sorts {
			sortList = append(sortList, marshal.SortEntry{Field: s.Field, Order: s.Order})
		}
		result.Sort = sortList
	}

	return result
}

// renderQuery converts a query to a typed marshal structure.
func (r *Renderer) renderQuery(q lucene.Query) (any, error) {
	switch v := q.(type) {
	// Term-level queries
	case *lucene.TermQuery:
		return marshal.FieldQuery[marshal.TermInner]{
			QueryType: "term",
			Field:     v.Field(),
			Inner: marshal.TermInner{
				Value: v.Value(),
				Boost: v.BoostValue(),
			},
		}, nil

	case *lucene.TermsQuery:
		return r.renderTerms(v), nil

	case *lucene.RangeQuery:
		return marshal.FieldQuery[marshal.RangeInner]{
			QueryType: "range",
			Field:     v.Field(),
			Inner: marshal.RangeInner{
				Gt:     v.GtValue(),
				Gte:    v.GteValue(),
				Lt:     v.LtValue(),
				Lte:    v.LteValue(),
				Format: v.FormatValue(),
				Boost:  v.BoostValue(),
			},
		}, nil

	case *lucene.ExistsQuery:
		return marshal.SimpleQuery[marshal.ExistsInner]{
			QueryType: "exists",
			Inner:     marshal.ExistsInner{Field: v.Field()},
		}, nil

	case *lucene.IDsQuery:
		return marshal.SimpleQuery[marshal.IDsInner]{
			QueryType: "ids",
			Inner:     marshal.IDsInner{Values: v.IDValues()},
		}, nil

	case *lucene.PrefixQuery:
		return marshal.FieldQuery[marshal.PrefixInner]{
			QueryType: "prefix",
			Field:     v.Field(),
			Inner: marshal.PrefixInner{
				Value:           v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Rewrite:         v.RewriteValue(),
				CaseInsensitive: v.CaseInsensitiveValue(),
				Boost:           v.BoostValue(),
			},
		}, nil

	case *lucene.WildcardQuery:
		return marshal.FieldQuery[marshal.WildcardInner]{
			QueryType: "wildcard",
			Field:     v.Field(),
			Inner: marshal.WildcardInner{
				Value:           v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Rewrite:         v.RewriteValue(),
				CaseInsensitive: v.CaseInsensitiveValue(),
				Boost:           v.BoostValue(),
			},
		}, nil

	case *lucene.RegexpQuery:
		return marshal.FieldQuery[marshal.RegexpInner]{
			QueryType: "regexp",
			Field:     v.Field(),
			Inner: marshal.RegexpInner{
				Value:           v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Flags:           v.FlagsValue(),
				Rewrite:         v.RewriteValue(),
				CaseInsensitive: v.CaseInsensitiveValue(),
				Boost:           v.BoostValue(),
			},
		}, nil

	case *lucene.FuzzyQuery:
		return marshal.FieldQuery[marshal.FuzzyInner]{
			QueryType: "fuzzy",
			Field:     v.Field(),
			Inner: marshal.FuzzyInner{
				Value:          v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Fuzziness:      v.FuzzinessValue(),
				PrefixLength:   v.PrefixLengthValue(),
				MaxExpansions:  v.MaxExpansionsValue(),
				Transpositions: v.TranspositionsValue(),
				Rewrite:        v.RewriteValue(),
				Boost:          v.BoostValue(),
			},
		}, nil

	// Full-text queries
	case *lucene.MatchQuery:
		return marshal.FieldQuery[marshal.MatchInner]{
			QueryType: "match",
			Field:     v.Field(),
			Inner: marshal.MatchInner{
				Query:     v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Fuzziness: v.FuzzinessValue(),
				Operator:  v.OperatorValue(),
				Analyzer:  v.AnalyzerValue(),
				Boost:     v.BoostValue(),
			},
		}, nil

	case *lucene.MatchPhraseQuery:
		return marshal.FieldQuery[marshal.MatchPhraseInner]{
			QueryType: "match_phrase",
			Field:     v.Field(),
			Inner: marshal.MatchPhraseInner{
				Query:    v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Slop:     v.SlopValue(),
				Analyzer: v.AnalyzerValue(),
				Boost:    v.BoostValue(),
			},
		}, nil

	case *lucene.MatchPhrasePrefixQuery:
		return marshal.FieldQuery[marshal.MatchPhrasePrefixInner]{
			QueryType: "match_phrase_prefix",
			Field:     v.Field(),
			Inner: marshal.MatchPhrasePrefixInner{
				Query:         v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Slop:          v.SlopValue(),
				MaxExpansions: v.MaxExpansionsValue(),
				Analyzer:      v.AnalyzerValue(),
				Boost:         v.BoostValue(),
			},
		}, nil

	case *lucene.MultiMatchQuery:
		return marshal.SimpleQuery[marshal.MultiMatchInner]{
			QueryType: "multi_match",
			Inner: marshal.MultiMatchInner{
				Query:      v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Fields:     v.Fields(),
				Type:       v.TypeValue(),
				TieBreaker: v.TieBreakerValue(),
				Fuzziness:  v.FuzzinessValue(),
				Operator:   v.OperatorValue(),
				Analyzer:   v.AnalyzerValue(),
				Boost:      v.BoostValue(),
			},
		}, nil

	case *lucene.QueryStringQuery:
		return marshal.SimpleQuery[marshal.QueryStringInner]{
			QueryType: "query_string",
			Inner: marshal.QueryStringInner{
				Query:                v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				DefaultField:         v.DefaultFieldValue(),
				DefaultOperator:      v.DefaultOperatorValue(),
				Analyzer:             v.AnalyzerValue(),
				AllowLeadingWildcard: v.AllowLeadingWildcardValue(),
				Fuzziness:            v.FuzzinessValue(),
				Boost:                v.BoostValue(),
			},
		}, nil

	case *lucene.SimpleQueryStringQuery:
		return marshal.SimpleQuery[marshal.SimpleQueryStringInner]{
			QueryType: "simple_query_string",
			Inner: marshal.SimpleQueryStringInner{
				Query:           v.Value().(string), //nolint:errcheck // Type guaranteed by builder
				Fields:          v.FieldsValue(),
				DefaultOperator: v.DefaultOperatorValue(),
				Analyzer:        v.AnalyzerValue(),
				Flags:           v.FlagsValue(),
				Boost:           v.BoostValue(),
			},
		}, nil

	// Compound queries
	case *lucene.BoolQuery:
		return r.renderBool(v)

	case *lucene.MatchAllQuery:
		return marshal.SimpleQuery[marshal.MatchAllInner]{
			QueryType: "match_all",
			Inner:     marshal.MatchAllInner{Boost: v.BoostValue()},
		}, nil

	case *lucene.MatchNoneQuery:
		return marshal.SimpleQuery[marshal.MatchNoneInner]{
			QueryType: "match_none",
			Inner:     marshal.MatchNoneInner{},
		}, nil

	case *lucene.BoostingQuery:
		return r.renderBoosting(v)

	case *lucene.DisMaxQuery:
		return r.renderDisMax(v)

	case *lucene.ConstantScoreQuery:
		return r.renderConstantScore(v)

	// Nested queries
	case *lucene.NestedQuery:
		return r.renderNested(v)

	case *lucene.HasChildQuery:
		return r.renderHasChild(v)

	case *lucene.HasParentQuery:
		return r.renderHasParent(v)

	// Vector queries
	case *lucene.KnnQuery:
		return r.renderKnn(v)

	// Geo queries
	case *lucene.GeoDistanceQuery:
		return r.renderGeoDistance(v), nil

	case *lucene.GeoBoundingBoxQuery:
		return r.renderGeoBoundingBox(v), nil

	default:
		return nil, fmt.Errorf("unsupported query type: %T", q)
	}
}

// renderTerms handles the special case of terms query where field maps to values array.
func (r *Renderer) renderTerms(q *lucene.TermsQuery) any {
	inner := map[string]any{
		q.Field(): q.Values(),
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"terms": inner}
}

func (r *Renderer) renderBool(q *lucene.BoolQuery) (any, error) {
	inner := marshal.BoolInner{
		MinimumShouldMatch: q.MinimumShouldMatchValue(),
		Boost:              q.BoostValue(),
	}

	if must := q.MustClauses(); len(must) > 0 {
		clauses, err := r.renderClauses(must)
		if err != nil {
			return nil, err
		}
		inner.Must = clauses
	}

	if should := q.ShouldClauses(); len(should) > 0 {
		clauses, err := r.renderClauses(should)
		if err != nil {
			return nil, err
		}
		inner.Should = clauses
	}

	if mustNot := q.MustNotClauses(); len(mustNot) > 0 {
		clauses, err := r.renderClauses(mustNot)
		if err != nil {
			return nil, err
		}
		inner.MustNot = clauses
	}

	if filter := q.FilterClauses(); len(filter) > 0 {
		clauses, err := r.renderClauses(filter)
		if err != nil {
			return nil, err
		}
		inner.Filter = clauses
	}

	return marshal.SimpleQuery[marshal.BoolInner]{
		QueryType: "bool",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderClauses(queries []lucene.Query) ([]any, error) {
	result := make([]any, 0, len(queries))
	for _, q := range queries {
		rendered, err := r.renderQuery(q)
		if err != nil {
			return nil, err
		}
		result = append(result, rendered)
	}
	return result, nil
}

func (r *Renderer) renderBoosting(q *lucene.BoostingQuery) (any, error) {
	inner := marshal.BoostingInner{
		NegativeBoost: q.NegativeBoostValue(),
	}

	if pos := q.PositiveQuery(); pos != nil {
		rendered, err := r.renderQuery(pos)
		if err != nil {
			return nil, err
		}
		inner.Positive = rendered
	}

	if neg := q.NegativeQuery(); neg != nil {
		rendered, err := r.renderQuery(neg)
		if err != nil {
			return nil, err
		}
		inner.Negative = rendered
	}

	return marshal.SimpleQuery[marshal.BoostingInner]{
		QueryType: "boosting",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderDisMax(q *lucene.DisMaxQuery) (any, error) {
	inner := marshal.DisMaxInner{
		TieBreaker: q.TieBreakerValue(),
		Boost:      q.BoostValue(),
	}

	if queries := q.Queries(); len(queries) > 0 {
		rendered, err := r.renderClauses(queries)
		if err != nil {
			return nil, err
		}
		inner.Queries = rendered
	}

	return marshal.SimpleQuery[marshal.DisMaxInner]{
		QueryType: "dis_max",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderConstantScore(q *lucene.ConstantScoreQuery) (any, error) {
	inner := marshal.ConstantScoreInner{
		Boost: q.BoostValue(),
	}

	if filter := q.FilterQuery(); filter != nil {
		rendered, err := r.renderQuery(filter)
		if err != nil {
			return nil, err
		}
		inner.Filter = rendered
	}

	return marshal.SimpleQuery[marshal.ConstantScoreInner]{
		QueryType: "constant_score",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderNested(q *lucene.NestedQuery) (any, error) {
	inner := marshal.NestedInner{
		Path:           q.Path(),
		ScoreMode:      q.ScoreModeValue(),
		IgnoreUnmapped: q.IgnoreUnmappedValue(),
	}

	if iq := q.InnerQuery(); iq != nil {
		rendered, err := r.renderQuery(iq)
		if err != nil {
			return nil, err
		}
		inner.Query = rendered
	}

	return marshal.SimpleQuery[marshal.NestedInner]{
		QueryType: "nested",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderHasChild(q *lucene.HasChildQuery) (any, error) {
	inner := marshal.HasChildInner{
		Type:           q.ChildType(),
		ScoreMode:      q.ScoreModeValue(),
		MinChildren:    q.MinChildrenValue(),
		MaxChildren:    q.MaxChildrenValue(),
		IgnoreUnmapped: q.IgnoreUnmappedValue(),
	}

	if iq := q.InnerQuery(); iq != nil {
		rendered, err := r.renderQuery(iq)
		if err != nil {
			return nil, err
		}
		inner.Query = rendered
	}

	return marshal.SimpleQuery[marshal.HasChildInner]{
		QueryType: "has_child",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderHasParent(q *lucene.HasParentQuery) (any, error) {
	inner := marshal.HasParentInner{
		ParentType:     q.ParentType(),
		Score:          q.ScoreValue(),
		IgnoreUnmapped: q.IgnoreUnmappedValue(),
	}

	if iq := q.InnerQuery(); iq != nil {
		rendered, err := r.renderQuery(iq)
		if err != nil {
			return nil, err
		}
		inner.Query = rendered
	}

	return marshal.SimpleQuery[marshal.HasParentInner]{
		QueryType: "has_parent",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderKnn(q *lucene.KnnQuery) (any, error) {
	inner := marshal.KnnInnerES{
		Field:         q.Field(),
		Vector:        q.Vector(),
		K:             q.KValue(),
		NumCandidates: q.NumCandidatesValue(),
		Boost:         q.BoostValue(),
	}

	if filter := q.FilterQuery(); filter != nil {
		rendered, err := r.renderQuery(filter)
		if err != nil {
			return nil, err
		}
		inner.Filter = rendered
	}

	return marshal.SimpleQuery[marshal.KnnInnerES]{
		QueryType: "knn",
		Inner:     inner,
	}, nil
}

func (r *Renderer) renderGeoDistance(q *lucene.GeoDistanceQuery) any {
	// GeoDistance has dynamic field name for the point
	inner := map[string]any{
		q.Field(): marshal.GeoPoint{
			Lat: q.Lat(),
			Lon: q.Lon(),
		},
	}
	if v := q.DistanceValue(); v != nil {
		inner["distance"] = *v
	}
	if v := q.DistanceTypeValue(); v != nil {
		inner["distance_type"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"geo_distance": inner}
}

func (r *Renderer) renderGeoBoundingBox(q *lucene.GeoBoundingBoxQuery) any {
	// GeoBoundingBox has dynamic field name for the box
	fieldInner := make(map[string]any)
	if lat := q.TopLeftLat(); lat != nil {
		if lon := q.TopLeftLon(); lon != nil {
			fieldInner["top_left"] = marshal.GeoPoint{Lat: *lat, Lon: *lon}
		}
	}
	if lat := q.BottomRightLat(); lat != nil {
		if lon := q.BottomRightLon(); lon != nil {
			fieldInner["bottom_right"] = marshal.GeoPoint{Lat: *lat, Lon: *lon}
		}
	}

	inner := map[string]any{
		q.Field(): fieldInner,
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"geo_bounding_box": inner}
}
