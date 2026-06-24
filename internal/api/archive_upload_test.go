package api

import (
	"context"
	"strings"
	"testing"
)

func TestUploadDocumentPostsWSDLFieldsAndParsesResult(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "UploadDocument", uploadDocumentResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:sessionID>session-1</they:sessionID>",
			"<they:fileName>invoice.pdf</they:fileName>",
			"<they:data>JVBERg==</they:data>",
			"<they:folder>1</they:folder>",
			"<they:administrationID>admin-1</they:administrationID>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.UploadDocument(context.Background(), "session-1", ArchiveUploadOptions{
		FileName:         "invoice.pdf",
		DataBase64:       "JVBERg==",
		FolderID:         1,
		AdministrationID: "admin-1",
	})
	if err != nil {
		t.Fatalf("UploadDocument: %v", err)
	}
	if result.DocumentID != "doc-basic" || result.Operation != "UploadDocument" {
		t.Fatalf("result = %#v", result)
	}
}

func TestUploadDocumentWithDataPostsPostmanAndWSDLFields(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "UploadDocumentWithData", uploadDocumentWithDataResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:fileName>receipt.jpg</they:fileName>",
			"<they:data>anBn</they:data>",
			"<they:folder>7</they:folder>",
			"<they:currency>EUR</they:currency>",
			"<they:amount>42.50</they:amount>",
			"<they:costCategory>meals</they:costCategory>",
			"<they:paymentMethod>4</they:paymentMethod>",
			"<they:project>OPS</they:project>",
			"<they:remarks>card receipt</they:remarks>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.UploadDocumentWithData(context.Background(), "session-1", ArchiveUploadOptions{
		FileName:         "receipt.jpg",
		DataBase64:       "anBn",
		FolderID:         7,
		AdministrationID: "admin-1",
		Currency:         "EUR",
		Amount:           "42.50",
		CostCategory:     "meals",
		PaymentMethod:    4,
		Project:          "OPS",
		Remarks:          "card receipt",
	})
	if err != nil {
		t.Fatalf("UploadDocumentWithData: %v", err)
	}
	if result.DocumentID != "doc-data" || result.Amount != "42.50" || result.PaymentMethod != 4 {
		t.Fatalf("result = %#v", result)
	}
}

func TestUploadDocumentWithAttachmentPostsPostmanFields(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "UploadDocumentWithAttachment", uploadDocumentWithAttachmentResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:fileName1>soda.xml</they:fileName1>",
			"<they:data1>PHhtbD4=</they:data1>",
			"<they:fileName2>soda.pdf</they:fileName2>",
			"<they:data2>JVBERg==</they:data2>",
			"<they:folder>1</they:folder>",
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:currency>EUR</they:currency>",
			"<they:amount>0</they:amount>",
			"<they:costCategory>office</they:costCategory>",
			"<they:paymentMethod>0</they:paymentMethod>",
			"<they:project>OPS</they:project>",
			"<they:remarks>merged evidence</they:remarks>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.UploadDocumentWithAttachment(context.Background(), "session-1", ArchiveAttachmentUploadOptions{
		FileName1:        "soda.xml",
		Data1Base64:      "PHhtbD4=",
		FileName2:        "soda.pdf",
		Data2Base64:      "JVBERg==",
		FolderID:         1,
		AdministrationID: "admin-1",
		Currency:         "EUR",
		Amount:           "0",
		CostCategory:     "office",
		PaymentMethod:    0,
		Project:          "OPS",
		Remarks:          "merged evidence",
	})
	if err != nil {
		t.Fatalf("UploadDocumentWithAttachment: %v", err)
	}
	if result.DocumentID != "doc-attachment" || result.AttachmentName != "soda.pdf" {
		t.Fatalf("result = %#v", result)
	}
}

const uploadDocumentResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UploadDocumentResponse xmlns="http://www.theyukicompany.com/">
      <UploadDocumentResult>doc-basic</UploadDocumentResult>
    </UploadDocumentResponse>
  </soap:Body>
</soap:Envelope>`

const uploadDocumentWithDataResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UploadDocumentWithDataResponse xmlns="http://www.theyukicompany.com/">
      <UploadDocumentWithDataResult>doc-data</UploadDocumentWithDataResult>
    </UploadDocumentWithDataResponse>
  </soap:Body>
</soap:Envelope>`

const uploadDocumentWithAttachmentResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UploadDocumentWithAttachmentResponse xmlns="http://www.theyukicompany.com/">
      <UploadDocumentWithAttachmentResult>doc-attachment</UploadDocumentWithAttachmentResult>
    </UploadDocumentWithAttachmentResponse>
  </soap:Body>
</soap:Envelope>`
