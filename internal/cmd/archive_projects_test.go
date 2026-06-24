package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestArchiveProjectsListJSONPrintsProjects(t *testing.T) {
	client := &cmdFakeClient{
		sessionID: "session-1",
		archiveProjects: []api.AccountingProject{{
			HID:         "5",
			Code:        "ARCHIVE",
			Description: "Archive Project",
			Company:     "Highpro BV",
		}},
	}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"archive", "projects", "list",
		"--administration", "admin-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.archiveProjectsAdminID != "admin-1" {
		t.Fatalf("archiveProjectsAdminID = %q", client.archiveProjectsAdminID)
	}
	var projects []api.AccountingProject
	if err := json.Unmarshal(out.Bytes(), &projects); err != nil {
		t.Fatalf("json: %v", err)
	}
	if len(projects) != 1 || projects[0].Code != "ARCHIVE" {
		t.Fatalf("projects = %#v", projects)
	}
}
