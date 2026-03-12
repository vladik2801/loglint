package logcheck

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSuggestedFixes(t *testing.T) {
	old := enableFix
	enableFix = true
	t.Cleanup(func() { enableFix = old })

	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), Analyzer, "fixpkg")
}
