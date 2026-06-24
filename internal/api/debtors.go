package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

type DebtorItemsOptions struct {
	AdministrationID        string
	StartDate               string
	EndDate                 string
	IncludeBankTransactions bool
	SortOrder               string
}

func (c *Client) OutstandingDebtorItems(ctx context.Context, sessionID string, opts DebtorItemsOptions) ([]DebtorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingDebtorItems", params)
	if err != nil {
		return nil, err
	}
	var env outstandingDebtorItemsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingDebtorItems response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

type outstandingDebtorItemsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []DebtorItem `xml:"Item"`
				} `xml:"OutstandingDebtorItems"`
			} `xml:"OutstandingDebtorItemsResult"`
		} `xml:"OutstandingDebtorItemsResponse"`
	} `xml:"Body"`
}
