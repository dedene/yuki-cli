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

type projectsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Projects []AccountingProject `xml:"Project"`
			} `xml:"GetProjectsResult"`
		} `xml:"GetProjectsResponse"`
	} `xml:"Body"`
}
