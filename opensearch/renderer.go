// Package opensearch provides an OpenSearch-specific query renderer.
package opensearch

import (
	"encoding/json"
	"fmt"

	"github.com/zoobzio/lucene"
)

// Version represents an OpenSearch version.
type Version int

const (
	// V1 targets OpenSearch 1.x.
	V1 Version = 1
	// V2 targets OpenSearch 2.x.
	V2 Version = 2
)

// Renderer converts lucene queries to OpenSearch JSON.
type Renderer struct {
	version Version
}

// NewRenderer creates a new OpenSearch renderer for the specified version.
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

	result := make(map[string]any)

	// Query
	if q := s.QueryValue(); q != nil {
		rendered, err := r.renderQuery(q)
		if err != nil {
			return nil, err
		}
		result["query"] = rendered
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
		result["aggs"] = aggsResult
	}

	// Size
	if v := s.SizeValue(); v != nil {
		result["size"] = *v
	}

	// From
	if v := s.FromValue(); v != nil {
		result["from"] = *v
	}

	// Sort
	if sorts := s.SortValue(); len(sorts) > 0 {
		result["sort"] = r.renderSort(sorts)
	}

	// Source filtering
	if includes := s.SourceIncludesValue(); len(includes) > 0 {
		if excludes := s.SourceExcludesValue(); len(excludes) > 0 {
			result["_source"] = map[string]any{
				"includes": includes,
				"excludes": excludes,
			}
		} else {
			result["_source"] = includes
		}
	} else if excludes := s.SourceExcludesValue(); len(excludes) > 0 {
		result["_source"] = map[string]any{
			"excludes": excludes,
		}
	}

	// Highlight
	if h := s.HighlightValue(); h != nil {
		result["highlight"] = r.renderHighlight(h)
	}

	// Track total hits
	if v := s.TrackTotalHitsValue(); v != nil {
		result["track_total_hits"] = v
	}

	// Min score
	if v := s.MinScoreValue(); v != nil {
		result["min_score"] = *v
	}

	// Timeout
	if v := s.TimeoutValue(); v != nil {
		result["timeout"] = *v
	}

	return json.Marshal(result)
}

func (r *Renderer) renderSort(sorts []lucene.SortField) []any {
	result := make([]any, 0, len(sorts))
	for _, s := range sorts {
		if s.Order == "" {
			result = append(result, s.Field)
		} else {
			result = append(result, map[string]any{
				s.Field: map[string]any{"order": s.Order},
			})
		}
	}
	return result
}

func (r *Renderer) renderHighlight(h *lucene.Highlight) map[string]any {
	result := make(map[string]any)

	// Global settings
	if v := h.PreTagsValue(); len(v) > 0 {
		result["pre_tags"] = v
	}
	if v := h.PostTagsValue(); len(v) > 0 {
		result["post_tags"] = v
	}
	if v := h.EncoderValue(); v != nil {
		result["encoder"] = *v
	}
	if v := h.FragmentSizeValue(); v != nil {
		result["fragment_size"] = *v
	}
	if v := h.NumFragmentsValue(); v != nil {
		result["number_of_fragments"] = *v
	}
	if v := h.OrderValue(); v != nil {
		result["order"] = *v
	}
	if v := h.HighlighterValue(); v != nil {
		result["type"] = *v
	}

	// Fields
	if fields := h.FieldsValue(); len(fields) > 0 {
		fieldsMap := make(map[string]any)
		for _, f := range fields {
			fieldsMap[f.Name] = r.renderHighlightField(f)
		}
		result["fields"] = fieldsMap
	}

	return result
}

