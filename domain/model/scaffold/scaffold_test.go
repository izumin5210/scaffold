package scaffold

import (
	"testing"
)

func Test_Scaffold(t *testing.T) {
	tests := []struct {
		path string
		meta *Meta
		name string
	}{
		{
			path: "/app/.scaffold/foo",
			meta: &Meta{},
			name: "foo",
		},
		{
			path: "/app/.scaffold/foo",
			meta: &Meta{Name: "bar"},
			name: "bar",
		},
		{
			path: "/app/.scaffold/foo",
			meta: &Meta{
				Name:     "bar",
				Synopsis: "Generate bar",
				Help:     "Usage: scaffold g bar <name>\n\n",
			},
			name: "bar",
		},
	}

	for _, test := range tests {
		scff := NewScaffold(test.path, test.meta)

		if actual, expected := scff.Name(), test.name; actual != expected {
			t.Errorf("Name() returns %s, want %s", actual, expected)
		}

		if actual, expected := scff.Synopsis(), test.meta.Synopsis; actual != expected {
			t.Errorf("Synopsis() returns %s, want %s", actual, expected)
		}

		if actual, expected := scff.Help(), test.meta.Help; actual != expected {
			t.Errorf("Help() returns %s, want %s", actual, expected)
		}
	}
}
