# Search Query Builder Design

> `lucene` — provider-agnostic query builder for OpenSearch/Elasticsearch

## Overview

A type-safe, schema-validated query builder following zoobzio patterns (vecna). Generates JSON Query DSL for OpenSearch and Elasticsearch.

**Dependencies:** `github.com/zoobzio/sentinel` for schema extraction.

---

## 1. Core Types

### Query Interface

All query types implement a common interface:

```go
type Query interface {
    Op() Op
    Err() error
    query()  // sealed interface
}
```

### Base Query Node

```go
type query struct {
    op       Op
    field    string
    value    any
    children []Query
    err      error
}

func (q *query) Op() Op     { return q.op }
func (q *query) Err() error { return q.err }
func (q *query) query()     {}
```

### Query-Specific Types

Each query type embeds the base and adds type-safe param methods:

```go
// Full-text query with chained params
type MatchQuery struct {
    query
    fuzziness string
    operator  string
    analyzer  string
}

func (q *MatchQuery) Fuzziness(f string) *MatchQuery { q.fuzziness = f; return q }
func (q *MatchQuery) Operator(o string) *MatchQuery  { q.operator = o; return q }
func (q *MatchQuery) Analyzer(a string) *MatchQuery  { q.analyzer = a; return q }

// Range query with bound methods
type RangeQuery struct {
    query
    gt, gte, lt, lte any
    format           string
}

func (q *RangeQuery) Gt(v any) *RangeQuery     { q.gt = v; return q }
func (q *RangeQuery) Gte(v any) *RangeQuery    { q.gte = v; return q }
func (q *RangeQuery) Lt(v any) *RangeQuery     { q.lt = v; return q }
func (q *RangeQuery) Lte(v any) *RangeQuery    { q.lte = v; return q }
func (q *RangeQuery) Format(f string) *RangeQuery { q.format = f; return q }

// Bool query with clause methods
type BoolQuery struct {
    query
    must, should, mustNot, filter []Query
    minimumShouldMatch             int
}

func (q *BoolQuery) Must(queries ...Query) *BoolQuery    { q.must = append(q.must, queries...); return q }
func (q *BoolQuery) Should(queries ...Query) *BoolQuery  { q.should = append(q.should, queries...); return q }
func (q *BoolQuery) MustNot(queries ...Query) *BoolQuery { q.mustNot = append(q.mustNot, queries...); return q }
func (q *BoolQuery) Filter(queries ...Query) *BoolQuery  { q.filter = append(q.filter, queries...); return q }
func (q *BoolQuery) MinimumShouldMatch(n int) *BoolQuery { q.minimumShouldMatch = n; return q }

// Vector query with kNN params
type KnnQuery struct {
    query
    vector        []float32
    k             int
    numCandidates int
}

func (q *KnnQuery) K(k int) *KnnQuery             { q.k = k; return q }
func (q *KnnQuery) NumCandidates(n int) *KnnQuery { q.numCandidates = n; return q }
```

### Operators (Op enum)

```go
type Op uint8

const (
    // Full-text
    OpMatch Op = iota
    OpMatchPhrase
    OpMatchPhrasePrefix
    OpMultiMatch
    OpQueryString
    OpSimpleQueryString

    // Term-level
    OpTerm
    OpTerms
    OpRange
    OpPrefix
    OpWildcard
    OpRegexp
    OpFuzzy
    OpExists
    OpIds

    // Compound
    OpBool
    OpBoosting
    OpConstantScore
    OpDisMax
    OpFunctionScore

    // Special
    OpMatchAll
    OpMatchNone
    OpNested
    OpHasChild
    OpHasParent

    // Vector
    OpKnn

    // Geo
    OpGeoDistance
    OpGeoBoundingBox

    // Logical (internal)
    OpMust
    OpShould
    OpMustNot
    OpFilter
)
```

---

## 2. Schema Extraction

Following the vecna pattern, schema extraction happens once at builder initialization:

```go
type Builder[T any] struct {
    spec   *Spec
    fields map[string]*FieldSpec  // O(1) lookup
}

type Spec struct {
    Fields []FieldSpec
}

type FieldSpec struct {
    Name string    // resolved name (json tag or Go name)
    Type string    // Go type string
    Kind FieldKind // categorized type
}

type FieldKind uint8

const (
    KindUnknown FieldKind = iota
    KindString
    KindInt
    KindFloat
    KindBool
    KindTime
    KindSlice
    KindVector  // []float32, []float64 for embeddings
)
```

