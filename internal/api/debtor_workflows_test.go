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

func TestOutstandingDebtorItemsByDateParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingDebtorItemsByDate", outstandingDebtorItemsByDateResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>true</they:includeBankTransactions>",
			"<they:sortOrder>DateDesc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-01-31</they:endDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	items, err := client.OutstandingDebtorItemsByDate(context.Background(), "session-1", DebtorItemsOptions{
		AdministrationID:        "admin-1",
		StartDate:               "2020-01-01",
		EndDate:                 "2020-01-31",
		IncludeBankTransactions: true,
		SortOrder:               "DateDesc",
	})
	if err != nil {
		t.Fatalf("OutstandingDebtorItemsByDate: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != "ef1a1588-fd5d-49a2-a478-2609012ddae6" ||
		items[0].Contact != "Apple Sales International" ||
		items[0].AddressLine1 != "Bergweg 25" ||
		items[0].Reference != "XX-12534" {
		t.Fatalf("item = %#v", items[0])
	}
}

func TestOutstandingDebtorItemsByDateOutstandingParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingDebtorItemsByDateOutstanding", outstandingDebtorItemsByDateOutstandingResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>false</they:includeBankTransactions>",
			"<they:sortOrder>DateAsc</they:sortOrder>",
			"<they:dateOutstanding>2020-01-31</they:dateOutstanding>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	items, err := client.OutstandingDebtorItemsByDateOutstanding(context.Background(), "session-1", DebtorItemsOptions{
		AdministrationID:        "admin-1",
		DateOutstanding:         "2020-01-31",
		IncludeBankTransactions: false,
		SortOrder:               "DateAsc",
	})
	if err != nil {
		t.Fatalf("OutstandingDebtorItemsByDateOutstanding: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].ID != "ef1a1588-fd5d-49a2-a478-2609012ddae6" ||
		items[0].Reference != "XX-12534" {
		t.Fatalf("item = %#v", items[0])
	}
}

func TestOutstandingDebtorItemsWithLanguageParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingDebtorItemsWithLanguage", outstandingDebtorItemsWithLanguageResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>true</they:includeBankTransactions>",
			"<they:sortOrder>DateDesc</they:sortOrder>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	items, err := client.OutstandingDebtorItemsWithLanguage(context.Background(), "session-1", DebtorItemsOptions{
		AdministrationID:        "admin-1",
		IncludeBankTransactions: true,
		SortOrder:               "DateDesc",
	})
	if err != nil {
		t.Fatalf("OutstandingDebtorItemsWithLanguage: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].LayoutLanguage != "en" || items[0].Reference != "XX-12534" {
		t.Fatalf("item = %#v", items[0])
	}
}

