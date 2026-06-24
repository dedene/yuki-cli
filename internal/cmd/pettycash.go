package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type PettyCashCmd struct {
	Statement   PettyCashStatementCmd   `cmd:"" help:"Import petty cash CSV statements."`
	Line        PettyCashLineCmd        `cmd:"" help:"Import a single petty cash line."`
	ProjectLine PettyCashProjectLineCmd `cmd:"" name:"project-line" help:"Import a single petty cash project line."`
}

type PettyCashStatementCmd struct {
	Import PettyCashStatementImportCmd `cmd:"" help:"Import a Yuki-format petty cash CSV statement."`
}

type PettyCashStatementImportCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	File           string `name:"file" required:"" help:"Petty cash CSV file to import." type:"existingfile"`
	DryRun         bool   `name:"dry-run" help:"Validate and preview the statement without authenticating or sending it."`
}

type pettyCashStatementFile struct {
	Path    string
	Content string
	Bytes   int
}

func (c *PettyCashStatementImportCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}
	statement, err := readPettyCashStatement(c.File)
	if err != nil {
		return err
	}
	if c.DryRun {
		return renderPettyCashResult(rt, globals, api.PettyCashImportResult{
			Operation:        "ImportStatement",
			AdministrationID: administrationID,
			DryRun:           true,
			Message:          fmt.Sprintf("dry run; no petty cash statement sent; bytes=%d", statement.Bytes),
		})
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: petty-cash statement import")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.ImportPettyCashStatement(rt.Context, sessionID, api.PettyCashStatementImportOptions{
		AdministrationID: administrationID,
		StatementText:    statement.Content,
	})
	if err != nil {
		return err
	}
	return renderPettyCashResult(rt, globals, result)
}

func readPettyCashStatement(path string) (pettyCashStatementFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return pettyCashStatementFile{}, fmt.Errorf("read %s: %w", path, err)
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	if len(bytes.TrimSpace(data)) == 0 {
		return pettyCashStatementFile{}, fmt.Errorf("validate %s: statement file is empty", path)
	}
	return pettyCashStatementFile{
		Path:    path,
		Content: string(data),
		Bytes:   len(data),
	}, nil
}

type PettyCashLineCmd struct {
	Import PettyCashLineImportCmd `cmd:"" help:"Import a single petty cash statement line."`
}

type PettyCashLineImportCmd struct {
	AccountGLCode          string `name:"account-gl-code" required:"" help:"Petty cash GL account code, e.g. 570000."`
	TransactionCode        string `name:"transaction-code" help:"Transaction code for processing rules."`
	OffsetAccount          string `name:"offset-account" help:"Offset account or contact code."`
	OffsetName             string `name:"offset-name" help:"Offset account/contact name."`
	TransactionDate        string `name:"transaction-date" required:"" help:"Transaction date, YYYY-MM-DD."`
	TransactionDescription string `name:"description" help:"Transaction description."`
	Amount                 string `name:"amount" required:"" help:"Transaction amount; costs negative, revenues positive."`
	DryRun                 bool   `name:"dry-run" help:"Preview the line without authenticating or sending it."`
}

func (c *PettyCashLineImportCmd) Run(rt *Runtime, globals *Globals) error {
	opts := api.PettyCashLineImportOptions{
		AccountGLCode:          c.AccountGLCode,
		TransactionCode:        c.TransactionCode,
		OffsetAccount:          c.OffsetAccount,
		OffsetName:             c.OffsetName,
		TransactionDate:        c.TransactionDate,
		TransactionDescription: c.TransactionDescription,
		Amount:                 c.Amount,
	}
	if c.DryRun {
		return renderPettyCashResult(rt, globals, pettyCashLineDryRun("ImportSingleStatementLine", opts))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: petty-cash line import")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.ImportPettyCashLine(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderPettyCashResult(rt, globals, result)
}

type PettyCashProjectLineCmd struct {
	Import PettyCashProjectLineImportCmd `cmd:"" help:"Import a single petty cash statement line with project fields."`
}

type PettyCashProjectLineImportCmd struct {
	AccountGLCode          string `name:"account-gl-code" required:"" help:"Petty cash GL account code, e.g. 570000."`
	TransactionCode        string `name:"transaction-code" help:"Transaction code for processing rules."`
	OffsetAccount          string `name:"offset-account" help:"Offset account or contact code."`
	OffsetName             string `name:"offset-name" help:"Offset account/contact name."`
	TransactionDate        string `name:"transaction-date" required:"" help:"Transaction date, YYYY-MM-DD."`
	TransactionDescription string `name:"description" help:"Transaction description."`
	Amount                 string `name:"amount" required:"" help:"Transaction amount; costs negative, revenues positive."`
	ProjectCode            string `name:"project-code" help:"Project code."`
	ProjectName            string `name:"project-name" help:"Project name."`
	DryRun                 bool   `name:"dry-run" help:"Preview the project line without authenticating or sending it."`
}

func (c *PettyCashProjectLineImportCmd) Run(rt *Runtime, globals *Globals) error {
	opts := api.PettyCashLineImportOptions{
		AccountGLCode:          c.AccountGLCode,
		TransactionCode:        c.TransactionCode,
		OffsetAccount:          c.OffsetAccount,
		OffsetName:             c.OffsetName,
		TransactionDate:        c.TransactionDate,
		TransactionDescription: c.TransactionDescription,
		Amount:                 c.Amount,
		ProjectCode:            c.ProjectCode,
		ProjectName:            c.ProjectName,
	}
	if c.DryRun {
		return renderPettyCashResult(rt, globals, pettyCashLineDryRun("ImportSingleStatementProjectLine", opts))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: petty-cash project-line import")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.ImportPettyCashProjectLine(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderPettyCashResult(rt, globals, result)
}

func pettyCashLineDryRun(operation string, opts api.PettyCashLineImportOptions) api.PettyCashImportResult {
	return api.PettyCashImportResult{
		Operation:       operation,
		AccountGLCode:   opts.AccountGLCode,
		TransactionDate: opts.TransactionDate,
		Amount:          opts.Amount,
		ProjectCode:     opts.ProjectCode,
		ProjectName:     opts.ProjectName,
		DryRun:          true,
		Message:         "dry run; no petty cash line sent",
	}
}

func renderPettyCashResult(rt *Runtime, globals *Globals, result api.PettyCashImportResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"OPERATION", "ADMINISTRATION", "ACCOUNT", "DATE", "AMOUNT", "PROJECT", "DOCUMENT", "DRY RUN", "MESSAGE"}, [][]string{{
		result.Operation,
		result.AdministrationID,
		result.AccountGLCode,
		result.TransactionDate,
		result.Amount,
		firstNonEmptyString(result.ProjectCode, result.ProjectName),
		result.DocumentID,
		output.Bool(result.DryRun),
		result.Message,
	}})
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
