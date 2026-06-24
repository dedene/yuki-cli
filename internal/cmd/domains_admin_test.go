package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestDomainsCreateDryRunSkipsAuthentication(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"domains", "create",
		"--administration-name", "Highpro BV",
		"--domain-name", "Highpro",
		"--default-language", "nl-be",
		"--dry-run",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.accessKey != "" || client.domainCreateOperation != "" {
		t.Fatalf("dry run authenticated or sent operation: access=%q operation=%q", client.accessKey, client.domainCreateOperation)
	}
	var result api.DomainAdminResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json: %v", err)
	}
	if result.Operation != "CreateDomain" || !result.DryRun || result.DomainName != "Highpro" {
		t.Fatalf("result = %#v", result)
	}
}

func TestDomainsCreateReadonlyBlocksBeforeAuthentication(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--readonly",
		"domains", "create",
		"--administration-name", "Highpro BV",
		"--domain-name", "Highpro",
		"--default-language", "nl-be",
	}, Runtime{
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "--readonly blocks mutating command: domains create") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" || client.domainCreateOperation != "" {
		t.Fatalf("readonly authenticated or sent operation: access=%q operation=%q", client.accessKey, client.domainCreateOperation)
	}
}

func TestDomainsCreateTrialSendsFields(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		domainAdminResult: api.DomainAdminResult{
			Operation: "CreateTrialDomain",
			Message:   "ok",
		},
	}

	err := Execute(context.Background(), []string{
		"domains", "create-trial",
		"--administration-name", "Highpro BV",
		"--domain-name", "Trial Highpro",
		"--default-language", "en",
	}, Runtime{
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.domainCreateOperation != "CreateTrialDomain" ||
		client.domainCreateOpts.AdministrationName != "Highpro BV" ||
		client.domainCreateOpts.DomainName != "Trial Highpro" ||
		client.domainCreateOpts.DefaultLanguage != "en" {
		t.Fatalf("domain create = %q %#v", client.domainCreateOperation, client.domainCreateOpts)
	}
}

func TestDomainsAddUserSendsFields(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"domains", "add-user",
		"--domain", "domain-1",
		"--email", "peter@example.com",
		"--name", "Peter Dedene",
		"--roles", "Backoffice",
		"--administrations", "admin-1",
		"--send-message",
		"--custom-message", "Welcome",
		"--language", "en",
	}, Runtime{
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.domainUserAddOpts.DomainID != "domain-1" ||
		client.domainUserAddOpts.Email != "peter@example.com" ||
		client.domainUserAddOpts.Roles != "Backoffice" ||
		!client.domainUserAddOpts.SendMessage ||
		client.domainUserAddOpts.CustomMessage != "Welcome" ||
		client.domainUserAddOpts.Language != "en" {
		t.Fatalf("domainUserAddOpts = %#v", client.domainUserAddOpts)
	}
}

func TestDomainsAddUserDryRunJSONIncludesAllFields(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"domains", "add-user",
		"--domain", "domain-1",
		"--email", "peter@example.com",
		"--name", "Peter Dedene",
		"--roles", "Backoffice",
		"--administrations", "admin-1",
		"--send-message",
		"--custom-message", "Welcome",
		"--language", "en",
		"--dry-run",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.accessKey != "" || client.domainUserAddOpts.DomainID != "" {
		t.Fatalf("dry run authenticated or sent operation: access=%q opts=%#v", client.accessKey, client.domainUserAddOpts)
	}
	var result api.DomainAdminResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json: %v", err)
	}
	if result.Administrations != "admin-1" ||
		result.SendMessage == nil ||
		!*result.SendMessage ||
		result.CustomMessage != "Welcome" ||
		result.Language != "en" {
		t.Fatalf("result = %#v", result)
	}
}

func TestDomainsLyantheSendsEnabled(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"domains", "lyanthe",
		"--domain", "domain-1",
		"--enabled", "true",
	}, Runtime{
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.lyantheDomainID != "domain-1" || !client.lyantheEnabled {
		t.Fatalf("lyanthe = %q %v", client.lyantheDomainID, client.lyantheEnabled)
	}
}
