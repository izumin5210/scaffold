package scaffold

import "testing"

func Test_Compile(t *testing.T) {
	testcases := []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "",
			text:     "",
			expected: "",
		},
		{
			name:     "test",
			text:     "name = \"{{name}}\"\nsynopsis = \"\"\"This is {{name}}.\"\"\"",
			expected: "name = \"test\"\nsynopsis = \"\"\"This is test.\"\"\"",
		},
		{
			name:     "testcase",
			text:     "{{name | toUpper}}",
			expected: "TESTCASE",
		},
		{
			name:     "TestCase",
			text:     "{{name | toLower}}",
			expected: "testcase",
		},
		{
			name:     "TestCase",
			text:     "{{name | underscored}}",
			expected: "test_case",
		},
		{
			name:     "TestCase",
			text:     "{{name | dasherize}}",
			expected: "test-case",
		},
		{
			name:     "test_case",
			text:     "{{name | camelize}}",
			expected: "testCase",
		},
		{
			name:     "test_case",
			text:     "{{name | pascalize}}",
			expected: "TestCase",
		},
		{
			name:     "test_case",
			text:     "{{name | pluralize}}",
			expected: "test_cases",
		},
		{
			name:     "test_cases",
			text:     "{{name | singularize}}",
			expected: "test_case",
		},
		{
			name:     "test_cases",
			text:     "{{name | firstChild}}",
			expected: "t",
		},
		{
			name:     "test_cases",
			text:     "{{name | firstChild | toUpper}}",
			expected: "T",
		},
	}

	for _, c := range testcases {
		actual, err := NewTemplate(c.name).Compile(c.text)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual != c.expected {
			t.Errorf("NewTemplate(%q).Compile(%q) returns %s, want %s", c.name, c.text, actual, c.expected)
		}
	}
}

func Test_Compile_Reuse(t *testing.T) {
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
