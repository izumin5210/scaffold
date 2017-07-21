package scaffold

import "testing"

func Test_IsParentOf(t *testing.T) {
	entry := &Entry{path: "/app/.scaffold"}

	testcases := []struct {
		in  *Entry
		out bool
	}{
		{in: &Entry{path: "/"}, out: false},
		{in: &Entry{path: "/app"}, out: false},
		{in: &Entry{path: "/foo"}, out: false},
		{in: &Entry{path: "/app/.scaffold"}, out: false},
		{in: &Entry{path: "/app/foo"}, out: false},
		{in: &Entry{path: "/app/.scaffold/bar"}, out: true},
		{in: &Entry{path: "/app/foo/bar"}, out: false},
		{in: &Entry{path: "app"}, out: false},
	}

	for _, tc := range testcases {
		if actual, expected := entry.IsParentOf(tc.in), tc.out; actual != expected {
			t.Errorf("%v.IsParentOf(%v) returns. got %v: want %v", entry, tc.in, actual, expected)
		}
	}
}

func Test_IsChildOf(t *testing.T) {
	entry := &Entry{path: "/app/.scaffold"}

	testcases := []struct {
		in  *Entry
		out bool
	}{
		{in: &Entry{path: "/"}, out: false},
		{in: &Entry{path: "/app"}, out: true},
		{in: &Entry{path: "/foo"}, out: false},
		{in: &Entry{path: "/app/.scaffold"}, out: false},
		{in: &Entry{path: "/app/foo"}, out: false},
		{in: &Entry{path: "/app/.scaffold/bar"}, out: false},
		{in: &Entry{path: "/app/.scaffold/bar/baz"}, out: false},
		{in: &Entry{path: "/app/foo/bar"}, out: false},
		{in: &Entry{path: "app"}, out: false},
	}

	for _, tc := range testcases {
		if actual, expected := entry.IsChildOf(tc.in), tc.out; actual != expected {
			t.Errorf("%v.IsChildOf(%v) returns. got %v: want %v", entry, tc.in, actual, expected)
		}
	}
}
