package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"
)

type Settings struct {
	DefaultTimezone string `json:"default_timezone"`
	DefaultModel    string `json:"default_model"`
	DefaultPrompt   string `json:"default_prompt"`
	BulkConcurrency int    `json:"bulk_concurrency"`
	LogRetention    int    `json:"log_retention"`
}

type Plan struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Group            string    `json:"group,omitempty"`
	Enabled          bool      `json:"enabled"`
	AuthIndex        string    `json:"auth_index"`
	AccountName      string    `json:"account_name,omitempty"`
	AccountEmail     string    `json:"account_email,omitempty"`
	Model            string    `json:"model"`
	Cron             string    `json:"cron"`
	Timezone         string    `json:"timezone"`
	Prompt           string    `json:"prompt"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	LastRunAt        time.Time `json:"last_run_at,omitempty"`
	LastStatus       string    `json:"last_status,omitempty"`
	LastHTTPStatus   int       `json:"last_http_status,omitempty"`
	LastLatencyMS    int64     `json:"last_latency_ms,omitempty"`
	LastError        string    `json:"last_error,omitempty"`
	LastScheduledKey string    `json:"last_scheduled_key,omitempty"`
}

type PersistentState struct {
	SchemaVersion int      `json:"schema_version"`
	Settings      Settings `json:"settings"`
	Plans         []Plan   `json:"plans"`
}

type BulkPlanSpec struct {
	NamePrefix string `json:"name_prefix"`
	Group      string `json:"group,omitempty"`
	Model      string `json:"model"`
	Cron       string `json:"cron"`
	Timezone   string `json:"timezone"`
	Prompt     string `json:"prompt"`
	Enabled    bool   `json:"enabled"`
}

type BulkPlanResult struct {
	TotalAccounts int    `json:"total_accounts"`
	Created       int    `json:"created"`
	Skipped       int    `json:"skipped"`
	CreatedPlans  []Plan `json:"created_plans"`
}

type AuthFile struct {
	ID            string `json:"id,omitempty"`
	AuthIndex     string `json:"auth_index,omitempty"`
	Name          string `json:"name"`
	Type          string `json:"type,omitempty"`
	Provider      string `json:"provider,omitempty"`
	Label         string `json:"label,omitempty"`
	Email         string `json:"email,omitempty"`
	Account       string `json:"account,omitempty"`
	Status        string `json:"status,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
	Disabled      bool   `json:"disabled,omitempty"`
	Unavailable   bool   `json:"unavailable,omitempty"`
	RuntimeOnly   bool   `json:"runtime_only,omitempty"`
}

type LogEntry struct {
	ID                    string    `json:"id"`
	At                    time.Time `json:"at"`
	Trigger               string    `json:"trigger"`
	PlanID                string    `json:"plan_id,omitempty"`
	PlanName              string    `json:"plan_name,omitempty"`
	AuthIndex             string    `json:"auth_index,omitempty"`
	AccountName           string    `json:"account_name"`
	AccountEmail          string    `json:"account_email,omitempty"`
	Disabled              bool      `json:"disabled,omitempty"`
	Unavailable           bool      `json:"unavailable,omitempty"`
	Model                 string    `json:"model"`
	Status                string    `json:"status"`
	HTTPStatus            int       `json:"http_status,omitempty"`
	FiveHourStatus        string    `json:"five_hour_status,omitempty"`
	FiveHourUsedPercent   *float64  `json:"five_hour_used_percent,omitempty"`
	FiveHourResetAt       string    `json:"five_hour_reset_at,omitempty"`
	FiveHourWindowMinutes int       `json:"five_hour_window_minutes,omitempty"`
	LatencyMS             int64     `json:"latency_ms"`
	Error                 string    `json:"error,omitempty"`
}

