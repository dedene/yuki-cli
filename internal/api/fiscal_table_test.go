package api

import (
	"context"
	"strings"
	"testing"
)

func TestFiscalTableParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "FiscalTable", "GetFiscalTable", fiscalTableResponse, "sessionId", func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:sessionId>session-1</they:sessionId>",
			"<they:companyId>company-1</they:companyId>",
			"<they:year>2023</they:year>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	totals, err := client.FiscalTable(context.Background(), "session-1", "company-1", 2023)
	if err != nil {
		t.Fatalf("FiscalTable: %v", err)
	}
	if totals.CompanyID != "company-1" ||
		totals.Year != 2023 ||
		totals.RevenueTotal != "1000.00" ||
		totals.GrossMarginTotal != "800.00" ||
		totals.ProfessionalCostsTotal != "300.00" ||
		totals.SocialContributionsTotal != "120.00" {
		t.Fatalf("totals = %#v", totals)
	}
}

const fiscalTableResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetFiscalTableResponse xmlns="http://www.theyukicompany.com/">
      <GetFiscalTableResult>
        <RevenueTotal>1000.00</RevenueTotal>
        <GrossMarginTotal>800.00</GrossMarginTotal>
        <ProfessionalCostsTotal>300.00</ProfessionalCostsTotal>
        <SocialContributionsTotal>120.00</SocialContributionsTotal>
      </GetFiscalTableResult>
    </GetFiscalTableResponse>
  </soap:Body>
</soap:Envelope>`
