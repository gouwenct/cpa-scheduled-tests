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
		`<th>分组</th>`,
		`id="planGroup"`,
		`id="planModel"`,
		`id="planTimezone"`,
		`<div class="title"><h1>CPA Scheduled 5H</h1>`,
		`id="planPagination"`,
		`onclick="togglePlan(`,
		`function fmtTime(s){`,
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("panel missing %q", want)
		}
	}
	if strings.Contains(html, "快速发送") || strings.Contains(html, "quickAccount") {
		t.Fatal("quick send panel should be removed")
	}
	if strings.Contains(html, "top-actions") {
		t.Fatal("legacy top layout should be restored")
	}
}
