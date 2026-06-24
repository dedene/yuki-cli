package cmd

import (
	"errors"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type CreditorItemsCmd struct {
	All                  CreditorItemsAllCmd                  `cmd:"" help:"List all outstanding creditor purchase invoice items."`
	List                 CreditorItemsListCmd                 `cmd:"" help:"List outstanding creditor purchase invoice items."`
	ByOutstandingDate    CreditorItemsByOutstandingDateCmd    `cmd:"" name:"by-outstanding-date" help:"List creditor items open on an outstanding date."`
	WithPaymentReference CreditorItemsWithPaymentReferenceCmd `cmd:"" name:"with-payment-reference" help:"List outstanding creditor items with payment references."`
}

type CreditorItemsAllCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateAsc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *CreditorItemsAllCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingCreditorItems(rt.Context, sessionID, api.CreditorItemsOptions{
		AdministrationID:        administrationID,
		IncludeBankTransactions: c.IncludeBankTransactions,
		SortOrder:               c.SortOrder,
	})
	if err != nil {
		return err
	}
	if c.PaymentMethod != "" {
		items = filterCreditorItemsByPaymentMethod(items, c.PaymentMethod)
	}
	return renderCreditorItems(rt, globals, items)
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
	return renderCreditorItems(rt, globals, items)
}

type CreditorItemsByOutstandingDateCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	Date                    string `name:"date" required:"" help:"Outstanding date, YYYY-MM-DD."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateDesc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *CreditorItemsByOutstandingDateCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingCreditorItemsByDateOutstanding(rt.Context, sessionID, api.CreditorItemsOptions{
		AdministrationID:        administrationID,
		DateOutstanding:         c.Date,
		IncludeBankTransactions: c.IncludeBankTransactions,
		SortOrder:               c.SortOrder,
	})
	if err != nil {
		return err
	}
	if c.PaymentMethod != "" {
		items = filterCreditorItemsByPaymentMethod(items, c.PaymentMethod)
	}
	return renderCreditorItems(rt, globals, items)
}

type CreditorItemsWithPaymentReferenceCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	From                    string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To                      string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateDesc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *CreditorItemsWithPaymentReferenceCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingCreditorWithPaymentReference(rt.Context, sessionID, api.CreditorItemsOptions{
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
	return renderCreditorItemsWithPaymentReference(rt, globals, items)
}

func renderCreditorItems(rt *Runtime, globals *Globals, items []api.CreditorItem) error {
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

func renderCreditorItemsWithPaymentReference(rt *Runtime, globals *Globals, items []api.CreditorItem) error {
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
			item.PaymentReference,
			item.DocumentID,
			item.Description,
		})
	}
	return output.Table(rt.Out, []string{"DATE", "CONTACT", "ORIGINAL", "OPEN", "PAYMENT", "REFERENCE", "PAYMENT REF", "DOCUMENT", "DESCRIPTION"}, rows)
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
