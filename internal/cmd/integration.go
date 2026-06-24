package cmd

import (
	"github.com/dedene/yuki-cli/internal/output"
)

type IntegrationCmd struct {
	AdministrationData IntegrationAdministrationDataCmd `cmd:"" name:"administration-data" help:"Get administration integration data."`
}

type IntegrationAdministrationDataCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *IntegrationAdministrationDataCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	data, err := client.AdministrationData(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, data)
	}

	rows := [][]string{{
		data.CompanyName,
		data.MainContactName,
		data.MainContactEmail,
		data.City,
		data.Country,
		data.IBAN,
		data.VATNumber,
	}}
	return output.Table(rt.Out, []string{"COMPANY", "CONTACT", "EMAIL", "CITY", "COUNTRY", "IBAN", "VAT"}, rows)
}
