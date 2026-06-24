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