### Initialization

```go
func New[T any]() (*Builder[T], error) {
    sentinel.Tag("json")
    sentinel.Tag("lucene")  // custom tag for field options

    metadata, err := sentinel.TryInspect[T]()
    if err != nil {
        return nil, err
    }

    spec := buildSpec(metadata)
    fields := make(map[string]*FieldSpec, len(spec.Fields))
    for i := range spec.Fields {
        fields[spec.Fields[i].Name] = &spec.Fields[i]
    }

    return &Builder[T]{spec: spec, fields: fields}, nil
}
```

### Field Resolution

```go
func (b *Builder[T]) resolveField(name string) (*FieldSpec, error) {
    if spec, ok := b.fields[name]; ok {
        return spec, nil
    }
    return nil, fmt.Errorf("unknown field: %s", name)
}
```

---

## 3. Query Categories

### 3.1 Full-Text Queries

| Query | Purpose | Key Params |
|-------|---------|------------|
| `Match` | Analyzed text search | fuzziness, operator, analyzer |
| `MatchPhrase` | Exact phrase | slop, analyzer |
| `MatchPhrasePrefix` | Autocomplete | max_expansions |
| `MultiMatch` | Search multiple fields | type, tie_breaker, fields |
| `QueryString` | Lucene syntax | default_field, default_operator |
| `SimpleQueryString` | User-friendly syntax | flags, default_operator |

### 3.2 Term-Level Queries

| Query | Purpose | Key Params |
|-------|---------|------------|
| `Term` | Exact value | boost |
| `Terms` | Multiple exact values | boost |
| `Range` | Numeric/date range | gt, gte, lt, lte, format |
| `Prefix` | Prefix match | rewrite |
| `Wildcard` | Pattern match (* ?) | case_insensitive |
| `Regexp` | Regex match | flags |
| `Fuzzy` | Edit distance | fuzziness, prefix_length |
| `Exists` | Field exists | - |
| `Ids` | Document IDs | - |

### 3.3 Compound Queries

| Query | Purpose | Clauses |
|-------|---------|---------|
| `Bool` | Combine queries | must, should, must_not, filter |
| `Boosting` | Demote matches | positive, negative, negative_boost |
| `ConstantScore` | Fixed score | filter, boost |
| `DisMax` | Best match wins | queries, tie_breaker |
| `FunctionScore` | Custom scoring | functions, score_mode, boost_mode |

### 3.4 Specialized Queries

| Query | Purpose |
|-------|---------|
| `Nested` | Query nested objects |
| `HasChild` | Parent by child match |
| `HasParent` | Child by parent match |
| `MatchAll` | Match everything |
| `MatchNone` | Match nothing |

### 3.5 Vector Queries (kNN)

| Query | Purpose | Key Params |
|-------|---------|------------|
| `Knn` | k-nearest neighbors | k, num_candidates, vector field |

**Provider differences:**
- Elasticsearch 8.x: `knn` at top level of search request
- OpenSearch 2.x: `knn` within query body

### 3.6 Geo Queries

| Query | Purpose | Key Params |
|-------|---------|------------|
| `GeoDistance` | Within radius | distance, distance_type |
| `GeoBoundingBox` | Within box | top_left, bottom_right |

---

## 4. Aggregations

### 4.1 Bucket Aggregations

| Agg | Purpose | Key Params |
|-----|---------|------------|
| `Terms` | Group by field value | size, order, min_doc_count |
| `Histogram` | Numeric buckets | interval, offset, min_doc_count |
| `DateHistogram` | Time buckets | calendar_interval, fixed_interval |
| `Range` | Custom ranges | ranges[] |
| `DateRange` | Date ranges | ranges[], format |
| `Filter` | Single filter bucket | filter query |
| `Filters` | Named filter buckets | filters map |
| `Nested` | Nested doc aggregation | path |
| `Missing` | Docs missing field | - |

### 4.2 Metric Aggregations

