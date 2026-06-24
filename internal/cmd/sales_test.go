package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/auth"
)

func TestSalesInvoiceSchemaPathJSONDoesNotAuthenticate(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID:              "session-1",
		salesInvoiceSchemaPath: "https://api.yukiworks.be/schemas/SalesInvoices.xsd",
	}

	err := Execute(context.Background(), []string{"--json", "sales", "invoices", "schema-path"}, Runtime{
		Out:       &out,
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
	var payload map[string]string
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if payload["schema_path"] != "https://api.yukiworks.be/schemas/SalesInvoices.xsd" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestSalesItemsListJSONUsesAdministration(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		salesItems: []api.SalesItem{{
			ID:          "item-1",
			Description: "Consulting",
		}},
	}

	err := Execute(context.Background(), []string{"--json", "sales", "items", "list", "--administration", "admin-1"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.administrationID != "admin-1" || client.accessKey != "stored-key" {
		t.Fatalf("administrationID/accessKey = %q/%q", client.administrationID, client.accessKey)
	}
	var items []api.SalesItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].Description != "Consulting" {
		t.Fatalf("items = %#v", items)
	}
}

func TestSalesInvoicesCreateDryRunSkipsAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeSalesInvoiceXMLFixture(t)

	err := Execute(context.Background(), []string{
		"--json", "sales", "invoices", "create",
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
		payload["operation"] != "ProcessSalesInvoices" ||
		payload["administration_id"] != "admin-1" ||
		payload["root"] != "SalesInvoices" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestSalesInvoicesCreateReadonlyBlocksBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeSalesInvoiceXMLFixture(t)
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--readonly", "sales", "invoices", "create",
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

func TestSalesInvoicesCreateSendsXMLFile(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeSalesInvoiceXMLFixture(t)
	client := &cmdFakeClient{
		sessionID: "session-1",
		salesInvoiceImportResult: api.SalesInvoiceImportResponse{
			TimeStamp:        "2025-10-23",
			AdministrationID: "admin-1",
			TotalSucceeded:   1,
			Invoices: []api.SalesInvoiceImportInvoice{{
				Succeeded: true,
				Processed: true,
				Reference: "VF-0001",
			}},
		},
	}

	err := Execute(context.Background(), []string{
		"--json", "sales", "invoices", "create",
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
	if client.salesInvoiceImportOperation != "ProcessSalesInvoices" ||
		client.salesInvoiceImportOpts.AdministrationID != "admin-1" ||
		!strings.Contains(client.salesInvoiceImportOpts.XMLDoc, "<Reference>VF-0001</Reference>") {
		t.Fatalf("opts/operation = %#v/%s", client.salesInvoiceImportOpts, client.salesInvoiceImportOperation)
	}
	if strings.Contains(client.salesInvoiceImportOpts.XMLDoc, "<?xml") {
		t.Fatalf("XML declaration was not stripped from nested xmlDoc: %q", client.salesInvoiceImportOpts.XMLDoc)
	}
	var result api.SalesInvoiceImportResponse
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if result.TotalSucceeded != 1 || len(result.Invoices) != 1 || result.Invoices[0].Reference != "VF-0001" {
		t.Fatalf("result = %#v", result)
	}
}

func TestSalesRecognizedInvoicesCreateUsesRecognizedOperation(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeSalesInvoiceXMLFixture(t)
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"sales", "recognized-invoices", "create",
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
	if client.salesInvoiceImportOperation != "ProcessRecognizedSalesInvoices" ||
		client.salesInvoiceImportOpts.AdministrationID != "admin-1" {
		t.Fatalf("opts/operation = %#v/%s", client.salesInvoiceImportOpts, client.salesInvoiceImportOperation)
	}
}

func TestSalesInvoicesCreateRejectsInvalidXMLBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	xmlPath := writeFile(t, "broken.xml", []byte(`<SalesInvoices>`))
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"sales", "invoices", "create",
		"--administration", "admin-1",
		"--file", xmlPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "validate") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}

func writeSalesInvoiceXMLFixture(t *testing.T) string {
	t.Helper()
	return writeFile(t, "sales-invoices.xml", []byte(`<?xml version="1.0" encoding="utf-8"?><SalesInvoices xmlns="urn:xmlns:http://www.theyukicompany.com:salesinvoices"><SalesInvoice><Reference>VF-0001</Reference></SalesInvoice></SalesInvoices>`))
}

func writeFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := t.TempDir() + string(os.PathSeparator) + name
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	return path
}
