package main

import (
	"encoding/base64"
	"testing"
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
