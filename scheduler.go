package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
)

func schedulerLoop(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	// Also evaluate immediately after load; persisted minute keys prevent duplicate runs.
	evaluateSchedules(ctx, time.Now())
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			evaluateSchedules(ctx, now)
		}
	}
}

func evaluateSchedules(ctx context.Context, now time.Time) {
	state := snapshotPersistent()
	for _, p := range state.Plans {
		if !p.Enabled {
			continue
		}
		loc, err := time.LoadLocation(p.Timezone)
		if err != nil {
			continue
		}
		cron, err := parseCron(p.Cron)
		if err != nil {
			continue
		}
		local := now.In(loc)
		if !cron.matches(local) {
			continue
		}
		key := local.Format("2006-01-02T15:04")
		if !markScheduled(p.ID, key) {
			continue
		}
		if _, err := startPlanJob(p.ID, "scheduled"); err != nil {
			hostLog("warn", "scheduled test was not started", map[string]string{"plan_id": p.ID, "error": err.Error()})
		}
	}
}

func startPlanJob(planID, trigger string) (string, error) {
	p, ok := getPlan(planID)
	if !ok {
		return "", fmt.Errorf("plan not found")
	}
	job, err := beginJob("plan", p.ID, p.Model, "plan:"+p.ID)
	if err != nil {
		return "", err
	}
	go func() {
		errText := ""
		defer func() { finishJob(job.ID, "plan:"+p.ID, errText) }()

		a, err := findAccount(p.AuthIndex)
		if err != nil {
			errText = err.Error()
			entry := LogEntry{
				ID: newID("log"), At: time.Now(), Trigger: trigger,
				PlanID: p.ID, PlanName: p.Name, AuthIndex: p.AuthIndex,
				AccountName: firstNonEmpty(p.AccountName, p.AuthIndex), AccountEmail: p.AccountEmail,
				Model: p.Model, Status: "failed", Error: sanitizeError(errText),
			}
			appendLog(entry)
			updatePlanResult(p.ID, entry)
			updateJob(job.ID, func(j *Job) { j.Total = 1; j.Completed = 1; j.Failed = 1 })
			return
		}

		updateJob(job.ID, func(j *Job) { j.Total = 1 })
		entry := sendCodexRequest(a, p.Model, p.Prompt, trigger, p.ID, p.Name)
		appendLog(entry)
		updatePlanResult(p.ID, entry)
		updateJob(job.ID, func(j *Job) {
			j.Completed = 1
			classifyJobResult(j, entry)
		})
	}()
	return job.ID, nil
}

func startQuickJob(authIndex, model, prompt string) (string, error) {
	authIndex = strings.TrimSpace(authIndex)
	model = strings.TrimSpace(model)
	if authIndex == "" {
		return "", fmt.Errorf("account is required")
	}
	if model == "" {
		return "", fmt.Errorf("model is required")
	}
	job, err := beginJob("quick", "", model, "quick:"+authIndex)
	if err != nil {
		return "", err
	}
	go func() {
		errText := ""
		defer func() { finishJob(job.ID, "quick:"+authIndex, errText) }()
		a, err := findAccount(authIndex)
		if err != nil {
			errText = err.Error()
			updateJob(job.ID, func(j *Job) { j.Total = 1; j.Completed = 1; j.Failed = 1 })
			appendLog(LogEntry{ID: newID("log"), At: time.Now(), Trigger: "manual", AuthIndex: authIndex, AccountName: authIndex, Model: model, Status: "failed", Error: sanitizeError(errText)})
			return
		}
		updateJob(job.ID, func(j *Job) { j.Total = 1 })
		entry := sendCodexRequest(a, model, prompt, "manual", "", "")
		appendLog(entry)
		updateJob(job.ID, func(j *Job) { j.Completed = 1; classifyJobResult(j, entry) })
	}()
	return job.ID, nil
}

func startAllJob(model, prompt string) (string, error) {
	model = strings.TrimSpace(model)
	if model == "" {
		return "", fmt.Errorf("model is required")
	}
	job, err := beginJob("all-accounts", "", model, "all-accounts")
	if err != nil {
		return "", err
	}
	go func() {
		errText := ""
		defer func() { finishJob(job.ID, "all-accounts", errText) }()

		accounts, err := listCodexAccounts()
		if err != nil {
			errText = err.Error()
			updateJob(job.ID, func(j *Job) { j.Error = errText })
			return
		}
		updateJob(job.ID, func(j *Job) { j.Total = len(accounts) })
		if len(accounts) == 0 {
			return
		}

		rt.mu.RLock()
		concurrency := rt.state.Settings.BulkConcurrency
		rt.mu.RUnlock()
		if concurrency < 1 {
			concurrency = 1
		}

		sem := make(chan struct{}, concurrency)
		var wg sync.WaitGroup
		for _, account := range accounts {
			a := account
			wg.Add(1)
			go func() {
				defer wg.Done()
				sem <- struct{}{}
				entry := sendCodexRequest(a, model, prompt, "all-accounts", "", "")
				<-sem
				appendLog(entry)
				updateJob(job.ID, func(j *Job) {
					j.Completed++
					classifyJobResult(j, entry)
				})
			}()
		}
		wg.Wait()
	}()
	return job.ID, nil
}

func findAccount(authIndex string) (AuthFile, error) {
	accounts, err := listCodexAccounts()
	if err != nil {
		return AuthFile{}, err
	}
	for _, a := range accounts {
		if a.AuthIndex == authIndex {
			return a, nil
		}
	}
	return AuthFile{}, fmt.Errorf("account %s not found", authIndex)
}

func classifyJobResult(j *Job, entry LogEntry) {
	switch entry.Status {
	case "success":
		j.Succeeded++
	case "limited":
		j.Limited++
	default:
		j.Failed++
	}
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return "-"
}

func hostLog(level, message string, fields map[string]string) {
	_ = callHost("host.log", map[string]any{"level": level, "message": message, "fields": fields}, nil)
}
