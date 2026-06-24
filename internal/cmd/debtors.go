package cmd

import (
	"errors"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type DebtorItemsCmd struct {
	All  DebtorItemsAllCmd  `cmd:"" help:"List all outstanding debtor sales invoice items."`
	List DebtorItemsListCmd `cmd:"" help:"List outstanding debtor sales invoice items."`
}

type DebtorItemsAllCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateAsc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *DebtorItemsAllCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingDebtorItems(rt.Context, sessionID, api.DebtorItemsOptions{
		AdministrationID:        administrationID,
		IncludeBankTransactions: c.IncludeBankTransactions,
		SortOrder:               c.SortOrder,
	})
	if err != nil {
		return err
	}
	if c.PaymentMethod != "" {
		items = filterDebtorItemsByPaymentMethod(items, c.PaymentMethod)
	}
	return renderDebtorItems(rt, globals, items)
}

type DebtorItemsListCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	From                    string `name:"from" help:"Start date, YYYY-MM-DD."`
	To                      string `name:"to" help:"End date, YYYY-MM-DD."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateDesc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *DebtorItemsListCmd) Run(rt *Runtime, globals *Globals) error {
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
	items, err := client.OutstandingDebtorItemsByDate(rt.Context, sessionID, api.DebtorItemsOptions{
		AdministrationID:        administrationID,
		StartDate:               c.From,
		EndDate:                 c.To,
		IncludeBankTransactions: c.IncludeBankTransactions,
		SortOrder:               c.SortOrder,
	})
	if err != nil {
		return err
	}
	if c.PaymentMethod != "" {
		items = filterDebtorItemsByPaymentMethod(items, c.PaymentMethod)
	}
	return renderDebtorItems(rt, globals, items)
}

func renderDebtorItems(rt *Runtime, globals *Globals, items []api.DebtorItem) error {
	if globals.JSON {
		return output.JSON(rt.Out, items)
	}

	rows := make([][]string, 0, len(items))
	for _, item := range items {
		rows = append(rows, []string{
			item.Date,
			item.Contact,
			item.OriginalAmount,
			item.OpenAmount,
			item.PaymentMethod,
			item.Reference,
			item.DocumentID,
			item.Description,
		})
	}
	return output.Table(rt.Out, []string{"DATE", "CONTACT", "ORIGINAL", "OPEN", "PAYMENT", "REFERENCE", "DOCUMENT", "DESCRIPTION"}, rows)
}

func filterDebtorItemsByPaymentMethod(items []api.DebtorItem, paymentMethod string) []api.DebtorItem {
	filtered := make([]api.DebtorItem, 0, len(items))
	for _, item := range items {
		if strings.EqualFold(item.PaymentMethod, paymentMethod) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
