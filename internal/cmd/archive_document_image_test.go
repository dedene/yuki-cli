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

func TestArchiveDocumentsImageJSONPrintsBase64ImageData(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentImageData: api.DocumentImageData{
			DocumentID:      "doc-1",
			MaxWidth:        800,
			MaxHeight:       1200,
			ImageDataBase64: "aW1hZ2UtYnl0ZXM=",
		},
	}

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "image",
		"--document", "doc-1",
		"--max-width", "800",
		"--max-height", "1200",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.documentID != "doc-1" || client.maxWidth != 800 || client.maxHeight != 1200 {
		t.Fatalf("image call = document %q maxWidth %d maxHeight %d", client.documentID, client.maxWidth, client.maxHeight)
	}
	var data api.DocumentImageData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if data.DocumentID != "doc-1" || data.ImageDataBase64 != "aW1hZ2UtYnl0ZXM=" {
		t.Fatalf("data = %#v", data)
	}
}

func TestArchiveDocumentsImageWritesDecodedImageFile(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		documentImageData: api.DocumentImageData{
			DocumentID:      "doc-1",
			MaxWidth:        800,
			MaxHeight:       1200,
			ImageDataBase64: "aW1hZ2UtYnl0ZXM=",
		},
	}
	outputPath := filepath.Join(t.TempDir(), "document-image.bin")

	err := Execute(context.Background(), []string{
		"archive", "documents", "image",
		"--document", "doc-1",
		"--max-width", "800",
		"--max-height", "1200",
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
	if string(data) != "image-bytes" {
		t.Fatalf("decoded file = %q", data)
	}
	if !strings.Contains(out.String(), "Wrote "+outputPath) {
		t.Fatalf("output missing write confirmation:\n%s", out.String())
	}
}

func TestArchiveDocumentsImageRejectsNonPositiveDimensions(t *testing.T) {
	var out bytes.Buffer
	err := Execute(context.Background(), []string{
		"archive", "documents", "image",
		"--document", "doc-1",
		"--max-width", "0",
		"--max-height", "1200",
		"--output", "doc.png",
	}, Runtime{Out: &out})
	if err == nil {
		t.Fatal("Execute returned nil, want dimension error")
	}
}
