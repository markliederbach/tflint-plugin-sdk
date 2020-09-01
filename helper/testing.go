package helper

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/markliederbach/tflint-plugin-sdk/tflint"
)

// TestRunner returns a mock Runner for testing.
// You can pass the map of file names and their contents in the second argument.
func TestRunner(t *testing.T, files map[string]string) *Runner {
	runner := &Runner{Files: map[string]*hcl.File{}, Issues: Issues{}}
	parser := hclparse.NewParser()

	for name, src := range files {
		file, diags := parser.ParseHCL([]byte(src), name)
		if diags.HasErrors() {
			t.Fatal(diags)
		}
		runner.Files[name] = file
	}

	return runner
}

// AssertIssues is an assertion helper for comparing issues.
func AssertIssues(t *testing.T, expected Issues, actual Issues) {
	opts := []cmp.Option{
		// Byte field will be ignored because it's not important in tests such as positions
		cmpopts.IgnoreFields(hcl.Pos{}, "Byte"),
		ruleComparer(),
	}
	if !cmp.Equal(expected, actual, opts...) {
		t.Fatalf("Expected issues are not matched:\n %s\n", cmp.Diff(expected, actual, opts...))
	}
}

// AssertIssuesWithoutRange is an assertion helper for comparing issues except for range.
func AssertIssuesWithoutRange(t *testing.T, expected Issues, actual Issues) {
	opts := []cmp.Option{
		cmpopts.IgnoreFields(Issue{}, "Range"),
		ruleComparer(),
	}
	if !cmp.Equal(expected, actual, opts...) {
		t.Fatalf("Expected issues are not matched:\n %s\n", cmp.Diff(expected, actual, opts...))
	}
}

// ruleComparer returns a Comparer func that checks that two rule interfaces
// have the same underlying type. It does not compare struct fields.
func ruleComparer() cmp.Option {
	return cmp.Comparer(func(x, y tflint.Rule) bool {
		return reflect.TypeOf(x) == reflect.TypeOf(y)
	})
}
