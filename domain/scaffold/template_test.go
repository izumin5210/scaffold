package scaffold

import "testing"

func Test_Compile(t *testing.T) {
	name := "test"
	tmpl := NewTemplate(name)
	testcases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			text:     "",
			expected: "",
		},
		{
			text:     "name = \"{{name}}\"\nsynopsis = \"\"\"This is {{name}}.\"\"\"",
			expected: "name = \"test\"\nsynopsis = \"\"\"This is test.\"\"\"",
		},
		{
			text:     "",
			expected: "",
		},
	}

	for _, c := range testcases {
		actual, err := tmpl.Compile(c.text)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual != c.expected {
			t.Errorf("NewTemplate(%q).Compile(%q) returns %s, want %s", name, c.text, actual, c.expected)
		}
	}
}
