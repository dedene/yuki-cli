package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/output"
)

type AccountingCmd struct {
	GLAccounts GLAccountsCmd `cmd:"" name:"gl-accounts" help:"Inspect GL accounts."`
}

type GLAccountsCmd struct {
	List GLAccountsListCmd `cmd:"" help:"List GL accounts for an administration."`
}

type GLAccountsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *GLAccountsListCmd) Run(rt *Runtime, globals *Globals) error {
	profile, err := loadProfile(globals)
	if err != nil {
		return err
	}
	administrationID := c.Administration
	if administrationID == "" {
		administrationID = profile.AdministrationID
	}
	if administrationID == "" {
		return errors.New("missing --administration; set it in config or pass --administration <id>")
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
