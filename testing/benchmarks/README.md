# Benchmarks

Performance benchmarks for lucene query building operations.

## Running

```bash
make test-bench
```

## Writing Benchmarks

Place benchmark files in this directory with the `_test.go` suffix.

```go
//go:build testing

package benchmarks

import "testing"

func BenchmarkQueryBuild(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Build query
    }
}

func BenchmarkQueryBuild_Complex(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Build complex query with nested clauses
    }
}
```

## Benchmark Guidelines

- Use `b.ReportAllocs()` for memory allocation tracking
- Name benchmarks descriptively: `BenchmarkOperation_Variant`
- Include both simple and complex scenarios
- Document expected performance characteristics
