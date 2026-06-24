package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/auth"
)

func TestArchiveDocumentsUploadBasicDryRunSkipsAuth(t *testing.T) {
	path := writeFile(t, "invoice.pdf", []byte("%PDF"))
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"--default-administration", "admin-default",
		"archive", "documents", "upload-basic",
		"--file", path,
		"--folder", "1",
		"--dry-run",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{err: auth.ErrAccessKeyNotFound},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.accessKey != "" || client.archiveUploadOperation != "" {
		t.Fatalf("dry-run authenticated or sent: accessKey=%q operation=%q", client.accessKey, client.archiveUploadOperation)
	}
	var result api.ArchiveUploadResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json: %v", err)
	}
	if result.Operation != "UploadDocument" || result.FileName != "invoice.pdf" || result.FolderID != 1 || !result.DryRun {
		t.Fatalf("result = %#v", result)
	}
}

func TestArchiveDocumentsUploadDataReadonlyBlocksBeforeAuth(t *testing.T) {
	path := writeFile(t, "invoice.pdf", []byte("%PDF"))
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--readonly",
		"archive", "documents", "upload-data",
		"--file", path,
		"--folder", "7",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "--readonly blocks mutating command: archive documents upload-data") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" || client.archiveUploadOperation != "" {
		t.Fatalf("readonly authenticated or sent: accessKey=%q operation=%q", client.accessKey, client.archiveUploadOperation)
	}
}

func TestArchiveDocumentsUploadDataSendsBase64AndFlags(t *testing.T) {
	path := writeFile(t, "local.pdf", []byte("%PDF"))
	client := &cmdFakeClient{
		sessionID: "session-1",
		archiveUploadResult: api.ArchiveUploadResult{
			Operation:  "UploadDocumentWithData",
			DocumentID: "doc-1",
		},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "documents", "upload-data",
		"--file", path,
		"--file-name", "invoice.pdf",
		"--folder", "7",
		"--administration", "admin-1",
		"--currency", "EUR",
		"--amount", "42.50",
		"--cost-category", "meals",
		"--payment-method", "4",
		"--project", "OPS",
		"--remarks", "card receipt",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	opts := client.archiveUploadOpts
	if client.archiveUploadOperation != "UploadDocumentWithData" ||
		opts.FileName != "invoice.pdf" ||
		opts.DataBase64 != base64.StdEncoding.EncodeToString([]byte("%PDF")) ||
		opts.FolderID != 7 ||
		opts.AdministrationID != "admin-1" ||
		opts.Currency != "EUR" ||
		opts.Amount != "42.50" ||
		opts.CostCategory != "meals" ||
		opts.PaymentMethod != 4 ||
		opts.Project != "OPS" ||
		opts.Remarks != "card receipt" {
		t.Fatalf("upload opts = %#v operation=%q", opts, client.archiveUploadOperation)
	}
	var result api.ArchiveUploadResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("json: %v", err)
	}
	if result.DocumentID != "doc-1" {
		t.Fatalf("result = %#v", result)
	}
}

func TestArchiveDocumentsUploadAttachmentSendsBothFiles(t *testing.T) {
	file := writeFile(t, "soda.xml", []byte("<xml/>"))
	attachment := writeFile(t, "soda.pdf", []byte("%PDF"))
	client := &cmdFakeClient{
		sessionID: "session-1",
		archiveUploadResult: api.ArchiveUploadResult{
			Operation:  "UploadDocumentWithAttachment",
			DocumentID: "doc-attachment",
		},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"archive", "documents", "upload-attachment",
		"--file", file,
		"--attachment", attachment,
		"--folder", "1",
		"--administration", "admin-1",
		"--amount", "0",
		"--payment-method", "0",
		"--remarks", "linked files",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	opts := client.archiveAttachmentUploadOpts
	if client.archiveUploadOperation != "UploadDocumentWithAttachment" ||
		opts.FileName1 != "soda.xml" ||
		opts.Data1Base64 != base64.StdEncoding.EncodeToString([]byte("<xml/>")) ||
		opts.FileName2 != "soda.pdf" ||
		opts.Data2Base64 != base64.StdEncoding.EncodeToString([]byte("%PDF")) ||
		opts.FolderID != 1 ||
		opts.AdministrationID != "admin-1" ||
		opts.Remarks != "linked files" {
		t.Fatalf("attachment opts = %#v operation=%q", opts, client.archiveUploadOperation)
	}
}

func TestArchiveDocumentsUploadRejectsEmptyFileBeforeAuth(t *testing.T) {
	path := writeFile(t, "empty.pdf", nil)
	client := &cmdFakeClient{sessionID: "session-1"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"archive", "documents", "upload-basic",
		"--file", path,
		"--folder", "1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "upload file is empty") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("authenticated before rejecting empty file: %q", client.accessKey)
	}
}
