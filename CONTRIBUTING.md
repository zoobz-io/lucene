# Contributing to lucene

Thank you for your interest in contributing to lucene.

## Development Setup

```bash
# Clone the repository
git clone https://github.com/zoobzio/lucene.git
cd lucene

# Install development tools
make install-tools

# Install git hooks
make install-hooks

# Run all checks
make check
```

## Available Commands

Run `make help` to see all available commands:

```bash
make help
```

## Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run `make check` to verify your changes
5. Commit your changes (`git commit -m 'feat: add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## Commit Messages

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` — New features
- `fix:` — Bug fixes
- `perf:` — Performance improvements
- `docs:` — Documentation changes
- `test:` — Test additions or modifications
- `refactor:` — Code refactoring
- `chore:` — Maintenance tasks

## Code Quality

All contributions must pass:

- `make lint` — Linting with golangci-lint
- `make test` — Unit tests with race detector
- `make security` — Security scanning with gosec

## Testing

- Maintain 1:1 relationship between source and test files
- New code must have 80% coverage minimum
- Place test helpers in `testing/helpers.go`
- Integration tests go in `testing/integration/`
- Benchmarks go in `testing/benchmarks/`

## Questions?

Open an issue for discussion before starting significant work.
