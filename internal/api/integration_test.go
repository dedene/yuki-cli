package api

import (
	"context"
	"strings"
	"testing"
)

func TestAdministrationDataParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "Integration", "GetAdministrationData", administrationDataResponse, "sessionId", func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:sessionId>session-1</they:sessionId>",
			"<they:administrationId>admin-1</they:administrationId>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	data, err := client.AdministrationData(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("AdministrationData: %v", err)
	}
	if data.CompanyName != "Highpro BV" ||
		data.MainContactEmail != "connections@yuki.be" ||
		data.AddressLine1 != "Orteliuskaai 2" ||
		data.City != "Antwerpen" ||
		data.Country != "BE" ||
		data.CompanyLogoB64 != "base64-logo" ||
		data.IBAN != "BExxxxxxxxxx" ||
		data.BankAccountName != "Highpro NV" ||
		data.CoCNumber != "xxxxxxxxxx" ||
		data.VATNumber != "BExxxx.xxx.xxx" {
		t.Fatalf("data = %#v", data)
	}
}

const administrationDataResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetAdministrationDataResponse xmlns="http://www.theyukicompany.com/">
      <GetAdministrationDataResult>
        <CompanyName>Highpro BV</CompanyName>
        <MainContactName>katrien 2</MainContactName>
        <MainContactEmail>connections@yuki.be</MainContactEmail>
        <AddressLine_1>Orteliuskaai 2</AddressLine_1>
        <Postcode>2000</Postcode>
        <City>Antwerpen</City>
        <Country>BE</Country>
        <EmailOutgoingInvoices>dhkatrien@outlook.com</EmailOutgoingInvoices>
        <CompanyLogoB64 xsi:type="xsd:base64Binary">base64-logo</CompanyLogoB64>
        <IBAN>BExxxxxxxxxx</IBAN>
        <BankAccountName>Highpro NV</BankAccountName>
        <CoCNumber>xxxxxxxxxx</CoCNumber>
        <VATNumber>BExxxx.xxx.xxx</VATNumber>
      </GetAdministrationDataResult>
    </GetAdministrationDataResponse>
  </soap:Body>
</soap:Envelope>`
