package main

import (
	"strings"
	"testing"
)

func TestPanelGitHubLink(t *testing.T) {
	html := panelHTML()
	for _, want := range []string{
		`href="https://github.com/gouwenct/cpa-scheduled-tests" target="_blank" rel="noopener noreferrer"`,
		`GitHub · Star</a>`,
		`onclick="openBulkPlan()"`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("panel missing %q", want)
		}
	}
}
