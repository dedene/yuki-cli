package cmd

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type SalesCmd struct {
	Invoices           SalesInvoicesCmd           `cmd:"" help:"Manage sales invoice imports."`
	RecognizedInvoices SalesRecognizedInvoicesCmd `cmd:"" name:"recognized-invoices" help:"Manage recognized sales invoice imports."`
	Items              SalesItemsCmd              `cmd:"" help:"Inspect sales items."`
}

type SalesInvoicesCmd struct {
	SchemaPath SalesInvoiceSchemaPathCmd `cmd:"" name:"schema-path" help:"Print the sales invoice XML schema URL."`
	Create     SalesInvoicesCreateCmd    `cmd:"" help:"Import sales invoices from an XML file."`
}

type SalesInvoiceSchemaPathCmd struct{}

func (c *SalesInvoiceSchemaPathCmd) Run(rt *Runtime, globals *Globals) error {
	profile, err := loadProfile(globals)
	if err != nil {
		return err
	}
	client := rt.client(globals, profile)
	path, err := client.SalesInvoiceSchemaPath(rt.Context)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, map[string]string{"schema_path": path})
	}
	_, err = fmt.Fprintln(rt.Out, path)
	return err
}

type SalesInvoicesCreateCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	File           string `name:"file" required:"" help:"SalesInvoices XML file to import." type:"existingfile"`
	DryRun         bool   `name:"dry-run" help:"Validate and preview the import without authenticating or sending it."`
}

func (c *SalesInvoicesCreateCmd) Run(rt *Runtime, globals *Globals) error {
	return runSalesInvoiceImport(rt, globals, salesInvoiceImportCommand{
		Command:        "sales invoices create",
		Operation:      "ProcessSalesInvoices",
		Administration: c.Administration,
		File:           c.File,
		DryRun:         c.DryRun,
		Call: func(client Client, sessionID string, opts api.SalesInvoiceImportOptions) (api.SalesInvoiceImportResponse, error) {
			return client.ProcessSalesInvoices(rt.Context, sessionID, opts)
		},
	})
}

type SalesRecognizedInvoicesCmd struct {
	Create SalesRecognizedInvoicesCreateCmd `cmd:"" help:"Import recognized sales invoices from an XML file."`
}

type SalesRecognizedInvoicesCreateCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	File           string `name:"file" required:"" help:"Recognized SalesInvoices XML file to import." type:"existingfile"`
	DryRun         bool   `name:"dry-run" help:"Validate and preview the import without authenticating or sending it."`
}

func (c *SalesRecognizedInvoicesCreateCmd) Run(rt *Runtime, globals *Globals) error {
	return runSalesInvoiceImport(rt, globals, salesInvoiceImportCommand{
		Command:        "sales recognized-invoices create",
		Operation:      "ProcessRecognizedSalesInvoices",
		Administration: c.Administration,
		File:           c.File,
		DryRun:         c.DryRun,
		Call: func(client Client, sessionID string, opts api.SalesInvoiceImportOptions) (api.SalesInvoiceImportResponse, error) {
			return client.ProcessRecognizedSalesInvoices(rt.Context, sessionID, opts)
		},
	})
}

type SalesItemsCmd struct {
	List SalesItemsListCmd `cmd:"" help:"List sales items for an administration."`
}

type SalesItemsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *SalesItemsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.SalesItems(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, items)
	}
	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{item.ID, item.Description})
	}
	return output.Table(rt.Out, []string{"ID", "DESCRIPTION"}, rows)
}

type salesInvoiceImportCommand struct {
	Command        string
	Operation      string
	Administration string
	File           string
	DryRun         bool
	Call           func(Client, string, api.SalesInvoiceImportOptions) (api.SalesInvoiceImportResponse, error)
}

type salesInvoiceXMLDocument struct {
	Path    string
	Content string
	Bytes   int
	Root    string
}

type salesInvoiceImportDryRun struct {
	DryRun           bool   `json:"dry_run"`
	Operation        string `json:"operation"`
	AdministrationID string `json:"administration_id"`
	File             string `json:"file"`
	Bytes            int    `json:"bytes"`
	Root             string `json:"root"`
	Message          string `json:"message"`
}

