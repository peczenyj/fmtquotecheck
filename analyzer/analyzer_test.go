package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/peczenyj/fmtquotecheck/analyzer"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()

	t.Run("run with default flag test", func(t *testing.T) {
		instance := analyzer.New()

		analysistest.Run(t, testdata, instance, "a")
	})

	t.Run("run with default flag test and suggested fixes", func(t *testing.T) {
		instance := analyzer.New()

		analysistest.RunWithSuggestedFixes(t, testdata, instance, "fix")
	})

	t.Run("run with restricted printf functions", func(t *testing.T) {
		instance := analyzer.New()

		err := instance.Flags.Set("funcs", "fmt.Printf")
		if err != nil {
			t.Fatalf("unexpected error while set flag '-funcs': %v", err)
		}

		analysistest.Run(t, testdata, instance, "b")
	})
}
