package pnp

import (
	"strconv"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	t.Parallel()
	tr := New[string]()
	tr.Set("foo", "A")
	tr.Set("foobar", "B")
	tr.Set("bar", "C")

	tests := []struct {
		key      string
		want     string
		wantOkay bool
	}{
		{"foo", "A", true},
		{"foobar", "B", true},
		{"bar", "C", true},
		{"fo", "", false},
		{"baz", "", false},
	}

	for _, tc := range tests {
		got, ok := tr.Get(tc.key)
		if ok != tc.wantOkay || (ok && got != tc.want) {
			t.Fatalf("Get(%q) = (%q,%v), want (%q,%v)", tc.key, got, ok, tc.want, tc.wantOkay)
		}
	}
}

func TestGetAncestorValue_Basic(t *testing.T) {
	t.Parallel()
	tr := New[string]()
	tr.Set("foo", "A")
	tr.Set("foobar", "B")

	tests := []struct {
		key      string
		want     string
		wantOkay bool
	}{
		{"foobar", "B", true},
		{"fo", "", false},
		{"foz", "", false},
		{"foobaz", "A", true},
		{"x", "", false},
	}

	for _, tc := range tests {
		got, ok := tr.GetAncestorValue(tc.key)
		if ok != tc.wantOkay || (ok && got != tc.want) {
			t.Fatalf("GetAncestorValue(%q) = (%q,%v), want (%q,%v)",
				tc.key, got, ok, tc.want, tc.wantOkay)
		}
	}
}

func TestGetAncestorValue_Unicode(t *testing.T) {
	t.Parallel()
	tr := New[string]()
	tr.Set("한", "H1")
	tr.Set("한글", "H2")
	tr.Set("한글날", "H3")

	tests := []struct {
		key      string
		want     string
		wantOkay bool
	}{
		{"한", "H1", true},
		{"한글", "H2", true},
		{"한글날", "H3", true},
		{"한글나라", "H2", true},
		{"한국", "H1", true},
	}

	for _, tc := range tests {
		got, ok := tr.GetAncestorValue(tc.key)
		if ok != tc.wantOkay || (ok && got != tc.want) {
			t.Fatalf("Unicode GetAncestorValue(%q) = (%q,%v), want (%q,%v)",
				tc.key, got, ok, tc.want, tc.wantOkay)
		}
	}
}

func TestRootValue_AsAncestor(t *testing.T) {
	t.Parallel()
	tr := New[string]()
	tr.Set("", "ROOT")
	tr.Set("a", "A")

	tests := []struct {
		key      string
		want     string
		wantOkay bool
	}{
		{"", "ROOT", true},
		{"x", "ROOT", true},
		{"ab", "A", true},
	}

	for _, tc := range tests {
		got, ok := tr.GetAncestorValue(tc.key)
		if ok != tc.wantOkay || (ok && got != tc.want) {
			t.Fatalf("Root GetAncestorValue(%q) = (%q,%v), want (%q,%v)",
				tc.key, got, ok, tc.want, tc.wantOkay)
		}
	}
}

func TestZeroValue_Int(t *testing.T) {
	t.Parallel()
	tr := New[int]()
	tr.Set("zero", 0)
	tr.Set("one", 1)

	got, ok := tr.Get("zero")
	if !ok || got != 0 {
		t.Fatalf("Get(zero) = (%d,%v), want (0,true)", got, ok)
	}

	got2, ok2 := tr.GetAncestorValue("zero-suffix")
	if !ok2 || got2 != 0 {
		t.Fatalf("GetAncestorValue(zero-suffix) = (%d,%v), want (0,true)", got2, ok2)
	}
}

func TestGetAncestorValue_NoAncestor(t *testing.T) {
	t.Parallel()
	tr := New[string]()

	if _, ok := tr.GetAncestorValue("anything"); ok {
		t.Fatalf("expected no ancestor value, but got ok=true")
	}
}

func TestGetAncestorValue_LongestPrefixWins(t *testing.T) {
	t.Parallel()
	tr := New[string]()
	// "a", "ab", "abc", ..., "abcdefghij"
	prefix := ""
	for i := range 10 {
		prefix += "a"
		tr.Set(prefix, "V"+strconv.Itoa(i))
	}
	// "aaaaaaaaax"
	key := prefix + "x"

	got, ok := tr.GetAncestorValue(key)
	if !ok || got != "V9" {
		t.Fatalf("LongestPrefix GetAncestorValue = (%q,%v), want (\"V9\",true)", got, ok)
	}
}

func BenchmarkGetAncestorValue(b *testing.B) {
	tr := New[string]()
	for i := range 1000 {
		tr.Set("key-"+strconv.Itoa(i), "V"+strconv.Itoa(i))
	}
	tr.Set("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "LONG")

	b.ResetTimer()
	for range b.N {
		_, _ = tr.GetAncestorValue("key-123-extra")
		_, _ = tr.GetAncestorValue("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaax")
		_, _ = tr.GetAncestorValue("zzz")
	}
}
