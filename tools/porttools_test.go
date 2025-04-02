// File: assignment-2/tools/porttools_test.go
package tools

import (
	"os"
	"testing"
)

// TestGetServerPort verifies that GetServerPort returns a fallback port
// when the PORT environment variable is unset or empty, and returns the
// environment variable's value if set.

func TestGetServerPort(t *testing.T) {

	// We'll save any existing PORT to restore it after the test.
	originalPort := os.Getenv("PORT")
	// Use a defer to restore the original environment afterwards.
	defer os.Setenv("PORT", originalPort)

	t.Run("FallbackUsedWhenUnset", func(t *testing.T) {
		// Unset PORT to ensure the function falls back.
		os.Unsetenv("PORT")

		fallback := "8080"
		got := GetServerPort(fallback)
		if got != fallback {
			t.Errorf("Expected fallback '%s', got '%s'", fallback, got)
		}
	})

	t.Run("FallbackUsedWhenEmptyString", func(t *testing.T) {
		// Set PORT to an empty string to simulate an empty variable.
		os.Setenv("PORT", "")

		fallback := "6060"
		got := GetServerPort(fallback)
		if got != fallback {
			t.Errorf("Expected fallback '%s', got '%s'", fallback, got)
		}
	})

	t.Run("EnvironmentOverride", func(t *testing.T) {
		// Set PORT to a non-empty value.
		os.Setenv("PORT", "5000")

		fallback := "9999"
		got := GetServerPort(fallback)
		if got != "5000" {
			t.Errorf("Expected '5000' from env var, got '%s'", got)
		}
	})
}
