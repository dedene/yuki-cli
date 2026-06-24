package api

import (
	"context"
	"strings"
	"testing"
)

func TestOutstandingDebtorItemsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingDebtorItems", outstandingDebtorItemsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>false</they:includeBankTransactions>",
			"<they:sortOrder>DateAsc</they:sortOrder>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
		for _, notWant := range []string{"<they:startDate>", "<they:endDate>"} {
			if strings.Contains(body, notWant) {
				t.Fatalf("request body unexpectedly contains %q:\n%s", notWant, body)
			}
		}
	})

	items, err := client.OutstandingDebtorItems(context.Background(), "session-1", DebtorItemsOptions{
		AdministrationID:        "admin-1",
		IncludeBankTransactions: false,
		SortOrder:               "DateAsc",
	})
	if err != nil {
		t.Fatalf("OutstandingDebtorItems: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != "7407c636-3fd2-40d9-af28-a61b7d87ffef" ||
		items[0].DocumentID != "5d5a3213-34ed-4c96-bf2f-f6bb23804597" ||
		items[0].OriginalAmount != "1000.00" ||
		items[0].Type.ID != "21" {
		t.Fatalf("item = %#v", items[0])
	}
}

const outstandingDebtorItemsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingDebtorItemsResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingDebtorItemsResult>
        <OutstandingDebtorItems xmlns="">
          <Item ID="7407c636-3fd2-40d9-af28-a61b7d87ffef">
            <Date>2018-12-30</Date>
            <Description>test</Description>
            <Contact />
            <ContactID />
            <OpenAmount>1000.00</OpenAmount>
            <OriginalAmount>1000.00</OriginalAmount>
            <Type ID="21">Diverse posten boeking</Type>
            <Reference />
            <DueDate />
            <DocumentID>5d5a3213-34ed-4c96-bf2f-f6bb23804597</DocumentID>
            <PaymentMethod />
            <ContactCode />
            <CoCNumber />
            <VATNumber />
            <AddressLine_1 />
            <AddressLine_2 />
            <Postcode />
            <City />
            <MailAddressLine_1 />
            <MailAddressLine_2 />
            <MailPostcode />
            <MailCity />
            <Country />
            <RecipientEmail />
            <PhoneHome />
            <PhoneWork />
            <EmailHome />
            <EmailWork />
          </Item>
        </OutstandingDebtorItems>
      </OutstandingDebtorItemsResult>
    </OutstandingDebtorItemsResponse>
  </soap:Body>
</soap:Envelope>`
