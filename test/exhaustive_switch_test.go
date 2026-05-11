package test

import (
	"testing"

	"github.com/nishanths/exhaustive"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExhaustiveSwitchCaseTest(t *testing.T) {
	_ = analysistest.Run(t, "../", exhaustive.Analyzer, "./src/...")
}
