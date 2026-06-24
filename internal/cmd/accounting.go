package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type AccountingCmd struct {
	GLAccounts     GLAccountsCmd               `cmd:"" name:"gl-accounts" help:"Inspect GL accounts."`
	CreditorItems  CreditorItemsCmd            `cmd:"" name:"creditor-items" help:"Inspect outstanding creditor purchase invoices."`
	Transactions   TransactionsCmd             `cmd:"" help:"Inspect accounting transactions."`
	PaymentMethods AccountingPaymentMethodsCmd `cmd:"" name:"payment-methods" help:"Inspect accounting payment methods."`
}

type GLAccountsCmd struct {
	List GLAccountsListCmd `cmd:"" help:"List GL accounts for an administration."`
}

type GLAccountsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *GLAccountsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	accounts, err := client.GLAccounts(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, accounts)
	}
	rows := make([][]string, 0, len(accounts))
	for _, account := range accounts {
		rows = append(rows, []string{account.Code, account.Type, account.Subtype, output.Bool(account.Enabled), account.Description})
	}
	return output.Table(rt.Out, []string{"CODE", "TYPE", "SUBTYPE", "ENABLED", "DESCRIPTION"}, rows)
}

func resolveAdministrationID(globals *Globals, explicit string) (string, error) {
	profile, err := loadProfile(globals)
	if err != nil {
		return "", err
	}
	administrationID := explicit
	if administrationID == "" {
		administrationID = profile.AdministrationID
	}
	if administrationID == "" {
		return "", errors.New("missing --administration; set it in config or pass --administration <id>")
	}
	return administrationID, nil
}

type AccountingPaymentMethodsCmd struct {
	List AccountingPaymentMethodsListCmd `cmd:"" help:"List custom payment methods for an administration."`
}

type AccountingPaymentMethodsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *AccountingPaymentMethodsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	methods, err := client.CustomPaymentMethods(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	return renderPaymentMethods(rt, globals, methods)
}

func renderPaymentMethods(rt *Runtime, globals *Globals, methods []api.PaymentMethod) error {
	if globals.JSON {
		return output.JSON(rt.Out, methods)
	}

	rows := make([][]string, 0, len(methods))
	for _, method := range methods {
		rows = append(rows, []string{method.ID, method.Description})
	}
	return output.Table(rt.Out, []string{"ID", "DESCRIPTION"}, rows)
}
