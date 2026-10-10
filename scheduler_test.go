package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupScheduledTest(t *testing.T, plan Plan) {
	t.Helper()
	rt.mu.Lock()
	oldDir, oldState, oldJobs, oldKeys := rt.dataDir, rt.state, rt.jobs, rt.activeKeys
	rt.dataDir = t.TempDir()
	rt.state = PersistentState{SchemaVersion: 1, Settings: defaultSettings(), Plans: []Plan{plan}}
	rt.jobs = make(map[string]*Job)
	rt.activeKeys = make(map[string]string)
	rt.mu.Unlock()
	t.Cleanup(func() {
		waitScheduledJobs(t)
		rt.mu.Lock()
		rt.dataDir, rt.state, rt.jobs, rt.activeKeys = oldDir, oldState, oldJobs, oldKeys
		rt.mu.Unlock()
	})
}

func waitScheduledJobs(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for {
		running := false
		for _, job := range jobsSnapshot() {
			running = running || job.Running
		}
		if !running {
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("scheduled jobs did not finish")
		}
		time.Sleep(time.Millisecond)
	}
}

func TestScheduledSendsAtCronMinuteAndNextMinute(t *testing.T) {
	for _, tc := range []struct {
		name string
		cron string
		at   string
	}{
		{"daily", "0 6,11,16,21 * * *", "2026-10-10T06:00:00+08:00"},
		{"hour boundary", "59 6 * * *", "2026-10-10T06:59:00+08:00"},
		{"year boundary", "59 23 31 12 *", "2026-12-31T23:59:00+08:00"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			at, _ := time.Parse(time.RFC3339, tc.at)
			setupScheduledTest(t, Plan{ID: "p", Enabled: true, AuthIndex: "missing", Model: "test", Cron: tc.cron, Timezone: "Asia/Taipei"})
			for _, delta := range []time.Duration{0, 10 * time.Second, 50 * time.Second, time.Minute, 70 * time.Second, 2 * time.Minute} {
				evaluateSchedules(context.Background(), at.Add(delta))
				waitScheduledJobs(t)
			}
			jobs := jobsSnapshot()
			if len(jobs) != 2 {
				t.Fatalf("got %d sends, want exactly 2", len(jobs))
			}
			logs, err := listLogs(10)
			if err != nil || len(logs) != 2 {
				t.Fatalf("logs = %#v, error = %v", logs, err)
			}
			if logs[0].Trigger != "scheduled+1min" || logs[1].Trigger != "scheduled" {
				t.Fatalf("unexpected triggers: %#v", logs)
			}
		})
	}
}

func TestAdjacentCronMinutesEachSendTwice(t *testing.T) {
	at := time.Date(2026, 10, 10, 6, 0, 0, 0, time.UTC)
	setupScheduledTest(t, Plan{ID: "p", Enabled: true, AuthIndex: "missing", Model: "test", Cron: "0,1 6 * * *", Timezone: "UTC"})
	for _, delta := range []time.Duration{0, time.Minute, 70 * time.Second, 2 * time.Minute, 130 * time.Second} {
		evaluateSchedules(context.Background(), at.Add(delta))
		waitScheduledJobs(t)
	}
	if got := len(jobsSnapshot()); got != 4 {
		t.Fatalf("got %d sends, want 4 for two Cron occurrences", got)
	}
}

func TestDisabledOrStoppedScheduleDoesNotSend(t *testing.T) {
	at := time.Date(2026, 10, 10, 6, 1, 0, 0, time.UTC)
	setupScheduledTest(t, Plan{ID: "p", Cron: "0 6 * * *", Timezone: "UTC"})
	evaluateSchedules(context.Background(), at)
	rt.mu.Lock()
	rt.state.Plans[0].Enabled = true
	rt.mu.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	evaluateSchedules(ctx, at)
	waitScheduledJobs(t)
	if len(jobsSnapshot()) != 0 {
		t.Fatal("disabled or stopped schedule dispatched a request")
	}
}

func TestScheduledRepeatSurvivesStateReloadAndEditing(t *testing.T) {
	at := time.Date(2026, 10, 10, 6, 0, 0, 0, time.UTC)
	setupScheduledTest(t, Plan{ID: "p", Enabled: true, AuthIndex: "missing", Model: "test", Cron: "0 6 * * *", Timezone: "UTC"})
	evaluateSchedules(context.Background(), at)
	waitScheduledJobs(t)
	evaluateSchedules(context.Background(), at.Add(time.Minute))
	waitScheduledJobs(t)
	raw, err := os.ReadFile(filepath.Join(rt.dataDir, "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	var saved PersistentState
	if err := json.Unmarshal(raw, &saved); err != nil {
		t.Fatal(err)
	}
	if saved.Plans[0].LastScheduledRepeatKey == "" {
		t.Fatal("second-send dedupe record was not persisted")
	}
	rt.mu.Lock()
	rt.state = saved
	rt.mu.Unlock()
	p := saved.Plans[0]
	p.LastScheduledRepeatKey = ""
	if _, err := upsertPlan(p); err != nil {
		t.Fatal(err)
	}
	evaluateSchedules(context.Background(), at.Add(70*time.Second))
	waitScheduledJobs(t)
	if got := len(jobsSnapshot()); got != 2 {
		t.Fatalf("state reload or editing repeated a send: %d jobs", got)
	}
}

func TestSecondSendDoesNotWaitForFirstJob(t *testing.T) {
	at := time.Date(2026, 10, 10, 6, 0, 0, 0, time.UTC)
	setupScheduledTest(t, Plan{ID: "p", Enabled: true, AuthIndex: "missing", Model: "test", Cron: "0 6 * * *", Timezone: "UTC"})
	key := "plan:p:scheduled:" + at.Format("2006-01-02T15:04")
	first, err := beginJob("plan", "p", "test", key)
	if err != nil {
		t.Fatal(err)
	}
	defer finishJob(first.ID, key, "")
	markScheduled("p", at.Format("2006-01-02T15:04"), false)
	evaluateSchedules(context.Background(), at.Add(time.Minute))
	if got := len(jobsSnapshot()); got != 2 {
		t.Fatalf("first job blocked second send: %d jobs", got)
	}
	finishJob(first.ID, key, "")
	waitScheduledJobs(t)
}

func TestSlowFirstSendDoesNotOverwriteSecondResult(t *testing.T) {
	at := time.Date(2026, 10, 10, 6, 0, 0, 0, time.UTC)
	setupScheduledTest(t, Plan{ID: "p"})
	updatePlanResult("p", LogEntry{At: at.Add(time.Minute), Status: "success"})
	updatePlanResult("p", LogEntry{At: at, Status: "failed"})
	p, _ := getPlan("p")
	if p.LastStatus != "success" || !p.LastRunAt.Equal(at.Add(time.Minute)) {
		t.Fatalf("slow first send replaced latest result: %#v", p)
	}
}
