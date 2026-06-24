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
