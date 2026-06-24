package api

import (
	"context"
	"encoding/xml"
	"fmt"
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

type projectsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Projects []AccountingProject `xml:"Project"`
			} `xml:"GetProjectsResult"`
		} `xml:"GetProjectsResponse"`
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

type projectBalanceEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Balances []ProjectBalance `xml:"ProjectBalance"`
			} `xml:"GetProjectBalanceResult"`
		} `xml:"GetProjectBalanceResponse"`
	} `xml:"Body"`
}