func runSalesInvoiceImport(rt *Runtime, globals *Globals, command salesInvoiceImportCommand) error {
	administrationID, err := resolveAdministrationID(globals, command.Administration)
	if err != nil {
		return err
	}
	doc, err := readSalesInvoiceXML(command.File)
	if err != nil {
		return err
	}
	if command.DryRun {
		return renderSalesInvoiceDryRun(rt, globals, salesInvoiceImportDryRun{
			DryRun:           true,
			Operation:        command.Operation,
			AdministrationID: administrationID,
			File:             doc.Path,
			Bytes:            doc.Bytes,
			Root:             doc.Root,
			Message:          "dry run; no invoice import sent",
		})
	}
	if globals.Readonly {
		return fmt.Errorf("--readonly blocks mutating command: %s", command.Command)
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := command.Call(client, sessionID, api.SalesInvoiceImportOptions{
		AdministrationID: administrationID,
		XMLDoc:           doc.Content,
	})
	if err != nil {
		return err
	}
	return renderSalesInvoiceImportResponse(rt, globals, result)
}

func readSalesInvoiceXML(path string) (salesInvoiceXMLDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return salesInvoiceXMLDocument{}, fmt.Errorf("read %s: %w", path, err)
	}
	root, err := validateXMLDocument(data)
	if err != nil {
		return salesInvoiceXMLDocument{}, fmt.Errorf("validate %s: %w", path, err)
	}
	data = stripXMLDeclaration(data)
	return salesInvoiceXMLDocument{
		Path:    path,
		Content: string(data),
		Bytes:   len(data),
		Root:    root,
	}, nil
}

func validateXMLDocument(data []byte) (string, error) {
	if len(bytes.TrimSpace(data)) == 0 {
		return "", errors.New("XML file is empty")
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	depth := 0
	root := ""
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", err
		}
		switch t := token.(type) {
		case xml.StartElement:
			if depth == 0 {
				if root != "" {
					return "", errors.New("XML document has multiple root elements")
				}
				root = t.Name.Local
			}
			depth++
		case xml.EndElement:
			if depth > 0 {
				depth--
			}
		case xml.CharData:
			if depth == 0 && strings.TrimSpace(string(t)) != "" {
				return "", errors.New("XML document has non-whitespace content outside the root element")
			}
		}
	}
	if root == "" {
		return "", errors.New("XML document has no root element")
	}
	if depth != 0 {
		return "", errors.New("XML document ended before the root element closed")
	}
	return root, nil
}

func stripXMLDeclaration(data []byte) []byte {
	trimmed := bytes.TrimSpace(data)
	trimmed = bytes.TrimPrefix(trimmed, []byte("\xef\xbb\xbf"))
	if !bytes.HasPrefix(trimmed, []byte("<?xml")) {
		return trimmed
	}
	end := bytes.Index(trimmed, []byte("?>"))
	if end == -1 {
		return trimmed
	}
	return bytes.TrimSpace(trimmed[end+2:])
}

func renderSalesInvoiceDryRun(rt *Runtime, globals *Globals, result salesInvoiceImportDryRun) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"OPERATION", "ADMINISTRATION", "FILE", "BYTES", "ROOT", "MESSAGE"}, [][]string{{
		result.Operation,
		result.AdministrationID,
		result.File,
		fmt.Sprintf("%d", result.Bytes),
		result.Root,
		result.Message,
	}})
}

func renderSalesInvoiceImportResponse(rt *Runtime, globals *Globals, result api.SalesInvoiceImportResponse) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	rows := make([][]string, 0, max(1, len(result.Invoices)))
	if len(result.Invoices) == 0 {
		rows = append(rows, salesInvoiceImportRow(result, api.SalesInvoiceImportInvoice{}))
	} else {
		for _, invoice := range result.Invoices {
			rows = append(rows, salesInvoiceImportRow(result, invoice))
		}
	}
	return output.Table(rt.Out, []string{"TIME", "ADMINISTRATION", "OK", "FAILED", "SKIPPED", "REFERENCE", "PROCESSED", "EMAIL", "MESSAGE"}, rows)
}

func salesInvoiceImportRow(result api.SalesInvoiceImportResponse, invoice api.SalesInvoiceImportInvoice) []string {
	return []string{
		result.TimeStamp,
		result.AdministrationID,
		fmt.Sprintf("%d", result.TotalSucceeded),
		fmt.Sprintf("%d", result.TotalFailed),
		fmt.Sprintf("%d", result.TotalSkipped),
		invoice.Reference,
		output.Bool(invoice.Processed),
		output.Bool(invoice.EmailSent),
		invoice.Message,
	}
}
