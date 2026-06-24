package cmd

import (
	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ProjectsCmd struct {
	List       ProjectsListCmd       `cmd:"" help:"List projects for an administration."`
	ListWithID ProjectsListWithIDCmd `cmd:"" name:"list-with-id" help:"List projects with Yuki IDs for an administration."`
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
