package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type AccountingCmd struct {
	GLAccounts     GLAccountsCmd               `cmd:"" name:"gl-accounts" help:"Inspect GL accounts."`
	Revenue        RevenueCmd                  `cmd:"" help:"Inspect revenue reports."`
	CreditorItems  CreditorItemsCmd            `cmd:"" name:"creditor-items" help:"Inspect outstanding creditor purchase invoices."`
	Transactions   TransactionsCmd             `cmd:"" help:"Inspect accounting transactions."`
	PaymentMethods AccountingPaymentMethodsCmd `cmd:"" name:"payment-methods" help:"Inspect accounting payment methods."`
}

type GLAccountsCmd struct {
	List           GLAccountsListCmd           `cmd:"" help:"List GL accounts for an administration."`
	Balance        GLAccountsBalanceCmd        `cmd:"" help:"List GL account balances at a transaction date."`
	BalanceFiscal  GLAccountsBalanceFiscalCmd  `cmd:"" name:"balance-fiscal" help:"List fiscal GL account balances at a transaction date."`
	BalanceYearEnd GLAccountsBalanceYearEndCmd `cmd:"" name:"balance-year-end" help:"List year-end GL account balances at a transaction date."`
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

type GLAccountsBalanceCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Date           string `name:"date" required:"" help:"Transaction date, YYYY-MM-DD."`
}

func (c *GLAccountsBalanceCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.GLAccountBalance(rt.Context, sessionID, api.GLAccountBalanceOptions{
		AdministrationID: administrationID,
		TransactionDate:  c.Date,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	return renderGLAccountBalances(rt, balances)
}

type GLAccountsBalanceFiscalCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Date           string `name:"date" required:"" help:"Transaction date, YYYY-MM-DD."`
}

func (c *GLAccountsBalanceFiscalCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.GLAccountBalanceFiscal(rt.Context, sessionID, api.GLAccountBalanceOptions{
		AdministrationID: administrationID,
		TransactionDate:  c.Date,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	return renderGLAccountBalances(rt, balances)
}

type GLAccountsBalanceYearEndCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Date           string `name:"date" required:"" help:"Transaction date, YYYY-MM-DD."`
}

func (c *GLAccountsBalanceYearEndCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.GLAccountBalanceYearEnd(rt.Context, sessionID, api.GLAccountBalanceOptions{
		AdministrationID: administrationID,
		TransactionDate:  c.Date,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	return renderGLAccountBalances(rt, balances)
}

func renderGLAccountBalances(rt *Runtime, balances []api.GLAccountBalanceItem) error {
	rows := make([][]string, 0, len(balances))
	for _, balance := range balances {
		rows = append(rows, []string{balance.Code, balance.BalanceType, balance.Amount, balance.Description})
	}
	return output.Table(rt.Out, []string{"CODE", "TYPE", "AMOUNT", "DESCRIPTION"}, rows)
}

type RevenueCmd struct {
	Net       RevenueNetCmd       `cmd:"" help:"Get net revenue for a date range."`
	NetFiscal RevenueNetFiscalCmd `cmd:"" name:"net-fiscal" help:"Get fiscal net revenue for a date range."`
}

type RevenueNetCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *RevenueNetCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	report, err := client.NetRevenue(rt.Context, sessionID, api.RevenueOptions{
		AdministrationID: administrationID,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, report)
	}
	return renderRevenueReport(rt, report)
}

type RevenueNetFiscalCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *RevenueNetFiscalCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	report, err := client.NetRevenueFiscal(rt.Context, sessionID, api.RevenueOptions{
		AdministrationID: administrationID,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, report)
	}
	return renderRevenueReport(rt, report)
}

func renderRevenueReport(rt *Runtime, report api.RevenueReport) error {
	return output.Table(rt.Out, []string{"FROM", "TO", "AMOUNT"}, [][]string{{report.StartDate, report.EndDate, report.Amount}})
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
