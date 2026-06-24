package cmd

import (
	"github.com/dedene/yuki-cli/internal/output"
)

type FiscalTableCmd struct {
	Totals FiscalTableTotalsCmd `cmd:"" help:"Get fiscal table totals for a company."`
}

type FiscalTableTotalsCmd struct {
	Company string `name:"company" help:"Company/administration ID. Defaults to profile/global administration."`
	Year    int    `name:"year" required:"" help:"Fiscal table year."`
}

func (c *FiscalTableTotalsCmd) Run(rt *Runtime, globals *Globals) error {
	companyID, err := resolveAdministrationID(globals, c.Company)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	totals, err := client.FiscalTable(rt.Context, sessionID, companyID, c.Year)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, totals)
	}

	rows := [][]string{{
		totals.RevenueTotal,
		totals.GrossMarginTotal,
		totals.ProfessionalCostsTotal,
		totals.SocialContributionsTotal,
	}}
	return output.Table(rt.Out, []string{"REVENUE", "GROSS MARGIN", "PRO COSTS", "SOCIAL CONTRIBUTIONS"}, rows)
}