type Job struct {
	ID         string    `json:"id"`
	Kind       string    `json:"kind"`
	PlanID     string    `json:"plan_id,omitempty"`
	Model      string    `json:"model"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	Running    bool      `json:"running"`
	Total      int       `json:"total"`
	Completed  int       `json:"completed"`
	Succeeded  int       `json:"succeeded"`
	Failed     int       `json:"failed"`
	Limited    int       `json:"limited"`
	Error      string    `json:"error,omitempty"`
}

type runtimeState struct {
	mu          sync.RWMutex
	initialized bool
	dataDir     string
	state       PersistentState
	jobs        map[string]*Job
	activeKeys  map[string]string
	ctx         context.Context
	cancel      context.CancelFunc
}

var rt runtimeState
var logMu sync.Mutex

const logRetentionWindow = 48 * time.Hour

func defaultSettings() Settings {
	return Settings{
		DefaultTimezone: "Asia/Shanghai",
		DefaultModel:    "gpt-5.6-luna",
		DefaultPrompt:   "你好",
		BulkConcurrency: 4,
		LogRetention:    5000,
	}
}

func ensureRuntime(pluginDir string) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if rt.initialized {
		return nil
	}

	dataDir, err := resolveDataDir(pluginDir)
	if err != nil {
		return err
	}

	state := PersistentState{SchemaVersion: 1, Settings: defaultSettings(), Plans: []Plan{}}
	statePath := filepath.Join(dataDir, "state.json")
	if raw, err := os.ReadFile(statePath); err == nil {
		if err := json.Unmarshal(raw, &state); err != nil {
			return fmt.Errorf("decode %s: %w", statePath, err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("read %s: %w", statePath, err)
	}
	normalizeState(&state)

	ctx, cancel := context.WithCancel(context.Background())
	rt.dataDir = dataDir
	rt.state = state
	rt.jobs = make(map[string]*Job)
	rt.activeKeys = make(map[string]string)
	rt.ctx = ctx
	rt.cancel = cancel
	rt.initialized = true

	go schedulerLoop(ctx)
	hostLog("info", "CPA Scheduled Tests started", map[string]string{"data_dir": dataDir})
	return nil
}

func stopRuntime() {
	rt.mu.Lock()
	if rt.cancel != nil {
		rt.cancel()
	}
	rt.cancel = nil
	rt.ctx = nil
	rt.initialized = false
	rt.mu.Unlock()
}

func resolveDataDir(pluginDir string) (string, error) {
	candidates := make([]string, 0, 3)
	if strings.TrimSpace(pluginDir) != "" {
		candidates = append(candidates, filepath.Join(pluginDir, ".data", pluginID))
	}
	if cfg, err := os.UserConfigDir(); err == nil && cfg != "" {
		candidates = append(candidates, filepath.Join(cfg, "CLIProxyAPI", "plugin-data", pluginID))
	}
	if exe, err := os.Executable(); err == nil && exe != "" {
		candidates = append(candidates, filepath.Join(filepath.Dir(exe), "plugin-data", pluginID))
	}
	if len(candidates) == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		candidates = append(candidates, filepath.Join(cwd, "plugin-data", pluginID))
	}
	var lastErr error
	for _, dir := range candidates {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			lastErr = err
			continue
		}
		probe := filepath.Join(dir, ".write-test")
		if err := os.WriteFile(probe, []byte("ok"), 0o600); err != nil {
			lastErr = err
			continue
		}
		_ = os.Remove(probe)
		return dir, nil
	}
	return "", fmt.Errorf("no writable plugin data directory: %w", lastErr)
}

func normalizeState(s *PersistentState) {
	if s.SchemaVersion == 0 {
		s.SchemaVersion = 1
	}
	ds := defaultSettings()
	if strings.TrimSpace(s.Settings.DefaultTimezone) == "" {
		s.Settings.DefaultTimezone = ds.DefaultTimezone
	}
	if strings.TrimSpace(s.Settings.DefaultModel) == "" {
		s.Settings.DefaultModel = ds.DefaultModel
	}
	if strings.TrimSpace(s.Settings.DefaultPrompt) == "" {
		s.Settings.DefaultPrompt = ds.DefaultPrompt
	}
	if s.Settings.BulkConcurrency < 1 || s.Settings.BulkConcurrency > 32 {
		s.Settings.BulkConcurrency = ds.BulkConcurrency
	}
	if s.Settings.LogRetention < 100 || s.Settings.LogRetention > 100000 {
		s.Settings.LogRetention = ds.LogRetention
	}
	if s.Plans == nil {
		s.Plans = []Plan{}
	}
	for i := range s.Plans {
		if strings.TrimSpace(s.Plans[i].Group) == "" {
			group := strings.TrimSpace(s.Plans[i].Name)
			if before, _, ok := strings.Cut(group, " · "); ok {
				group = before
			}
			if group == "" {
				group = "默认"
			}
			s.Plans[i].Group = group
		}
	}
}

func snapshotPersistent() PersistentState {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	out := rt.state
	out.Plans = append([]Plan(nil), rt.state.Plans...)
	return out
}

func saveStateLocked() error {
	raw, err := json.MarshalIndent(rt.state, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(rt.dataDir, "state.json")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func saveSettings(settings Settings) error {
	settings.DefaultTimezone = strings.TrimSpace(settings.DefaultTimezone)
	settings.DefaultModel = strings.TrimSpace(settings.DefaultModel)
	settings.DefaultPrompt = strings.TrimSpace(settings.DefaultPrompt)
	if settings.DefaultTimezone == "" {
		settings.DefaultTimezone = "Asia/Shanghai"
	}
	if _, err := time.LoadLocation(settings.DefaultTimezone); err != nil {
		return fmt.Errorf("invalid timezone: %s", settings.DefaultTimezone)
	}
	if settings.DefaultModel == "" {
		return fmt.Errorf("default model is required")
	}
	if settings.DefaultPrompt == "" {
		settings.DefaultPrompt = "你好"
	}
	if settings.BulkConcurrency < 1 || settings.BulkConcurrency > 32 {
		return fmt.Errorf("bulk_concurrency must be between 1 and 32")
	}
	if settings.LogRetention < 100 || settings.LogRetention > 100000 {
		return fmt.Errorf("log_retention must be between 100 and 100000")
	}

	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.state.Settings = settings
	return saveStateLocked()
}

func validatePlan(p *Plan) error {
	p.Name = strings.TrimSpace(p.Name)
	p.Group = strings.TrimSpace(p.Group)
	p.AuthIndex = strings.TrimSpace(p.AuthIndex)
	p.Model = strings.TrimSpace(p.Model)
	p.Cron = strings.TrimSpace(p.Cron)
	p.Timezone = strings.TrimSpace(p.Timezone)
	p.Prompt = strings.TrimSpace(p.Prompt)

	if p.Group == "" {
		p.Group = "默认"
	}
	if p.Name == "" {
		p.Name = p.Group
	}
	if p.AuthIndex == "" {
		return fmt.Errorf("account is required")
	}
	if p.Model == "" {
		return fmt.Errorf("model is required")
	}
	if p.Cron == "" {
		return fmt.Errorf("cron expression is required")
	}
	if p.Timezone == "" {
		p.Timezone = defaultSettings().DefaultTimezone
	}
	if _, err := time.LoadLocation(p.Timezone); err != nil {
		return fmt.Errorf("invalid timezone %q", p.Timezone)
	}
	if _, err := parseCron(p.Cron); err != nil {
		return err
	}
	if p.Prompt == "" {
		p.Prompt = "你好"
	}
	return nil
}

func upsertPlan(p Plan) (Plan, error) {
	if err := validatePlan(&p); err != nil {
		return Plan{}, err
	}
	now := time.Now()

	rt.mu.Lock()
	defer rt.mu.Unlock()

	if p.ID == "" {
		p.ID = newID("plan")
		p.CreatedAt = now
	}
	p.UpdatedAt = now

	for i := range rt.state.Plans {
		if rt.state.Plans[i].ID == p.ID {
			old := rt.state.Plans[i]
			if p.CreatedAt.IsZero() {
				p.CreatedAt = old.CreatedAt
			}
			p.LastRunAt = old.LastRunAt
			p.LastStatus = old.LastStatus
			p.LastHTTPStatus = old.LastHTTPStatus
			p.LastLatencyMS = old.LastLatencyMS
			p.LastError = old.LastError
			p.LastScheduledKey = old.LastScheduledKey
			rt.state.Plans[i] = p
			if err := saveStateLocked(); err != nil {
				return Plan{}, err
			}
			return p, nil
		}
	}

	rt.state.Plans = append(rt.state.Plans, p)
	if err := saveStateLocked(); err != nil {
		return Plan{}, err
	}
	return p, nil
}

func planDedupeKey(p Plan) string {
	return strings.Join([]string{
		strings.TrimSpace(p.AuthIndex),
		strings.ToLower(strings.TrimSpace(p.Model)),
		strings.TrimSpace(p.Cron),
		strings.ToLower(strings.TrimSpace(p.Timezone)),
		strings.TrimSpace(p.Prompt),
	}, "\x1f")
}

func bulkCreatePlans(accounts []AuthFile, spec BulkPlanSpec) (BulkPlanResult, error) {
	spec.Group = strings.TrimSpace(spec.Group)
	spec.NamePrefix = strings.TrimSpace(spec.NamePrefix)
	if spec.Group == "" {
		spec.Group = spec.NamePrefix
	}
	spec.Model = strings.TrimSpace(spec.Model)
	spec.Cron = strings.TrimSpace(spec.Cron)
	spec.Timezone = strings.TrimSpace(spec.Timezone)
	spec.Prompt = strings.TrimSpace(spec.Prompt)
	if spec.Group == "" {
		spec.Group = "全账号定时"
	}
	if spec.Timezone == "" {
		spec.Timezone = defaultSettings().DefaultTimezone
	}
	if spec.Prompt == "" {
		spec.Prompt = "你好"
	}

	// Validate the shared schedule once before mutating persistent state.
	template := Plan{
		Name:      spec.Group,
		Group:     spec.Group,
		AuthIndex: "validation-placeholder",
		Model:     spec.Model,
		Cron:      spec.Cron,
		Timezone:  spec.Timezone,
		Prompt:    spec.Prompt,
		Enabled:   spec.Enabled,
	}
	if err := validatePlan(&template); err != nil {
		return BulkPlanResult{}, err
	}

	result := BulkPlanResult{TotalAccounts: len(accounts), CreatedPlans: []Plan{}}
	now := time.Now()

	rt.mu.Lock()
	defer rt.mu.Unlock()

	existing := make(map[string]struct{}, len(rt.state.Plans)+len(accounts))
	for _, p := range rt.state.Plans {
		existing[planDedupeKey(p)] = struct{}{}
	}

	for _, a := range accounts {
		authIndex := strings.TrimSpace(a.AuthIndex)
		if authIndex == "" {
			result.Skipped++
			continue
		}
		label := accountDisplayName(a)
		if strings.TrimSpace(label) == "" {
			label = authIndex
		}
		p := Plan{
			ID:           newID("plan"),
			Name:         spec.Group + " · " + label,
			Group:        spec.Group,
			Enabled:      spec.Enabled,
			AuthIndex:    authIndex,
			AccountName:  label,
			AccountEmail: a.Email,
			Model:        spec.Model,
			Cron:         spec.Cron,
			Timezone:     spec.Timezone,
			Prompt:       spec.Prompt,
			CreatedAt:    now,
			UpdatedAt:    now,
		}
		if err := validatePlan(&p); err != nil {
			return BulkPlanResult{}, err
		}
		key := planDedupeKey(p)
		if _, ok := existing[key]; ok {
			result.Skipped++
			continue
		}
		existing[key] = struct{}{}
		rt.state.Plans = append(rt.state.Plans, p)
		result.CreatedPlans = append(result.CreatedPlans, p)
		result.Created++
	}

	if result.Created > 0 {
		if err := saveStateLocked(); err != nil {
			return BulkPlanResult{}, err
		}
	}
	return result, nil
}

func deletePlan(id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("plan id is required")
	}
	rt.mu.Lock()
	defer rt.mu.Unlock()
	out := rt.state.Plans[:0]
	found := false
	for _, p := range rt.state.Plans {
		if p.ID == id {
			found = true
			continue
		}
		out = append(out, p)
	}
	if !found {
		return fmt.Errorf("plan not found")
	}
	rt.state.Plans = append([]Plan(nil), out...)
	return saveStateLocked()
}

func getPlan(id string) (Plan, bool) {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	for _, p := range rt.state.Plans {
		if p.ID == id {
			return p, true
		}
	}
	return Plan{}, false
}

func updatePlanResult(planID string, entry LogEntry) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	for i := range rt.state.Plans {
		if rt.state.Plans[i].ID != planID {
			continue
		}
		p := &rt.state.Plans[i]
		p.LastRunAt = entry.At
		p.LastStatus = entry.Status
		p.LastHTTPStatus = entry.HTTPStatus
		p.LastLatencyMS = entry.LatencyMS
		p.LastError = entry.Error
		p.UpdatedAt = time.Now()
		_ = saveStateLocked()
		return
	}
}

func markScheduled(planID, key string) bool {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	for i := range rt.state.Plans {
		p := &rt.state.Plans[i]
		if p.ID != planID {
			continue
		}
		if p.LastScheduledKey == key {
			return false
		}
		p.LastScheduledKey = key
		_ = saveStateLocked()
		return true
	}
	return false
}

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d-%06d", prefix, time.Now().UnixNano(), rand.Intn(1000000))
}

func logPath() string {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	return filepath.Join(rt.dataDir, "runs.jsonl")
}

func appendLog(entry LogEntry) {
	logMu.Lock()
	defer logMu.Unlock()
	if entry.ID == "" {
		entry.ID = newID("log")
	}
	if entry.At.IsZero() {
		entry.At = time.Now()
	}

	raw, err := json.Marshal(entry)
	if err != nil {
		return
	}
	path := logPath()
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err == nil {
		_, _ = f.Write(append(raw, '\n'))
		_ = f.Close()
	}
	compactLogsIfNeededUnlocked()
}

func listLogs(limit int) ([]LogEntry, error) {
	logMu.Lock()
	defer logMu.Unlock()
	return listLogsUnlocked(limit)
}

func listLogsUnlocked(limit int) ([]LogEntry, error) {
	if limit <= 0 {
		limit = 200
	}
	if limit > 2000 {
		limit = 2000
	}
	path := logPath()
	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return []LogEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	logs := make([]LogEntry, 0, limit)
	cutoff := time.Now().Add(-logRetentionWindow)
	scanner := bufio.NewScanner(f)
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		var e LogEntry
		if json.Unmarshal(scanner.Bytes(), &e) != nil {
			continue
		}
		if e.At.Before(cutoff) {
			continue
		}
		logs = append(logs, e)
		if len(logs) > limit {
			copy(logs, logs[len(logs)-limit:])
			logs = logs[:limit]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}
	return logs, nil
}

func clearLogs() error {
	logMu.Lock()
	defer logMu.Unlock()
	path := logPath()
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}

func compactLogsIfNeededUnlocked() {
	rt.mu.RLock()
	retention := rt.state.Settings.LogRetention
	rt.mu.RUnlock()
	if retention < 100 {
		retention = 5000
	}
	path := logPath()
	info, err := os.Stat(path)
	if err != nil || info.Size() < 8*1024*1024 {
		return
	}
	logs, err := listLogsUnlocked(retention)
	if err != nil {
		return
	}
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}
	tmp := path + ".compact"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	enc := json.NewEncoder(f)
	for _, e := range logs {
		_ = enc.Encode(e)
	}
	_ = f.Close()
	_ = os.Rename(tmp, path)
}

func jobsSnapshot() []Job {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	out := make([]Job, 0, len(rt.jobs))
	for _, j := range rt.jobs {
		out = append(out, *j)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].StartedAt.After(out[j].StartedAt) })
	if len(out) > 20 {
		out = out[:20]
	}
	return out
}

func beginJob(kind, planID, model, activeKey string) (*Job, error) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if activeKey != "" {
		if existing, ok := rt.activeKeys[activeKey]; ok {
			return nil, fmt.Errorf("a job is already running (%s)", existing)
		}
	}
	id := newID("job")
	j := &Job{ID: id, Kind: kind, PlanID: planID, Model: model, StartedAt: time.Now(), Running: true}
	rt.jobs[id] = j
	if activeKey != "" {
		rt.activeKeys[activeKey] = id
	}
	trimJobsLocked()
	return j, nil
}

func finishJob(jobID, activeKey, errText string) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if j, ok := rt.jobs[jobID]; ok {
		j.Running = false
		j.FinishedAt = time.Now()
		j.Error = errText
	}
	if activeKey != "" {
		if rt.activeKeys[activeKey] == jobID {
			delete(rt.activeKeys, activeKey)
		}
	}
}

func updateJob(jobID string, fn func(*Job)) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	if j, ok := rt.jobs[jobID]; ok {
		fn(j)
	}
}

func trimJobsLocked() {
	if len(rt.jobs) <= 50 {
		return
	}
	type pair struct {
		id      string
		at      time.Time
		running bool
	}
	xs := make([]pair, 0, len(rt.jobs))
	for id, j := range rt.jobs {
		xs = append(xs, pair{id: id, at: j.StartedAt, running: j.Running})
	}
	sort.Slice(xs, func(i, j int) bool { return xs[i].at.Before(xs[j].at) })
	for _, x := range xs {
		if len(rt.jobs) <= 35 {
			break
		}
		if !x.running {
			delete(rt.jobs, x.id)
		}
	}
}