| Agg | Purpose |
|-----|---------|
| `Avg` | Average value |
| `Sum` | Sum of values |
| `Min` | Minimum value |
| `Max` | Maximum value |
| `Count` | Value count |
| `Cardinality` | Distinct count |
| `Stats` | All basic stats |
| `ExtendedStats` | Stats + variance |
| `Percentiles` | Percentile values |
| `TopHits` | Top matching docs |

### 4.3 Pipeline Aggregations

| Agg | Purpose |
|-----|---------|
| `AvgBucket` | Avg of bucket values |
| `SumBucket` | Sum of bucket values |
| `MaxBucket` | Max bucket value |
| `MinBucket` | Min bucket value |
| `Derivative` | Rate of change |
| `CumulativeSum` | Running total |
| `MovingAvg` | Smoothed average |

---

## 5. Builder API

### 5.1 Initialization

```go
b, err := lucene.New[Product]()
if err != nil {
    // struct inspection failed
}
```

### 5.2 Query Builders

Each builder method returns a query-specific type with only relevant param methods:

```go
// Full-text — returns *MatchQuery with Fuzziness(), Operator(), Analyzer()
b.Match("title", "search term").Fuzziness("AUTO").Operator("and")
b.MatchPhrase("title", "exact phrase").Slop(2)
b.MultiMatch("search term", "title", "description").Type("best_fields").TieBreaker(0.3)

// Term-level — returns *TermQuery, *RangeQuery, etc.
b.Term("status", "active").Boost(1.5)
b.Terms("category", "shoes", "boots", "sandals")
b.Range("price").Gte(10).Lt(100)
b.Range("created_at").Gte("2024-01-01").Format("yyyy-MM-dd")
b.Prefix("sku", "ABC").Rewrite("constant_score")
b.Wildcard("name", "jo*n").CaseInsensitive(true)
b.Regexp("code", "[A-Z]{3}[0-9]+").Flags("ALL")
b.Fuzzy("name", "john").Fuzziness("AUTO").PrefixLength(2)
b.Exists("email")
b.Ids("doc1", "doc2", "doc3")

// Compound — returns *BoolQuery, *BoostingQuery, etc.
b.Bool().
    Must(q1, q2).
    Should(q3).
    Filter(q4).
    MustNot(q5).
    MinimumShouldMatch(1)

b.Boosting().
    Positive(q1).
    Negative(q2).
    NegativeBoost(0.5)

b.DisMax(q1, q2, q3).TieBreaker(0.3)

b.ConstantScore(q).Boost(1.5)

// Nested — returns *NestedQuery
b.Nested("comments", q).ScoreMode("avg").IgnoreUnmapped(true)

// Vector — returns *KnnQuery
b.Knn("embedding", vector).K(10).NumCandidates(100)

// Geo — returns *GeoDistanceQuery, *GeoBoundingBoxQuery
b.GeoDistance("location", 40.73, -73.93).Distance("10km").DistanceType("arc")
b.GeoBoundingBox("location").TopLeft(40.73, -74.1).BottomRight(40.01, -71.12)
```

### 5.3 Aggregation Builders

```go
// Bucket
b.TermsAgg("by_brand", "brand").Size(10).MinDocCount(1)
b.DateHistogram("by_month", "created_at").CalendarInterval("month")
b.RangeAgg("price_ranges", "price").
    AddRange(0, 50).
    AddRange(50, 100).
    AddRange(100, nil)

// Metric
b.Avg("avg_price", "price")
b.Stats("price_stats", "price")
b.TopHits("top_products").Size(3).Sort("score", "desc")

// Nested aggregations
b.TermsAgg("by_category", "category").
    SubAgg(b.Avg("avg_price", "price")).
    SubAgg(b.Max("max_price", "price"))
```

### 5.4 Search Request Builder

```go
search := lucene.Search().
    Query(query).
    Aggs(agg1, agg2).
    Size(20).
    From(0).
    Sort(b.SortField("price", "asc")).
    Source("title", "price", "brand").
    Highlight(b.Highlight("title", "description")).
    TrackTotalHits(true)
```

### 5.5 Error Handling

Errors are deferred until explicitly checked (vecna pattern):

```go
q := b.Term("invalid_field", "value")  // stores error, does not panic
q := b.Bool().Must(q1, q2)             // propagates child errors

// Check before rendering
if err := q.Err(); err != nil {
    return err
}

// Or let Render check
json, err := renderer.Render(q)  // returns error if q.Err() != nil
```

