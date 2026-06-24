package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"strings"
)

type ProjectsOptions struct {
	AdministrationID string
	SearchOption     string
	SearchValue      string
}

type ProjectBalanceOptions struct {
	AdministrationID string
	GLAccountCode    string
	ProjectCode      string
	StartDate        string
	EndDate          string
}

type ProjectUpdateOptions struct {
	AdministrationID string        `json:"administration_id,omitempty"`
	Project          ProjectUpdate `json:"project"`
	DryRun           bool          `json:"dry_run,omitempty"`
}

type ProjectUpdate struct {
	Description      string `json:"description,omitempty"`
	Code             string `json:"code,omitempty"`
	Company          string `json:"company,omitempty"`
	Manager          string `json:"manager,omitempty"`
	Contact          string `json:"contact,omitempty"`
	Notes            string `json:"notes,omitempty"`
	SecurityLevel    string `json:"security_level,omitempty"`
	AllowOCRMatching string `json:"allow_ocr_matching,omitempty"`
	StartDate        string `json:"start_date,omitempty"`
	EndDate          string `json:"end_date,omitempty"`
	BudgetRevenue    string `json:"budget_revenue,omitempty"`
	BudgetCosts      string `json:"budget_costs,omitempty"`
}

type ProjectUpdateResult struct {
	AdministrationID string        `json:"administration_id,omitempty"`
	Project          ProjectUpdate `json:"project"`
	DryRun           bool          `json:"dry_run,omitempty"`
	Message          string        `json:"message,omitempty"`
}

func (c *Client) Projects(ctx context.Context, sessionID string, opts ProjectsOptions) ([]AccountingProject, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "searchOption", Value: opts.SearchOption},
		{Name: "searchValue", Value: opts.SearchValue},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetProjects", params)
	if err != nil {
		return nil, err
	}
	var env projectsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetProjects response: %w", err)
	}
	return env.Body.Response.Result.Projects, nil
}

func (c *Client) ProjectsAndID(ctx context.Context, sessionID string, opts ProjectsOptions) ([]AccountingProject, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "searchOption", Value: opts.SearchOption},
		{Name: "searchValue", Value: opts.SearchValue},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetProjectsAndID", params)
	if err != nil {
		return nil, err
	}
	var env projectsAndIDEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetProjectsAndID response: %w", err)
	}
	return env.Body.Response.Result.Projects, nil
}

func (c *Client) ArchiveProjects(ctx context.Context, sessionID, administrationID string) ([]AccountingProject, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "Archive", "Projects", params)
	if err != nil {
		return nil, err
	}
	var env archiveProjectsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Projects response: %w", err)
	}
	return env.Body.Response.Result.projects(), nil
}

func (c *Client) ProjectBalance(ctx context.Context, sessionID string, opts ProjectBalanceOptions) ([]ProjectBalance, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "GLAccountCode", Value: opts.GLAccountCode},
		{Name: "projectCode", Value: opts.ProjectCode},
		{Name: "StartDate", Value: opts.StartDate},
		{Name: "EndDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetProjectBalance", params)
	if err != nil {
		return nil, err
	}
	var env projectBalanceEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetProjectBalance response: %w", err)
	}
	return env.Body.Response.Result.Balances, nil
}

func (c *Client) UpdateProject(ctx context.Context, sessionID string, opts ProjectUpdateOptions) (ProjectUpdateResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "project", Value: projectUpdateXML(opts.Project), Raw: true},
	}
	data, err := c.call(ctx, "Projects", "UpdateProject", params)
	if err != nil {
		return ProjectUpdateResult{}, err
	}
	var env updateProjectEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return ProjectUpdateResult{}, fmt.Errorf("parse UpdateProject response: %w", err)
	}
	return ProjectUpdateResult{
		AdministrationID: opts.AdministrationID,
		Project:          opts.Project,
		Message:          "project upserted",
	}, nil
}

func projectUpdateXML(project ProjectUpdate) string {
	var b strings.Builder
	writeProjectField(&b, "Description", project.Description)
	writeProjectField(&b, "Code", project.Code)
	writeProjectField(&b, "Company", project.Company)
	writeProjectField(&b, "Manager", project.Manager)
	writeProjectField(&b, "Contact", project.Contact)
	writeProjectField(&b, "Notes", project.Notes)
	writeProjectField(&b, "SecurityLevel", project.SecurityLevel)
	writeProjectField(&b, "AllowOCRMatching", project.AllowOCRMatching)
	writeProjectField(&b, "StartDate", project.StartDate)
	writeProjectField(&b, "EndDate", project.EndDate)
	writeProjectField(&b, "BudgetRevenue", project.BudgetRevenue)
	writeProjectField(&b, "BudgetCosts", project.BudgetCosts)
	return b.String()
}

func writeProjectField(b *strings.Builder, name string, value string) {
	if value == "" {
		return
	}
	b.WriteString("<they:")
	b.WriteString(name)
	b.WriteString(">")
	b.WriteString(html.EscapeString(value))
	b.WriteString("</they:")
	b.WriteString(name)
	b.WriteString(">")
}

type projectsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Projects []AccountingProject `xml:"Project"`
			} `xml:"GetProjectsResult"`
		} `xml:"GetProjectsResponse"`
	} `xml:"Body"`
}

type updateProjectEnvelope struct {
	Body struct {
		Response struct{} `xml:"UpdateProjectResponse"`
	} `xml:"Body"`
}

type projectsAndIDEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Projects []AccountingProject `xml:"Project"`
			} `xml:"GetProjectsAndIDResult"`
		} `xml:"GetProjectsAndIDResponse"`
	} `xml:"Body"`
}

type archiveProjectsEnvelope struct {
	Body struct {
		Response struct {
			Result archiveProjectsResult `xml:"ProjectsResult"`
		} `xml:"ProjectsResponse"`
	} `xml:"Body"`
}

type archiveProjectsResult struct {
	Direct []AccountingProject `xml:"Project"`
	Nested struct {
		Projects []AccountingProject `xml:"Project"`
	} `xml:"Projects"`
}

func (r archiveProjectsResult) projects() []AccountingProject {
	if len(r.Nested.Projects) > 0 {
		return r.Nested.Projects
	}
	return r.Direct
}

type projectBalanceEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Balances []ProjectBalance `xml:"ProjectBalance"`
			} `xml:"GetProjectBalanceResult"`
		} `xml:"GetProjectBalanceResponse"`
	} `xml:"Body"`
}
