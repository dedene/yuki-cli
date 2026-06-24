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

func TestAccountingJournalsProcessDryRunSkipsAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeJournalXMLFixture(t)

	err := Execute(context.Background(), []string{
		"--json", "accounting", "journals", "process",
		"--administration", "admin-1",
		"--file", xmlPath,
		"--dry-run",
	}, Runtime{
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
	if payload["dry_run"] != true ||
		payload["operation"] != "ProcessJournal" ||
		payload["administration_id"] != "admin-1" ||
		payload["root"] != "Journal" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestAccountingJournalsProcessReadonlyBlocksBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeJournalXMLFixture(t)
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--readonly", "accounting", "journals", "process",
		"--administration", "admin-1",
		"--file", xmlPath,
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

func TestAccountingJournalsProcessSendsXMLFile(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeJournalXMLFixture(t)
	client := &cmdFakeClient{
		sessionID: "session-1",
		journalProcessResult: api.JournalProcessResult{
			AdministrationID: "admin-1",
			DocumentID:       "9999e179-04bd-43a7-8fb7-403999a88fac",
		},
	}

	err := Execute(context.Background(), []string{
		"--json", "accounting", "journals", "process",
		"--administration", "admin-1",
		"--file", xmlPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.journalImportOpts.AdministrationID != "admin-1" ||
		!strings.Contains(client.journalImportOpts.XMLDoc, "<DocumentSubject>Payroll</DocumentSubject>") {
		t.Fatalf("opts = %#v", client.journalImportOpts)
	}
	if strings.Contains(client.journalImportOpts.XMLDoc, "<?xml") {
		t.Fatalf("XML declaration was not stripped from nested xmlDoc: %q", client.journalImportOpts.XMLDoc)
	}
	var result api.JournalProcessResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if result.DocumentID != "9999e179-04bd-43a7-8fb7-403999a88fac" {
		t.Fatalf("result = %#v", result)
	}
}

func TestAccountingJournalsProcessRejectsWrongRootBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeFile(t, "sales-invoices.xml", []byte(`<SalesInvoices/>`))
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"accounting", "journals", "process",
		"--administration", "admin-1",
		"--file", xmlPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "expected root element Journal") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}

func writeJournalXMLFixture(t *testing.T) string {
	t.Helper()
	return writeFile(t, "journal.xml", []byte(`<?xml version="1.0" encoding="utf-8"?><Journal xmlns="urn:xmlns:http://www.theyukicompany.com:journal"><AdministrationID>admin-1</AdministrationID><DocumentSubject>Payroll</DocumentSubject><JournalType>GeneralJournal</JournalType><JournalEntry><EntryDate>2026-01-31</EntryDate><GLAccount>400000</GLAccount><Amount>100.00</Amount><Description>Payroll</Description></JournalEntry><JournalEntry><EntryDate>2026-01-31</EntryDate><GLAccount>456000</GLAccount><Amount>-100.00</Amount><Description>Payroll</Description></JournalEntry></Journal>`))
}
