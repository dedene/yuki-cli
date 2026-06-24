package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ContactsUpsertCmd struct {
	DomainID string `name:"domain" required:"" help:"Domain ID."`
	File     string `name:"file" required:"" help:"Contact XML file to create or update." type:"existingfile"`
	DryRun   bool   `name:"dry-run" help:"Validate and preview the contact XML without authenticating or sending it."`
}

type contactXMLDocument struct {
	Path    string
	Content string
	Bytes   int
	Root    string
}

type contactUpsertDryRun struct {
	DryRun    bool   `json:"dry_run"`
	Operation string `json:"operation"`
	DomainID  string `json:"domain_id"`
	File      string `json:"file"`
	Bytes     int    `json:"bytes"`
	Root      string `json:"root"`
	Message   string `json:"message"`
}

func (c *ContactsUpsertCmd) Run(rt *Runtime, globals *Globals) error {
	doc, err := readContactXML(c.File)
	if err != nil {
		return err
	}
	if c.DryRun {
		return renderContactUpsertDryRun(rt, globals, contactUpsertDryRun{
			DryRun:    true,
			Operation: "UpdateContact",
			DomainID:  c.DomainID,
			File:      doc.Path,
			Bytes:     doc.Bytes,
			Root:      doc.Root,
			Message:   "dry run; no contact sent",
		})
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: contacts upsert")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.UpdateContact(rt.Context, sessionID, api.ContactUpdateOptions{
		DomainID: c.DomainID,
		XMLDoc:   doc.Content,
	})
	if err != nil {
		return err
	}
	return renderContactUpdateResult(rt, globals, result)
}

func readContactXML(path string) (contactXMLDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return contactXMLDocument{}, fmt.Errorf("read %s: %w", path, err)
	}
	root, err := validateXMLDocument(data)
	if err != nil {
		return contactXMLDocument{}, fmt.Errorf("validate %s: %w", path, err)
	}
	if root != "Contact" {
		return contactXMLDocument{}, fmt.Errorf("validate %s: expected root element Contact, got %s", path, root)
	}
	data = stripXMLDeclaration(data)
	return contactXMLDocument{
		Path:    path,
		Content: string(data),
		Bytes:   len(data),
		Root:    root,
	}, nil
}

func renderContactUpsertDryRun(rt *Runtime, globals *Globals, result contactUpsertDryRun) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"OPERATION", "DOMAIN", "FILE", "BYTES", "ROOT", "MESSAGE"}, [][]string{{
		result.Operation,
		result.DomainID,
		result.File,
		strconv.Itoa(result.Bytes),
		result.Root,
		result.Message,
	}})
}

func renderContactUpdateResult(rt *Runtime, globals *Globals, result api.ContactUpdateResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"DOMAIN", "TIMESTAMP", "SUCCEEDED", "MESSAGE"}, [][]string{{
		result.DomainID,
		result.Timestamp,
		result.Succeeded,
		result.Message,
	}})
}
