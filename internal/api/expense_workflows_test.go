package api

import (
	"context"
	"strings"
	"testing"
)

func TestOutstandingCreditorItemsByDateParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingCreditorItemsByDate", creditorItemsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>false</they:includeBankTransactions>",
			"<they:sortOrder>DateAsc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-01-31</they:endDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	items, err := client.OutstandingCreditorItemsByDate(context.Background(), "session-1", CreditorItemsOptions{
		AdministrationID:        "admin-1",
		StartDate:               "2020-01-01",
		EndDate:                 "2020-01-31",
		IncludeBankTransactions: false,
		SortOrder:               "DateAsc",
	})
	if err != nil {
		t.Fatalf("OutstandingCreditorItemsByDate: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].DocumentID != "c5057bb0-652e-4f8a-ab71-7ecf0e00b82f" ||
		items[0].PaymentMethod != "Creditcard" ||
		items[0].Type.Text != "Aankoopfactuur" ||
		items[0].Type.ID != "2" {
		t.Fatalf("item = %#v", items[0])
	}
}

func TestOutstandingCreditorItemsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingCreditorItems", outstandingCreditorItemsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:includeBankTransactions>true</they:includeBankTransactions>",
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

	items, err := client.OutstandingCreditorItems(context.Background(), "session-1", CreditorItemsOptions{
		AdministrationID:        "admin-1",
		IncludeBankTransactions: true,
		SortOrder:               "DateAsc",
	})
	if err != nil {
		t.Fatalf("OutstandingCreditorItems: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].DocumentID != "a76c5a71-818e-4d94-b28c-1adc71acd285" ||
		items[0].Reference != "S05233212" ||
		items[0].PaymentMethod != "Overschrijving" {
		t.Fatalf("item = %#v", items[0])
	}
}

func TestOutstandingCreditorWithPaymentReferenceParsesWSDLCompatibleResponse(t *testing.T) {
	client := fixtureClientForService(t, "Accounting", "OutstandingCreditorWithPaymentReference", outstandingCreditorWithPaymentReferenceResponse, func(t *testing.T, body string) {
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

	items, err := client.OutstandingCreditorWithPaymentReference(context.Background(), "session-1", CreditorItemsOptions{
		AdministrationID:        "admin-1",
		StartDate:               "2020-01-01",
		EndDate:                 "2020-01-31",
		IncludeBankTransactions: true,
		SortOrder:               "DateDesc",
	})
	if err != nil {
		t.Fatalf("OutstandingCreditorWithPaymentReference: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(items))
	}
	if items[0].PaymentReference != "RF18539007547034" ||
		items[0].DocumentID != "c5057bb0-652e-4f8a-ab71-7ecf0e00b82f" ||
		items[0].PaymentMethod != "Creditcard" {
		t.Fatalf("item = %#v", items[0])
	}
}

func TestTransactionDetailsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetTransactionDetails", transactionDetailsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:GLAccountCode>400000</they:GLAccountCode>",
			"<they:StartDate>2020-08-01</they:StartDate>",
			"<they:EndDate>2020-08-12</they:EndDate>",
			"<they:financialMode>1</they:financialMode>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	transactions, err := client.TransactionDetails(context.Background(), "session-1", TransactionDetailsOptions{
		AdministrationID: "admin-1",
		GLAccountCode:    "400000",
		StartDate:        "2020-08-01",
		EndDate:          "2020-08-12",
		FinancialMode:    "1",
	})
	if err != nil {
		t.Fatalf("TransactionDetails: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("len(transactions) = %d, want 1", len(transactions))
	}
	if transactions[0].ID != "dbce2622-bddf-42de-9d9b-2aff4319f592" ||
		transactions[0].DocumentID != "b3b9d1de-3f29-44e7-b6a4-6b79d296d0e2" ||
		transactions[0].GLAccountCode != "400000" {
		t.Fatalf("transaction = %#v", transactions[0])
	}
}

func TestTransactionDocumentParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetTransactionDocument", transactionDocumentResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:transactionID>tx-1</they:transactionID>") {
			t.Fatalf("request body missing transaction ID:\n%s", body)
		}
	})

	document, err := client.TransactionDocument(context.Background(), "session-1", "admin-1", "tx-1")
	if err != nil {
		t.Fatalf("TransactionDocument: %v", err)
	}
	if document.FileName != "Invoice 201900001.pdf" || document.FileData != "JVBERg==" {
		t.Fatalf("document = %#v", document)
	}
}

func TestFindDocumentParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "FindDocument", findDocumentResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	document, err := client.FindDocument(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("FindDocument: %v", err)
	}
	if document.ID != "150458c9-fe55-4658-911d-055428ccac69" ||
		document.FileName != "Invoice A1040.pdf" ||
		document.Folder.Text != "Verkoop" ||
		document.Folder.ID != "2" {
		t.Fatalf("document = %#v", document)
	}
}

func TestDocumentFileParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentFile", documentFileResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	file, err := client.DocumentFile(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentFile: %v", err)
	}
	if file.ID != "e9c1f89a-a970-4368-b09b-82e61fd56b4d" ||
		file.FileName != "Invoice NV2018/147.pdf" ||
		file.FileData != "JVBERg==" {
		t.Fatalf("file = %#v", file)
	}
}

func TestDocumentBinaryDataParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentBinaryData", documentBinaryDataResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	data, err := client.DocumentBinaryData(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentBinaryData: %v", err)
	}
	if data.DocumentID != "doc-1" || data.FileData != "JVBERg==" {
		t.Fatalf("data = %#v", data)
	}
}

func TestDocumentImageCountParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentImageCount", documentImageCountResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	count, err := client.DocumentImageCount(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentImageCount: %v", err)
	}
	if count.DocumentID != "doc-1" || count.ImageCount != 0 {
		t.Fatalf("count = %#v", count)
	}
}

func TestDocumentXMLDataAsStringParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentXMLDataAsString", documentXMLDataAsStringResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	data, err := client.DocumentXMLDataAsString(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentXMLDataAsString: %v", err)
	}
	want := "<SalesInvoice><Reference>A1040</Reference><Subject>Testfactuur - 1</Subject></SalesInvoice>"
	if data.DocumentID != "doc-1" || data.XML != want {
		t.Fatalf("data = %#v", data)
	}
}

const creditorItemsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingCreditorItemsByDateResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingCreditorItemsByDateResult>
        <OutstandingCreditorItems xmlns="">
          <Item ID="3506683a-f5d1-4904-8957-01c3bc8f8879">
            <Date>2020-01-03</Date>
            <Description>Factuur van AD Delhaize, Goodwill</Description>
            <Contact>AD Delhaize</Contact>
            <ContactID>6249d031-e9b1-429a-b417-f21cfb0e5fb0</ContactID>
            <OpenAmount>242.00</OpenAmount>
            <OriginalAmount>242.00</OriginalAmount>
            <Type ID="2">Aankoopfactuur</Type>
            <Reference>test</Reference>
            <DueDate>2020-01-03</DueDate>
            <DocumentID>c5057bb0-652e-4f8a-ab71-7ecf0e00b82f</DocumentID>
            <PaymentMethod>Creditcard</PaymentMethod>
            <ContactCode>062</ContactCode>
            <VATNumber>BE0402.206.045</VATNumber>
            <Country>BE</Country>
          </Item>
        </OutstandingCreditorItems>
      </OutstandingCreditorItemsByDateResult>
    </OutstandingCreditorItemsByDateResponse>
  </soap:Body>
</soap:Envelope>`

const outstandingCreditorItemsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingCreditorItemsResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingCreditorItemsResult>
        <OutstandingCreditorItems xmlns="">
          <Item ID="eaa937e9-951a-46f1-82e4-fc4e7ad9f07f">
            <Date>2015-03-09</Date>
            <Description>Factuur van Belgian Shell S.A., Brandstoffen bedrijfswagens</Description>
            <Contact>Belgian Shell S.A.</Contact>
            <ContactID>a9fbcd16-32f2-4372-adf4-f186eb515c60</ContactID>
            <OpenAmount>75.86</OpenAmount>
            <OriginalAmount>75.86</OriginalAmount>
            <Type ID="2">Aankoopfactuur</Type>
            <Reference>S05233212</Reference>
            <DueDate>2015-03-09</DueDate>
            <PaymentMethod>Overschrijving</PaymentMethod>
            <DocumentID>a76c5a71-818e-4d94-b28c-1adc71acd285</DocumentID>
            <ContactCode />
            <CoCNumber>0403048262</CoCNumber>
            <VATNumber>BE0403.048.262</VATNumber>
            <Country>BE</Country>
          </Item>
        </OutstandingCreditorItems>
      </OutstandingCreditorItemsResult>
    </OutstandingCreditorItemsResponse>
  </soap:Body>
</soap:Envelope>`

const outstandingCreditorWithPaymentReferenceResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <OutstandingCreditorWithPaymentReferenceResponse xmlns="http://www.theyukicompany.com/">
      <OutstandingCreditorWithPaymentReferenceResult>
        <OutstandingCreditorItems xmlns="">
          <Item ID="3506683a-f5d1-4904-8957-01c3bc8f8879">
            <Date>2020-01-03</Date>
            <Description>Factuur van AD Delhaize, Goodwill</Description>
            <Contact>AD Delhaize</Contact>
            <ContactID>6249d031-e9b1-429a-b417-f21cfb0e5fb0</ContactID>
            <OpenAmount>242.00</OpenAmount>
            <OriginalAmount>242.00</OriginalAmount>
            <Type ID="2">Aankoopfactuur</Type>
            <Reference>test</Reference>
            <PaymentReference>RF18539007547034</PaymentReference>
            <DueDate>2020-01-03</DueDate>
            <DocumentID>c5057bb0-652e-4f8a-ab71-7ecf0e00b82f</DocumentID>
            <PaymentMethod>Creditcard</PaymentMethod>
          </Item>
        </OutstandingCreditorItems>
      </OutstandingCreditorWithPaymentReferenceResult>
    </OutstandingCreditorWithPaymentReferenceResponse>
  </soap:Body>
