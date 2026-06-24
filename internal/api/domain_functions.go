package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) DomainFunctions(ctx context.Context, sessionID, domainID string) ([]DomainFunctionAssignment, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domain", Value: domainID},
	}
	data, err := c.call(ctx, "Domains", "GetDomainFunctions", params)
	if err != nil {
		return nil, err
	}
	var env domainFunctionsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetDomainFunctions response: %w", err)
	}
	assignments := env.Body.Response.Result.Functions.assignments()
	for i := range assignments {
		assignments[i].DomainID = domainID
	}
	return assignments, nil
}

type domainFunctionsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Functions domainFunctions `xml:"DomainFunctions"`
			} `xml:"GetDomainFunctionsResult"`
		} `xml:"GetDomainFunctionsResponse"`
	} `xml:"Body"`
}

type domainFunctions struct {
	BOResponsible    domainFunctionUser `xml:"BOResponsible"`
	BOBackup         domainFunctionUser `xml:"BOBackup"`
	BOController     domainFunctionUser `xml:"BOController"`
	BOAccountManager domainFunctionUser `xml:"BOAccountManager"`
}

func (f domainFunctions) assignments() []DomainFunctionAssignment {
	return []DomainFunctionAssignment{
		f.BOResponsible.assignment("BOResponsible"),
		f.BOBackup.assignment("BOBackup"),
		f.BOController.assignment("BOController"),
		f.BOAccountManager.assignment("BOAccountManager"),
	}
}

type domainFunctionUser struct {
	FullName string `xml:"FullName"`
	Login    string `xml:"Login"`
}

func (u domainFunctionUser) assignment(function string) DomainFunctionAssignment {
	return DomainFunctionAssignment{
		Function: function,
		FullName: u.FullName,
		Login:    u.Login,
	}
}
