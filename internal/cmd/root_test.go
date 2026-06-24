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

func TestGLAccountsBalanceJSONUsesDate(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		balances: []api.GLAccountBalanceItem{{
			Code:        "100000",
			BalanceType: "B",
			Amount:      "-1222.22",
			Description: "Share capital",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "balance",
		"--administration", "admin-1",
		"--date", "2020-12-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.balanceOpts.AdministrationID != "admin-1" ||
		client.balanceOpts.TransactionDate != "2020-12-31" {
		t.Fatalf("balanceOpts = %#v", client.balanceOpts)
	}
	var balances []api.GLAccountBalanceItem
	if err := json.Unmarshal(out.Bytes(), &balances); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(balances) != 1 || balances[0].Code != "100000" {
		t.Fatalf("balances = %#v", balances)
	}
}

func TestGLAccountsBalanceFiscalJSONUsesDate(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		balances: []api.GLAccountBalanceItem{{
			Code:        "100000",
			BalanceType: "B",
			Amount:      "-1222.22",
			Description: "Geplaatst kapitaal",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "balance-fiscal",
		"--administration", "admin-1",
		"--date", "2020-12-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.balanceOpts.AdministrationID != "admin-1" ||
		client.balanceOpts.TransactionDate != "2020-12-31" {
		t.Fatalf("balanceOpts = %#v", client.balanceOpts)
	}
	var balances []api.GLAccountBalanceItem
	if err := json.Unmarshal(out.Bytes(), &balances); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(balances) != 1 || balances[0].Description != "Geplaatst kapitaal" {
		t.Fatalf("balances = %#v", balances)
	}
}

func TestGLAccountsBalanceYearEndJSONUsesDate(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		balances: []api.GLAccountBalanceItem{{
			Code:        "140000",
			BalanceType: "B",
			Amount:      "-1454.14",
			Description: "Overgedragen winst",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "balance-year-end",
		"--administration", "admin-1",
		"--date", "2020-12-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.balanceOpts.AdministrationID != "admin-1" ||
		client.balanceOpts.TransactionDate != "2020-12-31" {
		t.Fatalf("balanceOpts = %#v", client.balanceOpts)
	}
	var balances []api.GLAccountBalanceItem
	if err := json.Unmarshal(out.Bytes(), &balances); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(balances) != 1 || balances[0].Code != "140000" {
		t.Fatalf("balances = %#v", balances)
	}
}

func TestGLAccountsTransactionsJSONUsesDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		glTransactions: []api.GLAccountTransaction{{
			ID:            "tx-1",
			Date:          "2020-01-01",
			Description:   "Factuur voor Quentin test",
			Amount:        "-47.45",
			Contact:       "Quentin test",
			ContactID:     "contact-1",
			Project:       api.GLTransactionProject{Code: "WELLNESS", Text: "Wellness Event"},
			GLAccountCode: "700000",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "transactions",
		"--administration", "admin-1",
		"--gl-account", "700000",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.glTransactionOpts.AdministrationID != "admin-1" ||
		client.glTransactionOpts.GLAccountCode != "700000" ||
		client.glTransactionOpts.StartDate != "2020-01-01" ||
		client.glTransactionOpts.EndDate != "2020-01-31" {
		t.Fatalf("glTransactionOpts = %#v", client.glTransactionOpts)
	}
	var transactions []api.GLAccountTransaction
	if err := json.Unmarshal(out.Bytes(), &transactions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(transactions) != 1 || transactions[0].Project.Code != "WELLNESS" {
		t.Fatalf("transactions = %#v", transactions)
	}
}

func TestGLAccountsTransactionsFiscalJSONUsesDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		glTransactions: []api.GLAccountTransaction{{
			ID:            "tx-1",
			Date:          "2020-01-01",
			Description:   "Factuur voor Quentin test",
			Amount:        "-47.45",
			Contact:       "Quentin test",
			Project:       api.GLTransactionProject{Code: "WELLNESS", Text: "Wellness Event"},
			GLAccountCode: "700000",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "transactions-fiscal",
		"--administration", "admin-1",
		"--gl-account", "700000",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.glTransactionOpts.AdministrationID != "admin-1" ||
		client.glTransactionOpts.GLAccountCode != "700000" ||
		client.glTransactionOpts.StartDate != "2020-01-01" ||
		client.glTransactionOpts.EndDate != "2020-01-31" {
		t.Fatalf("glTransactionOpts = %#v", client.glTransactionOpts)
	}
	var transactions []api.GLAccountTransaction
	if err := json.Unmarshal(out.Bytes(), &transactions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(transactions) != 1 || transactions[0].Project.Code != "WELLNESS" {
		t.Fatalf("transactions = %#v", transactions)
	}
}

func TestGLAccountsTransactionsWithContactJSONUsesDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		glTransactions: []api.GLAccountTransaction{{
			ID:            "tx-1",
			Date:          "2020-01-01",
			Description:   "Factuur voor Quentin test",
			Amount:        "-47.45",
			Contact:       "Quentin test",
			ContactID:     "contact-1",
			Project:       api.GLTransactionProject{Text: "WELLNESS"},
			GLAccountCode: "700000",
			FileName:      "Invoice 2017/149.pdf",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "transactions-with-contact",
		"--administration", "admin-1",
		"--gl-account", "700000",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.glTransactionOpts.AdministrationID != "admin-1" ||
		client.glTransactionOpts.GLAccountCode != "700000" ||
		client.glTransactionOpts.StartDate != "2020-01-01" ||
		client.glTransactionOpts.EndDate != "2020-01-31" {
		t.Fatalf("glTransactionOpts = %#v", client.glTransactionOpts)
	}
	var transactions []api.GLAccountTransaction
	if err := json.Unmarshal(out.Bytes(), &transactions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(transactions) != 1 || transactions[0].FileName != "Invoice 2017/149.pdf" {
		t.Fatalf("transactions = %#v", transactions)
	}
}

func TestRevenueNetJSONUsesDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		revenueReport: api.RevenueReport{
			AdministrationID: "admin-1",
			StartDate:        "2020-01-01",
			EndDate:          "2020-01-31",
			Amount:           "1868.36",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "revenue", "net",
		"--administration", "admin-1",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.revenueOpts.AdministrationID != "admin-1" ||
		client.revenueOpts.StartDate != "2020-01-01" ||
		client.revenueOpts.EndDate != "2020-01-31" {
		t.Fatalf("revenueOpts = %#v", client.revenueOpts)
	}
	var report api.RevenueReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if report.Amount != "1868.36" {
		t.Fatalf("report = %#v", report)
	}
}

func TestRevenueNetFiscalJSONUsesDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		revenueReport: api.RevenueReport{
			AdministrationID: "admin-1",
			StartDate:        "2020-01-01",
			EndDate:          "2020-01-31",
			Amount:           "1868.36",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "revenue", "net-fiscal",
		"--administration", "admin-1",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.revenueOpts.AdministrationID != "admin-1" ||
		client.revenueOpts.StartDate != "2020-01-01" ||
		client.revenueOpts.EndDate != "2020-01-31" {
		t.Fatalf("revenueOpts = %#v", client.revenueOpts)
	}
	var report api.RevenueReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if report.Amount != "1868.36" {
		t.Fatalf("report = %#v", report)
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

func TestCreditorItemsAllJSONUsesFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		creditorItems: []api.CreditorItem{{
			Date:           "2026-01-03",
			Contact:        "Belgian Shell S.A.",
			OriginalAmount: "75.86",
			OpenAmount:     "75.86",
			PaymentMethod:  "Overschrijving",
			Reference:      "S05233212",
			DocumentID:     "doc-1",
			Description:    "Factuur van Belgian Shell S.A.",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "creditor-items", "all",
		"--administration", "admin-1",
		"--include-bank-transactions",
		"--sort-order", "DateDesc",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.creditorOpts.AdministrationID != "admin-1" ||
		!client.creditorOpts.IncludeBankTransactions ||
		client.creditorOpts.SortOrder != "DateDesc" ||
		client.creditorOpts.StartDate != "" ||
		client.creditorOpts.EndDate != "" {
		t.Fatalf("creditorOpts = %#v", client.creditorOpts)
	}
	var items []api.CreditorItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].Reference != "S05233212" {
		t.Fatalf("items = %#v", items)
	}
}

func TestCreditorItemsWithPaymentReferenceJSONUsesFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		creditorItems: []api.CreditorItem{{
			Date:             "2026-01-03",
			Contact:          "AD Delhaize",
			OriginalAmount:   "242.00",
			OpenAmount:       "242.00",
			PaymentMethod:    "Creditcard",
			Reference:        "test",
			PaymentReference: "RF18539007547034",
			DocumentID:       "doc-1",
			Description:      "Factuur van AD Delhaize",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "creditor-items", "with-payment-reference",
		"--administration", "admin-1",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
		"--include-bank-transactions",
		"--sort-order", "DateDesc",
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
		client.creditorOpts.StartDate != "2020-01-01" ||
		client.creditorOpts.EndDate != "2020-01-31" ||
		!client.creditorOpts.IncludeBankTransactions ||
		client.creditorOpts.SortOrder != "DateDesc" {
		t.Fatalf("creditorOpts = %#v", client.creditorOpts)
	}
	var items []api.CreditorItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].PaymentReference != "RF18539007547034" {
		t.Fatalf("items = %#v", items)
	}
}

func TestDebtorItemsAllJSONUsesFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		debtorItems: []api.DebtorItem{{
			Date:           "2026-01-03",
			Contact:        "Customer",
			OriginalAmount: "1000.00",
			OpenAmount:     "1000.00",
			PaymentMethod:  "Creditcard",
			Reference:      "INV-1",
			DocumentID:     "doc-1",
			Description:    "Sales invoice",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "debtor-items", "all",
		"--administration", "admin-1",
		"--include-bank-transactions",
		"--sort-order", "DateDesc",
		"--payment-method", "Creditcard",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.debtorOpts.AdministrationID != "admin-1" ||
		!client.debtorOpts.IncludeBankTransactions ||
		client.debtorOpts.SortOrder != "DateDesc" ||
		client.debtorOpts.StartDate != "" ||
		client.debtorOpts.EndDate != "" {
		t.Fatalf("debtorOpts = %#v", client.debtorOpts)
	}
	var items []api.DebtorItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].Reference != "INV-1" {
		t.Fatalf("items = %#v", items)
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

func TestArchiveDocumentsImageCountJSONPrintsCount(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentImageCount: api.DocumentImageCount{
			DocumentID: "doc-1",
			ImageCount: 3,
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "image-count",
		"--document", "doc-1",
	}, Runtime{
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
	var count api.DocumentImageCount
	if err := json.Unmarshal(out.Bytes(), &count); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if count.DocumentID != "doc-1" || count.ImageCount != 3 {
		t.Fatalf("count = %#v", count)
	}
}

func TestArchiveDocumentsBinaryJSONPrintsBase64(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentBinaryData: api.DocumentBinaryData{
			DocumentID: "doc-1",
			FileData:   "JVBERg==",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "binary",
		"--document", "doc-1",
	}, Runtime{
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
	var data api.DocumentBinaryData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if data.DocumentID != "doc-1" || data.FileData != "JVBERg==" {
		t.Fatalf("data = %#v", data)
	}
}

func TestArchiveDocumentsXMLJSONPrintsXMLData(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentXMLData: api.DocumentXMLData{
			DocumentID: "doc-1",
			XML:        "<SalesInvoice><Reference>A1040</Reference></SalesInvoice>",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "xml",
		"--document", "doc-1",
	}, Runtime{
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
	var data api.DocumentXMLData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if data.DocumentID != "doc-1" || data.XML != "<SalesInvoice><Reference>A1040</Reference></SalesInvoice>" {
		t.Fatalf("data = %#v", data)
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
	sessionID             string
	accessKey             string
	administrationID      string
	documentID            string
	transactionID         string
	domains               []api.Domain
	accounts              []api.GLAccount
	balances              []api.GLAccountBalanceItem
	balanceOpts           api.GLAccountBalanceOptions
	glTransactions        []api.GLAccountTransaction
	glTransactionOpts     api.GLAccountTransactionsOptions
	revenueReport         api.RevenueReport
	revenueOpts           api.RevenueOptions
	creditorItems         []api.CreditorItem
	creditorOpts          api.CreditorItemsOptions
	debtorItems           []api.DebtorItem
	debtorOpts            api.DebtorItemsOptions
	richTransactions      []api.Transaction
	transactionsOpts      api.TransactionsOptions
	transactions          []api.TransactionInfo
	transactionOpts       api.TransactionDetailsOptions
	customMethods         []api.PaymentMethod
	archiveMethods        []api.PaymentMethod
	folders               []api.DocumentFolder
	tabs                  []api.DocumentFolderTab
	currencies            []api.Currency
	costCategories        []api.CostCategory
	menuEntries           []api.MenuEntry
	folderID              string
	documents             []api.Document
	documentsOpts         api.DocumentsOptions
	documentsInFolderOpts api.DocumentsInFolderOptions
	documentsInTabOpts    api.DocumentsInTabOptions
	documentsByTypeOpts   api.DocumentsByTypeOptions
	modifiedInFolderOpts  api.ModifiedDocumentsInFolderOptions
	modifiedByTypeOpts    api.ModifiedDocumentsByTypeOptions
	searchDocuments       []api.Document
	searchDocumentsOpts   api.SearchDocumentsOptions
	document              api.Document
	documentBundle        []api.Document
	documentFile          api.DocumentFile
	documentBinaryData    api.DocumentBinaryData
	documentImageData     api.DocumentImageData
	documentImageCount    api.DocumentImageCount
	documentXMLBinaryData api.DocumentXMLBinaryData
	documentXMLData       api.DocumentXMLData
	transactionDocument   api.TransactionDocument
	maxWidth              int
	maxHeight             int
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

func (c *cmdFakeClient) GLAccountBalance(_ context.Context, _ string, opts api.GLAccountBalanceOptions) ([]api.GLAccountBalanceItem, error) {
	c.balanceOpts = opts
	return c.balances, nil
}

func (c *cmdFakeClient) GLAccountBalanceFiscal(_ context.Context, _ string, opts api.GLAccountBalanceOptions) ([]api.GLAccountBalanceItem, error) {
	c.balanceOpts = opts
	return c.balances, nil
}

func (c *cmdFakeClient) GLAccountBalanceYearEnd(_ context.Context, _ string, opts api.GLAccountBalanceOptions) ([]api.GLAccountBalanceItem, error) {
	c.balanceOpts = opts
	return c.balances, nil
}

func (c *cmdFakeClient) GLAccountTransactions(_ context.Context, _ string, opts api.GLAccountTransactionsOptions) ([]api.GLAccountTransaction, error) {
	c.glTransactionOpts = opts
	return c.glTransactions, nil
}

func (c *cmdFakeClient) GLAccountTransactionsFiscal(_ context.Context, _ string, opts api.GLAccountTransactionsOptions) ([]api.GLAccountTransaction, error) {
	c.glTransactionOpts = opts
	return c.glTransactions, nil
}

func (c *cmdFakeClient) GLAccountTransactionsAndContact(_ context.Context, _ string, opts api.GLAccountTransactionsOptions) ([]api.GLAccountTransaction, error) {
	c.glTransactionOpts = opts
	return c.glTransactions, nil
}

func (c *cmdFakeClient) NetRevenue(_ context.Context, _ string, opts api.RevenueOptions) (api.RevenueReport, error) {
	c.revenueOpts = opts
	return c.revenueReportWithDefaults(opts), nil
}

func (c *cmdFakeClient) NetRevenueFiscal(_ context.Context, _ string, opts api.RevenueOptions) (api.RevenueReport, error) {
	c.revenueOpts = opts
	return c.revenueReportWithDefaults(opts), nil
}

func (c *cmdFakeClient) revenueReportWithDefaults(opts api.RevenueOptions) api.RevenueReport {
	if c.revenueReport.AdministrationID == "" {
		c.revenueReport.AdministrationID = opts.AdministrationID
	}
	if c.revenueReport.StartDate == "" {
		c.revenueReport.StartDate = opts.StartDate
	}
	if c.revenueReport.EndDate == "" {
		c.revenueReport.EndDate = opts.EndDate
	}
	return c.revenueReport
}

func (c *cmdFakeClient) OutstandingCreditorItems(_ context.Context, _ string, opts api.CreditorItemsOptions) ([]api.CreditorItem, error) {
	c.creditorOpts = opts
	return c.creditorItems, nil
}

func (c *cmdFakeClient) OutstandingCreditorItemsByDate(_ context.Context, _ string, opts api.CreditorItemsOptions) ([]api.CreditorItem, error) {
	c.creditorOpts = opts
	return c.creditorItems, nil
}

func (c *cmdFakeClient) OutstandingCreditorWithPaymentReference(_ context.Context, _ string, opts api.CreditorItemsOptions) ([]api.CreditorItem, error) {
	c.creditorOpts = opts
	return c.creditorItems, nil
}

func (c *cmdFakeClient) OutstandingDebtorItems(_ context.Context, _ string, opts api.DebtorItemsOptions) ([]api.DebtorItem, error) {
	c.debtorOpts = opts
	return c.debtorItems, nil
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

func (c *cmdFakeClient) Documents(_ context.Context, _ string, opts api.DocumentsOptions) ([]api.Document, error) {
	c.documentsOpts = opts
	return c.documents, nil
}

func (c *cmdFakeClient) DocumentsInFolder(_ context.Context, _ string, opts api.DocumentsInFolderOptions) ([]api.Document, error) {
	c.documentsInFolderOpts = opts
	return c.documents, nil
}

func (c *cmdFakeClient) DocumentsInTab(_ context.Context, _ string, opts api.DocumentsInTabOptions) ([]api.Document, error) {
	c.documentsInTabOpts = opts
	return c.documents, nil
}

func (c *cmdFakeClient) DocumentsByType(_ context.Context, _ string, opts api.DocumentsByTypeOptions) ([]api.Document, error) {
	c.documentsByTypeOpts = opts
	return c.documents, nil
}

func (c *cmdFakeClient) ModifiedDocumentsInFolder(_ context.Context, _ string, opts api.ModifiedDocumentsInFolderOptions) ([]api.Document, error) {
	c.modifiedInFolderOpts = opts
	return c.documents, nil
}

func (c *cmdFakeClient) ModifiedDocumentsByType(_ context.Context, _ string, opts api.ModifiedDocumentsByTypeOptions) ([]api.Document, error) {
	c.modifiedByTypeOpts = opts
	return c.documents, nil
}

func (c *cmdFakeClient) SearchDocuments(_ context.Context, _ string, opts api.SearchDocumentsOptions) ([]api.Document, error) {
	c.searchDocumentsOpts = opts
	return c.searchDocuments, nil
}

func (c *cmdFakeClient) FindDocument(_ context.Context, _ string, documentID string) (api.Document, error) {
	c.documentID = documentID
	return c.document, nil
}

func (c *cmdFakeClient) DocumentBundle(_ context.Context, _ string, documentID string) ([]api.Document, error) {
	c.documentID = documentID
	return c.documentBundle, nil
}

func (c *cmdFakeClient) DocumentFile(_ context.Context, _ string, documentID string) (api.DocumentFile, error) {
	c.documentID = documentID
	return c.documentFile, nil
}

func (c *cmdFakeClient) DocumentBinaryData(_ context.Context, _ string, documentID string) (api.DocumentBinaryData, error) {
	c.documentID = documentID
	if c.documentBinaryData.DocumentID == "" {
		c.documentBinaryData.DocumentID = documentID
	}
	return c.documentBinaryData, nil
}

func (c *cmdFakeClient) DocumentImage(_ context.Context, _ string, documentID string, maxWidth int, maxHeight int) (api.DocumentImageData, error) {
	c.documentID = documentID
	c.maxWidth = maxWidth
	c.maxHeight = maxHeight
	if c.documentImageData.DocumentID == "" {
		c.documentImageData.DocumentID = documentID
	}
	return c.documentImageData, nil
}

func (c *cmdFakeClient) DocumentImageCount(_ context.Context, _ string, documentID string) (api.DocumentImageCount, error) {
	c.documentID = documentID
	if c.documentImageCount.DocumentID == "" {
		c.documentImageCount.DocumentID = documentID
	}
	return c.documentImageCount, nil
}

func (c *cmdFakeClient) DocumentXMLDataAsBinary(_ context.Context, _ string, documentID string) (api.DocumentXMLBinaryData, error) {
	c.documentID = documentID
	if c.documentXMLBinaryData.DocumentID == "" {
		c.documentXMLBinaryData.DocumentID = documentID
	}
	return c.documentXMLBinaryData, nil
}

func (c *cmdFakeClient) DocumentXMLDataAsString(_ context.Context, _ string, documentID string) (api.DocumentXMLData, error) {
	c.documentID = documentID
	if c.documentXMLData.DocumentID == "" {
		c.documentXMLData.DocumentID = documentID
	}
	return c.documentXMLData, nil
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

func (c *cmdFakeClient) Menu(context.Context, string) ([]api.MenuEntry, error) {
	return c.menuEntries, nil
}
