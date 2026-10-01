package main

import "testing"

func TestBulkCreatePlansCreatesEveryAccountAndDedupes(t *testing.T) {
	stopRuntime()
	rt.mu.Lock()
	rt.initialized = true
	rt.dataDir = t.TempDir()
	rt.state = PersistentState{SchemaVersion: 1, Settings: defaultSettings(), Plans: []Plan{}}
	rt.jobs = make(map[string]*Job)
	rt.activeKeys = make(map[string]string)
	rt.mu.Unlock()
	defer func() {
		rt.mu.Lock()
		rt.initialized = false
		rt.state = PersistentState{}
		rt.jobs = nil
		rt.activeKeys = nil
		rt.dataDir = ""
		rt.mu.Unlock()
	}()

	accounts := []AuthFile{
		{AuthIndex: "a1", Name: "a.json", Email: "a@example.com"},
		{AuthIndex: "b2", Name: "b.json", Email: "b@example.com", Disabled: true, Unavailable: true},
	}
	spec := BulkPlanSpec{
		NamePrefix: "全账号定时",
		Model:      "gpt-test",
		Cron:       "0 6,11,16,21 * * *",
		Timezone:   "Asia/Shanghai",
		Prompt:     "ping",
		Enabled:    true,
	}

	first, err := bulkCreatePlans(accounts, spec)
	if err != nil {
		t.Fatal(err)
	}
	if first.TotalAccounts != 2 || first.Created != 2 || first.Skipped != 0 {
		t.Fatalf("unexpected first result: %#v", first)
	}
	state := snapshotPersistent()
	if len(state.Plans) != 2 {
		t.Fatalf("plans = %d, want 2", len(state.Plans))
	}
	if state.Plans[0].AuthIndex == state.Plans[1].AuthIndex {
		t.Fatalf("expected one plan per account: %#v", state.Plans)
	}

	second, err := bulkCreatePlans(accounts, spec)
	if err != nil {
		t.Fatal(err)
	}
	if second.Created != 0 || second.Skipped != 2 {
		t.Fatalf("unexpected dedupe result: %#v", second)
	}
	if len(snapshotPersistent().Plans) != 2 {
		t.Fatal("duplicate plans were created")
	}
}
