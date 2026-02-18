# Testing Infrastructure

This directory contains shared testing utilities for the lucene package.

## Structure

```
testing/
├── README.md           # This file
├── helpers.go          # Domain-specific test helpers
├── helpers_test.go     # Tests for helpers themselves
├── integration/        # Integration tests
│   └── README.md
└── benchmarks/         # Performance benchmarks
    └── README.md
```

## Conventions

### Test Helpers

All helpers in `helpers.go`:

- Call `t.Helper()` as the first line
- Accept `*testing.T` as the first parameter
- Are domain-specific to lucene (query building, AST validation)
- Use the `testing` build tag

### Running Tests

```bash
# All tests
make test

# Unit tests only (fast)
make test-unit

# Integration tests
make test-integration

# Benchmarks
make test-bench
```

### Coverage

Target coverage levels:

- Project: 70% minimum
- New code (patch): 80% minimum

Generate coverage report:

```bash
make coverage
```
