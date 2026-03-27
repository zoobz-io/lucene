# lucene

[![CI Status](https://github.com/zoobz-io/lucene/workflows/CI/badge.svg)](https://github.com/zoobz-io/lucene/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/zoobz-io/lucene/graph/badge.svg?branch=main)](https://codecov.io/gh/zoobz-io/lucene)
[![Go Report Card](https://goreportcard.com/badge/github.com/zoobz-io/lucene)](https://goreportcard.com/report/github.com/zoobz-io/lucene)
[![CodeQL](https://github.com/zoobz-io/lucene/workflows/CodeQL/badge.svg)](https://github.com/zoobz-io/lucene/security/code-scanning)
[![Go Reference](https://pkg.go.dev/badge/github.com/zoobz-io/lucene.svg)](https://pkg.go.dev/github.com/zoobz-io/lucene)
[![License](https://img.shields.io/github/license/zoobz-io/lucene)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/zoobz-io/lucene)](go.mod)
[![Release](https://img.shields.io/github/v/release/zoobz-io/lucene)](https://github.com/zoobz-io/lucene/releases)

Type-safe search queries for Elasticsearch and OpenSearch. Compile-time field validation ensures your queries reference fields that actually exist.

## Your Struct, Your Schema

```go
type Product struct {
    Title    string  `json:"title"`
    Category string  `json:"category"`
    Price    float64 `json:"price"`
}

b := lucene.New[Product]()

// This compiles - "title" exists
query := b.Match("title", "laptop")

// This fails at build time - "titl" doesn't exist
query := b.Match("titl", "laptop")  // unknown field: titl
```

Your Go struct becomes the source of truth. No more runtime surprises from typos in field names.

## Install

```bash
go get github.com/zoobz-io/lucene
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/zoobz-io/lucene"
    "github.com/zoobz-io/lucene/elasticsearch"
)

type Article struct {
    Title     string   `json:"title"`
    Body      string   `json:"body"`
    Author    string   `json:"author"`
    Published string   `json:"published"`
    Views     int      `json:"views"`
}

func main() {
    // Create a type-safe builder
    b := lucene.New[Article]()

    // Build a search request
    search := lucene.NewSearch().
        Query(
            b.Bool().
                Must(b.Match("title", "golang")).
                Filter(b.Range("views").Gte(1000)).
                Should(b.Term("author", "gopher")),
        ).
        Aggs(b.TermsAgg("by_author", "author").Size(10)).
        Size(20)

    // Render to Elasticsearch JSON
    renderer := elasticsearch.NewRenderer(elasticsearch.V8)
    json, err := renderer.Render(search)
    if err != nil {
        panic(err)
    }

    fmt.Println(string(json))
}
```

## Capabilities

| Feature | Description |
|---------|-------------|
| **Full-text queries** | Match, match_phrase, multi_match, query_string |
| **Term-level queries** | Term, terms, range, prefix, wildcard, regexp, fuzzy, exists |
| **Compound queries** | Bool, boosting, constant_score, dis_max |
| **Joining queries** | Nested, has_child, has_parent |
| **Geo queries** | Geo_distance, geo_bounding_box |
| **Vector search** | k-NN with filter support |
| **Aggregations** | Terms, histogram, date_histogram, range, metrics, pipeline |
| **Search features** | Sort, pagination, source filtering, highlighting |

## Why lucene?

- **Catch errors early** - Field validation happens when you build the query, not when Elasticsearch rejects it
- **Chain naturally** - Fluent builder methods return typed results; check `.Err()` once at the end
- **Target both engines** - Same query AST renders to Elasticsearch or OpenSearch JSON
- **Cover the DSL** - Bool queries, aggregations, geo, vectors, highlights - it's all there

## The Zoobzio Ecosystem

lucene works alongside other zoobzio packages:

| Package | Purpose |
|---------|---------|
| [sentinel](https://github.com/zoobz-io/sentinel) | Struct metadata extraction (powers lucene's schema) |

## Documentation

**Learn**
- [Overview](docs/learn/overview.md) - Core concepts and architecture
- [Quickstart](docs/learn/quickstart.md) - Get running in 5 minutes

**Guides**
- [Query Building](docs/guides/queries.md) - All query types explained
- [Aggregations](docs/guides/aggregations.md) - Bucket and metric aggregations
- [Rendering](docs/guides/rendering.md) - ES vs OpenSearch output

**Reference**
- [API Reference](https://pkg.go.dev/github.com/zoobz-io/lucene) - Full package documentation

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](LICENSE)
