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

func TestConvertNumberToBase62(t *testing.T) {
	testcases := []struct {
		in   int64
		want string
	}{
		{0, "0"},
		{123, "1z"},
		{29290292, "1ytk4"},
	}

	for _, tc := range testcases {
		result := idToBase62(tc.in)
		if result != tc.want {
			t.Errorf("Expected %d to be converted to %s, got %s", tc.in, tc.want, result)
		}
	}
}
