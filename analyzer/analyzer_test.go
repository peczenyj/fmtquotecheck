package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/peczenyj/fmtquotecheck/analyzer"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer.New(), "a")
}