</soap:Envelope>`

const transactionDetailsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetTransactionDetailsResponse xmlns="http://www.theyukicompany.com/">
      <GetTransactionDetailsResult>
        <TransactionInfo>
          <id>dbce2622-bddf-42de-9d9b-2aff4319f592</id>
          <hID>55613</hID>
          <transactionDate>2020-08-06T00:00:00</transactionDate>
          <description>Factuur voor JOris</description>
          <transactionAmount>363.00</transactionAmount>
          <transactionAmountForeignCurrency>363.00</transactionAmountForeignCurrency>
          <currencyRate>1.000000</currencyRate>
          <currency>EUR</currency>
          <fullName>JOris</fullName>
          <contactCountry>DK</contactCountry>
          <glAccountCode>400000</glAccountCode>
          <documentID>b3b9d1de-3f29-44e7-b6a4-6b79d296d0e2</documentID>
          <documentReference>test</documentReference>
          <documentType>TRMSales invoice (Prescanned)</documentType>
          <documentFolder>TRMSales</documentFolder>
          <documentFolderTab>TRMInvoices</documentFolderTab>
          <contactID>99dacb8b-7daa-44aa-8a87-223e7d70ce75</contactID>
          <periodId>202008</periodId>
          <company>Highpro BV</company>
          <mutationUser>yuki</mutationUser>
        </TransactionInfo>
      </GetTransactionDetailsResult>
    </GetTransactionDetailsResponse>
  </soap:Body>
</soap:Envelope>`

const transactionDocumentResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetTransactionDocumentResponse xmlns="http://www.theyukicompany.com/">
      <GetTransactionDocumentResult>
        <fileName>Invoice 201900001.pdf</fileName>
        <filedata>JVBERg==</filedata>
      </GetTransactionDocumentResult>
    </GetTransactionDocumentResponse>
  </soap:Body>
</soap:Envelope>`

const findDocumentResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <FindDocumentResponse xmlns="http://www.theyukicompany.com/">
      <FindDocumentResult>
        <Documents xmlns="">
          <Document ID="150458c9-fe55-4658-911d-055428ccac69">
            <Subject>Testfactuur - 1</Subject>
            <DocumentDate>2021-09-01</DocumentDate>
            <Amount>861.89</Amount>
            <Folder ID="2">Verkoop</Folder>
            <Tab ID="201">Facturen</Tab>
            <Type>6</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <FileName>Invoice A1040.pdf</FileName>
            <ContentType>application/pdf</ContentType>
            <FileSize>95084</FileSize>
            <ContactName>James Bond</ContactName>
            <Created>2021-09-16T16:46:15</Created>
            <Creator>yuki</Creator>
            <Modified>2021-09-16T16:46:15</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </FindDocumentResult>
    </FindDocumentResponse>
  </soap:Body>
</soap:Envelope>`

const documentFileResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentFileResponse xmlns="http://www.theyukicompany.com/">
      <DocumentFileResult>
        <Document ID="e9c1f89a-a970-4368-b09b-82e61fd56b4d" xmlns="">
          <FileName>Invoice NV2018/147.pdf</FileName>
          <FileSize>110254.00</FileSize>
          <FileData>JVBERg==</FileData>
        </Document>
      </DocumentFileResult>
    </DocumentFileResponse>
  </soap:Body>
</soap:Envelope>`

const documentBinaryDataResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentBinaryDataResponse xmlns="http://www.theyukicompany.com/">
      <DocumentBinaryDataResult>JVBERg==</DocumentBinaryDataResult>
    </DocumentBinaryDataResponse>
  </soap:Body>
</soap:Envelope>`

const documentImageCountResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentImageCountResponse xmlns="http://www.theyukicompany.com/">
      <DocumentImageCountResult>0</DocumentImageCountResult>
    </DocumentImageCountResponse>
  </soap:Body>
</soap:Envelope>`

const documentXMLDataAsStringResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentXMLDataAsStringResponse xmlns="http://www.theyukicompany.com/">
      <DocumentXMLDataAsStringResult>&lt;SalesInvoice&gt;&lt;Reference&gt;A1040&lt;/Reference&gt;&lt;Subject&gt;Testfactuur - 1&lt;/Subject&gt;&lt;/SalesInvoice&gt;</DocumentXMLDataAsStringResult>
    </DocumentXMLDataAsStringResponse>
  </soap:Body>
</soap:Envelope>`
