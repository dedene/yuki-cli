package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type TransactionsCmd struct {
	Details  TransactionDetailsCmd  `cmd:"" help:"List transaction details for GL accounts."`
	Document TransactionDocumentCmd `cmd:"" help:"Fetch the document attached to a transaction."`
}

type TransactionDetailsCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	GLAccount      string `name:"gl-account" help:"GL account code. Pass an empty value to include all accounts."`
	From           string `name:"from" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" help:"End date, YYYY-MM-DD."`
	FinancialMode  string `name:"financial-mode" default:"1" help:"Yuki financial mode."`
}

func (c *TransactionDetailsCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}
	if c.From == "" || c.To == "" {
		return errors.New("missing --from/--to; pass a Yuki date range like --from 2026-01-01 --to 2026-01-31")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	transactions, err := client.TransactionDetails(rt.Context, sessionID, api.TransactionDetailsOptions{
		AdministrationID: administrationID,
		GLAccountCode:    c.GLAccount,
		StartDate:        c.From,
		EndDate:          c.To,
		FinancialMode:    c.FinancialMode,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, transactions)
	}

	rows := make([][]string, 0, len(transactions))
	for _, tx := range transactions {
		rows = append(rows, []string{
			tx.TransactionDate,
			tx.GLAccountCode,
			tx.TransactionAmount,
			tx.Currency,
			tx.FullName,
			tx.DocumentID,
			tx.Description,
		})
	}
	return output.Table(rt.Out, []string{"DATE", "GL", "AMOUNT", "CCY", "CONTACT", "DOCUMENT", "DESCRIPTION"}, rows)
}

type TransactionDocumentCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Transaction    string `name:"transaction" required:"" help:"Transaction ID."`
	Output         string `name:"output" short:"o" help:"Write decoded file bytes to this path."`
}

func (c *TransactionDocumentCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	document, err := client.TransactionDocument(rt.Context, sessionID, administrationID, c.Transaction)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, document)
	}
	if c.Output == "" {
		return errors.New("missing --output; pass --output <path> or use --json to print the base64 payload")
	}
	return writeBase64File(rt.Out, c.Output, document.FileName, document.FileData)
}
