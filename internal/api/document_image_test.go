package api

import (
	"context"
	"strings"
	"testing"
)

func TestDocumentImageParsesWSDLBase64Response(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "DocumentImage", documentImageResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:documentID>doc-1</they:documentID>",
			"<they:maxWidth>800</they:maxWidth>",
			"<they:maxHeight>1200</they:maxHeight>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	data, err := client.DocumentImage(context.Background(), "session-1", "doc-1", 800, 1200)
	if err != nil {
		t.Fatalf("DocumentImage: %v", err)
	}
	if data.DocumentID != "doc-1" ||
		data.MaxWidth != 800 ||
		data.MaxHeight != 1200 ||
		data.ImageDataBase64 != "aW1hZ2UtYnl0ZXM=" {
		t.Fatalf("data = %#v", data)
	}
}

const documentImageResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DocumentImageResponse xmlns="http://www.theyukicompany.com/">
      <DocumentImageResult>aW1hZ2UtYnl0ZXM=</DocumentImageResult>
    </DocumentImageResponse>
  </soap:Body>
</soap:Envelope>`
