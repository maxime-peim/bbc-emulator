package utils

import (
	"testing"
)

func TestSegment(t *testing.T) {
	s := NewSegment(0x1000, 0xA000)
	a := NewSegment(0x0100, 0x2000)
	b := NewSegment(0xB000, 0xFFF0)

	if s.Size() != 0x9001 {
		t.Fatalf("Wrong size (%d != %d)", s.Size(), 0x9001)
	} else if !s.Intersect(a) || !a.Intersect(s) {
		t.Fatal("Intersection fails")
	} else if s.Intersect(b) || b.Intersect(s) {
		t.Fatal("No intersection fails")
	} else if !s.IsIn(0x5123) {
		t.Fatal("0x5123 is in s")
	} else if s.IsIn(0x0000) {
		t.Fatal("0x0000 is not in s")
	}
}
