package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type planView struct {
	Plan
	NextRun *time.Time `json:"next_run,omitempty"`
}

type stateResponse struct {
	Version  string     `json:"version"`
	DataDir  string     `json:"data_dir"`
	Settings Settings   `json:"settings"`
	Plans    []planView `json:"plans"`
	Accounts []AuthFile `json:"accounts"`
	Jobs     []Job      `json:"jobs"`
}

type runRequest struct {
	PlanID    string `json:"plan_id,omitempty"`
	AuthIndex string `json:"auth_index,omitempty"`
	Model     string `json:"model,omitempty"`
	Prompt    string `json:"prompt,omitempty"`
}

type deletePlanRequest struct {
	ID string `json:"id"`
}

func handleManagement(raw []byte) []byte {
	var req managementRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		return failEnvelope("invalid_request", "invalid management request")
	}
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	path := strings.TrimSpace(req.Path)

	switch {
	case method == http.MethodGet && strings.Contains(path, "/v0/resource/plugins/") && strings.HasSuffix(path, "/panel"):
		return okEnvelope(htmlResponse(panelHTML()))

	case method == http.MethodGet && strings.HasSuffix(path, "/plugins/"+pluginID+"/state"):
		state, err := buildStateResponse()
		if err != nil {
			return managementError(500, err)
		}
		return okEnvelope(jsonResponse(200, state))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/plans/save"):
		var p Plan
		if err := json.Unmarshal(req.Body, &p); err != nil {
			return managementError(400, fmt.Errorf("invalid plan JSON"))
		}
		accounts, _ := listCodexAccounts()
		for _, a := range accounts {
			if a.AuthIndex == p.AuthIndex {
				p.AccountName = accountDisplayName(a)
				p.AccountEmail = a.Email
				break
			}
		}
		saved, err := upsertPlan(p)
		if err != nil {
			return managementError(400, err)
		}
		return okEnvelope(jsonResponse(200, saved))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/plans/bulk-create"):
		var spec BulkPlanSpec
		if json.Unmarshal(req.Body, &spec) != nil {
			return managementError(400, fmt.Errorf("invalid bulk plan JSON"))
		}
		accounts, err := listCodexAccounts()
		if err != nil {
			return managementError(500, err)
		}
		if len(accounts) == 0 {
			return managementError(400, fmt.Errorf("no Codex accounts found"))
		}
		result, err := bulkCreatePlans(accounts, spec)
		if err != nil {
			return managementError(400, err)
		}
		return okEnvelope(jsonResponse(200, result))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/plans/delete"):
		var body deletePlanRequest
		if json.Unmarshal(req.Body, &body) != nil {
			return managementError(400, fmt.Errorf("invalid request"))
		}
		if err := deletePlan(body.ID); err != nil {
			return managementError(404, err)
		}
		return okEnvelope(jsonResponse(200, map[string]any{"deleted": true, "id": body.ID}))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/settings"):
		var settings Settings
		if json.Unmarshal(req.Body, &settings) != nil {
			return managementError(400, fmt.Errorf("invalid settings JSON"))
		}
		if err := saveSettings(settings); err != nil {
			return managementError(400, err)
		}
		return okEnvelope(jsonResponse(200, settings))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/run"):
		var body runRequest
		if json.Unmarshal(req.Body, &body) != nil {
			return managementError(400, fmt.Errorf("invalid run request"))
		}
		var jobID string
		var err error
		if strings.TrimSpace(body.PlanID) != "" {
			jobID, err = startPlanJob(body.PlanID, "manual")
		} else {
			jobID, err = startQuickJob(body.AuthIndex, body.Model, body.Prompt)
		}
		if err != nil {
			return managementError(409, err)
		}
		return okEnvelope(jsonResponse(202, map[string]any{"accepted": true, "job_id": jobID}))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/run-all"):
		var body runRequest
		if json.Unmarshal(req.Body, &body) != nil {
			return managementError(400, fmt.Errorf("invalid run-all request"))
		}
		jobID, err := startAllJob(body.Model, body.Prompt)
		if err != nil {
			return managementError(409, err)
		}
		return okEnvelope(jsonResponse(202, map[string]any{"accepted": true, "job_id": jobID, "includes_disabled": true, "includes_unavailable": true}))

	case method == http.MethodGet && strings.HasSuffix(path, "/plugins/"+pluginID+"/logs"):
		limit := 200
		if values := req.Query["limit"]; len(values) > 0 {
			if n, err := strconv.Atoi(values[0]); err == nil {
				limit = n
			}
		}
		logs, err := listLogs(limit)
		if err != nil {
			return managementError(500, err)
		}
		return okEnvelope(jsonResponse(200, map[string]any{"logs": logs}))

	case method == http.MethodPost && strings.HasSuffix(path, "/plugins/"+pluginID+"/logs/clear"):
		if err := clearLogs(); err != nil {
			return managementError(500, err)
		}
		return okEnvelope(jsonResponse(200, map[string]any{"cleared": true}))

	default:
		return okEnvelope(jsonResponse(404, map[string]string{"error": "not found"}))
	}
}

func buildStateResponse() (stateResponse, error) {
	state := snapshotPersistent()
	accounts, err := listCodexAccounts()
	if err != nil {
		return stateResponse{}, err
	}
	sort.Slice(accounts, func(i, j int) bool { return accountDisplayName(accounts[i]) < accountDisplayName(accounts[j]) })

	views := make([]planView, 0, len(state.Plans))
	now := time.Now()
	for _, p := range state.Plans {
		view := planView{Plan: p}
		if p.Enabled {
			if t, err := nextCronTime(p.Cron, p.Timezone, now); err == nil {
				view.NextRun = &t
			}
		}
		views = append(views, view)
	}
	sort.Slice(views, func(i, j int) bool { return views[i].CreatedAt.Before(views[j].CreatedAt) })

	rt.mu.RLock()
	dataDir := rt.dataDir
	rt.mu.RUnlock()

	return stateResponse{
		Version:  pluginVersion,
		DataDir:  dataDir,
		Settings: state.Settings,
		Plans:    views,
		Accounts: accounts,
		Jobs:     jobsSnapshot(),
	}, nil
}

func jsonResponse(status int, body any) managementResponse {
	raw, _ := json.Marshal(body)
	return managementResponse{
		StatusCode: status,
		Headers: map[string][]string{
			"content-type":  {"application/json; charset=utf-8"},
			"cache-control": {"no-store"},
		},
		Body: raw,
	}
}

func htmlResponse(body string) managementResponse {
	return managementResponse{
		StatusCode: 200,
		Headers: map[string][]string{
			"content-type":  {"text/html; charset=utf-8"},
			"cache-control": {"no-store"},
		},
		Body: []byte(body),
	}
}

func managementError(status int, err error) []byte {
	return okEnvelope(jsonResponse(status, map[string]string{"error": err.Error()}))
}
