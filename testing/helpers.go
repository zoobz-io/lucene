//go:build testing

// Package testing provides test helpers for the lucene package.
package testing

import "testing"

// AssertQueryString validates that a query produces the expected string output.
func AssertQueryString(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("query string mismatch:\n  got:  %q\n  want: %q", got, want)
	}
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error, got nil")
	}
}
