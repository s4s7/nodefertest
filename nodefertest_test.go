package nodefertest_test

import (
	"testing"

	"github.com/gostaticanalysis/testutil"
	"github.com/s4s7/nodefertest"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {
	testdata := testutil.WithModules(t, analysistest.TestData(), nil)
	analysistest.Run(t, testdata, nodefertest.Analyzer, "a")
}
