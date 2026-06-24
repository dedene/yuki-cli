package api

import (
	"context"
	"strings"
	"testing"
)

func TestTransactionsParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetTransactions", transactionsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:glAccountCode>550002</they:glAccountCode>",
			"<they:startDate>2021-10-25</they:startDate>",
			"<they:endDate>2021-10-25</they:endDate>",
			"<they:financialMode>0</they:financialMode>",
			"<they:dataGroups>documentprocessed,document,documentmatching</they:dataGroups>",
			"<they:numberOfRecords>100</they:numberOfRecords>",
			"<they:startRecord>1</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	transactions, err := client.Transactions(context.Background(), "session-1", TransactionsOptions{
		AdministrationID: "admin-1",
		GLAccountCode:    "550002",
		StartDate:        "2021-10-25",
		EndDate:          "2021-10-25",
		FinancialMode:    "0",
		DataGroups:       "documentprocessed,document,documentmatching",
		NumberOfRecords:  100,
		StartRecord:      1,
	})
	if err != nil {
		t.Fatalf("Transactions: %v", err)
	}
	if len(transactions) != 1 {
		t.Fatalf("len(transactions) = %d, want 1", len(transactions))
	}
	tx := transactions[0]
	if tx.ID != "tx-1" ||
		tx.GLAccountCode != "550002" ||
		tx.Document == nil ||
		tx.Document.ID != "doc-1" ||
		tx.Document.Reference != "MC-2021-10-25" ||
		tx.DocumentMatched == nil ||
		tx.DocumentMatched.MatchedBy != "yuki" ||
		tx.Contact == nil ||
		tx.Contact.FullName != "Apple Sales International" {
		t.Fatalf("transaction = %#v", tx)
	}
}

func TestCustomPaymentMethodsParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetCustomPaymentMethods", customPaymentMethodsResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationID>admin-1</they:administrationID>") {
			t.Fatalf("request body missing administration ID:\n%s", body)
		}
	})

	methods, err := client.CustomPaymentMethods(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("CustomPaymentMethods: %v", err)
	}
	if len(methods) != 2 || methods[0].ID != "5" || methods[0].Description != "Creditcard" {
		t.Fatalf("methods = %#v", methods)
	}
}

func TestCustomPaymentMethodsParsesPostmanResponseLabel(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetCustomPaymentMethods", customPaymentMethodsPostmanResponse, nil)

	methods, err := client.CustomPaymentMethods(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("CustomPaymentMethods: %v", err)
	}
	if len(methods) != 1 || methods[0].ID != "5" || methods[0].Description != "Creditcard" {
		t.Fatalf("methods = %#v", methods)
	}
}

func TestSearchDocumentsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "SearchDocuments", searchDocumentsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:searchOption>All</they:searchOption>",
			"<they:searchText>test</they:searchText>",
			"<they:folderID>2</they:folderID>",
			"<they:tabID>2</they:tabID>",
			"<they:sortOrder>CreatedAsc</they:sortOrder>",
			"<they:startDate>2021-08-01</they:startDate>",
			"<they:endDate>2021-09-30</they:endDate>",
			"<they:numberOfRecords>5</they:numberOfRecords>",
			"<they:startRecord>1</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	documents, err := client.SearchDocuments(context.Background(), "session-1", SearchDocumentsOptions{
		SearchOption:    "All",
		SearchText:      "test",
		FolderID:        2,
		TabID:           2,
		SortOrder:       "CreatedAsc",
		StartDate:       "2021-08-01",
		EndDate:         "2021-09-30",
		NumberOfRecords: 5,
		StartRecord:     1,
	})
	if err != nil {
		t.Fatalf("SearchDocuments: %v", err)
	}
	if len(documents) != 2 {
		t.Fatalf("len(documents) = %d, want 2", len(documents))
	}
	if documents[0].ID != "5517bedb-57c6-499f-b87a-c2f3a151fe59" ||
		documents[0].ContactName != "Apple Sales International" ||
		documents[1].Amount != "111.79" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestDocumentsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "Documents", documentsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:sortOrder>CreatedAsc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-05-30</they:endDate>",
			"<they:numberOfRecords>10</they:numberOfRecords>",
			"<they:startRecord>0</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	documents, err := client.Documents(context.Background(), "session-1", DocumentsOptions{
		SortOrder:       "CreatedAsc",
		StartDate:       "2020-01-01",
		EndDate:         "2020-05-30",
		NumberOfRecords: 10,
		StartRecord:     0,
	})
	if err != nil {
		t.Fatalf("Documents: %v", err)
	}
	if len(documents) != 2 {
		t.Fatalf("len(documents) = %d, want 2", len(documents))
	}
	if documents[0].ID != "4e47ac57-c219-4c1f-a582-6f3ec94015b2" ||
		documents[0].Folder.Text != "Sales" ||
		documents[1].ContactName != " Engie ;Electrabel" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestDocumentsInFolderParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentsInFolder", documentsInFolderResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:folderID>2</they:folderID>",
			"<they:sortOrder>CreatedAsc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-01-31</they:endDate>",
			"<they:numberOfRecords>10</they:numberOfRecords>",
			"<they:startRecord>0</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	documents, err := client.DocumentsInFolder(context.Background(), "session-1", DocumentsInFolderOptions{
		FolderID:        2,
		SortOrder:       "CreatedAsc",
		StartDate:       "2020-01-01",
		EndDate:         "2020-01-31",
		NumberOfRecords: 10,
		StartRecord:     0,
	})
	if err != nil {
		t.Fatalf("DocumentsInFolder: %v", err)
	}
	if len(documents) != 2 {
		t.Fatalf("len(documents) = %d, want 2", len(documents))
	}
	if documents[0].ID != "9878e187-7541-4607-8339-f94d4791d735" ||
		documents[0].FileName != "Invoice XX-12534.pdf" ||
		documents[1].ContactName != "Belgian Shell S.A." {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestDocumentsInTabParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentsInTab", documentsInTabResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:tabID>201</they:tabID>",
			"<they:sortOrder>CreatedAsc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-12-31</they:endDate>",
			"<they:numberOfRecords>100</they:numberOfRecords>",
			"<they:startRecord>1</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	documents, err := client.DocumentsInTab(context.Background(), "session-1", DocumentsInTabOptions{
		TabID:           201,
		SortOrder:       "CreatedAsc",
		StartDate:       "2020-01-01",
		EndDate:         "2020-12-31",
		NumberOfRecords: 100,
		StartRecord:     1,
	})
	if err != nil {
		t.Fatalf("DocumentsInTab: %v", err)
	}
	if len(documents) != 3 {
		t.Fatalf("len(documents) = %d, want 3", len(documents))
	}
	if documents[0].ID != "8de9f559-e123-4715-b68e-6fcae2b592c4" ||
		documents[0].FileName != "YukiSmartForm_DeclarationForm.xml" ||
		documents[1].TypeDescription != "Aankoopfactuur" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestDocumentsByTypeParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentsByType", documentsByTypeResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:documentType>2</they:documentType>",
			"<they:sortOrder>DocumentDateAsc</they:sortOrder>",
			"<they:startDate>2020-01-01</they:startDate>",
			"<they:endDate>2020-01-31</they:endDate>",
			"<they:numberOfRecords>100</they:numberOfRecords>",
			"<they:startRecord>0</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	documents, err := client.DocumentsByType(context.Background(), "session-1", DocumentsByTypeOptions{
		DocumentType:    2,
		SortOrder:       "DocumentDateAsc",
		StartDate:       "2020-01-01",
		EndDate:         "2020-01-31",
		NumberOfRecords: 100,
		StartRecord:     0,
	})
	if err != nil {
		t.Fatalf("DocumentsByType: %v", err)
	}
	if len(documents) != 2 {
		t.Fatalf("len(documents) = %d, want 2", len(documents))
	}
	if documents[0].ID != "2e22b66c-0608-4bce-8f66-0e64ed024c3e" ||
		documents[0].TypeDescription != "Purchase invoice" ||
		documents[1].Amount != "242.00" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestModifiedDocumentsInFolderParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "ModifiedDocumentsInFolder", modifiedDocumentsInFolderResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:folderID>2</they:folderID>",
			"<they:sortOrder>CreatedAsc</they:sortOrder>",
			"<they:modifiedSince>2020-08-01</they:modifiedSince>",
			"<they:numberOfRecords>100</they:numberOfRecords>",
			"<they:startRecord>1</they:startRecord>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	documents, err := client.ModifiedDocumentsInFolder(context.Background(), "session-1", ModifiedDocumentsInFolderOptions{
		FolderID:        2,
		SortOrder:       "CreatedAsc",
		ModifiedSince:   "2020-08-01",
		NumberOfRecords: 100,
		StartRecord:     1,
	})
	if err != nil {
		t.Fatalf("ModifiedDocumentsInFolder: %v", err)
	}
	if len(documents) != 2 {
		t.Fatalf("len(documents) = %d, want 2", len(documents))
	}
	if documents[0].ID != "9a52c873-cf70-4807-9e74-b9820b0144ed" ||
		documents[0].Modified != "2021-08-23T17:20:18" ||
		documents[1].ContactName != "Apple Sales International" {
		t.Fatalf("documents = %#v", documents)
	}
}

const transactionsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetTransactionsResponse xmlns="http://www.theyukicompany.com/">
      <GetTransactionsResult>
        <Transaction>
          <id>tx-1</id>
          <hID>12345</hID>
          <transactionDate>2021-10-25T00:00:00</transactionDate>
          <description>Mastercard settlement</description>
          <amount>-42.50</amount>
          <glAccountCode>550002</glAccountCode>
          <contact id="contact-1">
            <HID>1144</HID>
            <fullName>Apple Sales International</fullName>
            <country>IE</country>
            <VATNumber>IE9700053D</VATNumber>
          </contact>
          <document id="doc-1">
            <HID>98765</HID>
            <reference>MC-2021-10-25</reference>
            <type>6</type>
            <typeDescription>Aankoopfactuur</typeDescription>
            <folderId>2</folderId>
            <folder>Inkoop</folder>
            <folderTabId>201</folderTabId>
            <folderTab>Facturen</folderTab>
            <created>2021-10-26T10:00:00</created>
            <modified>2021-10-26T10:05:00</modified>
            <uploadMethod>API</uploadMethod>
          </document>
          <documentProcessed>
            <processedDate>2021-10-26T10:06:00</processedDate>
            <processedBy>yuki</processedBy>
          </documentProcessed>
          <documentMatched>
            <matchDate>2021-10-26T10:07:00</matchDate>
            <matchedBy>yuki</matchedBy>
          </documentMatched>
          <foreignCurrency>
            <amountFC>-42.50</amountFC>
            <rate>1.000000</rate>
            <currency>EUR</currency>
          </foreignCurrency>
          <vat>
            <codeType>5</codeType>
            <codeDescription>21%</codeDescription>
            <codePercentage>21.00</codePercentage>
          </vat>
          <project>
            <code>OPS</code>
            <description>Operations</description>
          </project>
        </Transaction>
      </GetTransactionsResult>
    </GetTransactionsResponse>
  </soap:Body>
</soap:Envelope>`

const customPaymentMethodsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCustomPaymentMethodsResponse xmlns="http://www.theyukicompany.com/">
      <GetCustomPaymentMethodsResult>
        <PaymentMethod>
          <ID>5</ID>
          <Description>Creditcard</Description>
        </PaymentMethod>
        <PaymentMethod>
          <ID>10</ID>
          <Description>Privé betaald</Description>
        </PaymentMethod>
      </GetCustomPaymentMethodsResult>
    </GetCustomPaymentMethodsResponse>
  </soap:Body>
</soap:Envelope>`

const customPaymentMethodsPostmanResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetPaymentMethodsResponse xmlns="http://www.theyukicompany.com/">
      <GetPaymentMethodsResult>
        <PaymentMethod>
          <ID>5</ID>
          <Description>Creditcard</Description>
        </PaymentMethod>
      </GetPaymentMethodsResult>
    </GetPaymentMethodsResponse>
  </soap:Body>
</soap:Envelope>`

const documentsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentsResponse xmlns="http://www.theyukicompany.com/">
      <DocumentsResult>
        <Documents xmlns="">
          <Document ID="4e47ac57-c219-4c1f-a582-6f3ec94015b2">
            <Subject>Factuur voor  NMBS</Subject>
            <DocumentDate>2020-01-02</DocumentDate>
            <Amount>121.00</Amount>
            <Folder ID="2">Sales</Folder>
            <Tab ID="201">Invoices</Tab>
            <Type>31</Type>
            <TypeDescription>Sales invoice</TypeDescription>
            <FileName>Document 2 (2).pdf</FileName>
            <ContentType>application/pdf</ContentType>
            <FileSize>206896</FileSize>
            <ContactName> NMBS</ContactName>
            <Created>2020-01-02T11:12:22</Created>
            <Creator>yuki</Creator>
            <Modified>2020-01-02T11:13:55</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="83e15108-6243-41ca-b665-8490053211e1">
            <Subject>Invoice for  Engie Electrabel</Subject>
            <Amount>0.00</Amount>
            <Folder ID="2">Sales</Folder>
            <Tab ID="201">Invoices</Tab>
            <Type>6</Type>
            <TypeDescription>Sales invoice</TypeDescription>
            <ContactName> Engie ;Electrabel</ContactName>
            <Created>2020-01-02T11:53:27</Created>
            <Creator>yuki</Creator>
            <Modified>2020-01-02T11:53:27</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </DocumentsResult>
    </DocumentsResponse>
  </soap:Body>
</soap:Envelope>`

const documentsInFolderResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentsInFolderResponse xmlns="http://www.theyukicompany.com/">
      <DocumentsInFolderResult>
        <Documents xmlns="">
          <Document ID="9878e187-7541-4607-8339-f94d4791d735">
            <Subject>Testfactuur - 1</Subject>
            <DocumentDate>2020-01-31</DocumentDate>
            <Amount>29.76</Amount>
            <Type>6</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <FileName>Invoice XX-12534.pdf</FileName>
            <ContentType>application/pdf</ContentType>
            <FileSize>90576</FileSize>
            <ContactName>Apple Sales International</ContactName>
            <Created>2020-01-31T16:13:28</Created>
            <Creator>yuki</Creator>
            <Modified>2020-01-31T16:45:37</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="840af189-fe57-4c20-bb79-c9ba65ca5460">
            <Subject>Factuur voor Belgian Shell S.A. (Voor alle verrichtingen in België)</Subject>
            <DocumentDate>2019-01-31</DocumentDate>
            <Amount>300.00</Amount>
            <Type>31</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <FileName>document.htm</FileName>
            <ContentType>text/html</ContentType>
            <FileSize>0</FileSize>
            <ContactName>Belgian Shell S.A.</ContactName>
            <Created>2020-01-31T09:42:47</Created>
            <Creator>yuki</Creator>
            <Modified>2020-01-31T09:43:25</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </DocumentsInFolderResult>
    </DocumentsInFolderResponse>
  </soap:Body>
</soap:Envelope>`

const documentsInTabResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentsInTabResponse xmlns="http://www.theyukicompany.com/">
      <DocumentsInTabResult>
        <Documents xmlns="">
          <Document ID="8de9f559-e123-4715-b68e-6fcae2b592c4">
            <DocumentDate>2020-11-12</DocumentDate>
            <Amount>242.00</Amount>
            <Type>0</Type>
            <TypeDescription>Standaard</TypeDescription>
            <FileName>YukiSmartForm_DeclarationForm.xml</FileName>
            <ContentType>text/xml</ContentType>
            <FileSize>665</FileSize>
            <ContactName>katrien 2</ContactName>
            <Created>2020-11-12T15:38:11</Created>
            <Creator>yuki</Creator>
            <Modified>2021-03-10T09:41:51</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="669226dc-2cbc-4eff-916b-58ea790d82a7">
            <Subject>Factuur van  Sodexo Pass Belgium, Aankopen van handelsgoederen, borgtochten</Subject>
            <DocumentDate>2020-02-12</DocumentDate>
            <Amount>252.81</Amount>
            <Type>2</Type>
            <TypeDescription>Aankoopfactuur</TypeDescription>
            <FileName>document.htm</FileName>
            <ContentType>text/html</ContentType>
            <FileSize>0</FileSize>
            <ContactName> Sodexo Pass Belgium</ContactName>
            <Created>2020-02-12T11:14:04</Created>
            <Creator>yuki</Creator>
            <Modified>2020-02-12T11:15:51</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="7c88ab36-5637-4609-934c-74ed1d88a779">
            <DocumentDate>2019-10-23</DocumentDate>
            <Amount>0.00</Amount>
            <Type>0</Type>
            <TypeDescription>Standaard</TypeDescription>
            <FileName>Invoice 20193.pdf</FileName>
            <ContentType>application/pdf</ContentType>
            <FileSize>63361</FileSize>
            <Created>2019-10-23T14:33:37</Created>
            <Creator>yuki</Creator>
            <Modified>2019-10-23T14:33:37</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </DocumentsInTabResult>
    </DocumentsInTabResponse>
  </soap:Body>
</soap:Envelope>`

const documentsByTypeResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentsByTypeResponse xmlns="http://www.theyukicompany.com/">
      <DocumentsByTypeResult>
        <Documents xmlns="">
          <Document ID="2e22b66c-0608-4bce-8f66-0e64ed024c3e">
            <Subject>Factuur van AD Delhaize, Gebouwen, Rekening onbekend</Subject>
            <DocumentDate>2020-01-03</DocumentDate>
            <Amount>354520.32</Amount>
            <Folder ID="1">Purchase</Folder>
            <Tab ID="101">Invoices</Tab>
            <Type>2</Type>
            <TypeDescription>Purchase invoice</TypeDescription>
            <FileName>document.htm</FileName>
            <ContentType>text/html</ContentType>
            <FileSize>0</FileSize>
            <ContactName>AD Delhaize</ContactName>
            <Created>2020-01-03T15:56:52</Created>
            <Creator>yuki</Creator>
            <Modified>2020-01-03T16:04:56</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="4439aea3-9705-4119-acc0-3ed65342ee21">
            <Subject>Factuur van AD Delhaize, Gebouwen</Subject>
            <DocumentDate>2020-01-03</DocumentDate>
            <Amount>242.00</Amount>
            <Folder ID="1">Purchase</Folder>
            <Tab ID="101">Invoices</Tab>
            <Type>2</Type>
            <TypeDescription>Purchase invoice</TypeDescription>
            <FileName>document.htm</FileName>
            <ContentType>text/html</ContentType>
            <FileSize>0</FileSize>
            <ContactName>AD Delhaize</ContactName>
            <Created>2020-01-03T15:40:00</Created>
            <Creator>yuki</Creator>
            <Modified>2020-01-03T15:40:41</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </DocumentsByTypeResult>
    </DocumentsByTypeResponse>
  </soap:Body>
</soap:Envelope>`

const modifiedDocumentsInFolderResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ModifiedDocumentsInFolderResponse xmlns="http://www.theyukicompany.com/">
      <ModifiedDocumentsInFolderResult>
        <Documents xmlns="">
          <Document ID="9a52c873-cf70-4807-9e74-b9820b0144ed">
            <Subject>Testfactuur - 1</Subject>
            <DocumentDate>2021-07-05</DocumentDate>
            <Amount>111.79</Amount>
            <Type>6</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <FileName>Invoice B1240.pdf</FileName>
            <ContentType>application/pdf</ContentType>
            <FileSize>95360</FileSize>
            <ContactName>ES klant</ContactName>
            <Created>2021-07-05T10:15:33</Created>
            <Creator>yuki</Creator>
            <Modified>2021-08-23T17:20:18</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="5517bedb-57c6-499f-b87a-c2f3a151fe59">
            <Subject>Testfactuur - 1</Subject>
            <Amount>29.76</Amount>
            <Type>6</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <ContactName>Apple Sales International</ContactName>
            <Created>2021-08-01T04:10:22</Created>
            <Creator>yuki</Creator>
            <Modified>2021-08-01T04:10:23</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </ModifiedDocumentsInFolderResult>
    </ModifiedDocumentsInFolderResponse>
  </soap:Body>
</soap:Envelope>`

const searchDocumentsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <SearchDocumentsResponse xmlns="http://www.theyukicompany.com/">
      <SearchDocumentsResult>
        <Documents xmlns="">
          <Document ID="5517bedb-57c6-499f-b87a-c2f3a151fe59">
            <Subject>Testfactuur - 1</Subject>
            <Amount>29.76</Amount>
            <Type>6</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <ContactName>Apple Sales International</ContactName>
            <Created>2021-08-01T04:10:22</Created>
            <Creator>yuki</Creator>
            <Modified>2021-08-01T04:10:23</Modified>
            <Modifier>yuki</Modifier>
          </Document>
          <Document ID="8102afb2-5b1c-4917-bcbf-40f8875e684c">
            <Subject>Testfactuur - 1</Subject>
            <Amount>111.79</Amount>
            <Type>6</Type>
            <TypeDescription>Verkoopfactuur</TypeDescription>
            <ContactName>Molly Malone</ContactName>
            <Created>2021-08-18T16:13:50</Created>
            <Creator>yuki</Creator>
            <Modified>2021-08-18T16:13:50</Modified>
            <Modifier>yuki</Modifier>
          </Document>
        </Documents>
      </SearchDocumentsResult>
    </SearchDocumentsResponse>
  </soap:Body>
</soap:Envelope>`
