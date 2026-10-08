package main

import (
	"encoding/base64"
	"math"
	"testing"
	"time"
)

func TestParseAuthMaterialFlat(t *testing.T) {
	got, err := parseAuthMaterial([]byte(`{"access_token":"abc","account_id":"acct"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.AccessToken != "abc" || got.AccountID != "acct" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestParseAuthMaterialNested(t *testing.T) {
	got, err := parseAuthMaterial([]byte(`{"tokens":{"access_token":"abc","account_id":"acct"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.AccessToken != "abc" || got.AccountID != "acct" {
		t.Fatalf("unexpected: %#v", got)
	}
}

func TestCodexAccountDetection(t *testing.T) {
	if !isCodexAccount(AuthFile{Type: "codex"}) {
		t.Fatal("codex type not detected")
	}
	if !isCodexAccount(AuthFile{Name: "codex-user@example.com.json"}) {
		t.Fatal("codex filename not detected")
	}
	if isCodexAccount(AuthFile{Type: "claude"}) {
		t.Fatal("non-codex account detected")
	}
}

func TestParseAuthMaterialAccountIDFromJWT(t *testing.T) {
	payload := `{"https://api.openai.com/auth":{"chatgpt_account_id":"acct-from-jwt"}}`
	jwt := "x." + base64.RawURLEncoding.EncodeToString([]byte(payload)) + ".y"
	raw := []byte(`{"access_token":"abc","id_token":"` + jwt + `"}`)
	got, err := parseAuthMaterial(raw)
	if err != nil {
		t.Fatal(err)
	}
	if got.AccountID != "acct-from-jwt" {
		t.Fatalf("unexpected account id: %q", got.AccountID)
	}
}

func TestParseCodexQuotaHeaders(t *testing.T) {
	now := time.Unix(1699990000, 0)
	info := parseCodexQuotaHeaders(map[string][]string{
		"X-Codex-Primary-Used-Percent":   {"42.5"},
		"x-codex-primary-window-minutes": {"300"},
		"x-codex-primary-reset-at":       {"1700000000"},
	}, now)
	if info.Status != "open" || info.WindowMinutes != 300 || info.ResetAt == "" {
		t.Fatalf("unexpected quota info: %#v", info)
	}
	if info.UsedPercent == nil || math.Abs(*info.UsedPercent-42.5) > 0.001 {
		t.Fatalf("unexpected used percent: %#v", info.UsedPercent)
	}

	exhausted := parseCodexQuotaHeaders(map[string][]string{
		"x-codex-primary-used-percent":   {"100"},
		"x-codex-primary-window-minutes": {"300"},
		"x-codex-primary-reset-at":       {"1700000000"},
	}, now)
	if exhausted.Status != "exhausted" {
		t.Fatalf("expected exhausted status, got %#v", exhausted)
	}
}

func TestFiveHourEvidence(t *testing.T) {
	now := time.Unix(1699990000, 0)
	for _, tc := range []struct {
		name string
		body string
		want string
	}{
		{"active", `{"rate_limit":{"primary_window":{"limit_window_seconds":18000,"used_percent":1,"reset_at":1700000000}}}`, "open"},
		{"secondary five-hour", `{"rate_limit":{"primary_window":{"limit_window_seconds":604800,"used_percent":100,"reset_at":1700000000},"secondary_window":{"limit_window_seconds":18000,"used_percent":10,"reset_at":1700000000}}}`, "open"},
		{"exhausted", `{"rate_limit":{"primary_window":{"limit_window_seconds":18000,"used_percent":100,"reset_at":1700000000}}}`, "exhausted"},
		{"weekly only", `{"rate_limit":{"primary_window":{"limit_window_seconds":604800,"used_percent":100,"reset_at":1700000000}}}`, "unknown"},
		{"zero rounded usage", `{"rate_limit":{"primary_window":{"limit_window_seconds":18000,"used_percent":0,"reset_at":1700000000}}}`, "unknown"},
		{"expired", `{"rate_limit":{"primary_window":{"limit_window_seconds":18000,"used_percent":10,"reset_at":1699990000}}}`, "unknown"},
		{"missing usage", `{"rate_limit":{"primary_window":{"limit_window_seconds":18000,"reset_at":1700000000}}}`, "unknown"},
		{"invalid JSON", `not JSON`, "unknown"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := parseCodexUsage([]byte(tc.body), now); got.Status != tc.want {
				t.Fatalf("got %s, want %s", got.Status, tc.want)
			}
		})
	}
	for _, headers := range []map[string][]string{
		nil,
		{"x-codex-primary-used-percent": {"100"}},
		{"x-codex-primary-window-minutes": {"10080"}, "x-codex-primary-used-percent": {"100"}},
		{"x-codex-primary-window-minutes": {"300"}, "x-codex-primary-used-percent": {"NaN"}},
	} {
		if got := parseCodexQuotaHeaders(headers, now); got.Status != "unknown" {
			t.Fatalf("insufficient header evidence: %#v", got)
		}
	}
}
