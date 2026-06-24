package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/dedene/yuki-cli/internal/api"
)

func TestAuthSessionClientJSONUsesStoredAccessKey(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-client"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"auth", "session", "client",
		"--client-id", "client-1",
		"--client-secret", "secret-1",
	}, Runtime{
		Out:       &out,
		Store:     &cmdFakeStore{key: "stored-key"},
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.clientID != "client-1" || client.clientSecret != "secret-1" || client.accessKey != "stored-key" {
		t.Fatalf("client auth fields = clientID:%q secret:%q access:%q", client.clientID, client.clientSecret, client.accessKey)
	}
	var payload map[string]string
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v", err)
	}
	if payload["session_id"] != "session-client" {
		t.Fatalf("session_id = %q", payload["session_id"])
	}
}

func TestAuthSessionUsernameJSONUsesPasswordFlag(t *testing.T) {
	client := &cmdFakeClient{sessionID: "session-user"}
	var out bytes.Buffer

	err := Execute(context.Background(), []string{
		"--json",
		"auth", "session", "username",
		"--username", "peter@example.com",
		"--password", "secret-1",
	}, Runtime{
		Out:       &out,
		NewClient: func(api.Config) Client { return client },
	})
	if err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if client.userName != "peter@example.com" || client.password != "secret-1" {
		t.Fatalf("username auth fields = username:%q password:%q", client.userName, client.password)
	}
	var payload map[string]string
	if err := json.Unmarshal(out.Bytes(), &payload); err != nil {
		t.Fatalf("json: %v", err)
	}
	if payload["session_id"] != "session-user" {
		t.Fatalf("session_id = %q", payload["session_id"])
	}
}
