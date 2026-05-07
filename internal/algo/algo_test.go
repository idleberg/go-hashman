package algo

import (
	"testing"
)

func TestRegistryHasUniqueIDs(t *testing.T) {
	seen := make(map[string]bool)
	for _, a := range Registry {
		if seen[a.ID] {
			t.Errorf("duplicate algorithm ID: %s", a.ID)
		}
		seen[a.ID] = true
	}
}

func TestRegistryHasUniqueFlags(t *testing.T) {
	seen := make(map[string]bool)
	for _, a := range Registry {
		if seen[a.Flag] {
			t.Errorf("duplicate flag: %s", a.Flag)
		}
		seen[a.Flag] = true
	}
}

func TestNewHashReturnsNonNil(t *testing.T) {
	for _, a := range Registry {
		h := a.NewHash()
		if h == nil {
			t.Errorf("NewHash() returned nil for %s", a.ID)
		}
	}
}
