package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNormalizeStateMigratesLegacyPlanNameToGroup(t *testing.T) {
	s := PersistentState{Plans: []Plan{{Name: "日常测试 · account@example.com"}}}
	normalizeState(&s)
	if got := s.Plans[0].Group; got != "日常测试" {
		t.Fatalf("group = %q, want 日常测试", got)
	}
}

func TestListLogsKeepsOnlyRecentTwoDays(t *testing.T) {
	oldDataDir, oldState := rt.dataDir, rt.state
	defer func() { rt.dataDir, rt.state = oldDataDir, oldState }()
	rt.dataDir = t.TempDir()
	entries := []LogEntry{
		{ID: "old", At: time.Now().Add(-49 * time.Hour), AccountName: "old"},
		{ID: "new", At: time.Now().Add(-47 * time.Hour), AccountName: "new"},
	}
	f, err := os.Create(filepath.Join(rt.dataDir, "runs.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	enc := json.NewEncoder(f)
	for _, entry := range entries {
		if err := enc.Encode(entry); err != nil {
			t.Fatal(err)
		}
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	got, err := listLogs(200)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].ID != "new" {
		t.Fatalf("logs = %#v, want only recent entry", got)
	}
}
