package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/auth"
)

func TestPettyCashStatementImportDryRunSkipsAuth(t *testing.T) {
	var out bytes.Buffer
	path := writePettyCashStatementFixture(t)

	err := Execute(context.Background(), []string{
		"--json", "petty-cash", "statement", "import",
		"--administration", "admin-1",
		"--file", path,
		"--dry-run",
	}, Runtime{
		Out:   &out,
		Store: &cmdFakeStore{err: auth.ErrAccessKeyNotFound},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var result api.PettyCashImportResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if !result.DryRun ||
		result.Operation != "ImportStatement" ||
		result.AdministrationID != "admin-1" {
		t.Fatalf("result = %#v", result)
	}
}

func TestPettyCashStatementImportReadonlyBlocksBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	path := writePettyCashStatementFixture(t)
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--readonly", "petty-cash", "statement", "import",
		"--administration", "admin-1",
		"--file", path,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "--readonly blocks mutating command") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}

func TestPettyCashStatementImportSendsFile(t *testing.T) {
	var out bytes.Buffer
	path := writePettyCashStatementFixture(t)
	client := &cmdFakeClient{
		sessionID: "session-1",
		pettyCashImportResult: api.PettyCashImportResult{
			Operation:        "ImportStatement",
			AdministrationID: "admin-1",
			DocumentID:       "doc-1",
		},
	}

	err := Execute(context.Background(), []string{
		"--json", "petty-cash", "statement", "import",
		"--administration", "admin-1",
		"--file", path,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.pettyCashImportOperation != "ImportStatement" ||
		client.pettyCashStatementOpts.AdministrationID != "admin-1" ||
		!strings.Contains(client.pettyCashStatementOpts.StatementText, "Ledger account Petty Cash") {
		t.Fatalf("opts/operation = %#v/%s", client.pettyCashStatementOpts, client.pettyCashImportOperation)
	}
	var result api.PettyCashImportResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if result.DocumentID != "doc-1" {
		t.Fatalf("result = %#v", result)
	}
}

func TestPettyCashLineImportSendsFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--json", "petty-cash", "line", "import",
		"--account-gl-code", "570000",
		"--transaction-code", "REV6",
		"--offset-account", "700000",
		"--offset-name", "Revenue",
		"--transaction-date", "2021-01-15",
		"--description", "revenue 6%",
		"--amount", "100",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.pettyCashImportOperation != "ImportSingleStatementLine" ||
		client.pettyCashLineOpts.AccountGLCode != "570000" ||
		client.pettyCashLineOpts.TransactionDescription != "revenue 6%" {
		t.Fatalf("opts/operation = %#v/%s", client.pettyCashLineOpts, client.pettyCashImportOperation)
	}
}

func TestPettyCashProjectLineImportSendsProjectFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--json", "petty-cash", "project-line", "import",
		"--account-gl-code", "570000",
		"--transaction-code", "REV6",
		"--offset-account", "700000",
		"--offset-name", "Revenue",
		"--transaction-date", "2021-01-15",
		"--description", "revenue 6%",
		"--amount", "100",
		"--project-code", "PROJECTNEW",
		"--project-name", "New Project",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.pettyCashImportOperation != "ImportSingleStatementProjectLine" ||
		client.pettyCashLineOpts.ProjectCode != "PROJECTNEW" ||
		client.pettyCashLineOpts.ProjectName != "New Project" {
		t.Fatalf("opts/operation = %#v/%s", client.pettyCashLineOpts, client.pettyCashImportOperation)
	}
}

func writePettyCashStatementFixture(t *testing.T) string {
	t.Helper()
	return writeFile(t, "pettycash.csv", []byte("Ledger account Petty Cash;Petty cash description;Transaction code;Suspense account;Name suspense account;Date transaction;Description;Amount;Balance petty cash;Project code;Project name\n570001;Kas Highpro;REV6;700000;revenue;15/01/2021;revenue 6%;100;;PROJECTNEW;New Project\n"))
}
