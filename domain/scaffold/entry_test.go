package scaffold

import "testing"

func Test_IsParentOf(t *testing.T) {
	target := &entry{path: "/app/.scaffold"}

	testcases := []struct {
		in  Entry
		out bool
	}{
		{in: &entry{path: "/"}, out: false},
		{in: &entry{path: "/app"}, out: false},
		{in: &entry{path: "/foo"}, out: false},
		{in: &entry{path: "/app/.scaffold"}, out: false},
		{in: &entry{path: "/app/foo"}, out: false},
		{in: &entry{path: "/app/.scaffold/bar"}, out: true},
		{in: &entry{path: "/app/foo/bar"}, out: false},
		{in: &entry{path: "app"}, out: false},
	}

	for _, tc := range testcases {
		if actual, expected := target.IsParentOf(tc.in), tc.out; actual != expected {
			t.Errorf("%v.IsParentOf(%v) returns. got %v: want %v", target, tc.in, actual, expected)
		}
	}
}

func Test_IsChildOf(t *testing.T) {
	target := &entry{path: "/app/.scaffold"}

	testcases := []struct {
		in  Entry
		out bool
	}{
		{in: &entry{path: "/"}, out: false},
		{in: &entry{path: "/app"}, out: true},
		{in: &entry{path: "/foo"}, out: false},
		{in: &entry{path: "/app/.scaffold"}, out: false},
		{in: &entry{path: "/app/foo"}, out: false},
		{in: &entry{path: "/app/.scaffold/bar"}, out: false},
		{in: &entry{path: "/app/.scaffold/bar/baz"}, out: false},
		{in: &entry{path: "/app/foo/bar"}, out: false},
		{in: &entry{path: "app"}, out: false},
	}

	for _, tc := range testcases {
		if actual, expected := target.IsChildOf(tc.in), tc.out; actual != expected {
			t.Errorf("%v.IsChildOf(%v) returns. got %v: want %v", target, tc.in, actual, expected)
		}
	}
}
