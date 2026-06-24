package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestArchiveDocumentsDownloadURLJSONPrintsURL(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentDownloadURL: api.DocumentDownloadURL{
			DocumentID: "doc-1",
			URL:        "https://api.yukiworks.be/download/document/doc-1?token=abc",
		},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "download-url",
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
	var downloadURL api.DocumentDownloadURL
	if err := json.Unmarshal(out.Bytes(), &downloadURL); err != nil {
		t.Fatalf("json: %v", err)
	}
	if downloadURL.DocumentID != "doc-1" ||
		downloadURL.URL != "https://api.yukiworks.be/download/document/doc-1?token=abc" {
		t.Fatalf("downloadURL = %#v", downloadURL)
	}
}
