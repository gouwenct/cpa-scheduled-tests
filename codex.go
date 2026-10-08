package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const codexResponsesURL = "https://chatgpt.com/backend-api/codex/responses"
const codexUsageURL = "https://chatgpt.com/backend-api/wham/usage"

type hostHTTPRequest struct {
	Method  string
	URL     string
	Headers map[string][]string
	Body    []byte
}

type hostHTTPResponse struct {
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

type codexQuotaInfo struct {
	Status        string
	UsedPercent   *float64
	ResetAt       string
	WindowMinutes int
}

type authMaterial struct {
	AccessToken string
	AccountID   string
}

type codexBody struct {
	Model        string         `json:"model"`
	Instructions string         `json:"instructions"`
	Input        []codexMessage `json:"input"`
	Store        bool           `json:"store"`
	Stream       bool           `json:"stream"`
}

type codexMessage struct {
	Type    string      `json:"type"`
	Role    string      `json:"role"`
	Content []codexPart `json:"content"`
}

type codexPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func listCodexAccounts() ([]AuthFile, error) {
	var resp struct {
		Files []AuthFile `json:"files"`
	}
	if err := callHost("host.auth.list", map[string]any{}, &resp); err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	out := make([]AuthFile, 0)
	for _, a := range resp.Files {
		if !isCodexAccount(a) {
			continue
		}
		key := a.AuthIndex
		if key == "" {
			key = "name:" + a.Name
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, a)
	}
	return out, nil
}

func isCodexAccount(a AuthFile) bool {
	fields := []string{a.Provider, a.Type, a.Name, a.Label}
	for _, v := range fields {
		v = strings.ToLower(strings.TrimSpace(v))
		if v == "codex" || strings.Contains(v, "codex") || strings.Contains(v, "chatgpt") {
			return true
		}
	}
	return false
}

func getAuthJSON(authIndex string) ([]byte, error) {
	var resp struct {
		AuthIndex string          `json:"auth_index"`
		Name      string          `json:"name"`
		JSON      json.RawMessage `json:"json"`
	}
	if err := callHost("host.auth.get", map[string]any{"auth_index": authIndex}, &resp); err != nil {
		return nil, err
	}
	if len(resp.JSON) == 0 {
		return nil, fmt.Errorf("empty auth JSON")
	}
	return resp.JSON, nil
}

func parseAuthMaterial(raw []byte) (authMaterial, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal(raw, &root); err != nil {
		return authMaterial{}, fmt.Errorf("invalid auth JSON")
	}
	token := firstString(root, "access_token", "accessToken", "oauth_access_token", "oauthAccessToken", "token")
	accountID := firstString(root, "account_id", "chatgpt_account_id", "accountId", "chatgptAccountId")
	idToken := firstString(root, "id_token", "idToken")

	for _, key := range []string{"tokens", "credentials", "auth", "oauth", "session"} {
		nestedRaw, ok := root[key]
		if !ok {
			continue
		}
		var nested map[string]json.RawMessage
		if json.Unmarshal(nestedRaw, &nested) != nil {
			continue
		}
		if token == "" {
			token = firstString(nested, "access_token", "accessToken", "oauth_access_token", "oauthAccessToken", "token")
		}
		if accountID == "" {
			accountID = firstString(nested, "account_id", "chatgpt_account_id", "accountId", "chatgptAccountId")
		}
		if idToken == "" {
			idToken = firstString(nested, "id_token", "idToken")
		}
	}
	if accountID == "" && idToken != "" {
		accountID = accountIDFromJWT(idToken)
	}
	if token == "" {
		return authMaterial{}, fmt.Errorf("missing access token")
	}
	return authMaterial{AccessToken: token, AccountID: accountID}, nil
}

func accountIDFromJWT(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return ""
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		if padded := parts[1] + strings.Repeat("=", (4-len(parts[1])%4)%4); padded != parts[1] {
			payload, err = base64.URLEncoding.DecodeString(padded)
		}
	}
	if err != nil {
		return ""
	}
	var claims map[string]any
	if json.Unmarshal(payload, &claims) != nil {
		return ""
	}
	for _, key := range []string{"chatgpt_account_id", "account_id"} {
		if v, ok := claims[key].(string); ok && strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	if raw, ok := claims["https://api.openai.com/auth"].(map[string]any); ok {
		for _, key := range []string{"chatgpt_account_id", "account_id"} {
			if v, ok := raw[key].(string); ok && strings.TrimSpace(v) != "" {
				return strings.TrimSpace(v)
			}
		}
	}
	return ""
}

func firstString(m map[string]json.RawMessage, keys ...string) string {
	for _, k := range keys {
		raw, ok := m[k]
		if !ok {
			continue
		}
		var s string
		if json.Unmarshal(raw, &s) == nil && strings.TrimSpace(s) != "" {
			return strings.TrimSpace(s)
		}
	}
	return ""
}

func accountDisplayName(a AuthFile) string {
	if strings.TrimSpace(a.Label) != "" {
		return a.Label
	}
	if strings.TrimSpace(a.Email) != "" {
		return a.Email
	}
	if strings.TrimSpace(a.Name) != "" {
		return a.Name
	}
	if strings.TrimSpace(a.ID) != "" {
		return a.ID
	}
	return a.AuthIndex
}

func sendCodexRequest(a AuthFile, model, prompt, trigger, planID, planName string) LogEntry {
	started := time.Now()
	entry := LogEntry{
		ID:           newID("log"),
		At:           started,
		Trigger:      trigger,
		PlanID:       planID,
		PlanName:     planName,
		AuthIndex:    a.AuthIndex,
		AccountName:  accountDisplayName(a),
		AccountEmail: a.Email,
		Disabled:     a.Disabled,
		Unavailable:  a.Unavailable,
		Model:        strings.TrimSpace(model),
	}

	if a.AuthIndex == "" {
		entry.Status = "failed"
		entry.Error = "missing auth_index"
		entry.LatencyMS = time.Since(started).Milliseconds()
		return entry
	}
	if entry.Model == "" {
		entry.Status = "failed"
		entry.Error = "model is empty"
		entry.LatencyMS = time.Since(started).Milliseconds()
		return entry
	}
	if strings.TrimSpace(prompt) == "" {
		prompt = "你好"
	}

	rawAuth, err := getAuthJSON(a.AuthIndex)
	if err != nil {
		entry.Status = "auth_error"
		entry.Error = "auth get: " + sanitizeError(err.Error())
		entry.LatencyMS = time.Since(started).Milliseconds()
		return entry
	}
	material, err := parseAuthMaterial(rawAuth)
	if err != nil {
		entry.Status = "auth_error"
		entry.Error = sanitizeError(err.Error())
		entry.LatencyMS = time.Since(started).Milliseconds()
		return entry
	}

	body, _ := json.Marshal(codexBody{
		Model:        entry.Model,
		Instructions: "Return a minimal response.",
		Input: []codexMessage{{
			Type:    "message",
			Role:    "user",
			Content: []codexPart{{Type: "input_text", Text: prompt}},
		}},
		Store:  false,
		Stream: true,
	})

	headers := map[string][]string{
		"Accept":        {"text/event-stream"},
		"Authorization": {"Bearer " + material.AccessToken},
		"Content-Type":  {"application/json"},
		"OpenAI-Beta":   {"responses=v1"},
		"originator":    {"codex_cli_rs"},
		"User-Agent":    {"codex_cli_rs/0.76.0"},
	}
	if material.AccountID != "" {
		headers["Chatgpt-Account-Id"] = []string{material.AccountID}
	}

	var resp hostHTTPResponse
	err = callHost("host.http.do", hostHTTPRequest{
		Method:  "POST",
		URL:     codexResponsesURL,
		Headers: headers,
		Body:    body,
	}, &resp)
	entry.LatencyMS = time.Since(started).Milliseconds()
	if err != nil {
		entry.Status = "network_error"
		entry.Error = sanitizeError(err.Error())
		return entry
	}

	entry.HTTPStatus = resp.StatusCode
	detail := parseUpstreamMessage(resp.Body)
	quota := parseCodexQuotaHeaders(resp.Headers, time.Now())
	if quota.Status == "unknown" {
		usageHeaders := map[string][]string{"Accept": {"application/json"}, "Authorization": {"Bearer " + material.AccessToken}}
		if material.AccountID != "" {
			usageHeaders["Chatgpt-Account-Id"] = []string{material.AccountID}
		}
		var usage hostHTTPResponse
		if callHost("host.http.do", hostHTTPRequest{Method: "GET", URL: codexUsageURL, Headers: usageHeaders}, &usage) == nil && usage.StatusCode == 200 {
			quota = parseCodexUsage(usage.Body, time.Now())
		}
	}
	entry.FiveHourStatus = quota.Status
	entry.FiveHourUsedPercent = quota.UsedPercent
	entry.FiveHourResetAt = quota.ResetAt
	entry.FiveHourWindowMinutes = quota.WindowMinutes
	switch {
	case resp.StatusCode >= 200 && resp.StatusCode <= 299:
		entry.Status = "success"
	case resp.StatusCode == 401:
		entry.Status = "auth_error"
		entry.Error = fallbackDetail(detail, "HTTP 401 authentication failed")
	case resp.StatusCode == 403:
		entry.Status = "forbidden"
		entry.Error = fallbackDetail(detail, "HTTP 403 forbidden")
	case resp.StatusCode == 404:
		entry.Status = "model_not_found"
		entry.Error = fallbackDetail(detail, "HTTP 404 model/endpoint not found")
	case resp.StatusCode == 408 || resp.StatusCode == 425 || resp.StatusCode == 429:
		entry.Status = "limited"
		entry.Error = fallbackDetail(detail, fmt.Sprintf("HTTP %d rate/temporary limit", resp.StatusCode))
	case resp.StatusCode >= 500:
		entry.Status = "upstream_error"
		entry.Error = fallbackDetail(detail, fmt.Sprintf("HTTP %d upstream error", resp.StatusCode))
	default:
		entry.Status = "failed"
		entry.Error = fallbackDetail(detail, fmt.Sprintf("HTTP %d", resp.StatusCode))
	}
	return entry
}

func parseCodexQuotaHeaders(headers map[string][]string, now time.Time) codexQuotaInfo {
	for _, window := range []string{"primary", "secondary"} {
		prefix := "x-codex-" + window + "-"
		minutes, _ := strconv.Atoi(strings.TrimSpace(firstHeader(headers, prefix+"window-minutes")))
		if minutes != 300 {
			continue
		}
		used, err := strconv.ParseFloat(strings.TrimSpace(firstHeader(headers, prefix+"used-percent")), 64)
		if err != nil {
			return codexQuotaInfo{Status: "unknown"}
		}
		reset, _ := strconv.ParseInt(strings.TrimSpace(firstHeader(headers, prefix+"reset-at")), 10, 64)
		return codexWindowInfo(&used, reset, now)
	}
	return codexQuotaInfo{Status: "unknown"}
}

func parseCodexUsage(body []byte, now time.Time) codexQuotaInfo {
	type window struct {
		Seconds int      `json:"limit_window_seconds"`
		Used    *float64 `json:"used_percent"`
		Reset   int64    `json:"reset_at"`
	}
	var usage struct {
		RateLimit struct {
			Primary   *window `json:"primary_window"`
			Secondary *window `json:"secondary_window"`
		} `json:"rate_limit"`
	}
	if json.Unmarshal(body, &usage) == nil {
		for _, w := range []*window{usage.RateLimit.Primary, usage.RateLimit.Secondary} {
			if w != nil && w.Seconds == 18000 {
				return codexWindowInfo(w.Used, w.Reset, now)
			}
		}
	}
	return codexQuotaInfo{Status: "unknown"}
}

func codexWindowInfo(used *float64, reset int64, now time.Time) codexQuotaInfo {
	info := codexQuotaInfo{Status: "unknown", WindowMinutes: 300}
	if used == nil || math.IsNaN(*used) || math.IsInf(*used, 0) || *used < 0 || *used > 100 {
		return info
	}
	info.UsedPercent = used
	if reset > 0 {
		info.ResetAt = time.Unix(reset, 0).UTC().Format(time.RFC3339)
	}
	if reset > now.Unix() && *used > 0 {
		info.Status = "open"
		if *used == 100 {
			info.Status = "exhausted"
		}
	}
	return info
}

func firstHeader(headers map[string][]string, name string) string {
	for key, values := range headers {
		if !strings.EqualFold(key, name) || len(values) == 0 {
			continue
		}
		return values[0]
	}
	return ""
}

func parseUpstreamMessage(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var root map[string]any
	if json.Unmarshal(body, &root) != nil {
		return ""
	}
	candidates := []any{root["error"], root["message"], root["detail"]}
	for _, c := range candidates {
		switch v := c.(type) {
		case string:
			if strings.TrimSpace(v) != "" {
				return sanitizeError(v)
			}
		case map[string]any:
			for _, k := range []string{"message", "type", "code"} {
				if s, ok := v[k].(string); ok && strings.TrimSpace(s) != "" {
					return sanitizeError(s)
				}
			}
		}
	}
	return ""
}

func sanitizeError(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(s, "\r", " "), "\n", " "))
	if len(s) > 400 {
		s = s[:400] + "..."
	}
	return s
}

func fallbackDetail(detail, fallback string) string {
	if strings.TrimSpace(detail) != "" {
		return detail
	}
	return fallback
}
