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

func TestCreditorItemsListFiltersPaymentMethod(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		creditorItems: []api.CreditorItem{{
			Date:           "2026-01-03",
			Contact:        "AD Delhaize",
			OriginalAmount: "242.00",
			OpenAmount:     "242.00",
			PaymentMethod:  "Creditcard",
			Reference:      "test",
			DocumentID:     "doc-1",
			Description:    "Factuur van AD Delhaize",
		}, {
			Date:           "2026-01-04",
			Contact:        "Other Vendor",
			OriginalAmount: "10.00",
			OpenAmount:     "10.00",
			PaymentMethod:  "Transfer",
			DocumentID:     "doc-2",
			Description:    "Bank transfer",
		}},
	}

	err := Execute(context.Background(), []string{
		"accounting", "creditor-items", "list",
		"--administration", "admin-1",
		"--from", "2026-01-01",
		"--to", "2026-01-31",
		"--payment-method", "Creditcard",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.creditorOpts.AdministrationID != "admin-1" ||
		client.creditorOpts.StartDate != "2026-01-01" ||
		client.creditorOpts.EndDate != "2026-01-31" {
		t.Fatalf("creditorOpts = %#v", client.creditorOpts)
	}
	got := out.String()
	for _, want := range []string{"DATE", "AD Delhaize", "Creditcard", "doc-1"} {
		if !strings.Contains(got, want) {
			t.Fatalf("creditor-items output missing %q in:\n%s", want, got)
		}
	}
	if strings.Contains(got, "Other Vendor") {
		t.Fatalf("creditor-items output was not filtered:\n%s", got)
	}
}

func TestTransactionDetailsJSONUsesFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		transactions: []api.TransactionInfo{{
			ID:                "tx-1",
			TransactionDate:   "2026-01-03",
			GLAccountCode:     "400000",
			TransactionAmount: "-242.00",
			DocumentID:        "doc-1",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "transactions", "details",
		"--administration", "admin-1",
		"--gl-account", "400000",
		"--from", "2026-01-01",
		"--to", "2026-01-31",
		"--financial-mode", "1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.transactionOpts.AdministrationID != "admin-1" ||
		client.transactionOpts.GLAccountCode != "400000" ||
		client.transactionOpts.StartDate != "2026-01-01" ||
		client.transactionOpts.FinancialMode != "1" {
		t.Fatalf("transactionOpts = %#v", client.transactionOpts)
	}
	var transactions []api.TransactionInfo
	if err := json.Unmarshal(out.Bytes(), &transactions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(transactions) != 1 || transactions[0].ID != "tx-1" {
		t.Fatalf("transactions = %#v", transactions)
	}
}

func TestArchiveDocumentsFindPrintsDocument(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		document: api.Document{
			ID:              "doc-1",
			Subject:         "Invoice",
			DocumentDate:    "2026-01-03",
			Amount:          "242.00",
			TypeDescription: "Aankoopfactuur",
			FileName:        "invoice.pdf",
			ContactName:     "AD Delhaize",
		},
	}

	err := Execute(context.Background(), []string{"archive", "documents", "find", "--document", "doc-1"}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.documentID != "doc-1" {
		t.Fatalf("documentID = %q", client.documentID)
	}
	got := out.String()
	for _, want := range []string{"ID", "doc-1", "Invoice", "invoice.pdf"} {
		if !strings.Contains(got, want) {
			t.Fatalf("archive document output missing %q in:\n%s", want, got)
		}
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
	sessionID           string
	accessKey           string
	administrationID    string
	documentID          string
	transactionID       string
	domains             []api.Domain
	accounts            []api.GLAccount
	creditorItems       []api.CreditorItem
	creditorOpts        api.CreditorItemsOptions
	richTransactions    []api.Transaction
	transactionsOpts    api.TransactionsOptions
	transactions        []api.TransactionInfo
	transactionOpts     api.TransactionDetailsOptions
	customMethods       []api.PaymentMethod
	archiveMethods      []api.PaymentMethod
	folders             []api.DocumentFolder
	tabs                []api.DocumentFolderTab
	currencies          []api.Currency
	costCategories      []api.CostCategory
	folderID            string
	searchDocuments     []api.Document
	searchDocumentsOpts api.SearchDocumentsOptions
	document            api.Document
	documentFile        api.DocumentFile
	transactionDocument api.TransactionDocument
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

func (c *cmdFakeClient) OutstandingCreditorItemsByDate(_ context.Context, _ string, opts api.CreditorItemsOptions) ([]api.CreditorItem, error) {
	c.creditorOpts = opts
	return c.creditorItems, nil
}

func (c *cmdFakeClient) Transactions(_ context.Context, _ string, opts api.TransactionsOptions) ([]api.Transaction, error) {
	c.transactionsOpts = opts
	return c.richTransactions, nil
}

func (c *cmdFakeClient) TransactionDetails(_ context.Context, _ string, opts api.TransactionDetailsOptions) ([]api.TransactionInfo, error) {
	c.transactionOpts = opts
	return c.transactions, nil
}

func (c *cmdFakeClient) TransactionDocument(_ context.Context, _ string, administrationID string, transactionID string) (api.TransactionDocument, error) {
	c.administrationID = administrationID
	c.transactionID = transactionID
	return c.transactionDocument, nil
}

func (c *cmdFakeClient) CustomPaymentMethods(_ context.Context, _ string, administrationID string) ([]api.PaymentMethod, error) {
	c.administrationID = administrationID
	return c.customMethods, nil
}

func (c *cmdFakeClient) DocumentFolders(context.Context, string) ([]api.DocumentFolder, error) {
	return c.folders, nil
}

func (c *cmdFakeClient) DocumentFolderTabs(_ context.Context, _ string, folderID string) ([]api.DocumentFolderTab, error) {
	c.folderID = folderID
	return c.tabs, nil
}

func (c *cmdFakeClient) SearchDocuments(_ context.Context, _ string, opts api.SearchDocumentsOptions) ([]api.Document, error) {
	c.searchDocumentsOpts = opts
	return c.searchDocuments, nil
}

func (c *cmdFakeClient) FindDocument(_ context.Context, _ string, documentID string) (api.Document, error) {
	c.documentID = documentID
	return c.document, nil
}

func (c *cmdFakeClient) DocumentFile(_ context.Context, _ string, documentID string) (api.DocumentFile, error) {
	c.documentID = documentID
	return c.documentFile, nil
}

func (c *cmdFakeClient) PaymentMethods(context.Context, string) ([]api.PaymentMethod, error) {
	return c.archiveMethods, nil
}

func (c *cmdFakeClient) Currencies(context.Context, string) ([]api.Currency, error) {
	return c.currencies, nil
}

func (c *cmdFakeClient) CostCategories(context.Context, string) ([]api.CostCategory, error) {
	return c.costCategories, nil
}
