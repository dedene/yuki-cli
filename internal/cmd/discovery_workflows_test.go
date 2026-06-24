package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestTransactionsListJSONUsesPagingFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		richTransactions: []api.Transaction{{
			ID:              "tx-1",
			TransactionDate: "2026-01-03",
			GLAccountCode:   "550002",
			Amount:          "-242.00",
			Document: &api.TransactionDocumentReference{
				ID:        "doc-1",
				Reference: "MC-2026-01-03",
			},
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "transactions", "list",
		"--administration", "admin-1",
		"--gl-account", "550002",
		"--from", "2026-01-01",
		"--to", "2026-01-31",
		"--financial-mode", "0",
		"--data-groups", "documentprocessed,document,documentmatching",
		"--limit", "50",
		"--start-record", "51",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.transactionsOpts.AdministrationID != "admin-1" ||
		client.transactionsOpts.GLAccountCode != "550002" ||
		client.transactionsOpts.StartDate != "2026-01-01" ||
		client.transactionsOpts.DataGroups != "documentprocessed,document,documentmatching" ||
		client.transactionsOpts.NumberOfRecords != 50 ||
		client.transactionsOpts.StartRecord != 51 {
		t.Fatalf("transactionsOpts = %#v", client.transactionsOpts)
	}
	var transactions []api.Transaction
	if err := json.Unmarshal(out.Bytes(), &transactions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(transactions) != 1 || transactions[0].Document == nil || transactions[0].Document.ID != "doc-1" {
		t.Fatalf("transactions = %#v", transactions)
	}
}

func TestArchiveDocumentsSearchUsesDefaults(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		searchDocuments: []api.Document{{
			ID:          "doc-1",
			Subject:     "Apple invoice",
			Amount:      "29.76",
			ContactName: "Apple Sales International",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "search",
		"--search-text", "apple",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	opts := client.searchDocumentsOpts
	if opts.SearchOption != "All" ||
		opts.SearchText != "apple" ||
		opts.FolderID != -1 ||
		opts.TabID != -1 ||
		opts.SortOrder != "CreatedDesc" ||
		opts.StartDate != "0001-01-01" ||
		opts.EndDate != "0001-01-01" ||
		opts.NumberOfRecords != 25 ||
		opts.StartRecord != 1 {
		t.Fatalf("searchDocumentsOpts = %#v", opts)
	}
	var documents []api.Document
	if err := json.Unmarshal(out.Bytes(), &documents); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(documents) != 1 || documents[0].ID != "doc-1" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestArchiveFoldersListPrintsRows(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		folders: []api.DocumentFolder{{
			ID:              "1",
			Description:     "Purchase",
			Icon:            "DocumentFolder_red_label.png",
			ProcessedByYuki: true,
		}},
	}

	err := Execute(context.Background(), []string{"archive", "folders", "list"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	got := out.String()
	for _, want := range []string{"ID", "DESCRIPTION", "PROCESSED", "1", "Purchase", "true"} {
		if !strings.Contains(got, want) {
			t.Fatalf("archive folders output missing %q in:\n%s", want, got)
		}
	}
}

func TestArchiveFoldersTabsPrintsRows(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		tabs: []api.DocumentFolderTab{{
			ID:              "303",
			Description:     "Credit cards",
			ProcessedByYuki: true,
		}},
	}

	err := Execute(context.Background(), []string{"archive", "folders", "tabs", "--folder", "3"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.folderID != "3" {
		t.Fatalf("folderID = %q", client.folderID)
	}
	got := out.String()
	for _, want := range []string{"ID", "DESCRIPTION", "PROCESSED", "303", "Credit cards", "true"} {
		if !strings.Contains(got, want) {
			t.Fatalf("archive folder tabs output missing %q in:\n%s", want, got)
		}
	}
}

func TestAccountingPaymentMethodsListPrintsRows(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		customMethods: []api.PaymentMethod{{
			ID:          "5",
			Description: "Creditcard",
		}},
	}

	err := Execute(context.Background(), []string{
		"accounting", "payment-methods", "list",
		"--administration", "admin-1",
	}, Runtime{
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
	got := out.String()
	for _, want := range []string{"ID", "DESCRIPTION", "5", "Creditcard"} {
		if !strings.Contains(got, want) {
			t.Fatalf("payment methods output missing %q in:\n%s", want, got)
		}
	}
}

func TestArchivePaymentMethodsListPrintsRows(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		archiveMethods: []api.PaymentMethod{{
			ID:          "4",
			Description: "Zakelijke Bancontact",
		}},
	}

	err := Execute(context.Background(), []string{"archive", "payment-methods", "list"}, Runtime{
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
	for _, want := range []string{"ID", "DESCRIPTION", "4", "Zakelijke Bancontact"} {
		if !strings.Contains(got, want) {
			t.Fatalf("archive payment methods output missing %q in:\n%s", want, got)
		}
	}
}

func TestArchiveCurrenciesListPrintsRows(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		currencies: []api.Currency{{
			ID:          "EUR",
			Default:     true,
			Description: "Euro (EUR)",
		}},
	}

	err := Execute(context.Background(), []string{"archive", "currencies", "list"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	got := out.String()
	for _, want := range []string{"ID", "DEFAULT", "DESCRIPTION", "EUR", "true", "Euro (EUR)"} {
		if !strings.Contains(got, want) {
			t.Fatalf("archive currencies output missing %q in:\n%s", want, got)
		}
	}
}
