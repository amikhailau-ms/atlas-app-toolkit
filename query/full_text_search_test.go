package query

import "testing"

func TestFullTextSearch(t *testing.T) {
	fts := FullTextSearch{Query: ""}
	if fts.GoString() != "" {
		t.Errorf("invalid string representation: %v - expected: %s", fts, "")
	}

	fts = FullTextSearch{Query: "some query"}
	if fts.GoString() != "some query" {
		t.Errorf("invalid string representation: %v - expected: %s", fts, "some query")
	}
}
