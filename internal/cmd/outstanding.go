package cmd

import (
	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type OutstandingCmd struct {
	Check      OutstandingCheckCmd      `cmd:"" help:"Check outstanding items by reference."`
	CheckAdmin OutstandingCheckAdminCmd `cmd:"" name:"check-admin" help:"Check outstanding items by administration and reference."`
}

type OutstandingCheckCmd struct {
	Reference string `name:"reference" required:"" help:"Outstanding item reference."`
}

func (c *OutstandingCheckCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.CheckOutstandingItem(rt.Context, sessionID, c.Reference)
	if err != nil {
		return err
	}
	return renderOutstandingItems(rt, globals, items)
}

type OutstandingCheckAdminCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Reference      string `name:"reference" required:"" help:"Outstanding item reference."`
}

func (c *OutstandingCheckAdminCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	items, err := client.CheckOutstandingItemAdmin(rt.Context, sessionID, administrationID, c.Reference)
	if err != nil {
		return err
	}
	return renderOutstandingItems(rt, globals, items)
}

func renderOutstandingItems(rt *Runtime, globals *Globals, items []api.OutstandingItem) error {
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