func (r *Renderer) renderHighlightField(f lucene.HighlightField) map[string]any {
	result := make(map[string]any)

	if v := f.FragmentSize; v != nil {
		result["fragment_size"] = *v
	}
	if v := f.NumFragments; v != nil {
		result["number_of_fragments"] = *v
	}
	if v := f.PreTags; len(v) > 0 {
		result["pre_tags"] = v
	}
	if v := f.PostTags; len(v) > 0 {
		result["post_tags"] = v
	}
	if v := f.MatchedFields; len(v) > 0 {
		result["matched_fields"] = v
	}
	if v := f.FragmentOffset; v != nil {
		result["fragment_offset"] = *v
	}
	if v := f.NoMatchSize; v != nil {
		result["no_match_size"] = *v
	}
	if v := f.RequireFieldMatch; v != nil {
		result["require_field_match"] = *v
	}
	if q := f.HighlightQuery; q != nil {
		if rendered, err := r.renderQuery(q); err == nil {
			result["highlight_query"] = rendered
		}
	}

	return result
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

func (r *Renderer) renderAgg(a lucene.Aggregation) (map[string]any, error) {
	var inner map[string]any
	var aggKey string

	switch v := a.(type) {
	case *lucene.TermsAgg:
		aggKey = "terms"
		inner = r.renderTermsAgg(v)
	case *lucene.HistogramAgg:
		aggKey = "histogram"
		inner = r.renderHistogramAgg(v)
	case *lucene.DateHistogramAgg:
		aggKey = "date_histogram"
		inner = r.renderDateHistogramAgg(v)
	case *lucene.RangeAgg:
		aggKey = "range"
		inner = r.renderRangeAgg(v)
	case *lucene.DateRangeAgg:
		aggKey = "date_range"
		inner = r.renderDateRangeAgg(v)
	case *lucene.FilterAgg:
		return r.renderFilterAgg(v)
	case *lucene.FiltersAgg:
		return r.renderFiltersAgg(v)
	case *lucene.NestedAgg:
		aggKey = "nested"
		inner = map[string]any{"path": v.Path()}
	case *lucene.MissingAgg:
		aggKey = "missing"
		inner = map[string]any{"field": v.Field()}
	case *lucene.AvgAgg:
		aggKey = "avg"
		inner = r.renderMetricAgg(v.Field(), v.MissingValue())
	case *lucene.SumAgg:
		aggKey = "sum"
		inner = r.renderMetricAgg(v.Field(), v.MissingValue())
	case *lucene.MinAgg:
		aggKey = "min"
		inner = r.renderMetricAgg(v.Field(), v.MissingValue())
	case *lucene.MaxAgg:
		aggKey = "max"
		inner = r.renderMetricAgg(v.Field(), v.MissingValue())
	case *lucene.CountAgg:
		aggKey = "value_count"
		inner = map[string]any{"field": v.Field()}
	case *lucene.CardinalityAgg:
		aggKey = "cardinality"
		inner = r.renderCardinalityAgg(v)
	case *lucene.StatsAgg:
		aggKey = "stats"
		inner = r.renderMetricAgg(v.Field(), v.MissingValue())
	case *lucene.ExtendedStatsAgg:
		aggKey = "extended_stats"
		inner = r.renderExtendedStatsAgg(v)
	case *lucene.PercentilesAgg:
		aggKey = "percentiles"
		inner = r.renderPercentilesAgg(v)
	case *lucene.TopHitsAgg:
		aggKey = "top_hits"
		inner = r.renderTopHitsAgg(v)
	case *lucene.AvgBucketAgg:
		aggKey = "avg_bucket"
		inner = r.renderPipelineAgg(&v.PipelineAgg)
	case *lucene.SumBucketAgg:
		aggKey = "sum_bucket"
		inner = r.renderPipelineAgg(&v.PipelineAgg)
	case *lucene.MaxBucketAgg:
		aggKey = "max_bucket"
		inner = r.renderPipelineAgg(&v.PipelineAgg)
	case *lucene.MinBucketAgg:
		aggKey = "min_bucket"
		inner = r.renderPipelineAgg(&v.PipelineAgg)
	case *lucene.DerivativeAgg:
		aggKey = "derivative"
		inner = r.renderDerivativeAgg(v)
	case *lucene.CumulativeSumAgg:
		aggKey = "cumulative_sum"
		inner = r.renderPipelineAgg(&v.PipelineAgg)
	case *lucene.MovingAvgAgg:
		aggKey = "moving_avg"
		inner = r.renderMovingAvgAgg(v)
	default:
		return nil, fmt.Errorf("unsupported aggregation type: %T", a)
	}

	result := map[string]any{aggKey: inner}

	if subs := a.SubAggs(); len(subs) > 0 {
		subResult := make(map[string]any)
		for _, sub := range subs {
			rendered, err := r.renderAgg(sub)
			if err != nil {
				return nil, err
			}
			subResult[sub.Name()] = rendered
		}
		result["aggs"] = subResult
	}

	return result, nil
}

func (r *Renderer) renderTermsAgg(a *lucene.TermsAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	if v := a.SizeValue(); v != nil {
		inner["size"] = *v
	}
	if v := a.MinDocCountValue(); v != nil {
		inner["min_doc_count"] = *v
	}
	if order := a.OrderValue(); len(order) > 0 {
		inner["order"] = order
	}
	return inner
}

func (r *Renderer) renderHistogramAgg(a *lucene.HistogramAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	if v := a.IntervalValue(); v != nil {
		inner["interval"] = *v
	}
	if v := a.OffsetValue(); v != nil {
		inner["offset"] = *v
	}
	if v := a.MinDocCountValue(); v != nil {
		inner["min_doc_count"] = *v
	}
	return inner
}

func (r *Renderer) renderDateHistogramAgg(a *lucene.DateHistogramAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	if v := a.CalendarIntervalValue(); v != nil {
		inner["calendar_interval"] = *v
	}
	if v := a.FixedIntervalValue(); v != nil {
		inner["fixed_interval"] = *v
	}
	if v := a.FormatValue(); v != nil {
		inner["format"] = *v
	}
	if v := a.TimeZoneValue(); v != nil {
		inner["time_zone"] = *v
	}
	if v := a.MinDocCountValue(); v != nil {
		inner["min_doc_count"] = *v
	}
	return inner
}

func (r *Renderer) renderRangeAgg(a *lucene.RangeAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	ranges := make([]map[string]any, 0, len(a.Ranges()))
	for _, rng := range a.Ranges() {
		rangeMap := make(map[string]any)
		if rng.Key != "" {
			rangeMap["key"] = rng.Key
		}
		if rng.From != nil {
			rangeMap["from"] = rng.From
		}
		if rng.To != nil {
			rangeMap["to"] = rng.To
		}
		ranges = append(ranges, rangeMap)
	}
	inner["ranges"] = ranges
	if v := a.KeyedValue(); v != nil {
		inner["keyed"] = *v
	}
	return inner
}

func (r *Renderer) renderDateRangeAgg(a *lucene.DateRangeAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	ranges := make([]map[string]any, 0, len(a.Ranges()))
	for _, rng := range a.Ranges() {
		rangeMap := make(map[string]any)
		if rng.Key != "" {
			rangeMap["key"] = rng.Key
		}
		if rng.From != nil {
			rangeMap["from"] = rng.From
		}
		if rng.To != nil {
			rangeMap["to"] = rng.To
		}
		ranges = append(ranges, rangeMap)
	}
	inner["ranges"] = ranges
	if v := a.FormatValue(); v != nil {
		inner["format"] = *v
	}
	if v := a.KeyedValue(); v != nil {
		inner["keyed"] = *v
	}
	return inner
}

func (r *Renderer) renderFiltersAgg(a *lucene.FiltersAgg) (map[string]any, error) {
	filters := make(map[string]any)
	for name, q := range a.Filters() {
		rendered, err := r.renderQuery(q)
		if err != nil {
			return nil, err
		}
		filters[name] = rendered
	}
	result := map[string]any{
		"filters": map[string]any{
			"filters": filters,
		},
	}
	if subs := a.SubAggs(); len(subs) > 0 {
		subResult := make(map[string]any)
		for _, sub := range subs {
			rendered, err := r.renderAgg(sub)
			if err != nil {
				return nil, err
			}
			subResult[sub.Name()] = rendered
		}
		result["aggs"] = subResult
	}
	return result, nil
}

func (r *Renderer) renderFilterAgg(a *lucene.FilterAgg) (map[string]any, error) {
	filter, err := r.renderQuery(a.FilterQuery())
	if err != nil {
		return nil, err
	}
	result := map[string]any{"filter": filter}
	if subs := a.SubAggs(); len(subs) > 0 {
		subResult := make(map[string]any)
		for _, sub := range subs {
			rendered, err := r.renderAgg(sub)
			if err != nil {
				return nil, err
			}
			subResult[sub.Name()] = rendered
		}
		result["aggs"] = subResult
	}
	return result, nil
}

func (r *Renderer) renderMetricAgg(field string, missing any) map[string]any {
	inner := map[string]any{"field": field}
	if missing != nil {
		inner["missing"] = missing
	}
	return inner
}

func (r *Renderer) renderCardinalityAgg(a *lucene.CardinalityAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	if v := a.PrecisionThresholdValue(); v != nil {
		inner["precision_threshold"] = *v
	}
	return inner
}

func (r *Renderer) renderExtendedStatsAgg(a *lucene.ExtendedStatsAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	if v := a.MissingValue(); v != nil {
		inner["missing"] = v
	}
	if v := a.SigmaValue(); v != nil {
		inner["sigma"] = *v
	}
	return inner
}

func (r *Renderer) renderPercentilesAgg(a *lucene.PercentilesAgg) map[string]any {
	inner := map[string]any{"field": a.Field()}
	if v := a.PercentsValue(); len(v) > 0 {
		inner["percents"] = v
	}
	if v := a.MissingValue(); v != nil {
		inner["missing"] = v
	}
	return inner
}

func (r *Renderer) renderTopHitsAgg(a *lucene.TopHitsAgg) map[string]any {
	inner := make(map[string]any)
	if v := a.SizeValue(); v != nil {
		inner["size"] = *v
	}
	if v := a.FromValue(); v != nil {
		inner["from"] = *v
	}
	if sorts := a.SortValue(); len(sorts) > 0 {
		sortList := make([]map[string]any, 0, len(sorts))
		for _, s := range sorts {
			sortList = append(sortList, map[string]any{s.Field: map[string]any{"order": s.Order}})
		}
		inner["sort"] = sortList
	}
	if src := a.SourceValue(); len(src) > 0 {
		inner["_source"] = src
	}
	return inner
}

func (r *Renderer) renderPipelineAgg(a *lucene.PipelineAgg) map[string]any {
	inner := map[string]any{"buckets_path": a.BucketsPath()}
	if v := a.GapPolicyValue(); v != nil {
		inner["gap_policy"] = *v
	}
	if v := a.FormatValue(); v != nil {
		inner["format"] = *v
	}
	return inner
}

func (r *Renderer) renderDerivativeAgg(a *lucene.DerivativeAgg) map[string]any {
	inner := r.renderPipelineAgg(&a.PipelineAgg)
	if v := a.UnitValue(); v != nil {
		inner["unit"] = *v
	}
	return inner
}

func (r *Renderer) renderMovingAvgAgg(a *lucene.MovingAvgAgg) map[string]any {
	inner := r.renderPipelineAgg(&a.PipelineAgg)
	if v := a.WindowValue(); v != nil {
		inner["window"] = *v
	}
	if v := a.ModelValue(); v != nil {
		inner["model"] = *v
	}
	if v := a.PredictValue(); v != nil {
		inner["predict"] = *v
	}
	return inner
}

// renderQuery converts a query to a map structure.
func (r *Renderer) renderQuery(q lucene.Query) (map[string]any, error) {
	switch v := q.(type) {
	// Term-level queries.
	case *lucene.TermQuery:
		return r.renderTerm(v), nil
	case *lucene.TermsQuery:
		return r.renderTerms(v), nil
	case *lucene.RangeQuery:
		return r.renderRange(v), nil
	case *lucene.ExistsQuery:
		return r.renderExists(v), nil
	case *lucene.IDsQuery:
		return r.renderIDs(v), nil
	case *lucene.PrefixQuery:
		return r.renderPrefix(v), nil
	case *lucene.WildcardQuery:
		return r.renderWildcard(v), nil
	case *lucene.RegexpQuery:
		return r.renderRegexp(v), nil
	case *lucene.FuzzyQuery:
		return r.renderFuzzy(v), nil
	// Full-text queries.
	case *lucene.MatchQuery:
		return r.renderMatch(v), nil
	case *lucene.MatchPhraseQuery:
		return r.renderMatchPhrase(v), nil
	case *lucene.MatchPhrasePrefixQuery:
		return r.renderMatchPhrasePrefix(v), nil
	case *lucene.MultiMatchQuery:
		return r.renderMultiMatch(v), nil
	case *lucene.QueryStringQuery:
		return r.renderQueryString(v), nil
	case *lucene.SimpleQueryStringQuery:
		return r.renderSimpleQueryString(v), nil
	// Compound queries.
	case *lucene.BoolQuery:
		return r.renderBool(v)
	case *lucene.MatchAllQuery:
		return r.renderMatchAll(v), nil
	case *lucene.MatchNoneQuery:
		return r.renderMatchNone(), nil
	case *lucene.BoostingQuery:
		return r.renderBoosting(v)
	case *lucene.DisMaxQuery:
		return r.renderDisMax(v)
	case *lucene.ConstantScoreQuery:
		return r.renderConstantScore(v)
	// Nested queries.
	case *lucene.NestedQuery:
		return r.renderNested(v)
	case *lucene.HasChildQuery:
		return r.renderHasChild(v)
	case *lucene.HasParentQuery:
		return r.renderHasParent(v)
	// Vector queries - OpenSearch uses different syntax.
	case *lucene.KnnQuery:
		return r.renderKnn(v)
	// Geo queries.
	case *lucene.GeoDistanceQuery:
		return r.renderGeoDistance(v), nil
	case *lucene.GeoBoundingBoxQuery:
		return r.renderGeoBoundingBox(v), nil
	default:
		return nil, fmt.Errorf("unsupported query type: %T", q)
	}
}

func (r *Renderer) renderTerm(q *lucene.TermQuery) map[string]any {
	inner := map[string]any{"value": q.Value()}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"term": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderTerms(q *lucene.TermsQuery) map[string]any {
	result := map[string]any{"terms": map[string]any{q.Field(): q.Values()}}
	if b := q.BoostValue(); b != nil {
		result["terms"].(map[string]any)["boost"] = *b
	}
	return result
}

func (r *Renderer) renderRange(q *lucene.RangeQuery) map[string]any {
	inner := make(map[string]any)
	if v := q.GtValue(); v != nil {
		inner["gt"] = v
	}
	if v := q.GteValue(); v != nil {
		inner["gte"] = v
	}
	if v := q.LtValue(); v != nil {
		inner["lt"] = v
	}
	if v := q.LteValue(); v != nil {
		inner["lte"] = v
	}
	if f := q.FormatValue(); f != nil {
		inner["format"] = *f
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"range": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderExists(q *lucene.ExistsQuery) map[string]any {
	return map[string]any{"exists": map[string]any{"field": q.Field()}}
}

func (r *Renderer) renderIDs(q *lucene.IDsQuery) map[string]any {
	return map[string]any{"ids": map[string]any{"values": q.IDValues()}}
}

func (r *Renderer) renderPrefix(q *lucene.PrefixQuery) map[string]any {
	inner := map[string]any{"value": q.Value()}
	if v := q.RewriteValue(); v != nil {
		inner["rewrite"] = *v
	}
	if v := q.CaseInsensitiveValue(); v != nil {
		inner["case_insensitive"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"prefix": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderWildcard(q *lucene.WildcardQuery) map[string]any {
	inner := map[string]any{"value": q.Value()}
	if v := q.RewriteValue(); v != nil {
		inner["rewrite"] = *v
	}
	if v := q.CaseInsensitiveValue(); v != nil {
		inner["case_insensitive"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"wildcard": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderRegexp(q *lucene.RegexpQuery) map[string]any {
	inner := map[string]any{"value": q.Value()}
	if v := q.FlagsValue(); v != nil {
		inner["flags"] = *v
	}
	if v := q.RewriteValue(); v != nil {
		inner["rewrite"] = *v
	}
	if v := q.CaseInsensitiveValue(); v != nil {
		inner["case_insensitive"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"regexp": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderFuzzy(q *lucene.FuzzyQuery) map[string]any {
	inner := map[string]any{"value": q.Value()}
	if v := q.FuzzinessValue(); v != nil {
		inner["fuzziness"] = *v
	}
	if v := q.PrefixLengthValue(); v != nil {
		inner["prefix_length"] = *v
	}
	if v := q.MaxExpansionsValue(); v != nil {
		inner["max_expansions"] = *v
	}
	if v := q.TranspositionsValue(); v != nil {
		inner["transpositions"] = *v
	}
	if v := q.RewriteValue(); v != nil {
		inner["rewrite"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"fuzzy": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderMatch(q *lucene.MatchQuery) map[string]any {
	inner := map[string]any{"query": q.Value()}
	if f := q.FuzzinessValue(); f != nil {
		inner["fuzziness"] = *f
	}
	if o := q.OperatorValue(); o != nil {
		inner["operator"] = *o
	}
	if a := q.AnalyzerValue(); a != nil {
		inner["analyzer"] = *a
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"match": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderMatchPhrase(q *lucene.MatchPhraseQuery) map[string]any {
	inner := map[string]any{"query": q.Value()}
	if s := q.SlopValue(); s != nil {
		inner["slop"] = *s
	}
	if a := q.AnalyzerValue(); a != nil {
		inner["analyzer"] = *a
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"match_phrase": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderMatchPhrasePrefix(q *lucene.MatchPhrasePrefixQuery) map[string]any {
	inner := map[string]any{"query": q.Value()}
	if s := q.SlopValue(); s != nil {
		inner["slop"] = *s
	}
	if m := q.MaxExpansionsValue(); m != nil {
		inner["max_expansions"] = *m
	}
	if a := q.AnalyzerValue(); a != nil {
		inner["analyzer"] = *a
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"match_phrase_prefix": map[string]any{q.Field(): inner}}
}

func (r *Renderer) renderMultiMatch(q *lucene.MultiMatchQuery) map[string]any {
	inner := map[string]any{"query": q.Value(), "fields": q.Fields()}
	if t := q.TypeValue(); t != nil {
		inner["type"] = *t
	}
	if t := q.TieBreakerValue(); t != nil {
		inner["tie_breaker"] = *t
	}
	if f := q.FuzzinessValue(); f != nil {
		inner["fuzziness"] = *f
	}
	if o := q.OperatorValue(); o != nil {
		inner["operator"] = *o
	}
	if a := q.AnalyzerValue(); a != nil {
		inner["analyzer"] = *a
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"multi_match": inner}
}

func (r *Renderer) renderQueryString(q *lucene.QueryStringQuery) map[string]any {
	inner := map[string]any{"query": q.Value()}
	if v := q.DefaultFieldValue(); v != nil {
		inner["default_field"] = *v
	}
	if v := q.DefaultOperatorValue(); v != nil {
		inner["default_operator"] = *v
	}
	if v := q.AnalyzerValue(); v != nil {
		inner["analyzer"] = *v
	}
	if v := q.AllowLeadingWildcardValue(); v != nil {
		inner["allow_leading_wildcard"] = *v
	}
	if v := q.FuzzinessValue(); v != nil {
		inner["fuzziness"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"query_string": inner}
}

func (r *Renderer) renderSimpleQueryString(q *lucene.SimpleQueryStringQuery) map[string]any {
	inner := map[string]any{"query": q.Value()}
	if fields := q.FieldsValue(); len(fields) > 0 {
		inner["fields"] = fields
	}
	if v := q.DefaultOperatorValue(); v != nil {
		inner["default_operator"] = *v
	}
	if v := q.AnalyzerValue(); v != nil {
		inner["analyzer"] = *v
	}
	if v := q.FlagsValue(); v != nil {
		inner["flags"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"simple_query_string": inner}
}

func (r *Renderer) renderBool(q *lucene.BoolQuery) (map[string]any, error) {
	inner := make(map[string]any)
	if must := q.MustClauses(); len(must) > 0 {
		clauses, err := r.renderClauses(must)
		if err != nil {
			return nil, err
		}
		inner["must"] = clauses
	}
	if should := q.ShouldClauses(); len(should) > 0 {
		clauses, err := r.renderClauses(should)
		if err != nil {
			return nil, err
		}
		inner["should"] = clauses
	}
	if mustNot := q.MustNotClauses(); len(mustNot) > 0 {
		clauses, err := r.renderClauses(mustNot)
		if err != nil {
			return nil, err
		}
		inner["must_not"] = clauses
	}
	if filter := q.FilterClauses(); len(filter) > 0 {
		clauses, err := r.renderClauses(filter)
		if err != nil {
			return nil, err
		}
		inner["filter"] = clauses
	}
	if m := q.MinimumShouldMatchValue(); m != nil {
		inner["minimum_should_match"] = *m
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"bool": inner}, nil
}

func (r *Renderer) renderClauses(queries []lucene.Query) ([]map[string]any, error) {
	result := make([]map[string]any, 0, len(queries))
	for _, q := range queries {
		rendered, err := r.renderQuery(q)
		if err != nil {
			return nil, err
		}
		result = append(result, rendered)
	}
	return result, nil
}

func (r *Renderer) renderMatchAll(q *lucene.MatchAllQuery) map[string]any {
	inner := make(map[string]any)
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"match_all": inner}
}

func (r *Renderer) renderMatchNone() map[string]any {
	return map[string]any{"match_none": map[string]any{}}
}

func (r *Renderer) renderBoosting(q *lucene.BoostingQuery) (map[string]any, error) {
	inner := make(map[string]any)
	if pos := q.PositiveQuery(); pos != nil {
		rendered, err := r.renderQuery(pos)
		if err != nil {
			return nil, err
		}
		inner["positive"] = rendered
	}
	if neg := q.NegativeQuery(); neg != nil {
		rendered, err := r.renderQuery(neg)
		if err != nil {
			return nil, err
		}
		inner["negative"] = rendered
	}
	if v := q.NegativeBoostValue(); v != nil {
		inner["negative_boost"] = *v
	}
	return map[string]any{"boosting": inner}, nil
}

func (r *Renderer) renderDisMax(q *lucene.DisMaxQuery) (map[string]any, error) {
	inner := make(map[string]any)
	if queries := q.Queries(); len(queries) > 0 {
		rendered, err := r.renderClauses(queries)
		if err != nil {
			return nil, err
		}
		inner["queries"] = rendered
	}
	if v := q.TieBreakerValue(); v != nil {
		inner["tie_breaker"] = *v
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"dis_max": inner}, nil
}

func (r *Renderer) renderConstantScore(q *lucene.ConstantScoreQuery) (map[string]any, error) {
	inner := make(map[string]any)
	if filter := q.FilterQuery(); filter != nil {
		rendered, err := r.renderQuery(filter)
		if err != nil {
			return nil, err
		}
		inner["filter"] = rendered
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"constant_score": inner}, nil
}

func (r *Renderer) renderNested(q *lucene.NestedQuery) (map[string]any, error) {
	inner := map[string]any{"path": q.Path()}
	if iq := q.InnerQuery(); iq != nil {
		rendered, err := r.renderQuery(iq)
		if err != nil {
			return nil, err
		}
		inner["query"] = rendered
	}
	if v := q.ScoreModeValue(); v != nil {
		inner["score_mode"] = *v
	}
	if v := q.IgnoreUnmappedValue(); v != nil {
		inner["ignore_unmapped"] = *v
	}
	return map[string]any{"nested": inner}, nil
}

func (r *Renderer) renderHasChild(q *lucene.HasChildQuery) (map[string]any, error) {
	inner := map[string]any{"type": q.ChildType()}
	if iq := q.InnerQuery(); iq != nil {
		rendered, err := r.renderQuery(iq)
		if err != nil {
			return nil, err
		}
		inner["query"] = rendered
	}
	if v := q.ScoreModeValue(); v != nil {
		inner["score_mode"] = *v
	}
	if v := q.MinChildrenValue(); v != nil {
		inner["min_children"] = *v
	}
	if v := q.MaxChildrenValue(); v != nil {
		inner["max_children"] = *v
	}
	if v := q.IgnoreUnmappedValue(); v != nil {
		inner["ignore_unmapped"] = *v
	}
	return map[string]any{"has_child": inner}, nil
}

func (r *Renderer) renderHasParent(q *lucene.HasParentQuery) (map[string]any, error) {
	inner := map[string]any{"parent_type": q.ParentType()}
	if iq := q.InnerQuery(); iq != nil {
		rendered, err := r.renderQuery(iq)
		if err != nil {
			return nil, err
		}
		inner["query"] = rendered
	}
	if v := q.ScoreValue(); v != nil {
		inner["score"] = *v
	}
	if v := q.IgnoreUnmappedValue(); v != nil {
		inner["ignore_unmapped"] = *v
	}
	return map[string]any{"has_parent": inner}, nil
}

// renderKnn renders kNN for OpenSearch.
// OpenSearch 2.x uses knn within the query body (neural search plugin).
func (r *Renderer) renderKnn(q *lucene.KnnQuery) (map[string]any, error) {
	inner := map[string]any{
		q.Field(): map[string]any{
			"vector": q.Vector(),
			"k":      q.KValue(),
		},
	}
	if filter := q.FilterQuery(); filter != nil {
		rendered, err := r.renderQuery(filter)
		if err != nil {
			return nil, err
		}
		inner["filter"] = rendered
	}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"knn": inner}, nil
}

func (r *Renderer) renderGeoDistance(q *lucene.GeoDistanceQuery) map[string]any {
	inner := map[string]any{
		q.Field(): map[string]any{"lat": q.Lat(), "lon": q.Lon()},
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

func (r *Renderer) renderGeoBoundingBox(q *lucene.GeoBoundingBoxQuery) map[string]any {
	fieldInner := make(map[string]any)
	if lat := q.TopLeftLat(); lat != nil {
		if lon := q.TopLeftLon(); lon != nil {
			fieldInner["top_left"] = map[string]any{"lat": *lat, "lon": *lon}
		}
	}
	if lat := q.BottomRightLat(); lat != nil {
		if lon := q.BottomRightLon(); lon != nil {
			fieldInner["bottom_right"] = map[string]any{"lat": *lat, "lon": *lon}
		}
	}
	inner := map[string]any{q.Field(): fieldInner}
	if b := q.BoostValue(); b != nil {
		inner["boost"] = *b
	}
	return map[string]any{"geo_bounding_box": inner}
}
