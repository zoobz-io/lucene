# Integration Tests

Integration tests for lucene that verify behavior against real OpenSearch/Elasticsearch instances.

## Running

```bash
make test-integration
```

## Requirements

Integration tests may require:

- Running OpenSearch or Elasticsearch instance
- Environment variables for connection details

## Writing Integration Tests

1. Place tests in this directory
2. Use the `integration` build tag
3. Skip tests gracefully when dependencies are unavailable

```go
//go:build integration

package integration

import (
    "os"
    "testing"
)

func TestQueryExecution(t *testing.T) {
    endpoint := os.Getenv("OPENSEARCH_ENDPOINT")
    if endpoint == "" {
        t.Skip("OPENSEARCH_ENDPOINT not set")
    }
    // Test implementation
}
```
