package api

import (
	"context"
	"strings"
	"testing"
)

func TestDocumentXMLDataParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentXMLData", documentXMLDataResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	data, err := client.DocumentXMLData(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentXMLData: %v", err)
	}
	if data.DocumentID != "doc-1" ||
		!strings.Contains(data.XML, "<SalesInvoice") ||
		!strings.Contains(data.XML, "<Reference>A1040</Reference>") ||
		!strings.Contains(data.XML, "<Subject>Testfactuur - 1</Subject>") {
		t.Fatalf("data = %#v", data)
	}
}

func TestDocumentXMLDataAsBinaryParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentXMLDataAsBinary", documentXMLDataAsBinaryResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	data, err := client.DocumentXMLDataAsBinary(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentXMLDataAsBinary: %v", err)
	}
	if data.DocumentID != "doc-1" || data.XMLDataBase64 != "PFNhbGVzSW52b2ljZT48UmVmZXJlbmNlPkExMDQwPC9SZWZlcmVuY2U+PC9TYWxlc0ludm9pY2U+" {
		t.Fatalf("data = %#v", data)
	}
}

const documentXMLDataResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentXMLDataResponse xmlns="http://www.theyukicompany.com/">
      <DocumentXMLDataResult>
        <SalesInvoice xmlns="">
          <Reference>A1040</Reference>
          <Subject>Testfactuur - 1</Subject>
          <PaymentMethod>ElectronicTransfer</PaymentMethod>
        </SalesInvoice>
      </DocumentXMLDataResult>
    </DocumentXMLDataResponse>
  </soap:Body>
</soap:Envelope>`

const documentXMLDataAsBinaryResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentXMLDataAsBinaryResponse xmlns="http://www.theyukicompany.com/">
      <DocumentXMLDataAsBinaryResult>PFNhbGVzSW52b2ljZT48UmVmZXJlbmNlPkExMDQwPC9SZWZlcmVuY2U+PC9TYWxlc0ludm9pY2U+</DocumentXMLDataAsBinaryResult>
    </DocumentXMLDataAsBinaryResponse>
  </soap:Body>
</soap:Envelope>`
