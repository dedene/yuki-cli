package cmd

import (
	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type VATCmd struct {
	Codes   VATCodesCmd   `cmd:"" help:"Inspect VAT codes."`
	Returns VATReturnsCmd `cmd:"" help:"Inspect VAT returns."`
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

type VATReturnsCmd struct {
	List VATReturnsListCmd `cmd:"" help:"List VAT returns for an administration."`
}

type VATReturnsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Year           int    `name:"year" required:"" help:"VAT return year."`
	ModifiedAfter  string `name:"modified-after" required:"" help:"Return VAT returns modified after this date/datetime, e.g. 2021-01-01."`
}

func (c *VATReturnsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	returns, err := client.VATReturns(rt.Context, sessionID, api.VATReturnListOptions{
		AdministrationID: administrationID,
		Year:             c.Year,
		ModifiedAfter:    c.ModifiedAfter,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, returns)
	}

	rows := make([][]string, 0, len(returns))
	for _, vatReturn := range returns {
		rows = append(rows, []string{
			vatReturn.StartDate,
			vatReturn.EndDate,
			vatReturn.Status,
			vatReturn.SendDate,
			vatReturn.AcknowledgeDate,
			vatReturn.Modified,
			vatReturn.DocumentID,
		})
	}
	return output.Table(rt.Out, []string{"START", "END", "STATUS", "SENT", "ACK", "MODIFIED", "DOCUMENT"}, rows)
}
