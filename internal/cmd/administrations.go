package cmd

import "github.com/dedene/yuki-cli/internal/output"

type AdministrationsCmd struct {
	List AdministrationsListCmd `cmd:"" help:"List accessible administrations."`
}

type AdministrationsListCmd struct{}

func (c *AdministrationsListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	admins, err := client.Administrations(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, admins)
	}
	rows := make([][]string, 0, len(admins))
	for _, admin := range admins {
		rows = append(rows, []string{admin.ID, admin.Name, admin.Country, admin.VATNumber, admin.DomainID, output.Bool(admin.Active)})
	}
	return output.Table(rt.Out, []string{"ID", "NAME", "COUNTRY", "VAT", "DOMAIN", "ACTIVE"}, rows)
}
