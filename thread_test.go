package gmuch

import (
	"os"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestThreadIDValidation(t *testing.T) {
	g := New(dbPath, log.NewNopLogger())
	errTests := map[string]error{
		"123":              ErrIDNotFound,
		"abc":              ErrIDNotFound,
		"abc123":           ErrIDNotFound,
		"a-":               ErrIDInvalid,
		"thread:abc":       ErrIDInvalid,
		"thread:123":       ErrIDInvalid,
		"0000000000000001": nil,
		"000000000000001b": nil,
	}

	for id, want := range errTests {
		if _, got := g.Thread(id); want != got {
			t.Errorf("Thread(%q): want error %q got %q", id, want, got)
		}
	}
}

func TestThread(t *testing.T) {
	g := New(dbPath, log.NewLogfmtLogger(os.Stderr))
	ts, err := g.Query("*", 0, 1000)
	if err != nil {
		t.Fatal(err)
	}

	for _, thread := range ts {
		tr, err := g.Thread(thread.ID)
		if err != nil {
			t.Errorf("g.Thread(%q): got error: %s", thread.ID, err)
			continue
		}

		if want, got := thread.ID, tr.ID; want != got {
			t.Errorf("g.Thread(%q).ID: want %s got %s", thread.ID, want, got)
		}
	}
}
