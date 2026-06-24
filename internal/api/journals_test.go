package api

import (
	"context"
	"strings"
	"testing"
)

func TestProcessJournalPostsRawXMLAndParsesDocumentedResponse(t *testing.T) {
	xmlDoc := `<Journal xmlns="urn:xmlns:http://www.theyukicompany.com:journal"><AdministrationID>admin-1</AdministrationID><DocumentSubject>Payroll</DocumentSubject></Journal>`
	client := fixtureClientForService(t, "Accounting", "ProcessJournal", processJournalResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationID>admin-1</they:administrationID>") {
			t.Fatalf("request body missing administration ID:\n%s", body)
		}
		if !strings.Contains(body, "<they:xmlDoc>"+xmlDoc+"</they:xmlDoc>") {
			t.Fatalf("request body missing raw xmlDoc:\n%s", body)
		}
		if strings.Contains(body, "&lt;Journal") {
			t.Fatalf("request body escaped xmlDoc:\n%s", body)
		}
	})

	result, err := client.ProcessJournal(context.Background(), "session-1", JournalImportOptions{
		AdministrationID: "admin-1",
		XMLDoc:           xmlDoc,
	})
	if err != nil {
		t.Fatalf("ProcessJournal: %v", err)
	}
	if result.AdministrationID != "admin-1" ||
		result.DocumentID != "9999e179-04bd-43a7-8fb7-403999a88fac" {
		t.Fatalf("result = %#v", result)
	}
}

const processJournalResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ProcessJournalResponse xmlns="http://www.theyukicompany.com/">
      <ProcessJournalResult>9999e179-04bd-43a7-8fb7-403999a88fac</ProcessJournalResult>
    </ProcessJournalResponse>
  </soap:Body>
</soap:Envelope>`
