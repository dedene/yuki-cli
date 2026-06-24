package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/auth"
)

func TestExecuteVersionCommand(t *testing.T) {
	var out bytes.Buffer
	oldVersion := Version
	Version = "1.2.3"
	t.Cleanup(func() { Version = oldVersion })

	err := Execute(context.Background(), []string{"version"}, Runtime{Out: &out})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if got := strings.TrimSpace(out.String()); got != "yuki 1.2.3" {
		t.Fatalf("version output = %q", got)
	}
}

func TestAuthStatusJSONUsesEnvironmentAccessKey(t *testing.T) {
	t.Setenv("YUKI_ACCESS_KEY", "env-key")
	var out bytes.Buffer

	err := Execute(context.Background(), []string{"--json", "auth", "status"}, Runtime{
		Out:   &out,
		Store: &cmdFakeStore{err: auth.ErrAccessKeyNotFound},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if payload["authenticated"] != true || payload["source"] != string(auth.SourceEnv) {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestAuthStatusJSONUsesEnvironmentWithoutOpeningStore(t *testing.T) {
	t.Setenv("YUKI_ACCESS_KEY", "env-key")
	var out bytes.Buffer

	err := Execute(context.Background(), []string{"--json", "auth", "status"}, Runtime{Out: &out})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if payload["authenticated"] != true || payload["source"] != string(auth.SourceEnv) {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestAuthLogoutReportsMissingStoredAccessKey(t *testing.T) {
	var out bytes.Buffer

	err := Execute(context.Background(), []string{"auth", "logout"}, Runtime{
		Out:   &out,
		Store: &cmdFakeStore{err: auth.ErrAccessKeyNotFound},
	})
	if !errors.Is(err, auth.ErrAccessKeyNotFound) {
		t.Fatalf("err = %v, want ErrAccessKeyNotFound", err)
	}
	if out.Len() != 0 {
		t.Fatalf("logout wrote output on failure: %q", out.String())
	}
}

func TestDomainsListAuthenticatesAndPrintsTable(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		domains: []api.Domain{{
			ID:   "domain-1",
			Name: "Acme",
			URL:  "acme.yukiworks.be",
		}},
	}

	err := Execute(context.Background(), []string{"domains", "list"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.accessKey != "stored-key" {
		t.Fatalf("accessKey = %q", client.accessKey)
	}
	got := out.String()
	for _, want := range []string{"ID", "NAME", "URL", "domain-1", "Acme", "acme.yukiworks.be"} {
		if !strings.Contains(got, want) {
			t.Fatalf("domains output missing %q in:\n%s", want, got)
		}
	}
}

func TestGLAccountsListJSONUsesAdministrationFlag(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		accounts: []api.GLAccount{{
			Code:        "100000",
			Type:        "2",
			Subtype:     "0",
			Enabled:     true,
			Description: "Geplaatst kapitaal",
		}},
	}

	err := Execute(context.Background(), []string{"--json", "accounting", "gl-accounts", "list", "--administration", "admin-1"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.administrationID != "admin-1" {
		t.Fatalf("administrationID = %q", client.administrationID)
	}
	var accounts []api.GLAccount
	if err := json.Unmarshal(out.Bytes(), &accounts); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(accounts) != 1 || accounts[0].Code != "100000" {
		t.Fatalf("accounts = %#v", accounts)
	}
}

type cmdFakeStore struct {
	key string
	err error
}

func (s *cmdFakeStore) SetAccessKey(context.Context, string, string) error {
	return nil
}

func (s *cmdFakeStore) AccessKey(context.Context, string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	if s.key == "" {
		return "", auth.ErrAccessKeyNotFound
	}
	return s.key, nil
}

func (s *cmdFakeStore) DeleteAccessKey(context.Context, string) error {
	if s.err != nil {
		return s.err
	}
	return nil
}

type cmdFakeClient struct {
	sessionID        string
	accessKey        string
	administrationID string
	domains          []api.Domain
	accounts         []api.GLAccount
}

func (c *cmdFakeClient) Authenticate(_ context.Context, accessKey string) (string, error) {
	c.accessKey = accessKey
	return c.sessionID, nil
}

func (c *cmdFakeClient) Domains(context.Context, string) ([]api.Domain, error) {
	return c.domains, nil
}

func (c *cmdFakeClient) CurrentDomain(context.Context, string) (api.Domain, error) {
	return api.Domain{ID: "domain-1", Name: "Acme"}, nil
}

func (c *cmdFakeClient) Administrations(context.Context, string) ([]api.Administration, error) {
	return nil, nil
}

func (c *cmdFakeClient) Companies(context.Context, string) ([]api.Company, error) {
	return nil, nil
}

func (c *cmdFakeClient) GLAccounts(_ context.Context, _ string, administrationID string) ([]api.GLAccount, error) {
	c.administrationID = administrationID
	return c.accounts, nil
}
