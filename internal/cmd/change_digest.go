package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ChangeDigestCmd struct {
	Transactions ChangeDigestTransactionsCmd `cmd:"" help:"List updated and deleted transactions."`
	Detail       ChangeDigestDetailCmd       `cmd:"" help:"Get one changed transaction detail."`
}

type ChangeDigestTransactionsCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	From           string `name:"from" required:"" help:"Start datetime, e.g. 2025-07-23T00:00:00.00Z."`
	To             string `name:"to" required:"" help:"End datetime, e.g. 2025-08-23T13:00:00.00Z."`
	Limit          int    `name:"limit" default:"100" help:"Number of records to request."`
	StartRecord    int    `name:"start-record" default:"0" help:"Zero-based start record."`
}

func (c *ChangeDigestTransactionsCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}
	if c.Limit < 0 || c.StartRecord < 0 {
		return errors.New("--limit and --start-record must be zero or greater")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	transactions, err := client.UpdatedAndDeletedTransactions(rt.Context, sessionID, api.UpdatedAndDeletedTransactionsOptions{
		AdministrationID: administrationID,
		StartDate:        c.From,
		EndDate:          c.To,
		NumberOfRecords:  c.Limit,
		StartRecord:      c.StartRecord,
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
			tx.Updated,
			tx.Deleted,
			tx.TransactionDate,
			tx.GLAccountCode,
			tx.TransactionAmount,
			tx.Currency,
			tx.FullName,
			tx.DocumentID,
			tx.Description,
		})
	}
	return output.Table(rt.Out, []string{"UPDATED", "DELETED", "DATE", "GL", "AMOUNT", "CCY", "CONTACT", "DOCUMENT", "DESCRIPTION"}, rows)
}

type ChangeDigestDetailCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Transaction    string `name:"transaction" required:"" help:"Transaction ID."`
}

func (c *ChangeDigestDetailCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	tx, err := client.ChangeDigestTransactionDetail(rt.Context, sessionID, administrationID, c.Transaction)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, tx)
	}
	rows := [][]string{{
		tx.TransactionDate,
		tx.GLAccountCode,
		tx.TransactionAmount,
		tx.Currency,
		tx.FullName,
		tx.DocumentID,
		tx.Description,
	}}
	return output.Table(rt.Out, []string{"DATE", "GL", "AMOUNT", "CCY", "CONTACT", "DOCUMENT", "DESCRIPTION"}, rows)
}
