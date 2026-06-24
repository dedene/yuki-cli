package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestArchiveDocumentsBundleJSONPrintsBundledDocuments(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentBundle: []api.Document{{
			ID:              "bundle-doc-1",
			Subject:         "Bundled receipt",
			DocumentDate:    "2026-01-03",
			Amount:          "42.00",
			TypeDescription: "Receipt",
			FileName:        "receipt.pdf",
			ContactName:     "Vendor",
		}},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "bundle",
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
	var documents []api.Document
	if err := json.Unmarshal(out.Bytes(), &documents); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if len(documents) != 1 || documents[0].ID != "bundle-doc-1" {
		t.Fatalf("documents = %#v", documents)
	}
}
