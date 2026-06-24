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

func TestVATReturnsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Vat", "VATReturnList", vatReturnListResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:year>2023</they:year>",
			"<they:modifiedAfter>2021-01-01</they:modifiedAfter>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	returns, err := client.VATReturns(context.Background(), "session-1", VATReturnListOptions{
		AdministrationID: "admin-1",
		Year:             2023,
		ModifiedAfter:    "2021-01-01",
	})
	if err != nil {
		t.Fatalf("VATReturns: %v", err)
	}
	if len(returns) != 2 {
		t.Fatalf("len(returns) = %d, want 2", len(returns))
	}
	if returns[0].StartDate != "2023-07-01T00:00:00" ||
		returns[0].EndDate != "2023-07-31T00:00:00" ||
		returns[0].Status != "Draft" ||
		returns[0].SendDate != "" ||
		returns[0].AcknowledgeDate != "" ||
		returns[0].Modified != "2023-08-01T09:14:43.033" ||
		returns[1].StartDate != "2023-03-01T00:00:00" {
		t.Fatalf("returns = %#v", returns)
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

const vatReturnListResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://www.w3.org/2003/05/soap-envelope" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <VATReturnListResponse xmlns="http://www.theyukicompany.com/">
      <VATReturnListResult>
        <VATReturnInfo>
          <startDate>2023-07-01T00:00:00</startDate>
          <endDate>2023-07-31T00:00:00</endDate>
          <status>Draft</status>
          <sendDate xsi:nil="true"/>
          <acknowledgeDate xsi:nil="true"/>
          <modified>2023-08-01T09:14:43.033</modified>
        </VATReturnInfo>
        <VATReturnInfo>
          <startDate>2023-03-01T00:00:00</startDate>
          <endDate>2023-03-31T00:00:00</endDate>
          <status>Draft</status>
          <sendDate xsi:nil="true"/>
          <acknowledgeDate xsi:nil="true"/>
          <modified>2023-03-02T09:30:37.707</modified>
        </VATReturnInfo>
      </VATReturnListResult>
    </VATReturnListResponse>
  </soap:Body>
</soap:Envelope>`
