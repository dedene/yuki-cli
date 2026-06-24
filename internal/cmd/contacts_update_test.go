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

func TestContactsUpsertDryRunSkipsAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeContactXMLFixture(t)
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--json", "contacts", "upsert",
		"--domain", "domain-1",
		"--file", xmlPath,
		"--dry-run",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{err: auth.ErrAccessKeyNotFound},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.accessKey != "" || client.contactUpdateOpts.XMLDoc != "" {
		t.Fatalf("dry-run authenticated or sent: accessKey=%q opts=%#v", client.accessKey, client.contactUpdateOpts)
	}
	var payload map[string]any
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if payload["dry_run"] != true ||
		payload["operation"] != "UpdateContact" ||
		payload["domain_id"] != "domain-1" ||
		payload["root"] != "Contact" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestContactsUpsertReadonlyBlocksBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeContactXMLFixture(t)
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--readonly", "contacts", "upsert",
		"--domain", "domain-1",
		"--file", xmlPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "--readonly blocks mutating command: contacts upsert") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}

func TestContactsUpsertSendsXMLFile(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeContactXMLFixture(t)
	client := &cmdFakeClient{
		sessionID: "session-1",
		contactUpdateResult: api.ContactUpdateResult{
			DomainID:  "domain-1",
			Timestamp: "2021-03-08",
			Succeeded: "Succesfully updated Contact contact-1",
			Message:   "Succesfully updated Contact contact-1",
		},
	}

	err := Execute(context.Background(), []string{
		"--json", "contacts", "upsert",
		"--domain", "domain-1",
		"--file", xmlPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.contactUpdateOpts.DomainID != "domain-1" ||
		!strings.Contains(client.contactUpdateOpts.XMLDoc, "<FullName>A van B</FullName>") {
		t.Fatalf("opts = %#v", client.contactUpdateOpts)
	}
	if strings.Contains(client.contactUpdateOpts.XMLDoc, "<?xml") {
		t.Fatalf("XML declaration was not stripped from nested xmlDoc: %q", client.contactUpdateOpts.XMLDoc)
	}
	var result api.ContactUpdateResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if result.Succeeded != "Succesfully updated Contact contact-1" {
		t.Fatalf("result = %#v", result)
	}
}

func TestContactsUpsertRejectsWrongRootBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeFile(t, "journal.xml", []byte(`<Journal/>`))
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"contacts", "upsert",
		"--domain", "domain-1",
		"--file", xmlPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "expected root element Contact") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}

func writeContactXMLFixture(t *testing.T) string {
	t.Helper()
	return writeFile(t, "contact.xml", []byte(`<?xml version="1.0" encoding="utf-8"?><Contact xmlns="urn:xmlns:http://www.theyukicompany.com:contact"><ID/><Type>0</Type><Code>1</Code><FirstName>A</FirstName><MiddleName>van</MiddleName><LastName>B</LastName><FullName>A van B</FullName><EmailHome>support@yuki.nl</EmailHome></Contact>`))
}
