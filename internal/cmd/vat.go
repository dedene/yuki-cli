package cmd

import (
	"github.com/dedene/yuki-cli/internal/output"
)

type VATCmd struct {
	Codes VATCodesCmd `cmd:"" help:"Inspect VAT codes."`
}

type VATCodesCmd struct {
	Active VATCodesActiveCmd `cmd:"" help:"List active VAT codes for an administration."`
}

type VATCodesActiveCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *VATCodesActiveCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	codes, err := client.ActiveVATCodes(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, codes)
	}

	rows := make([][]string, 0, len(codes))
	for _, code := range codes {
		rows = append(rows, []string{
			code.Type,
			code.Percentage,
			code.Country,
			code.StartDate,
			code.EndDate,
			code.Description,
			code.TypeDescription,
		})
	}
	return output.Table(rt.Out, []string{"TYPE", "PCT", "COUNTRY", "START", "END", "DESCRIPTION", "TYPE DESCRIPTION"}, rows)
}
