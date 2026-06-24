package cmd

import (
	"errors"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type CreditorItemsCmd struct {
	List CreditorItemsListCmd `cmd:"" help:"List outstanding creditor purchase invoice items."`
}

type CreditorItemsListCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	From                    string `name:"from" help:"Start date, YYYY-MM-DD."`
	To                      string `name:"to" help:"End date, YYYY-MM-DD."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateAsc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *CreditorItemsListCmd) Run(rt *Runtime, globals *Globals) error {
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
	items, err := client.OutstandingCreditorItemsByDate(rt.Context, sessionID, api.CreditorItemsOptions{
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
		items = filterCreditorItemsByPaymentMethod(items, c.PaymentMethod)
	}
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

func filterCreditorItemsByPaymentMethod(items []api.CreditorItem, paymentMethod string) []api.CreditorItem {
	filtered := make([]api.CreditorItem, 0, len(items))
	for _, item := range items {
		if strings.EqualFold(item.PaymentMethod, paymentMethod) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
