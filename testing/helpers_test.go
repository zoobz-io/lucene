//go:build testing

package testing

import (
	"errors"
	"testing"
)

func TestAssertQueryString(t *testing.T) {
	t.Run("matching strings pass", func(t *testing.T) {
		// This should not fail
		AssertQueryString(t, "field:value", "field:value")
	})
}

func TestAssertNoError(t *testing.T) {
	t.Run("nil error passes", func(t *testing.T) {
		AssertNoError(t, nil)
	})
}

func TestAssertError(t *testing.T) {
	t.Run("non-nil error passes", func(t *testing.T) {
		AssertError(t, errors.New("expected error"))
	})
}
