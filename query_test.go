package gmuch

import (
	"testing"

	"github.com/go-kit/kit/log"
)

func TestQuery(t *testing.T) {
	g := New(dbPath, log.NewNopLogger())
	ts, err := g.Query("", 0, 100)
	if err != nil {
		t.Fatalf("g.Query(%q, %d, %d): got error: %s", "", 0, 100, err)
	}
	if want, got := 24, len(ts); want != got {
		t.Errorf("g.Query(%q, %d, %d): want %d threads got %d", "", 0, 100, want, got)
	}

	ts, err = g.Query("", 0, 20)
	if err != nil {
		t.Fatalf("g.Query(%q, %d, %d): got error: %s", "", 0, 20, err)
	}
	if want, got := 20, len(ts); want != got {
		t.Errorf("g.Query(%q, %d, %d): want %d threads got %d", "", 0, 20, want, got)
	}

	ts, err = g.Query("", 20, 100)
	if err != nil {
		t.Fatalf("g.Query(%q, %d, %d): got error: %s", "", 20, 100, err)
	}
	if want, got := 4, len(ts); want != got {
		t.Errorf("g.Query(%q, %d, %d): want %d threads got %d", "", 20, 100, want, got)
	}
}
