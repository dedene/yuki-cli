package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestArchiveDocumentsXMLBinaryJSONPrintsBase64XMLData(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentXMLBinaryData: api.DocumentXMLBinaryData{
			DocumentID:    "doc-1",
			XMLDataBase64: "PFNhbGVzSW52b2ljZS8+",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "xml-binary",
		"--document", "doc-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.documentID != "doc-1" {
		t.Fatalf("documentID = %q", client.documentID)
	}
	var data api.DocumentXMLBinaryData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if data.DocumentID != "doc-1" || data.XMLDataBase64 != "PFNhbGVzSW52b2ljZS8+" {
		t.Fatalf("data = %#v", data)
	}
}

func TestArchiveDocumentsXMLBinaryWritesDecodedXMLFile(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentXMLBinaryData: api.DocumentXMLBinaryData{
			DocumentID:    "doc-1",
			XMLDataBase64: "PFNhbGVzSW52b2ljZS8+",
		},
	}
	outputPath := filepath.Join(t.TempDir(), "invoice.xml")

	err := Execute(context.Background(), []string{
		"archive", "documents", "xml-binary",
		"--document", "doc-1",
		"--output", outputPath,
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if string(data) != "<SalesInvoice/>" {
		t.Fatalf("decoded file = %q", data)
	}
	if !strings.Contains(out.String(), "Wrote "+outputPath) {
		t.Fatalf("output missing write confirmation:\n%s", out.String())
	}
}
