package passwordhasher

import (
	"testing"
)

func TestGenerateAndCompare(t *testing.T) {
	input := "thisismypassword"

	h, err := Generate(input)
	if err != nil {
		t.Fatalf("Generate produced an error: %v", err)
	}

	if !Compare(h, input) {
		t.Errorf("Expected response to be true when comparing to correct password")
	}

	if Compare(h, "invalidpassword") {
		t.Errorf("Expected response to be false when comparing to incorrect password")
	}
}
