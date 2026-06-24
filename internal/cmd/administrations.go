package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type AdministrationsCmd struct {
	List             AdministrationsListCmd             `cmd:"" help:"List accessible administrations."`
	ID               AdministrationsIDCmd               `cmd:"" name:"id" help:"Resolve an administration ID by administration name."`
	WithCustomerCode AdministrationsWithCustomerCodeCmd `cmd:"" name:"with-customer-code" help:"List administrations including internal customer codes."`
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
	return renderAdministrations(rt, admins, false)
}

type AdministrationsIDCmd struct {
	Name string `name:"name" required:"" help:"Administration name."`
}

func (c *AdministrationsIDCmd) Run(rt *Runtime, globals *Globals) error {
	if c.Name == "" {
		return errors.New("missing --name; pass an administration name")
	}
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	id, err := client.AdministrationID(rt.Context, sessionID, c.Name)
	if err != nil {
		return err
	}
	result := administrationIDResult{Name: c.Name, ID: id}
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"NAME", "ID"}, [][]string{{result.Name, result.ID}})
}

type administrationIDResult struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type AdministrationsWithCustomerCodeCmd struct{}

func (c *AdministrationsWithCustomerCodeCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	admins, err := client.AdministrationsWithInternalCustomerCode(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, admins)
	}
	return renderAdministrations(rt, admins, true)
}

func renderAdministrations(rt *Runtime, admins []api.Administration, includeInternalCustomerCode bool) error {
	headers := []string{"ID", "NAME", "COUNTRY", "VAT", "DOMAIN", "ACTIVE"}
	if includeInternalCustomerCode {
		headers = append(headers, "INTERNAL CUSTOMER CODE")
	}
	rows := make([][]string, 0, len(admins))
	for _, admin := range admins {
		row := []string{admin.ID, admin.Name, admin.Country, admin.VATNumber, admin.DomainID, output.Bool(admin.Active)}
		if includeInternalCustomerCode {
			row = append(row, admin.InternalCustomerCode)
		}
		rows = append(rows, row)
	}
	return output.Table(rt.Out, headers, rows)
}
