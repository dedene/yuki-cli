package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

type CreditorItemsOptions struct {
	AdministrationID        string
	StartDate               string
	EndDate                 string
	IncludeBankTransactions bool
	SortOrder               string
}

func (c *Client) OutstandingCreditorItems(ctx context.Context, sessionID string, opts CreditorItemsOptions) ([]CreditorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingCreditorItems", params)
	if err != nil {
		return nil, err
	}
	var env outstandingCreditorItemsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingCreditorItems response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

func (c *Client) OutstandingCreditorItemsByDate(ctx context.Context, sessionID string, opts CreditorItemsOptions) ([]CreditorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingCreditorItemsByDate", params)
	if err != nil {
		return nil, err
	}
	var env creditorItemsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingCreditorItemsByDate response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

type outstandingCreditorItemsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []CreditorItem `xml:"Item"`
				} `xml:"OutstandingCreditorItems"`
			} `xml:"OutstandingCreditorItemsResult"`
		} `xml:"OutstandingCreditorItemsResponse"`
	} `xml:"Body"`
}

type creditorItemsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []CreditorItem `xml:"Item"`
				} `xml:"OutstandingCreditorItems"`
			} `xml:"OutstandingCreditorItemsByDateResult"`
		} `xml:"OutstandingCreditorItemsByDateResponse"`
	} `xml:"Body"`
}
