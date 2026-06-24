package api

import (
	"context"
	"strings"
	"testing"
)

func TestSearchContactsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Contact", "SearchContacts", searchContactsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:domainID>domain-1</they:domainID>",
			"<they:searchOption>ID</they:searchOption>",
			"<they:searchValue>a79d3806-xxxx-xxxx-xxxxx-469dd282c109</they:searchValue>",
			"<they:sortOrder>CreatedDesc</they:sortOrder>",
			"<they:modifiedAfter>2018-01-01</they:modifiedAfter>",
			"<they:active>Active</they:active>",
			"<they:pageNumber>1</they:pageNumber>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	contacts, err := client.SearchContacts(context.Background(), "session-1", ContactSearchOptions{
		DomainID:      "domain-1",
		SearchOption:  "ID",
		SearchValue:   "a79d3806-xxxx-xxxx-xxxxx-469dd282c109",
		SortOrder:     "CreatedDesc",
		ModifiedAfter: "2018-01-01",
		Active:        "Active",
		PageNumber:    1,
	})
	if err != nil {
		t.Fatalf("SearchContacts: %v", err)
	}
	if len(contacts) != 1 {
		t.Fatalf("len(contacts) = %d, want 1", len(contacts))
	}
	contact := contacts[0]
	if contact.ID != "a79d3806-xxxx-xxxx-xxxxx-469dd282c109" ||
		contact.Name != " Bol.com" ||
		contact.EmailWork != "klantenservice@bol.com" ||
		contact.VATNumber != "NL820471616B01" {
		t.Fatalf("contact = %#v", contact)
	}
}

func TestSuppliersAndCustomersParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Contact", "GetSuppliersAndCustomers", suppliersAndCustomersResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:domainID>domain-1</they:domainID>",
			"<they:searchOption>ID</they:searchOption>",
			"<they:sortOrder>CreatedAsc</they:sortOrder>",
			"<they:modifiedAfter>2024-01-01</they:modifiedAfter>",
			"<they:active>Both</they:active>",
			"<they:pageNumber>1</they:pageNumber>",
			"<they:contactType>Both</they:contactType>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	contacts, err := client.SuppliersAndCustomers(context.Background(), "session-1", ContactSearchOptions{
		DomainID:      "domain-1",
		SearchOption:  "ID",
		SearchValue:   "d8cxxxx-xxxx-xxxx-xxxx-xxxxf352b1a",
		SortOrder:     "CreatedAsc",
		ModifiedAfter: "2024-01-01",
		Active:        "Both",
		PageNumber:    1,
		ContactType:   "Both",
	})
	if err != nil {
		t.Fatalf("SuppliersAndCustomers: %v", err)
	}
	if len(contacts) != 1 {
		t.Fatalf("len(contacts) = %d, want 1", len(contacts))
	}
	contact := contacts[0]
	if contact.ID != "d8c438fa-6f15-4ce2-bcb5-562aaf352b1a" ||
		contact.Name != "Telenet" ||
		!contact.IsSupplier ||
		!contact.IsCustomer {
		t.Fatalf("contact = %#v", contact)
	}
}

const searchContactsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <SearchContactsResponse xmlns="http://www.theyukicompany.com/">
      <SearchContactsResult>
        <Contacts xmlns="">
          <Contact ID="a79d3806-xxxx-xxxx-xxxxx-469dd282c109">
            <Type>Company</Type>
            <HID>139</HID>
            <Code />
            <Name> Bol.com</Name>
            <AddressLine_1>Papendorpseweg 100</AddressLine_1>
            <AddressLine_2 />
            <Postcode>3528 BJ</Postcode>
            <City>UTRECHT</City>
            <MailAddressLine_1>Papendorpseweg 100</MailAddressLine_1>
            <MailAddressLine_2 />
            <MailPostcode>3528 BJ</MailPostcode>
            <MailCity>UTRECHT</MailCity>
            <Country>NL</Country>
            <PhoneHome>09002025438</PhoneHome>
            <EMailWork>klantenservice@bol.com</EMailWork>
            <Website>www.bol.com</Website>
            <VATNumber>NL820471616B01</VATNumber>
            <CoCNumber>0824.148.721</CoCNumber>
            <IncomeTaxNumber />
            <Created>4/10/2019 3:50:52 PM</Created>
            <Modified>12/18/2019 10:21:40 AM</Modified>
            <MainContactPerson>ba62b21c-xxxx-xxxxx-xxxxx-f7f3b2daab79</MainContactPerson>
            <Tags>F281_PA, Telecom</Tags>
            <BackofficeStatus>Inactief</BackofficeStatus>
          </Contact>
        </Contacts>
      </SearchContactsResult>
    </SearchContactsResponse>
  </soap:Body>
</soap:Envelope>`

const suppliersAndCustomersResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetSuppliersAndCustomersResponse xmlns="http://www.theyukicompany.com/">
      <GetSuppliersAndCustomersResult>
        <Contacts xmlns="">
          <Contact ID="d8c438fa-6f15-4ce2-bcb5-562aaf352b1a">
            <Type>Company</Type>
            <HID>19</HID>
            <Code>CL0005</Code>
            <Name>Telenet</Name>
            <AddressLine_1>Diksmuidelaan 25</AddressLine_1>
            <AddressLine_2 />
            <Postcode>2000</Postcode>
            <City>Antwerpen</City>
            <Country>BE</Country>
            <Created>3/9/2020 9:56:07 AM</Created>
            <Modified>9/12/2024 10:57:25 AM</Modified>
            <IsSupplier>True</IsSupplier>
            <IsCustomer>True</IsCustomer>
            <Tags>F281_PA, Telecom</Tags>
            <BackofficeStatus>Inactief</BackofficeStatus>
          </Contact>
        </Contacts>
      </GetSuppliersAndCustomersResult>
    </GetSuppliersAndCustomersResponse>
  </soap:Body>
</soap:Envelope>`
