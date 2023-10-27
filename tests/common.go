package tests

import (
	"sort"
	"testing"

	"github.com/metux/go-magicdict/api"
)

// explicitly skip nil's
func fetchEntry(t *testing.T, root api.Entry, k api.Key) api.Entry {
	v, err := root.Get(k)
	if err != nil {
		t.Fatalf("entry \"%s\" is not Entry: \"%s\"", k, err)
		return nil
	}
	return v
}

func checkStrs(t *testing.T, got []string, want []string) {
	sort.Strings(got)
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("checkStrs: size mismatch %d want %d -- %s vs %s", len(got), len(want), got, want)
	}

	for idx, s := range got {
		if s != want[idx] {
			t.Fatalf("IDX #%d mismatch: \"%s\" should be \"%s\"", idx, s, want[idx])
		}
	}
}

func checkKeys(t *testing.T, got []api.Key, want []string) {
	gotk := make([]string, len(got))
	for x, y := range got {
		gotk[x] = string(y)
	}

	sort.Strings(gotk)
	sort.Strings(want)

	for idx, s := range gotk {
		if s != want[idx] {
			t.Fatalf("IDX #%d mismatch: \"%s\" should be \"%s\"", idx, s, want[idx])
		}
	}
}
