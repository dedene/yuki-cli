package cmd

import (
	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ProjectsCmd struct {
	List       ProjectsListCmd       `cmd:"" help:"List projects for an administration."`
	ListWithID ProjectsListWithIDCmd `cmd:"" name:"list-with-id" help:"List projects with Yuki IDs for an administration."`
	Balance    ProjectsBalanceCmd    `cmd:"" help:"List project balances for a date range."`
}

type ProjectsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	SearchOption   string `name:"search-option" default:"All" help:"Yuki project search option: All, Company, Code, or Description."`
	SearchValue    string `name:"search-value" help:"Search value for the selected project search option."`
}

func (c *ProjectsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	projects, err := client.Projects(rt.Context, sessionID, api.ProjectsOptions{
		AdministrationID: administrationID,
		SearchOption:     c.SearchOption,
		SearchValue:      c.SearchValue,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, projects)
	}

	rows := make([][]string, 0, len(projects))
	for _, project := range projects {
		rows = append(rows, []string{
			project.HID,
			project.Code,
			project.Description,
			project.Company,
			project.Contact,
			project.StartDate,
			project.EndDate,
		})
	}
	return output.Table(rt.Out, []string{"HID", "CODE", "DESCRIPTION", "COMPANY", "CONTACT", "START", "END"}, rows)
}

type ProjectsListWithIDCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	SearchOption   string `name:"search-option" default:"All" help:"Yuki project search option: All, Company, Code, or Description."`
	SearchValue    string `name:"search-value" help:"Search value for the selected project search option."`
}

func (c *ProjectsListWithIDCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	projects, err := client.ProjectsAndID(rt.Context, sessionID, api.ProjectsOptions{
		AdministrationID: administrationID,
		SearchOption:     c.SearchOption,
		SearchValue:      c.SearchValue,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, projects)
	}

	rows := make([][]string, 0, len(projects))
	for _, project := range projects {
		rows = append(rows, []string{
			project.ID,
			project.HID,
			project.Code,
			project.Description,
			project.Company,
			project.Contact,
		})
	}
	return output.Table(rt.Out, []string{"ID", "HID", "CODE", "DESCRIPTION", "COMPANY", "CONTACT"}, rows)
}

type ProjectsBalanceCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	GLAccount      string `name:"gl-account" help:"Optional GL account code."`
	ProjectCode    string `name:"project-code" help:"Optional project code."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *ProjectsBalanceCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.ProjectBalance(rt.Context, sessionID, api.ProjectBalanceOptions{
		AdministrationID: administrationID,
		GLAccountCode:    c.GLAccount,
		ProjectCode:      c.ProjectCode,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}

	rows := make([][]string, 0, len(balances))
	for _, balance := range balances {
		rows = append(rows, []string{
			balance.GLAccountCode,
			balance.ProjectCode,
			balance.Project,
			balance.Amount,
		})
	}
	return output.Table(rt.Out, []string{"GL", "PROJECT CODE", "PROJECT", "AMOUNT"}, rows)
}
