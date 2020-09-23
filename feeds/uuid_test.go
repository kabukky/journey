package feeds

import (
	"testing"
)

func TestUUID(t *testing.T) {
	s := NewUUID()
	s2 := NewUUID()
	if len(s) != 16 {
		t.Errorf("Expecting len of 16, got %d\n", len(s))
	}
	if len(s.String()) != 36 {
		t.Errorf("Expecting uuid hex string len of 36, got %d\n", len(s.String()))
	}
	if s == s2 {
		t.Errorf("Expecting different UUIDs to be different, but they are the same.\n")
	}
}
