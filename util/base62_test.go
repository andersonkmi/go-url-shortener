package util

import (
	"testing"
)

func TestConvertZeroToBase62(t *testing.T) {
	result := idToBase62(0)
	if result != "0" {
		t.Errorf("Expected 0 to be converted to 0, got %s", result)
	}
}
