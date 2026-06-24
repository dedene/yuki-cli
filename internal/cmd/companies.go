package cmd

import (
	"strconv"

	"github.com/dedene/yuki-cli/internal/output"
)

type CompaniesCmd struct {
	List CompaniesListCmd `cmd:"" help:"List accessible companies."`
}

type CompaniesListCmd struct{}

func (c *CompaniesListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	companies, err := client.Companies(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, companies)
	}
	rows := make([][]string, 0, len(companies))
	for _, company := range companies {
		rows = append(rows, []string{company.ID, company.Name, strconv.FormatBool(company.Active)})
	}
	return output.Table(rt.Out, []string{"ID", "NAME", "ACTIVE"}, rows)
}
