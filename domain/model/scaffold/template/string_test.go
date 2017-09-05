package template

import (
	"testing"
)

func Test_String_Compile(t *testing.T) {
	testcases := []struct {
		value    interface{}
		text     String
		expected string
	}{
		{
			value:    nil,
			text:     "",
			expected: "",
		},
		{
			value:    struct{ Name, namespace string }{Name: "test", namespace: "ignored"},
			text:     "name = \"{{name}}\"\nsynopsis = \"\"\"This is {{name}}.\"\"\"",
			expected: "name = \"test\"\nsynopsis = \"\"\"This is test.\"\"\"",
		},
		{
			value:    struct{ Name, Namespace string }{Name: "test", Namespace: "template_string"},
			text:     "package {{namespace}}\n\n type {{name}} struct{}",
			expected: "package template_string\n\n type test struct{}",
		},
		{
			value:    struct{ Name string }{Name: "testcase"},
			text:     "{{name | toUpper}}",
			expected: "TESTCASE",
		},
		{
			value:    struct{ Name string }{Name: "TestCase"},
			text:     "{{name | toLower}}",
			expected: "testcase",
		},
		{
			value:    struct{ Name string }{Name: "TestCase"},
			text:     "{{name | underscored}}",
			expected: "test_case",
		},
		{
			value:    struct{ Name string }{Name: "TestCase"},
			text:     "{{name | dasherize}}",
			expected: "test-case",
		},
		{
			value:    struct{ Name string }{Name: "test_case"},
			text:     "{{name | camelize}}",
			expected: "testCase",
		},
		{
			value:    struct{ Name string }{Name: "test_case"},
			text:     "{{name | pascalize}}",
			expected: "TestCase",
		},
		{
			value:    struct{ Name string }{Name: "test_case"},
			text:     "{{name | pluralize}}",
			expected: "test_cases",
		},
		{
			value:    struct{ Name string }{Name: "test_cases"},
			text:     "{{name | singularize}}",
			expected: "test_case",
		},
		{
			value:    struct{ Name string }{Name: "test_cases"},
			text:     "{{name | firstChild}}",
			expected: "t",
		},
		{
			value:    struct{ Name string }{Name: "test_cases"},
			text:     "{{name | firstChild | toUpper}}",
			expected: "T",
		},
		{
			value:    &struct{ Name string }{Name: "test_case"},
			text:     "{{name}}",
			expected: "test_case",
		},
	}

	for _, c := range testcases {
		actual, err := c.text.Compile("template", c.value)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if actual != c.expected {
			t.Errorf("%q.Compile(%q) returns %s, want %s", c.text, c.value, actual, c.expected)
		}
	}
}

func Test_String_Compile_WithInvalidString(t *testing.T) {
	var ts String = "{{name}"
	s, err := ts.Compile("template", struct{ Name string }{Name: "foobar"})

	if actual, expected := s, string(ts); actual != expected {
		t.Errorf("Compile() returns %s, want %s", actual, expected)
	}

	if err == nil {
		t.Errorf("Should return an error")
	}
}

func Test_String_Compile_WithPrivateFields(t *testing.T) {
	var ts String = "package {{namespace}}\n\n type {{name}} struct{}"
	s, err := ts.Compile("template", struct{ Name, namespace string }{Name: "foobar", namespace: "baz"})

	if actual, expected := s, string(ts); actual != expected {
		t.Errorf("Compile() returns %s, want %s", actual, expected)
	}

	if err == nil {
		t.Errorf("Should return an error")
	}
}

func Test_String_Compile_WithUnsupportedTypeValue(t *testing.T) {
	var ts String = "{{name}}"
	s, err := ts.Compile("template", struct{ Name int }{Name: 1})

	if actual, expected := s, string(ts); actual != expected {
		t.Errorf("Compile() returns %s, want %s", actual, expected)
	}

	if err == nil {
		t.Errorf("Should return an error")
	}
}

func Test_String_Compile_WithNotStructObject(t *testing.T) {
	var ts String = "foobar"
	s, err := ts.Compile("template", 1)

	if actual, expected := s, string(ts); actual != expected {
		t.Errorf("Compile() returns %s, want %s", actual, expected)
	}

	if err == nil {
		t.Errorf("Should return an error")
	}
}

func Test_String_Compile_WithNoGivenValue(t *testing.T) {
	var ts String = "{{name}}"
	s, err := ts.Compile("template", struct{}{})

	if actual, expected := s, string(ts); actual != expected {
		t.Errorf("Compile() returns %s, want %s", actual, expected)
	}

	if err == nil {
		t.Errorf("Should return an error")
	}
}