func TestOutstandingDebtorWithPaymentReferenceParsesWSDLCompatibleResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingDebtorWithPaymentReference", outstandingDebtorWithPaymentReferenceResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>true</they:includeBankTransactions>",
			"<they:sortOrder>DateDesc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-01-31</they:endDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	items, err := client.OutstandingDebtorWithPaymentReference(context.Background(), "session-1", DebtorItemsOptions{
		AdministrationID:        "admin-1",
		StartDate:               "2020-01-01",
		EndDate:                 "2020-01-31",
		IncludeBankTransactions: true,
		SortOrder:               "DateDesc",
	})
	if err != nil {
		t.Fatalf("OutstandingDebtorWithPaymentReference: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].PaymentReference != "RF18539007547034" ||
		items[0].DocumentID != "c9cc2001-2ea4-41ee-a20e-48f40cdf4e38" ||
		items[0].PaymentMethod != "Overschrijving" {
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

const outstandingDebtorWithPaymentReferenceResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingDebtorWithPaymentReferenceResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingDebtorWithPaymentReferenceResult>
        <OutstandingDebtorItems xmlns="">
          <Item ID="ef1a1588-fd5d-49a2-a478-2609012ddae6">
            <Date>2020-01-31</Date>
            <Description>Testfactuur - 1</Description>
            <Contact>Apple Sales International</Contact>
            <ContactID>22f8a673-d5de-4b85-a845-814d15dd33cd</ContactID>
            <OpenAmount>29.76</OpenAmount>
            <OriginalAmount>29.76</OriginalAmount>
            <Type ID="6">Verkoopfactuur</Type>
            <Reference>XX-12534</Reference>
            <PaymentReference>RF18539007547034</PaymentReference>
            <DueDate>2012-07-22</DueDate>
            <DocumentID>c9cc2001-2ea4-41ee-a20e-48f40cdf4e38</DocumentID>
            <PaymentMethod>Overschrijving</PaymentMethod>
          </Item>
        </OutstandingDebtorItems>
      </OutstandingDebtorWithPaymentReferenceResult>
    </OutstandingDebtorWithPaymentReferenceResponse>
  </soap:Body>
</soap:Envelope>`

const outstandingDebtorItemsByDateResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingDebtorItemsByDateResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingDebtorItemsByDateResult>
        <OutstandingDebtorItems xmlns="">
          <Item ID="ef1a1588-fd5d-49a2-a478-2609012ddae6">
            <Date>2020-01-31</Date>
            <Description>Testfactuur - 1</Description>
            <Contact>Apple Sales International</Contact>
            <ContactID>22f8a673-d5de-4b85-a845-814d15dd33cd</ContactID>
            <OpenAmount>29.76</OpenAmount>
            <OriginalAmount>29.76</OriginalAmount>
            <Type ID="6">Verkoopfactuur</Type>
            <Reference>XX-12534</Reference>
            <DueDate>2012-07-22</DueDate>
            <DocumentID>c9cc2001-2ea4-41ee-a20e-48f40cdf4e38</DocumentID>
            <PaymentMethod>Overschrijving</PaymentMethod>
            <ContactCode>1122</ContactCode>
            <CoCNumber />
            <VATNumber />
            <AddressLine_1>Bergweg 25</AddressLine_1>
            <AddressLine_2 />
            <Postcode>1234 AA</Postcode>
            <City>Rotterdam</City>
            <MailAddressLine_1 />
            <MailAddressLine_2 />
            <MailPostcode />
            <MailCity />
            <Country>NL</Country>
            <RecipientEmail />
            <PhoneHome />
            <PhoneWork />
            <EmailHome />
            <EmailWork />
          </Item>
        </OutstandingDebtorItems>
      </OutstandingDebtorItemsByDateResult>
    </OutstandingDebtorItemsByDateResponse>
  </soap:Body>
</soap:Envelope>`

const outstandingDebtorItemsByDateOutstandingResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingDebtorItemsByDateOutstandingResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingDebtorItemsByDateOutstandingResult>
        <OutstandingDebtorItems xmlns="">
          <Item ID="ef1a1588-fd5d-49a2-a478-2609012ddae6">
            <Date>2020-01-31</Date>
            <Description>Testfactuur - 1</Description>
            <Contact>Apple Sales International</Contact>
            <OpenAmount>29.76</OpenAmount>
            <OriginalAmount>29.76</OriginalAmount>
            <Type ID="6">Verkoopfactuur</Type>
            <Reference>XX-12534</Reference>
            <DueDate>2012-07-22</DueDate>
            <DocumentID>c9cc2001-2ea4-41ee-a20e-48f40cdf4e38</DocumentID>
            <PaymentMethod>Overschrijving</PaymentMethod>
          </Item>
        </OutstandingDebtorItems>
      </OutstandingDebtorItemsByDateOutstandingResult>
    </OutstandingDebtorItemsByDateOutstandingResponse>
  </soap:Body>
</soap:Envelope>`

const outstandingDebtorItemsWithLanguageResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingDebtorItemsWithLanguageResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingDebtorItemsWithLanguageResult>
        <OutstandingDebtorItems xmlns="">
          <Item ID="ef1a1588-fd5d-49a2-a478-2609012ddae6">
            <Date>2020-01-31</Date>
            <Description>Testfactuur - 1</Description>
            <Contact>Apple Sales International</Contact>
            <OpenAmount>29.76</OpenAmount>
            <OriginalAmount>29.76</OriginalAmount>
            <Type ID="6">Verkoopfactuur</Type>
            <Reference>XX-12534</Reference>
            <DueDate>2012-07-22</DueDate>
            <DocumentID>c9cc2001-2ea4-41ee-a20e-48f40cdf4e38</DocumentID>
            <PaymentMethod>Overschrijving</PaymentMethod>
            <LayoutLanguage>en</LayoutLanguage>
          </Item>
        </OutstandingDebtorItems>
      </OutstandingDebtorItemsWithLanguageResult>
    </OutstandingDebtorItemsWithLanguageResponse>
  </soap:Body>
</soap:Envelope>`
