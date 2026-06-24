package api

import (
	"context"
	"strings"
	"testing"
)

func TestGLAccountBalanceParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "GLAccountBalance", glAccountBalanceResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:transactionDate>2020-12-31</they:transactionDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	balances, err := client.GLAccountBalance(context.Background(), "session-1", GLAccountBalanceOptions{
		AdministrationID: "admin-1",
		TransactionDate:  "2020-12-31",
	})
	if err != nil {
		t.Fatalf("GLAccountBalance: %v", err)
	}
	if len(balances) != 2 {
		t.Fatalf("len(balances) = %d, want 2", len(balances))
	}
	if balances[0].Code != "100000" ||
		balances[0].BalanceType != "B" ||
		balances[0].Description != "Share capital" ||
		balances[0].Amount != "-1222.22" {
		t.Fatalf("balance = %#v", balances[0])
	}
}

func TestNetRevenueParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "NetRevenue", netRevenueResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:StartDate>2020-01-01</they:StartDate>",
			"<they:EndDate>2020-01-31</they:EndDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	report, err := client.NetRevenue(context.Background(), "session-1", RevenueOptions{
		AdministrationID: "admin-1",
		StartDate:        "2020-01-01",
		EndDate:          "2020-01-31",
	})
	if err != nil {
		t.Fatalf("NetRevenue: %v", err)
	}
	if report.Amount != "1868.36" ||
		report.StartDate != "2020-01-01" ||
		report.EndDate != "2020-01-31" {
		t.Fatalf("report = %#v", report)
	}
}

func TestNetRevenueFiscalParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "NetRevenueFiscal", netRevenueFiscalResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:StartDate>2020-01-01</they:StartDate>",
			"<they:EndDate>2020-01-31</they:EndDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	report, err := client.NetRevenueFiscal(context.Background(), "session-1", RevenueOptions{
		AdministrationID: "admin-1",
		StartDate:        "2020-01-01",
		EndDate:          "2020-01-31",
	})
	if err != nil {
		t.Fatalf("NetRevenueFiscal: %v", err)
	}
	if report.Amount != "1868.36" ||
		report.StartDate != "2020-01-01" ||
		report.EndDate != "2020-01-31" {
		t.Fatalf("report = %#v", report)
	}
}

func TestGLAccountBalanceFiscalParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "GLAccountBalanceFiscal", glAccountBalanceFiscalResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:transactionDate>2020-12-31</they:transactionDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	balances, err := client.GLAccountBalanceFiscal(context.Background(), "session-1", GLAccountBalanceOptions{
		AdministrationID: "admin-1",
		TransactionDate:  "2020-12-31",
	})
	if err != nil {
		t.Fatalf("GLAccountBalanceFiscal: %v", err)
	}
	if len(balances) != 1 {
		t.Fatalf("len(balances) = %d, want 1", len(balances))
	}
	if balances[0].Code != "100000" ||
		balances[0].Description != "Geplaatst kapitaal" ||
		balances[0].Amount != "-1222.22" {
		t.Fatalf("balance = %#v", balances[0])
	}
}

func TestGLAccountBalanceYearEndParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "GLAccountBalanceYearEnd", glAccountBalanceYearEndResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:transactionDate>2020-12-31</they:transactionDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	balances, err := client.GLAccountBalanceYearEnd(context.Background(), "session-1", GLAccountBalanceOptions{
		AdministrationID: "admin-1",
		TransactionDate:  "2020-12-31",
	})
	if err != nil {
		t.Fatalf("GLAccountBalanceYearEnd: %v", err)
	}
	if len(balances) != 1 {
		t.Fatalf("len(balances) = %d, want 1", len(balances))
	}
	if balances[0].Code != "140000" ||
		balances[0].Description != "Overgedragen winst" ||
		balances[0].Amount != "-1454.14" {
		t.Fatalf("balance = %#v", balances[0])
	}
}

const glAccountBalanceResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GLAccountBalanceResponse xmlns="http://www.theyukicompany.com/">
      <GLAccountBalanceResult>
        <GLAccountBalance xmlns="">
          <GLAccount Code="100000" BalanceType="B">
            <Description>Share capital</Description>
            <Amount>-1222.22</Amount>
          </GLAccount>
          <GLAccount Code="100100" BalanceType="B">
            <Description>Capital De Herdt Katrien</Description>
            <Amount>555.81</Amount>
          </GLAccount>
        </GLAccountBalance>
      </GLAccountBalanceResult>
    </GLAccountBalanceResponse>
  </soap:Body>
</soap:Envelope>`

const netRevenueResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <NetRevenueResponse xmlns="http://www.theyukicompany.com/">
      <NetRevenueResult>1868.36</NetRevenueResult>
    </NetRevenueResponse>
  </soap:Body>
</soap:Envelope>`

const netRevenueFiscalResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <NetRevenueFiscalResponse xmlns="http://www.theyukicompany.com/">
      <NetRevenueFiscalResult>1868.36</NetRevenueFiscalResult>
    </NetRevenueFiscalResponse>
  </soap:Body>
</soap:Envelope>`

const glAccountBalanceFiscalResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GLAccountBalanceFiscalResponse xmlns="http://www.theyukicompany.com/">
      <GLAccountBalanceFiscalResult>
        <GLAccountBalance xmlns="">
          <GLAccount Code="100000" BalanceType="B">
            <Description>Geplaatst kapitaal</Description>
            <Amount>-1222.22</Amount>
          </GLAccount>
        </GLAccountBalance>
      </GLAccountBalanceFiscalResult>
    </GLAccountBalanceFiscalResponse>
  </soap:Body>
</soap:Envelope>`

const glAccountBalanceYearEndResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GLAccountBalanceYearEndResponse xmlns="http://www.theyukicompany.com/">
      <GLAccountBalanceYearEndResult>
        <GLAccountBalance xmlns="">
          <GLAccount Code="140000" BalanceType="B">
            <Description>Overgedragen winst</Description>
            <Amount>-1454.14</Amount>
          </GLAccount>
        </GLAccountBalance>
      </GLAccountBalanceYearEndResult>
    </GLAccountBalanceYearEndResponse>
  </soap:Body>
</soap:Envelope>`
