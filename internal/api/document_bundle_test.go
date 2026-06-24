package api

import (
	"context"
	"strings"
	"testing"
)

func TestDocumentBundleParsesDocumentedEmptyResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentBundle", documentBundleResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:documentID>doc-1</they:documentID>") {
			t.Fatalf("request body missing document ID:\n%s", body)
		}
	})

	documents, err := client.DocumentBundle(context.Background(), "session-1", "doc-1")
	if err != nil {
		t.Fatalf("DocumentBundle: %v", err)
	}
	if len(documents) != 0 {
		t.Fatalf("len(documents) = %d, want 0", len(documents))
	}
}

const documentBundleResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentBundleResponse xmlns="http://www.theyukicompany.com/">
      <DocumentBundleResult>
        <Documents xmlns="" />
      </DocumentBundleResult>
    </DocumentBundleResponse>
  </soap:Body>
</soap:Envelope>`
