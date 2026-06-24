package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/auth"
)

func TestProjectsUpsertDryRunSkipsAuth(t *testing.T) {
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json", "accounting", "projects", "upsert",
		"--administration", "admin-1",
		"--description", "New Project",
		"--code", "PROJECTNEW",
		"--allow-ocr-matching", "true",
		"--dry-run",
	}, Runtime{
		Out:   &out,
		Store: &cmdFakeStore{err: auth.ErrAccessKeyNotFound},
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	var result api.ProjectUpdateResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if !result.DryRun ||
		result.AdministrationID != "admin-1" ||
		result.Project.Description != "New Project" ||
		result.Project.AllowOCRMatching != "true" {
		t.Fatalf("result = %#v", result)
	}
}

func TestProjectsUpsertReadonlyBlocksBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"--readonly", "accounting", "projects", "upsert",
		"--administration", "admin-1",
		"--description", "New Project",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "--readonly blocks mutating command") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}

func TestProjectsUpsertSendsDocumentedFields(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{
		sessionID: "session-1",
		projectUpdateResult: api.ProjectUpdateResult{
			Message: "project upserted",
		},
	}

	err := Execute(context.Background(), []string{
		"--json", "accounting", "projects", "upsert",
		"--administration", "admin-1",
		"--description", "New Project",
		"--code", "PROJECTNEW",
		"--company", "admin-1",
		"--manager", "manager@example.com",
		"--contact", "contact-1",
		"--notes", "this is a new project",
		"--security-level", "1",
		"--allow-ocr-matching", "yes",
		"--start-date", "2020-01-20",
		"--end-date", "2022-12-31",
		"--budget-revenue", "3000",
		"--budget-costs", "1000",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.projectUpdateOpts.AdministrationID != "admin-1" ||
		client.projectUpdateOpts.Project.Description != "New Project" ||
		client.projectUpdateOpts.Project.SecurityLevel != "1" ||
		client.projectUpdateOpts.Project.AllowOCRMatching != "true" {
		t.Fatalf("opts = %#v", client.projectUpdateOpts)
	}
	var result api.ProjectUpdateResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	if result.AdministrationID != "admin-1" || result.Message != "project upserted" {
		t.Fatalf("result = %#v", result)
	}
}

func TestProjectsUpsertRejectsInvalidOCRMatchingBeforeAuth(t *testing.T) {
	var out bytes.Buffer
	client := &cmdFakeClient{sessionID: "session-1"}

	err := Execute(context.Background(), []string{
		"accounting", "projects", "upsert",
		"--administration", "admin-1",
		"--description", "New Project",
		"--allow-ocr-matching", "maybe",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err == nil || !strings.Contains(err.Error(), "invalid --allow-ocr-matching") {
		t.Fatalf("err = %v", err)
	}
	if client.accessKey != "" {
		t.Fatalf("accessKey = %q, want no authentication", client.accessKey)
	}
}
