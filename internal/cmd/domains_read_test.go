package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestDomainsNameJSONPrintsDomainName(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		domainNameResult: api.DomainNameResult{
			AdministrationName: "Highpro BV",
			DomainName:         "highpro.yukiworks.be",
		},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"domains", "name",
		"--administration-name", "Highpro BV",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.administrationDomainName != "Highpro BV" {
		t.Fatalf("administrationDomainName = %q", client.administrationDomainName)
	}
	var result api.DomainNameResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json: %v", err)
	}
	if result.DomainName != "highpro.yukiworks.be" {
		t.Fatalf("result = %#v", result)
	}
}

func TestDomainsUsersJSONPrintsUsers(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		domainUsers: []api.DomainUser{{
			ID:       "user-1",
			FullName: "Peter Dedene",
			Login:    "peter@example.com",
			Email:    "peter@example.com",
			Roles:    "Backoffice",
			Active:   "true",
		}},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"domains", "users",
		"--domain", "domain-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.domainID != "domain-1" {
		t.Fatalf("domainID = %q", client.domainID)
	}
	var users []api.DomainUser
	if err := json.Unmarshal(out.Bytes(), &users); err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(users) != 1 || users[0].Login != "peter@example.com" {
		t.Fatalf("users = %#v", users)
	}
}
