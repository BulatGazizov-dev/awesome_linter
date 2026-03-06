package linters

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestPluginExample(t *testing.T) {
	testdata := analysistest.TestData()

	tests := []struct {
		name         string
		settings     MySettings
		pattern      string
		suggestedFix bool
	}{
		{
			name:         "Rule: Lowercase",
			settings:     MySettings{DenyBeginUpper: true},
			pattern:      "a/lower",
			suggestedFix: true,
		},
		{
			name:         "Rule: English Only",
			settings:     MySettings{CheckEnglishOnly: true},
			pattern:      "a/english",
			suggestedFix: false,
		},
		{
			name:         "Rule: Sensitive Data",
			settings:     MySettings{CheckSensitive: true},
			pattern:      "a/sensitive",
			suggestedFix: true,
		},
		{
			name:         "Rule: Deny special symbols",
			settings:     MySettings{DenySpecialSymbols: true},
			pattern:      "a/symbols",
			suggestedFix: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PluginExample{settings: tt.settings}
			analyzers, _ := p.BuildAnalyzers()
			if tt.suggestedFix {
				analysistest.RunWithSuggestedFixes(t, testdata, analyzers[0], tt.pattern)
			} else {
				analysistest.Run(t, testdata, analyzers[0], tt.pattern)
			}
		})
	}
}
