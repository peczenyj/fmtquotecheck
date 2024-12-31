package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/peczenyj/fmtquotecheck/analyzer"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	t.Run("run with default flag test", func(t *testing.T) {
		instance, err := analyzer.New()
		if err != nil {
			t.Fatalf("unexpected error build a new analyzer: %v", err)
		}

		analysistest.Run(t, testdata, instance, "a")
	})

	t.Run("run with default flag test and suggested fixes", func(t *testing.T) {
		instance, err := analyzer.New()
		if err != nil {
			t.Fatalf("unexpected error build a new analyzer: %v", err)
		}

		analysistest.RunWithSuggestedFixes(t, testdata, instance, "fix")
	})

	t.Run("run with restricted printf functions via flags", func(t *testing.T) {
		instance, err := analyzer.New()
		if err != nil {
			t.Fatalf("unexpected error build a new analyzer: %v", err)
		}

		err = instance.Flags.Set("funcs", "fmt.Printf")
		if err != nil {
			t.Fatalf("unexpected error while set flag '-funcs': %v", err)
		}

		analysistest.Run(t, testdata, instance, "b")
	})

	t.Run("run with restricted printf functions via optons", func(t *testing.T) {
		instance, err := analyzer.New(analyzer.WithPrintfFuncs("fmt.Printf"))
		if err != nil {
			t.Fatalf("unexpected error build a new analyzer: %v", err)
		}

		analysistest.Run(t, testdata, instance, "b")
	})

	t.Run("should not accept empty printf funcs", func(t *testing.T) {
		instance, err := analyzer.New(analyzer.WithPrintfFuncs(""))
		require.EqualError(t, err, "invalid printf func: empty string")

		assert.Nil(t, instance)
	})
}
