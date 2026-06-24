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

func TestGLAccountsRGSSchemeJSONUsesAdministrationAndVersion(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		rgsEntries: []api.RGSEntry{{
			YukiCode:         "100000",
			YukiIsEnabled:    "True",
			YukiDescription:  "Geplaatst kapitaal",
			RGSReferenceCode: "BEivGokGea",
			RGSDescription:   "Normale aandelen aandelenkapitaal",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "rgs-scheme",
		"--administration", "admin-1",
		"--rgs-version", "2.0",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.rgsOpts.AdministrationID != "admin-1" ||
		client.rgsOpts.RGSVersion != "2.0" {
		t.Fatalf("rgsOpts = %#v", client.rgsOpts)
	}
	var entries []api.RGSEntry
	if err := json.Unmarshal(out.Bytes(), &entries); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(entries) != 1 || entries[0].RGSReferenceCode != "BEivGokGea" {
		t.Fatalf("entries = %#v", entries)
	}
}

func TestGLAccountsStartBalanceJSONUsesBookyearAndFinancialMode(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		startBalances: []api.GLAccountStartBalance{{
			AccountID:          "100000",
			StartBalance:       "-500.00",
			AccountDescription: "Share capital",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "gl-accounts", "start-balance",
		"--administration", "admin-1",
		"--bookyear", "2018",
		"--financial-mode", "1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.startBalanceOpts.AdministrationID != "admin-1" ||
		client.startBalanceOpts.Bookyear != 2018 ||
		client.startBalanceOpts.FinancialMode != 1 {
		t.Fatalf("startBalanceOpts = %#v", client.startBalanceOpts)
	}
	var balances []api.GLAccountStartBalance
	if err := json.Unmarshal(out.Bytes(), &balances); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(balances) != 1 || balances[0].AccountID != "100000" {
		t.Fatalf("balances = %#v", balances)
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

func TestPeriodsTableJSONUsesAdministrationAndYear(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		period: api.AdministrationPeriod{
			AdministrationID: "admin-1",
			YearID:           2020,
			Name:             "Highpro NV",
			Period:           "2021-01-02T00:00:00",
			WholePeriod:      "2021-01-02T00:00:00 2022-01-01T00:00:00",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "periods", "table",
		"--administration", "admin-1",
		"--year", "2020",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.periodOpts.AdministrationID != "admin-1" ||
		client.periodOpts.YearID != 2020 {
		t.Fatalf("periodOpts = %#v", client.periodOpts)
	}
	var period api.AdministrationPeriod
	if err := json.Unmarshal(out.Bytes(), &period); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if period.Name != "Highpro NV" || period.WholePeriod != "2021-01-02T00:00:00 2022-01-01T00:00:00" {
		t.Fatalf("period = %#v", period)
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

func TestDebtorItemsListJSONUsesDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		debtorItems: []api.DebtorItem{{
			Date:           "2020-01-31",
			Contact:        "Apple Sales International",
			OriginalAmount: "29.76",
			OpenAmount:     "29.76",
			PaymentMethod:  "Overschrijving",
			Reference:      "XX-12534",
			DocumentID:     "doc-1",
			Description:    "Testfactuur - 1",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "debtor-items", "list",
		"--administration", "admin-1",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
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
	if client.debtorOpts.AdministrationID != "admin-1" ||
		client.debtorOpts.StartDate != "2020-01-01" ||
		client.debtorOpts.EndDate != "2020-01-31" ||
		!client.debtorOpts.IncludeBankTransactions ||
		client.debtorOpts.SortOrder != "DateDesc" {
		t.Fatalf("debtorOpts = %#v", client.debtorOpts)
	}
	var items []api.DebtorItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].Reference != "XX-12534" {
		t.Fatalf("items = %#v", items)
	}
}

func TestDebtorItemsWithPaymentReferenceJSONUsesFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		debtorItems: []api.DebtorItem{{
			Date:             "2020-01-31",
			Contact:          "Apple Sales International",
			OriginalAmount:   "29.76",
			OpenAmount:       "29.76",
			PaymentMethod:    "Overschrijving",
			Reference:        "XX-12534",
			PaymentReference: "RF18539007547034",
			DocumentID:       "doc-1",
			Description:      "Testfactuur - 1",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "debtor-items", "with-payment-reference",
		"--administration", "admin-1",
		"--from", "2020-01-01",
		"--to", "2020-01-31",
		"--include-bank-transactions",
		"--sort-order", "DateDesc",
		"--payment-method", "Overschrijving",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.debtorOpts.AdministrationID != "admin-1" ||
		client.debtorOpts.StartDate != "2020-01-01" ||
		client.debtorOpts.EndDate != "2020-01-31" ||
		!client.debtorOpts.IncludeBankTransactions ||
		client.debtorOpts.SortOrder != "DateDesc" {
		t.Fatalf("debtorOpts = %#v", client.debtorOpts)
	}
	var items []api.DebtorItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].PaymentReference != "RF18539007547034" {
		t.Fatalf("items = %#v", items)
	}
}

func TestOutstandingCheckJSONUsesReference(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		outstandingItems: []api.OutstandingItem{{
			Date:           "2020-12-16",
			Contact:        "Bol.com",
			OriginalAmount: "91.96",
			OpenAmount:     "91.96",
			PaymentMethod:  "Electronic transfer",
			Reference:      "NV2018/156",
			Description:    "Factuur voor Bol.com",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "outstanding", "check",
		"--reference", "NV2018/156",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.outstandingReference != "NV2018/156" {
		t.Fatalf("outstandingReference = %q", client.outstandingReference)
	}
	var items []api.OutstandingItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].Reference != "NV2018/156" {
		t.Fatalf("items = %#v", items)
	}
}

func TestOutstandingCheckAdminJSONUsesAdministrationAndReference(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		outstandingItems: []api.OutstandingItem{{
			Date:           "2021-01-22",
			Contact:        "blabla 007",
			OriginalAmount: "91.80",
			OpenAmount:     "91.80",
			PaymentMethod:  "Electronic transfer",
			Reference:      "A1010",
			Description:    "Testfactuur - 1",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "outstanding", "check-admin",
		"--administration", "admin-1",
		"--reference", "A1010",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.outstandingAdministrationID != "admin-1" ||
		client.outstandingReference != "A1010" {
		t.Fatalf("outstanding administration/reference = %q/%q", client.outstandingAdministrationID, client.outstandingReference)
	}
	var items []api.OutstandingItem
	if err := json.Unmarshal(out.Bytes(), &items); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(items) != 1 || items[0].Reference != "A1010" {
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

func TestChangeDigestTransactionsJSONUsesDateRangeAndPaging(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		updatedTransactions: []api.UpdatedTransaction{{
			ID:                "tx-1",
			TransactionDate:   "2023-02-06T00:00:00",
			TransactionAmount: "-63.00",
			Currency:          "EUR",
			GLAccountCode:     "494100",
			FullName:          "Biovita Bvba",
			Updated:           "2023-07-25T13:37:58.3",
			Deleted:           "true",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "change-digest", "transactions",
		"--administration", "admin-1",
		"--from", "2025-07-23T00:00:00.00Z",
		"--to", "2025-08-23T13:00:00.00Z",
		"--limit", "100",
		"--start-record", "0",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.updatedTransactionsOpts.AdministrationID != "admin-1" ||
		client.updatedTransactionsOpts.StartDate != "2025-07-23T00:00:00.00Z" ||
		client.updatedTransactionsOpts.EndDate != "2025-08-23T13:00:00.00Z" ||
		client.updatedTransactionsOpts.NumberOfRecords != 100 ||
		client.updatedTransactionsOpts.StartRecord != 0 {
		t.Fatalf("updatedTransactionsOpts = %#v", client.updatedTransactionsOpts)
	}
	var transactions []api.UpdatedTransaction
	if err := json.Unmarshal(out.Bytes(), &transactions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(transactions) != 1 || transactions[0].Deleted != "true" {
		t.Fatalf("transactions = %#v", transactions)
	}
}

func TestChangeDigestDetailJSONUsesTransactionID(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		changeDigestTransaction: api.TransactionInfo{
			ID:                "tx-1",
			TransactionDate:   "2021-01-01T00:00:00",
			GLAccountCode:     "494190",
			TransactionAmount: "0.00",
			FullName:          "Topolino bvba",
			DocumentID:        "doc-1",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "change-digest", "detail",
		"--administration", "admin-1",
		"--transaction", "tx-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.administrationID != "admin-1" || client.transactionID != "tx-1" {
		t.Fatalf("administration/transaction = %q/%q", client.administrationID, client.transactionID)
	}
	var tx api.TransactionInfo
	if err := json.Unmarshal(out.Bytes(), &tx); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if tx.ID != "tx-1" || tx.DocumentID != "doc-1" {
		t.Fatalf("transaction = %#v", tx)
	}
}

func TestProjectsListJSONUsesSearchFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		projects: []api.AccountingProject{{
			HID:         "1",
			Code:        "WELLNESS",
			Description: "Wellness Event",
			Company:     "Highpro BV",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "projects", "list",
		"--administration", "admin-1",
		"--search-option", "Code",
		"--search-value", "WELLNESS",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.projectsOpts.AdministrationID != "admin-1" ||
		client.projectsOpts.SearchOption != "Code" ||
		client.projectsOpts.SearchValue != "WELLNESS" {
		t.Fatalf("projectsOpts = %#v", client.projectsOpts)
	}
	var projects []api.AccountingProject
	if err := json.Unmarshal(out.Bytes(), &projects); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(projects) != 1 || projects[0].Code != "WELLNESS" {
		t.Fatalf("projects = %#v", projects)
	}
}

func TestProjectsListWithIDJSONUsesSearchFlags(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		projects: []api.AccountingProject{{
			ID:          "project-1",
			HID:         "1",
			Code:        "WELLNESS",
			Description: "Wellness Event",
			Company:     "Highpro BV",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "projects", "list-with-id",
		"--administration", "admin-1",
		"--search-option", "Code",
		"--search-value", "WELLNESS",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.projectsAndIDOpts.AdministrationID != "admin-1" ||
		client.projectsAndIDOpts.SearchOption != "Code" ||
		client.projectsAndIDOpts.SearchValue != "WELLNESS" {
		t.Fatalf("projectsAndIDOpts = %#v", client.projectsAndIDOpts)
	}
	var projects []api.AccountingProject
	if err := json.Unmarshal(out.Bytes(), &projects); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(projects) != 1 || projects[0].ID != "project-1" {
		t.Fatalf("projects = %#v", projects)
	}
}

func TestProjectsBalanceJSONUsesScopeAndDateRange(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		projectBalances: []api.ProjectBalance{{
			GLAccountCode: "400000",
			Project:       "Dossier1",
			ProjectCode:   "DOS1",
			Amount:        "542.00",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"accounting", "projects", "balance",
		"--administration", "admin-1",
		"--gl-account", "400000",
		"--project-code", "DOS1",
		"--from", "2018-01-01",
		"--to", "2020-12-31",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.projectBalanceOpts.AdministrationID != "admin-1" ||
		client.projectBalanceOpts.GLAccountCode != "400000" ||
		client.projectBalanceOpts.ProjectCode != "DOS1" ||
		client.projectBalanceOpts.StartDate != "2018-01-01" ||
		client.projectBalanceOpts.EndDate != "2020-12-31" {
		t.Fatalf("projectBalanceOpts = %#v", client.projectBalanceOpts)
	}
	var balances []api.ProjectBalance
	if err := json.Unmarshal(out.Bytes(), &balances); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(balances) != 1 || balances[0].Amount != "542.00" {
		t.Fatalf("balances = %#v", balances)
	}
}

func TestVATCodesActiveJSONUsesAdministration(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		vatCodes: []api.VATCode{{
			Description:     "VAT 21%",
			Type:            "1",
			TypeDescription: "VAT high",
			Percentage:      "21.00",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"vat", "codes", "active",
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
	var codes []api.VATCode
	if err := json.Unmarshal(out.Bytes(), &codes); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(codes) != 1 || codes[0].Percentage != "21.00" {
		t.Fatalf("codes = %#v", codes)
	}
}

func TestVATReturnsListJSONUsesScope(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		vatReturns: []api.VATReturnInfo{{
			StartDate: "2023-07-01T00:00:00",
			EndDate:   "2023-07-31T00:00:00",
			Status:    "Draft",
			Modified:  "2023-08-01T09:14:43.033",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"vat", "returns", "list",
		"--administration", "admin-1",
		"--year", "2023",
		"--modified-after", "2021-01-01",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.vatReturnOpts.AdministrationID != "admin-1" ||
		client.vatReturnOpts.Year != 2023 ||
		client.vatReturnOpts.ModifiedAfter != "2021-01-01" {
		t.Fatalf("vatReturnOpts = %#v", client.vatReturnOpts)
	}
	var returns []api.VATReturnInfo
	if err := json.Unmarshal(out.Bytes(), &returns); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(returns) != 1 || returns[0].Status != "Draft" {
		t.Fatalf("returns = %#v", returns)
	}
}

func TestIntegrationAdministrationDataJSONUsesAdministration(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		administrationData: api.AdministrationIntegrationData{
			CompanyName:      "Highpro BV",
			MainContactEmail: "connections@yuki.be",
			City:             "Antwerpen",
			Country:          "BE",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"integration", "administration-data",
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
	var data api.AdministrationIntegrationData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if data.CompanyName != "Highpro BV" || data.MainContactEmail != "connections@yuki.be" {
		t.Fatalf("data = %#v", data)
	}
}

func TestFiscalTableTotalsJSONUsesCompanyAndYear(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		fiscalTableTotals: api.FiscalTableTotals{
			RevenueTotal:             "1000.00",
			GrossMarginTotal:         "800.00",
			ProfessionalCostsTotal:   "300.00",
			SocialContributionsTotal: "120.00",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"fiscal-table", "totals",
		"--company", "company-1",
		"--year", "2023",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.companyID != "company-1" || client.fiscalTableYear != 2023 {
		t.Fatalf("company/year = %q/%d", client.companyID, client.fiscalTableYear)
	}
	var totals api.FiscalTableTotals
	if err := json.Unmarshal(out.Bytes(), &totals); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if totals.RevenueTotal != "1000.00" || totals.SocialContributionsTotal != "120.00" {
		t.Fatalf("totals = %#v", totals)
	}
}

func TestBackofficeWorkflowJSONUsesAdministration(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		backofficeWorkflow: []api.BackofficeWorkflowDocument{{
			SubmitDate:   "2020-08-26T14:10:05",
			DocumentType: api.XMLText{ID: "2", Text: "Purchase invoice"},
			FileName:     "ININV-00004.pdf",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"backoffice", "workflow",
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
	var documents []api.BackofficeWorkflowDocument
	if err := json.Unmarshal(out.Bytes(), &documents); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(documents) != 1 || documents[0].FileName != "ININV-00004.pdf" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestBackofficeOutstandingQuestionsJSONUsesAdministration(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		backofficeQuestions: []api.BackofficeQuestion{{
			Date:        "2022-03-09T16:57:47",
			Type:        api.XMLText{ID: "29", Text: "Question"},
			Description: "vraag",
			From:        "Katrien",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"backoffice", "outstanding-questions",
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
	var questions []api.BackofficeQuestion
	if err := json.Unmarshal(out.Bytes(), &questions); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(questions) != 1 || questions[0].From != "Katrien" {
		t.Fatalf("questions = %#v", questions)
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
	if client.documentID != "doc-1" || client.documentXMLDataOperation != "string" {
		t.Fatalf("document call = %q/%q", client.documentID, client.documentXMLDataOperation)
	}
	var data api.DocumentXMLData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if data.DocumentID != "doc-1" || data.XML != "<SalesInvoice><Reference>A1040</Reference></SalesInvoice>" {
		t.Fatalf("data = %#v", data)
	}
}

func TestArchiveDocumentsXMLDataJSONPrintsEmbeddedXMLData(t *testing.T) {
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
		"archive", "documents", "xml-data",
		"--document", "doc-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.documentID != "doc-1" || client.documentXMLDataOperation != "raw" {
		t.Fatalf("document call = %q/%q", client.documentID, client.documentXMLDataOperation)
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
	sessionID                   string
	accessKey                   string
	administrationID            string
	companyID                   string
	documentID                  string
	transactionID               string
	domains                     []api.Domain
	accounts                    []api.GLAccount
	rgsEntries                  []api.RGSEntry
	rgsOpts                     api.RGSSchemeOptions
	startBalances               []api.GLAccountStartBalance
	startBalanceOpts            api.StartBalanceByGLAccountOptions
	balances                    []api.GLAccountBalanceItem
	balanceOpts                 api.GLAccountBalanceOptions
	glTransactions              []api.GLAccountTransaction
	glTransactionOpts           api.GLAccountTransactionsOptions
	revenueReport               api.RevenueReport
	revenueOpts                 api.RevenueOptions
	period                      api.AdministrationPeriod
	periodOpts                  api.PeriodDateTableOptions
	creditorItems               []api.CreditorItem
	creditorOpts                api.CreditorItemsOptions
	debtorItems                 []api.DebtorItem
	debtorOpts                  api.DebtorItemsOptions
	outstandingItems            []api.OutstandingItem
	outstandingReference        string
	outstandingAdministrationID string
	richTransactions            []api.Transaction
	transactionsOpts            api.TransactionsOptions
	transactions                []api.TransactionInfo
	transactionOpts             api.TransactionDetailsOptions
	updatedTransactions         []api.UpdatedTransaction
	updatedTransactionsOpts     api.UpdatedAndDeletedTransactionsOptions
	changeDigestTransaction     api.TransactionInfo
	projects                    []api.AccountingProject
	projectsOpts                api.ProjectsOptions
	projectsAndIDOpts           api.ProjectsOptions
	projectBalances             []api.ProjectBalance
	projectBalanceOpts          api.ProjectBalanceOptions
	customMethods               []api.PaymentMethod
	archiveMethods              []api.PaymentMethod
	folders                     []api.DocumentFolder
	folderCounts                []api.DocumentFolderCount
	folderCountYear             int
	tabs                        []api.DocumentFolderTab
	currencies                  []api.Currency
	costCategories              []api.CostCategory
	menuEntries                 []api.MenuEntry
	folderID                    string
	documents                   []api.Document
	documentsOpts               api.DocumentsOptions
	documentsInFolderOpts       api.DocumentsInFolderOptions
	documentsInTabOpts          api.DocumentsInTabOptions
	documentsByTypeOpts         api.DocumentsByTypeOptions
	modifiedInFolderOpts        api.ModifiedDocumentsInFolderOptions
	modifiedByTypeOpts          api.ModifiedDocumentsByTypeOptions
	searchDocuments             []api.Document
	searchDocumentsOpts         api.SearchDocumentsOptions
	document                    api.Document
	documentBundle              []api.Document
	documentFile                api.DocumentFile
	documentBinaryData          api.DocumentBinaryData
	documentImageData           api.DocumentImageData
	documentImageCount          api.DocumentImageCount
	documentXMLBinaryData       api.DocumentXMLBinaryData
	documentXMLData             api.DocumentXMLData
	documentXMLDataOperation    string
	transactionDocument         api.TransactionDocument
	vatCodes                    []api.VATCode
	vatReturns                  []api.VATReturnInfo
	vatReturnOpts               api.VATReturnListOptions
	administrationData          api.AdministrationIntegrationData
	fiscalTableTotals           api.FiscalTableTotals
	fiscalTableYear             int
	backofficeWorkflow          []api.BackofficeWorkflowDocument
	backofficeQuestions         []api.BackofficeQuestion
	maxWidth                    int
	maxHeight                   int
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

func (c *cmdFakeClient) RGSScheme(_ context.Context, _ string, opts api.RGSSchemeOptions) ([]api.RGSEntry, error) {
	c.rgsOpts = opts
	return c.rgsEntries, nil
}

func (c *cmdFakeClient) StartBalanceByGLAccount(_ context.Context, _ string, opts api.StartBalanceByGLAccountOptions) ([]api.GLAccountStartBalance, error) {
	c.startBalanceOpts = opts
	return c.startBalances, nil
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

func (c *cmdFakeClient) CheckOutstandingItem(_ context.Context, _ string, reference string) ([]api.OutstandingItem, error) {
	c.outstandingReference = reference
	return c.outstandingItems, nil
}

func (c *cmdFakeClient) CheckOutstandingItemAdmin(_ context.Context, _ string, administrationID string, reference string) ([]api.OutstandingItem, error) {
	c.outstandingAdministrationID = administrationID
	c.outstandingReference = reference
	return c.outstandingItems, nil
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

func (c *cmdFakeClient) PeriodDateTable(_ context.Context, _ string, opts api.PeriodDateTableOptions) (api.AdministrationPeriod, error) {
	c.periodOpts = opts
	if c.period.AdministrationID == "" {
		c.period.AdministrationID = opts.AdministrationID
	}
	if c.period.YearID == 0 {
		c.period.YearID = opts.YearID
	}
	return c.period, nil
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

func (c *cmdFakeClient) OutstandingDebtorItemsByDate(_ context.Context, _ string, opts api.DebtorItemsOptions) ([]api.DebtorItem, error) {
	c.debtorOpts = opts
	return c.debtorItems, nil
}

func (c *cmdFakeClient) OutstandingDebtorWithPaymentReference(_ context.Context, _ string, opts api.DebtorItemsOptions) ([]api.DebtorItem, error) {
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

func (c *cmdFakeClient) UpdatedAndDeletedTransactions(_ context.Context, _ string, opts api.UpdatedAndDeletedTransactionsOptions) ([]api.UpdatedTransaction, error) {
	c.updatedTransactionsOpts = opts
	return c.updatedTransactions, nil
}

func (c *cmdFakeClient) ChangeDigestTransactionDetail(_ context.Context, _ string, administrationID string, transactionID string) (api.TransactionInfo, error) {
	c.administrationID = administrationID
	c.transactionID = transactionID
	return c.changeDigestTransaction, nil
}

func (c *cmdFakeClient) Projects(_ context.Context, _ string, opts api.ProjectsOptions) ([]api.AccountingProject, error) {
	c.projectsOpts = opts
	return c.projects, nil
}

func (c *cmdFakeClient) ProjectsAndID(_ context.Context, _ string, opts api.ProjectsOptions) ([]api.AccountingProject, error) {
	c.projectsAndIDOpts = opts
	return c.projects, nil
}

func (c *cmdFakeClient) ProjectBalance(_ context.Context, _ string, opts api.ProjectBalanceOptions) ([]api.ProjectBalance, error) {
	c.projectBalanceOpts = opts
	return c.projectBalances, nil
}

func (c *cmdFakeClient) CustomPaymentMethods(_ context.Context, _ string, administrationID string) ([]api.PaymentMethod, error) {
	c.administrationID = administrationID
	return c.customMethods, nil
}

func (c *cmdFakeClient) DocumentFolders(context.Context, string) ([]api.DocumentFolder, error) {
	return c.folders, nil
}

func (c *cmdFakeClient) DocumentFolderCounts(_ context.Context, _ string, year int) ([]api.DocumentFolderCount, error) {
	c.folderCountYear = year
	return c.folderCounts, nil
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

func (c *cmdFakeClient) DocumentXMLData(_ context.Context, _ string, documentID string) (api.DocumentXMLData, error) {
	c.documentID = documentID
	c.documentXMLDataOperation = "raw"
	if c.documentXMLData.DocumentID == "" {
		c.documentXMLData.DocumentID = documentID
	}
	return c.documentXMLData, nil
}

func (c *cmdFakeClient) DocumentXMLDataAsString(_ context.Context, _ string, documentID string) (api.DocumentXMLData, error) {
	c.documentID = documentID
	c.documentXMLDataOperation = "string"
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

func (c *cmdFakeClient) ActiveVATCodes(_ context.Context, _ string, administrationID string) ([]api.VATCode, error) {
	c.administrationID = administrationID
	return c.vatCodes, nil
}

func (c *cmdFakeClient) VATReturns(_ context.Context, _ string, opts api.VATReturnListOptions) ([]api.VATReturnInfo, error) {
	c.vatReturnOpts = opts
	return c.vatReturns, nil
}

func (c *cmdFakeClient) AdministrationData(_ context.Context, _ string, administrationID string) (api.AdministrationIntegrationData, error) {
	c.administrationID = administrationID
	return c.administrationData, nil
}

func (c *cmdFakeClient) FiscalTable(_ context.Context, _ string, companyID string, year int) (api.FiscalTableTotals, error) {
	c.companyID = companyID
	c.fiscalTableYear = year
	return c.fiscalTableTotals, nil
}

func (c *cmdFakeClient) BackofficeWorkflow(_ context.Context, _ string, administrationID string) ([]api.BackofficeWorkflowDocument, error) {
	c.administrationID = administrationID
	return c.backofficeWorkflow, nil
}

func (c *cmdFakeClient) BackofficeOutstandingQuestions(_ context.Context, _ string, administrationID string) ([]api.BackofficeQuestion, error) {
	c.administrationID = administrationID
	return c.backofficeQuestions, nil
}
