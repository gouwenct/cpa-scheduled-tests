package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestDeleteAllPlansPersistsAndPreservesOtherData(t *testing.T) {
	setupScheduledTest(t, Plan{ID: "first", Enabled: true})
	for i := 1; i < 12; i++ {
		rt.state.Plans = append(rt.state.Plans, Plan{ID: fmt.Sprint(i), Group: fmt.Sprint(i % 2), Enabled: i%2 == 0})
	}
	settings := rt.state.Settings
	logPath := filepath.Join(rt.dataDir, "runs.jsonl")
	logData := []byte("preserved log\n")
	if err := os.WriteFile(logPath, logData, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := saveStateLocked(); err != nil {
		t.Fatal(err)
	}
	deleted, err := deleteAllPlans()
	if err != nil || deleted != 12 {
		t.Fatalf("deleteAllPlans() = %d, %v, want 12, nil", deleted, err)
	}
	if len(snapshotPersistent().Plans) != 0 || rt.state.Settings != settings {
		t.Fatal("plans were not cleared or settings changed")
	}
	raw, err := os.ReadFile(filepath.Join(rt.dataDir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	var persisted PersistentState
	if err := json.Unmarshal(raw, &persisted); err != nil {
		t.Fatal(err)
	}
	if persisted.Plans == nil || len(persisted.Plans) != 0 || persisted.Settings != settings {
		t.Fatalf("unexpected persisted state: %#v", persisted)
	}
	if raw, err := os.ReadFile(logPath); err != nil || string(raw) != string(logData) {
		t.Fatalf("logs changed: %q, %v", raw, err)
	}
	if markScheduled("first", "2026-10-23T06:00", false) || markScheduled("first", "2026-10-23T06:00", true) {
		t.Fatal("deleted plan can still be scheduled")
	}
	evaluateSchedules(context.Background(), time.Now())
	if len(jobsSnapshot()) != 0 {
		t.Fatal("deleted plans started jobs")
	}
	updatePlanResult("first", LogEntry{Status: "success"})
	if len(snapshotPersistent().Plans) != 0 {
		t.Fatal("finished request restored a deleted plan")
	}
	if deleted, err := deleteAllPlans(); err != nil || deleted != 0 {
		t.Fatalf("empty delete = %d, %v, want 0, nil", deleted, err)
	}
}

func TestDeleteAllPlansRollsBackWhenSaveFails(t *testing.T) {
	setupScheduledTest(t, Plan{ID: "preserved", Enabled: true})
	before := snapshotPersistent()
	// A directory at the temporary file path forces a portable write failure.
	if err := os.Mkdir(filepath.Join(rt.dataDir, "state.json.tmp"), 0o700); err != nil {
		t.Fatal(err)
	}
	if deleted, err := deleteAllPlans(); err == nil || deleted != 0 {
		t.Fatalf("failed save = %d, %v, want 0 and error", deleted, err)
	}
	if !reflect.DeepEqual(snapshotPersistent(), before) {
		t.Fatal("failed deletion changed the in-memory state")
	}
}

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
