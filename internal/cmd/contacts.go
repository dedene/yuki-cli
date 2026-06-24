package cmd

import (
	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ContactsCmd struct {
	Search             ContactsSearchCmd                `cmd:"" help:"Search contacts in a domain."`
	SuppliersCustomers ContactsSuppliersAndCustomersCmd `cmd:"" name:"suppliers-customers" help:"Search supplier and customer contacts in a domain."`
	Upsert             ContactsUpsertCmd                `cmd:"" help:"Create or update a contact from Yuki contact XML."`
}

type ContactsSearchCmd struct {
	DomainID      string `name:"domain" required:"" help:"Domain ID."`
	SearchOption  string `name:"search-option" default:"All" help:"Yuki contact search option."`
	SearchValue   string `name:"search-value" help:"Search value."`
	SortOrder     string `name:"sort-order" default:"CreatedDesc" help:"Yuki contact sort order."`
	ModifiedAfter string `name:"modified-after" default:"0001-01-01" help:"Return contacts modified after this date, YYYY-MM-DD."`
	Active        string `name:"active" default:"Both" help:"Active filter: Both, Active, or Inactive."`
	Page          int    `name:"page" default:"1" help:"One-based page number. Yuki returns max 100 records per page."`
}

func (c *ContactsSearchCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	contacts, err := client.SearchContacts(rt.Context, sessionID, c.options())
	if err != nil {
		return err
	}
	return renderContacts(rt, globals, contacts)
}

func (c *ContactsSearchCmd) options() api.ContactSearchOptions {
	return api.ContactSearchOptions{
		DomainID:      c.DomainID,
		SearchOption:  c.SearchOption,
		SearchValue:   c.SearchValue,
		SortOrder:     c.SortOrder,
		ModifiedAfter: c.ModifiedAfter,
		Active:        c.Active,
		PageNumber:    c.Page,
	}
}

type ContactsSuppliersAndCustomersCmd struct {
	DomainID      string `name:"domain" required:"" help:"Domain ID."`
	SearchOption  string `name:"search-option" default:"All" help:"Yuki contact search option."`
	SearchValue   string `name:"search-value" help:"Search value."`
	SortOrder     string `name:"sort-order" default:"CreatedAsc" help:"Yuki contact sort order."`
	ModifiedAfter string `name:"modified-after" default:"0001-01-01" help:"Return contacts modified after this date, YYYY-MM-DD."`
	Active        string `name:"active" default:"Both" help:"Active filter: Both, Active, or Inactive."`
	Page          int    `name:"page" default:"1" help:"One-based page number. Yuki returns max 100 records per page."`
	ContactType   string `name:"contact-type" default:"Both" help:"Contact type: Customer, Supplier, Both, or None."`
}

func (c *ContactsSuppliersAndCustomersCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	contacts, err := client.SuppliersAndCustomers(rt.Context, sessionID, c.options())
	if err != nil {
		return err
	}
	return renderContacts(rt, globals, contacts)
}

func (c *ContactsSuppliersAndCustomersCmd) options() api.ContactSearchOptions {
	return api.ContactSearchOptions{
		DomainID:      c.DomainID,
		SearchOption:  c.SearchOption,
		SearchValue:   c.SearchValue,
		SortOrder:     c.SortOrder,
		ModifiedAfter: c.ModifiedAfter,
		Active:        c.Active,
		PageNumber:    c.Page,
		ContactType:   c.ContactType,
	}
}

func renderContacts(rt *Runtime, globals *Globals, contacts []api.Contact) error {
	if globals.JSON {
		return output.JSON(rt.Out, contacts)
	}
	rows := make([][]string, 0, len(contacts))
	for _, contact := range contacts {
		rows = append(rows, []string{
			contact.ID,
			contact.HID,
			contact.Code,
			contact.Name,
			contact.City,
			contact.Country,
			contact.EmailWork,
			contact.VATNumber,
		})
	}
	return output.Table(rt.Out, []string{"ID", "HID", "CODE", "NAME", "CITY", "COUNTRY", "EMAIL", "VAT"}, rows)
}
