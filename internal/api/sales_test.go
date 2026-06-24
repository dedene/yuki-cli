package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSalesInvoiceSchemaPathParsesDocumentedResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/Sales.asmx" {
			t.Fatalf("path = %s, want /Sales.asmx", r.URL.Path)
		}
		if got := r.Header.Get("SOAPAction"); got != SOAPAction("SalesInvoiceSchemaPath") {
			t.Fatalf("SOAPAction = %q", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), "<they:SalesInvoiceSchemaPath>") {
			t.Fatalf("request body missing operation:\n%s", body)
		}
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		_, _ = w.Write([]byte(salesInvoiceSchemaPathResponse))
	}))
	defer srv.Close()

	client := New(Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	path, err := client.SalesInvoiceSchemaPath(context.Background())
	if err != nil {
		t.Fatalf("SalesInvoiceSchemaPath: %v", err)
	}
	if path != "https://api.yukiworks.be/schemas/SalesInvoices.xsd" {
		t.Fatalf("path = %q", path)
	}
}

func TestSalesItemsParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "Sales", "GetSalesItems", salesItemsResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationID>admin-1</they:administrationID>") {
			t.Fatalf("request body missing administration ID:\n%s", body)
		}
	})

	items, err := client.SalesItems(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("SalesItems: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if items[0].ID != "item-1" || items[0].Description != "Consulting" ||
		items[1].ID != "item-2" || items[1].Description != "Hosting" {
		t.Fatalf("items = %#v", items)
	}
}

func TestProcessSalesInvoicesPostsRawXMLAndParsesDocumentedResponse(t *testing.T) {
	xmlDoc := `<SalesInvoices xmlns="urn:xmlns:http://www.theyukicompany.com:salesinvoices"><SalesInvoice><Reference>VF-0001</Reference></SalesInvoice></SalesInvoices>`
	client := fixtureClientForServiceWithSessionElement(t, "Sales", "ProcessSalesInvoices", processSalesInvoicesResponse, "sessionId", func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationId>admin-1</they:administrationId>") {
			t.Fatalf("request body missing administration ID:\n%s", body)
		}
		if !strings.Contains(body, "<they:xmlDoc>"+xmlDoc+"</they:xmlDoc>") {
			t.Fatalf("request body missing raw xmlDoc:\n%s", body)
		}
		if strings.Contains(body, "&lt;SalesInvoices") {
			t.Fatalf("request body escaped xmlDoc:\n%s", body)
		}
	})

	result, err := client.ProcessSalesInvoices(context.Background(), "session-1", SalesInvoiceImportOptions{
		AdministrationID: "admin-1",
		XMLDoc:           xmlDoc,
	})
	if err != nil {
		t.Fatalf("ProcessSalesInvoices: %v", err)
	}
	if result.AdministrationID != "e5151219-a7bd-44f7-b9b2-a87d379e71f6" ||
		result.TotalSucceeded != 1 ||
		result.TotalFailed != 0 ||
		len(result.Invoices) != 1 ||
		!result.Invoices[0].Succeeded ||
		result.Invoices[0].Reference != "VF-0001" {
		t.Fatalf("result = %#v", result)
	}
}

func TestProcessRecognizedSalesInvoicesPostsWSDLOperation(t *testing.T) {
	xmlDoc := `<SalesInvoices xmlns="urn:xmlns:http://www.theyukicompany.com:salesinvoices"/>`
	client := fixtureClientForServiceWithSessionElement(t, "Sales", "ProcessRecognizedSalesInvoices", processRecognizedSalesInvoicesResponse, "sessionId", func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationId>admin-1</they:administrationId>") ||
			!strings.Contains(body, "<they:xmlDoc>"+xmlDoc+"</they:xmlDoc>") {
			t.Fatalf("request body missing recognized invoice params:\n%s", body)
		}
	})

	result, err := client.ProcessRecognizedSalesInvoices(context.Background(), "session-1", SalesInvoiceImportOptions{
		AdministrationID: "admin-1",
		XMLDoc:           xmlDoc,
	})
	if err != nil {
		t.Fatalf("ProcessRecognizedSalesInvoices: %v", err)
	}
	if result.TotalSkipped != 1 || len(result.Invoices) != 1 || result.Invoices[0].Message != "Already recognized" {
		t.Fatalf("result = %#v", result)
	}
}

const salesInvoiceSchemaPathResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <SalesInvoiceSchemaPathResponse xmlns="http://www.theyukicompany.com/">
      <SalesInvoiceSchemaPathResult>https://api.yukiworks.be/schemas/SalesInvoices.xsd</SalesInvoiceSchemaPathResult>
    </SalesInvoiceSchemaPathResponse>
  </soap:Body>
</soap:Envelope>`

const salesItemsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetSalesItemsResponse xmlns="http://www.theyukicompany.com/">
      <GetSalesItemsResult>
        <SalesItem>
          <id>item-1</id>
          <description>Consulting</description>
        </SalesItem>
        <SalesItem>
          <id>item-2</id>
          <description>Hosting</description>
        </SalesItem>
      </GetSalesItemsResult>
    </GetSalesItemsResponse>
  </soap:Body>
</soap:Envelope>`

const processSalesInvoicesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ProcessSalesInvoicesResponse xmlns="http://www.theyukicompany.com/">
      <ProcessSalesInvoicesResult>
        <SalesInvoicesImportResponse xmlns="urn:xmlns:http://www.theyukicompany.com:salesinvoicesresponse">
          <TimeStamp xmlns="">2025-10-23</TimeStamp>
          <AdministrationId xmlns="">e5151219-a7bd-44f7-b9b2-a87d379e71f6</AdministrationId>
          <TotalSucceeded xmlns="">1</TotalSucceeded>
          <TotalFailed xmlns="">0</TotalFailed>
          <TotalSkipped xmlns="">0</TotalSkipped>
          <Invoice xmlns="">
            <Succeeded>true</Succeeded>
            <Processed>true</Processed>
            <EmailSent>false</EmailSent>
            <Reference>VF-0001</Reference>
            <Subject>Testfactuur - 1</Subject>
            <Contact>blabla 007</Contact>
            <Message></Message>
          </Invoice>
        </SalesInvoicesImportResponse>
      </ProcessSalesInvoicesResult>
    </ProcessSalesInvoicesResponse>
  </soap:Body>
</soap:Envelope>`

const processRecognizedSalesInvoicesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ProcessRecognizedSalesInvoicesResponse xmlns="http://www.theyukicompany.com/">
      <ProcessRecognizedSalesInvoicesResult>
        <SalesInvoicesImportResponse xmlns="urn:xmlns:http://www.theyukicompany.com:salesinvoicesresponse">
          <TimeStamp xmlns="">2025-10-24</TimeStamp>
          <AdministrationId xmlns="">admin-1</AdministrationId>
          <TotalSucceeded xmlns="">0</TotalSucceeded>
          <TotalFailed xmlns="">0</TotalFailed>
          <TotalSkipped xmlns="">1</TotalSkipped>
          <Invoice xmlns="">
            <Succeeded>false</Succeeded>
            <Processed>false</Processed>
            <EmailSent>false</EmailSent>
            <Reference>VF-0002</Reference>
            <Subject>Recognized invoice</Subject>
            <Contact>Existing customer</Contact>
            <Message>Already recognized</Message>
          </Invoice>
        </SalesInvoicesImportResponse>
      </ProcessRecognizedSalesInvoicesResult>
    </ProcessRecognizedSalesInvoicesResponse>
  </soap:Body>
</soap:Envelope>`
