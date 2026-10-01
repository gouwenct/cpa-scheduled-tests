package main

import "testing"

func TestRegistrationHasRequiredMetadataAndManagementCapability(t *testing.T) {
	reg := registration(6)
	if reg.SchemaVersion != 6 {
		t.Fatalf("schema version = %d, want 6", reg.SchemaVersion)
	}
	if reg.Metadata.Name == "" || reg.Metadata.Version == "" || reg.Metadata.Author == "" || reg.Metadata.GitHubRepository == "" {
		t.Fatalf("required metadata missing: %#v", reg.Metadata)
	}
	if !reg.Capabilities.ManagementAPI {
		t.Fatal("management_api capability must be enabled")
	}
}

func TestRegistrationSchemaNegotiation(t *testing.T) {
	cases := []struct {
		host uint32
		want int
	}{
		{0, 1},
		{1, 1},
		{5, 5},
		{6, 6},
		{7, 6},
	}
	for _, tc := range cases {
		if got := registration(tc.host).SchemaVersion; got != tc.want {
			t.Fatalf("registration(%d).SchemaVersion = %d, want %d", tc.host, got, tc.want)
		}
	}
}

func TestManagementRegistrationIncludesBulkCreate(t *testing.T) {
	mg := managementRegistrationResult()
	want := "/plugins/" + pluginID + "/plans/bulk-create"
	for _, r := range mg.Routes {
		if r.Path == want {
			return
		}
	}
	t.Fatalf("bulk create route %q not registered", want)
}
