package cmd

import (
	"errors"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type DebtorItemsCmd struct {
	All                  DebtorItemsAllCmd                  `cmd:"" help:"List all outstanding debtor sales invoice items."`
	List                 DebtorItemsListCmd                 `cmd:"" help:"List outstanding debtor sales invoice items."`
	ByOutstandingDate    DebtorItemsByOutstandingDateCmd    `cmd:"" name:"by-outstanding-date" help:"List debtor items open on an outstanding date."`
	WithLanguage         DebtorItemsWithLanguageCmd         `cmd:"" name:"with-language" help:"List outstanding debtor items including layout language."`
	WithPaymentReference DebtorItemsWithPaymentReferenceCmd `cmd:"" name:"with-payment-reference" help:"List outstanding debtor items with payment references."`
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

type DebtorItemsByOutstandingDateCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	Date                    string `name:"date" required:"" help:"Outstanding date, YYYY-MM-DD."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateDesc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *DebtorItemsByOutstandingDateCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingDebtorItemsByDateOutstanding(rt.Context, sessionID, api.DebtorItemsOptions{
		AdministrationID:        administrationID,
		DateOutstanding:         c.Date,
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

type DebtorItemsWithLanguageCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateDesc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *DebtorItemsWithLanguageCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingDebtorItemsWithLanguage(rt.Context, sessionID, api.DebtorItemsOptions{
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
	return renderDebtorItemsWithLanguage(rt, globals, items)
}

type DebtorItemsWithPaymentReferenceCmd struct {
	Administration          string `help:"Administration ID. Defaults to profile/global administration."`
	From                    string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To                      string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
	IncludeBankTransactions bool   `name:"include-bank-transactions" help:"Include outstanding bank transactions."`
	SortOrder               string `name:"sort-order" default:"DateDesc" help:"Yuki sort order, e.g. DateAsc or DateDesc."`
	PaymentMethod           string `name:"payment-method" help:"Filter results by payment method, e.g. Creditcard."`
}

func (c *DebtorItemsWithPaymentReferenceCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.OutstandingDebtorWithPaymentReference(rt.Context, sessionID, api.DebtorItemsOptions{
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
	return renderDebtorItemsWithPaymentReference(rt, globals, items)
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

func renderDebtorItemsWithPaymentReference(rt *Runtime, globals *Globals, items []api.DebtorItem) error {
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

func renderDebtorItemsWithLanguage(rt *Runtime, globals *Globals, items []api.DebtorItem) error {
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
			debtorItemLanguage(item),
			item.Reference,
			item.DocumentID,
			item.Description,
		})
	}
	return output.Table(rt.Out, []string{"DATE", "CONTACT", "ORIGINAL", "OPEN", "PAYMENT", "LANGUAGE", "REFERENCE", "DOCUMENT", "DESCRIPTION"}, rows)
}

func debtorItemLanguage(item api.DebtorItem) string {
	if item.LayoutLanguage != "" {
		return item.LayoutLanguage
	}
	return item.Language
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
