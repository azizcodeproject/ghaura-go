package resi

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateFormat(t *testing.T) {
	resiNumber, err := Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	expectedPrefix := "GHR-" + time.Now().UTC().Format("060102") + "-"
	if !strings.HasPrefix(resiNumber, expectedPrefix) {
		t.Fatalf("unexpected prefix: got %s want prefix %s", resiNumber, expectedPrefix)
	}

	if len(resiNumber) != len(expectedPrefix)+6 {
		t.Fatalf("unexpected length: got %d (%s)", len(resiNumber), resiNumber)
	}
}

func TestGenerateUnique(t *testing.T) {
	seen := make(map[string]struct{})
	for i := 0; i < 50; i++ {
		resiNumber, err := Generate()
		if err != nil {
			t.Fatalf("Generate() error = %v", err)
		}
		if _, exists := seen[resiNumber]; exists {
			t.Fatalf("duplicate resi generated: %s", resiNumber)
		}
		seen[resiNumber] = struct{}{}
	}
}
