package api

import (
	"context"
	"strings"
	"testing"
)

func TestPeriodDateTableParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetPeriodDateTable", periodDateTableResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:yearID>2020</they:yearID>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	period, err := client.PeriodDateTable(context.Background(), "session-1", PeriodDateTableOptions{
		AdministrationID: "admin-1",
		YearID:           2020,
	})
	if err != nil {
		t.Fatalf("PeriodDateTable: %v", err)
	}
	if period.AdministrationID != "admin-1" ||
		period.YearID != 2020 ||
		period.Name != "Highpro NV" ||
		period.Period != "2021-01-02T00:00:00" ||
		period.WholePeriod != "2021-01-02T00:00:00 2022-01-01T00:00:00" ||
		period.ISO8601Period {
		t.Fatalf("period = %#v", period)
	}
}

func TestRGSSchemeParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetRGSScheme", rgsSchemeResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:rgsVersion>2.0</they:rgsVersion>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	entries, err := client.RGSScheme(context.Background(), "session-1", RGSSchemeOptions{
		AdministrationID: "admin-1",
		RGSVersion:       "2.0",
	})
	if err != nil {
		t.Fatalf("RGSScheme: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2", len(entries))
	}
	if entries[0].AdministrationID != "admin-1" ||
		entries[0].RGSVersion != "2.0" ||
		entries[0].YukiCode != "100000" ||
		entries[0].YukiIsEnabled != "True" ||
		entries[0].RGSReferenceCode != "BEivGokGea" ||
		entries[1].YukiCode != "101000" ||
		entries[1].RGSReferenceCode != "" {
		t.Fatalf("entries = %#v", entries)
	}
}

func TestStartBalanceByGLAccountParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetStartBalanceByGlAccount", startBalanceByGLAccountResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:bookyear>2018</they:bookyear>",
			"<they:financialMode>1</they:financialMode>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	balances, err := client.StartBalanceByGLAccount(context.Background(), "session-1", StartBalanceByGLAccountOptions{
		AdministrationID: "admin-1",
		Bookyear:         2018,
		FinancialMode:    1,
	})
	if err != nil {
		t.Fatalf("StartBalanceByGLAccount: %v", err)
	}
	if len(balances) != 2 {
		t.Fatalf("len(balances) = %d, want 2", len(balances))
	}
	if balances[0].AdministrationID != "admin-1" ||
		balances[0].Bookyear != 2018 ||
		balances[0].FinancialMode != 1 ||
		balances[0].AccountID != "100000" ||
		balances[0].StartBalance != "-500.00" ||
		balances[1].AccountID != "140000" {
		t.Fatalf("balances = %#v", balances)
	}
}

const periodDateTableResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetPeriodDateTableResponse xmlns="http://www.theyukicompany.com/">
      <GetPeriodDateTableResult>
        <name>Highpro NV</name>
        <period>2021-01-02T00:00:00</period>
        <wholePeriod>2021-01-02T00:00:00 2022-01-01T00:00:00</wholePeriod>
        <ISO8601Period>false</ISO8601Period>
      </GetPeriodDateTableResult>
    </GetPeriodDateTableResponse>
  </soap:Body>
</soap:Envelope>`

const rgsSchemeResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetRGSSchemeResponse xmlns="http://www.theyukicompany.com/">
      <GetRGSSchemeResult>
        <RGSEntry>
          <YukiCode>100000</YukiCode>
          <YukiIsEnabled>True</YukiIsEnabled>
          <YukiDescription>Geplaatst kapitaal</YukiDescription>
          <RgsReferenceCode>BEivGokGea</RgsReferenceCode>
          <RgsDescription>Normale aandelen aandelenkapitaal</RgsDescription>
        </RGSEntry>
        <RGSEntry>
          <YukiCode>101000</YukiCode>
          <YukiIsEnabled>True</YukiIsEnabled>
          <YukiDescription>Niet-opgevraagd kapitaal (-)</YukiDescription>
        </RGSEntry>
      </GetRGSSchemeResult>
    </GetRGSSchemeResponse>
  </soap:Body>
</soap:Envelope>`

const startBalanceByGLAccountResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetStartBalanceByGlAccountResponse xmlns="http://www.theyukicompany.com/">
      <GetStartBalanceByGlAccountResult>
        <AccountStartBalance>
          <accountID>100000</accountID>
          <startBalance>-500.00</startBalance>
          <accountDescription>Share capital</accountDescription>
        </AccountStartBalance>
        <AccountStartBalance>
          <accountID>140000</accountID>
          <startBalance>-1454.14</startBalance>
          <accountDescription>Retained earnings</accountDescription>
        </AccountStartBalance>
      </GetStartBalanceByGlAccountResult>
    </GetStartBalanceByGlAccountResponse>
  </soap:Body>
</soap:Envelope>`
