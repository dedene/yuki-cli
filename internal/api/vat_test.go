package api

import (
	"context"
	"strings"
	"testing"
)

func TestActiveVATCodesParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Vat", "ActiveVATCodesList", activeVATCodesResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	codes, err := client.ActiveVATCodes(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("ActiveVATCodes: %v", err)
	}
	if len(codes) != 3 {
		t.Fatalf("len(codes) = %d, want 3", len(codes))
	}
	if codes[0].Description != "VAT 12%" ||
		codes[0].Type != "21" ||
		codes[0].TypeDescription != "VAT medium" ||
		codes[0].Percentage != "12.00" ||
		codes[0].StartDate != "" ||
		codes[1].Country != "NL" ||
		codes[1].StartDate != "2021-07-01T00:00:00" ||
		codes[2].Type != "17" {
		t.Fatalf("codes = %#v", codes)
	}
}

const activeVATCodesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <soap:Body>
    <ActiveVATCodesListResponse xmlns="http://www.theyukicompany.com/">
      <ActiveVATCodesListResult>
        <VATCode>
          <description>VAT 12%</description>
          <type>21</type>
          <typeDescription>VAT medium</typeDescription>
          <percentage>12.00</percentage>
          <startDate xsi:nil="true"/>
          <endDate xsi:nil="true"/>
        </VATCode>
        <VATCode>
          <description>OSS NL</description>
          <type>40</type>
          <typeDescription>VAT digital services and/or distance selling (OSS)</typeDescription>
          <percentage>21.00</percentage>
          <country>NL</country>
          <startDate>2021-07-01T00:00:00</startDate>
          <endDate xsi:nil="true"/>
        </VATCode>
        <VATCode>
          <description>VAT reverse-charge</description>
          <type>17</type>
          <typeDescription>VAT reverse-charged, sales</typeDescription>
          <percentage>0.00</percentage>
          <startDate xsi:nil="true"/>
          <endDate xsi:nil="true"/>
        </VATCode>
      </ActiveVATCodesListResult>
    </ActiveVATCodesListResponse>
  </soap:Body>
</soap:Envelope>`
