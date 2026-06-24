package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestDefaultDomainGlobalSetsCurrentDomainBeforeCommand(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		domains: []api.Domain{{
			ID:   "domain-1",
			Name: "Acme",
		}},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"--default-domain", "domain-1",
		"domains", "list",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.setCurrentDomainID != "domain-1" {
		t.Fatalf("setCurrentDomainID = %q", client.setCurrentDomainID)
	}
}

func TestSessionLanguageGlobalSetsLanguageBeforeCommand(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		language:  "en",
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"--session-language", "en",
		"language", "current",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.setLanguage != "en" {
		t.Fatalf("setLanguage = %q", client.setLanguage)
	}
}

func TestDomainsSetCurrentJSONCallsSetCurrentDomain(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"domains", "set-current",
		"--domain", "domain-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.setCurrentDomainID != "domain-1" {
		t.Fatalf("setCurrentDomainID = %q", client.setCurrentDomainID)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v", err)
	}
	if payload["domain_id"] != "domain-1" || payload["message"] != "current domain set for this session" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestLanguageSetJSONCallsSetLanguage(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"language", "set",
		"--language", "en",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.setLanguage != "en" {
		t.Fatalf("setLanguage = %q", client.setLanguage)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v", err)
	}
	if payload["language"] != "en" || payload["message"] != "language set for this session" {
		t.Fatalf("payload = %#v", payload)
	}
}
