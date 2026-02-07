package split

import "testing"

func TestNewTemplate(t *testing.T) {
	tmpl := NewTemplate("t-1", "Super Mario 64", []string{"Bob-omb", "Whomp", "Jolly Roger"})

	if tmpl.ID != "t-1" {
		t.Fatalf("expected ID t-1, got %s", tmpl.ID)
	}

	if tmpl.Name != "Super Mario 64" {
		t.Fatalf("expected name Super Mario 64, got %s", tmpl.Name)
	}

	if len(tmpl.SegmentNames) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(tmpl.SegmentNames))
	}

	if tmpl.SegmentNames[0] != "Bob-omb" {
		t.Fatalf("expected Bob-omb, got %s", tmpl.SegmentNames[0])
	}

	if tmpl.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
}
