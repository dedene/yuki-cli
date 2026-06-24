package cmd

import (
	"errors"
	"strconv"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ProjectsCmd struct {
	List       ProjectsListCmd       `cmd:"" help:"List projects for an administration."`
	ListWithID ProjectsListWithIDCmd `cmd:"" name:"list-with-id" help:"List projects with Yuki IDs for an administration."`
	Upsert     ProjectsUpsertCmd     `cmd:"" help:"Create or update a project by description."`
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

type ProjectsUpsertCmd struct {
	Administration   string `help:"Administration ID. Defaults to profile/global administration."`
	Description      string `name:"description" required:"" help:"Unique project description/name."`
	Code             string `name:"code" help:"Project code."`
	Company          string `name:"company" help:"Company/administration ID to store on the project."`
	Manager          string `name:"manager" help:"Manager user email address."`
	Contact          string `name:"contact" help:"Contact ID to assign."`
	Notes            string `name:"notes" help:"Project notes."`
	SecurityLevel    int    `name:"security-level" help:"Security level code: 1 all users, 2 employees, 3 management, 4 manager and members."`
	AllowOCRMatching string `name:"allow-ocr-matching" help:"Set OCR project-code matching: true or false."`
	StartDate        string `name:"start-date" help:"Start date, YYYY-MM-DD."`
	EndDate          string `name:"end-date" help:"End date, YYYY-MM-DD."`
	BudgetRevenue    string `name:"budget-revenue" help:"Estimated project revenue."`
	BudgetCosts      string `name:"budget-costs" help:"Estimated project costs."`
	DryRun           bool   `name:"dry-run" help:"Preview the project update without authenticating or sending it."`
}

func (c *ProjectsUpsertCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}
	opts, err := c.options(administrationID)
	if err != nil {
		return err
	}
	if c.DryRun {
		result := api.ProjectUpdateResult{
			AdministrationID: administrationID,
			Project:          opts.Project,
			DryRun:           true,
			Message:          "dry run; no project update sent",
		}
		return renderProjectUpdateResult(rt, globals, result)
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: accounting projects upsert")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.UpdateProject(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderProjectUpdateResult(rt, globals, result)
}

func (c *ProjectsUpsertCmd) options(administrationID string) (api.ProjectUpdateOptions, error) {
	ocrMatching, err := normalizeOptionalBool(c.AllowOCRMatching)
	if err != nil {
		return api.ProjectUpdateOptions{}, err
	}
	project := api.ProjectUpdate{
		Description:      c.Description,
		Code:             c.Code,
		Company:          c.Company,
		Manager:          c.Manager,
		Contact:          c.Contact,
		Notes:            c.Notes,
		AllowOCRMatching: ocrMatching,
		StartDate:        c.StartDate,
		EndDate:          c.EndDate,
		BudgetRevenue:    c.BudgetRevenue,
		BudgetCosts:      c.BudgetCosts,
	}
	if c.SecurityLevel != 0 {
		project.SecurityLevel = strconv.Itoa(c.SecurityLevel)
	}
	return api.ProjectUpdateOptions{
		AdministrationID: administrationID,
		Project:          project,
	}, nil
}

func normalizeOptionalBool(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "":
		return "", nil
	case "true", "1", "yes":
		return "true", nil
	case "false", "0", "no":
		return "false", nil
	default:
		return "", errors.New("invalid --allow-ocr-matching; use true or false")
	}
}

func renderProjectUpdateResult(rt *Runtime, globals *Globals, result api.ProjectUpdateResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"ADMINISTRATION", "DESCRIPTION", "CODE", "DRY RUN", "MESSAGE"}, [][]string{{
		result.AdministrationID,
		result.Project.Description,
		result.Project.Code,
		output.Bool(result.DryRun),
		result.Message,
	}})
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
