// File: assignment-2/structs/features_test.go
package structs

import "testing"

// TestFeatures checks that NewFeatures() returns a Features struct
// with all booleans set to false and TargetCurrencies as an empty slice.
func TestFeatures(t *testing.T) {
	t.Run("NewFeatures_defaultValues", func(t *testing.T) {
		feats := NewFeatures()

		// Check that booleans are all false
		if feats.Temperature {
			t.Error("Expected Temperature to be false by default")
		}
		if feats.Precipitation {
			t.Error("Expected Precipitation to be false by default")
		}
		if feats.Capital {
			t.Error("Expected Capital to be false by default")
		}
		if feats.Coordinates {
			t.Error("Expected Coordinates to be false by default")
		}
		if feats.Population {
			t.Error("Expected Population to be false by default")
		}
		if feats.Area {
			t.Error("Expected Area to be false by default")
		}

		// Check TargetCurrencies is a non-nil, empty slice
		if feats.TargetCurrencies == nil {
			t.Error("Expected TargetCurrencies to be a non-nil empty slice, got nil")
		} else if len(feats.TargetCurrencies) != 0 {
			t.Errorf("Expected TargetCurrencies length = 0, got %d", len(feats.TargetCurrencies))
		}
	})
}