---

## 6. Renderer Interface

```go
type Renderer interface {
    Render(search *Search) ([]byte, error)
    RenderQuery(query *Query) ([]byte, error)
    RenderAggs(aggs []*Aggregation) ([]byte, error)
}
```

### Provider Implementations

```go
// Elasticsearch renderer
es := elasticsearch.NewRenderer(elasticsearch.V8)

// OpenSearch renderer
os := opensearch.NewRenderer(opensearch.V2)
```

### Provider Differences

| Feature | Elasticsearch | OpenSearch |
|---------|---------------|------------|
| kNN syntax | `knn` top-level | `knn` in query |
| Script syntax | Painless | Painless (compatible) |
| Version compat | 7.x, 8.x | 1.x, 2.x |

---

## 7. File Structure

```
lucene/
├── api.go                 # Query, Op, Search types
├── builder.go             # Builder[T], schema extraction
├── spec.go                # Spec, FieldSpec, FieldKind
├── search.go              # Search request builder
├── fulltext.go            # Match, MultiMatch, QueryString
├── term.go                # Term, Range, Prefix, Wildcard, etc.
├── compound.go            # Bool, Boosting, DisMax, ConstantScore
├── nested.go              # Nested, HasChild, HasParent
├── vector.go              # Knn
├── geo.go                 # GeoDistance, GeoBoundingBox
├── aggregation.go         # Agg interface, bucket/metric aggs
├── pipeline.go            # Pipeline aggregations
├── sort.go                # Sort builders
├── highlight.go           # Highlight builders
├── render.go              # Renderer interface
│
├── elasticsearch/
│   └── renderer.go        # ES-specific rendering
│
├── opensearch/
│   └── renderer.go        # OS-specific rendering
│
├── testing/
│   ├── helpers.go
│   ├── helpers_test.go
│   ├── integration/
│   └── benchmarks/
│
└── docs/
```

---

## 8. Implementation Phases

### Phase 1: Core

- [ ] Op enum with all operators
- [ ] Query interface and base query struct
- [ ] Spec, FieldSpec, FieldKind types
- [ ] Builder[T] with sentinel integration
- [ ] Field resolution and validation
- [ ] Renderer interface

### Phase 2: Basic Queries

- [ ] Term-level queries (Term, Terms, Range, Exists, Ids)
- [ ] Full-text queries (Match, MatchPhrase, MultiMatch)
- [ ] Compound queries (Bool)
- [ ] Special queries (MatchAll, MatchNone)
- [ ] Elasticsearch renderer (basic)

### Phase 3: Advanced Queries

- [ ] Remaining term-level (Prefix, Wildcard, Regexp, Fuzzy)
- [ ] Remaining full-text (QueryString, SimpleQueryString)
- [ ] Remaining compound (Boosting, DisMax, ConstantScore)
- [ ] Nested queries (Nested, HasChild, HasParent)
- [ ] Vector queries (Knn)
- [ ] Geo queries (GeoDistance, GeoBoundingBox)
- [ ] OpenSearch renderer

### Phase 4: Aggregations

- [ ] Bucket aggregations (Terms, Histogram, DateHistogram, Range)
- [ ] Metric aggregations (Avg, Sum, Min, Max, Stats, Cardinality)
- [ ] Nested aggregations
- [ ] Pipeline aggregations

### Phase 5: Search Features

- [ ] Search request builder (size, from, sort, source)
- [ ] Highlighting
- [ ] Track total hits
- [ ] Source filtering

---

## 9. Decisions

| Topic | Decision | Rationale |
|-------|----------|-----------|
| Schema extraction | sentinel | Established pattern, zero-cost after init |
| Error handling | Deferred (vecna pattern) | Enables fluent chaining |
| Query types | Per-query structs (soy pattern) | Type-safe params, IDE autocomplete, only relevant methods |
| Vector search | Included | OS/ES commonly used as vector store |
| Geo queries | Included | Common use case |
| Scripts | Excluded | Complex, provider-specific, defer to raw JSON |
| Suggesters | Excluded | Specialized, low priority |
| Percolate | Excluded | Niche feature |
