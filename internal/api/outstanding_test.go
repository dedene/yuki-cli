package api

import (
	"context"
	"strings"
	"testing"
)

func TestCheckOutstandingItemParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "CheckOutstandingItem", checkOutstandingItemResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:Reference>NV2018/156</they:Reference>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	items, err := client.CheckOutstandingItem(context.Background(), "session-1", "NV2018/156")
	if err != nil {
		t.Fatalf("CheckOutstandingItem: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != "99a9dbe2-b010-4f78-876c-46b113249096" ||
		items[0].Reference != "NV2018/156" ||
		items[0].PhoneHome != "09002025438" ||
		items[0].EmailWork != "klantenservice@bol.com" {
		t.Fatalf("item = %#v", items[0])
	}
}

const checkOutstandingItemResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CheckOutstandingItemResponse xmlns="http://www.theyukicompany.com/">
      <CheckOutstandingItemResult>
        <OutstandingItems xmlns="">
          <Item ID="99a9dbe2-b010-4f78-876c-46b113249096">
            <Date>2020-12-16</Date>
            <Description>Factuur voor  Bol.com</Description>
            <Contact> Bol.com</Contact>
            <ContactID>a79d3806-829e-48ea-833e-469dd282c109</ContactID>
            <OpenAmount>91.96</OpenAmount>
            <OriginalAmount>91.96</OriginalAmount>
            <Type ID="6">Sales invoice</Type>
            <Reference>NV2018/156</Reference>
            <DueDate>2020-12-30</DueDate>
            <PaymentMethod>Electronic transfer</PaymentMethod>
            <ContactCode />
            <CoCNumber>0824.148.721</CoCNumber>
            <VATNumber>NL820471616B01</VATNumber>
            <AddressLine_1>Papendorpseweg 100</AddressLine_1>
            <AddressLine_2 />
            <Postcode>3528 BJ</Postcode>
            <City>UTRECHT</City>
            <MailAddressLine_1>Papendorpseweg 100</MailAddressLine_1>
            <MailAddressLine_2 />
            <MailPostcode>3528 BJ</MailPostcode>
            <MailCity>UTRECHT</MailCity>
            <Country>NL</Country>
            <RecipientEmail />
            <PhoneHome>09002025438</PhoneHome>
            <PhoneWork />
            <EmailHome>klantenservice@bol.com</EmailHome>
            <EmailWork>klantenservice@bol.com</EmailWork>
          </Item>
        </OutstandingItems>
      </CheckOutstandingItemResult>
    </CheckOutstandingItemResponse>
  </soap:Body>
</soap:Envelope>`
