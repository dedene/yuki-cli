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
	DateOutstanding         string
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

func (c *Client) OutstandingDebtorItemsByDate(ctx context.Context, sessionID string, opts DebtorItemsOptions) ([]DebtorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingDebtorItemsByDate", params)
	if err != nil {
		return nil, err
	}
	var env outstandingDebtorItemsByDateEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingDebtorItemsByDate response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

func (c *Client) OutstandingDebtorItemsByDateOutstanding(ctx context.Context, sessionID string, opts DebtorItemsOptions) ([]DebtorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "dateOutstanding", Value: opts.DateOutstanding},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingDebtorItemsByDateOutstanding", params)
	if err != nil {
		return nil, err
	}
	var env outstandingDebtorItemsByDateOutstandingEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingDebtorItemsByDateOutstanding response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

func (c *Client) OutstandingDebtorItemsWithLanguage(ctx context.Context, sessionID string, opts DebtorItemsOptions) ([]DebtorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingDebtorItemsWithLanguage", params)
	if err != nil {
		return nil, err
	}
	var env outstandingDebtorItemsWithLanguageEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingDebtorItemsWithLanguage response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

func (c *Client) OutstandingDebtorWithPaymentReference(ctx context.Context, sessionID string, opts DebtorItemsOptions) ([]DebtorItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "includeBankTransactions", Value: boolString(opts.IncludeBankTransactions)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "Accounting", "OutstandingDebtorWithPaymentReference", params)
	if err != nil {
		return nil, err
	}
	var env outstandingDebtorWithPaymentReferenceEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse OutstandingDebtorWithPaymentReference response: %w", err)
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

type outstandingDebtorWithPaymentReferenceEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []DebtorItem `xml:"Item"`
				} `xml:"OutstandingDebtorItems"`
			} `xml:"OutstandingDebtorWithPaymentReferenceResult"`
		} `xml:"OutstandingDebtorWithPaymentReferenceResponse"`
	} `xml:"Body"`
}

type outstandingDebtorItemsByDateEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []DebtorItem `xml:"Item"`
				} `xml:"OutstandingDebtorItems"`
			} `xml:"OutstandingDebtorItemsByDateResult"`
		} `xml:"OutstandingDebtorItemsByDateResponse"`
	} `xml:"Body"`
}

type outstandingDebtorItemsByDateOutstandingEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []DebtorItem `xml:"Item"`
				} `xml:"OutstandingDebtorItems"`
			} `xml:"OutstandingDebtorItemsByDateOutstandingResult"`
		} `xml:"OutstandingDebtorItemsByDateOutstandingResponse"`
	} `xml:"Body"`
}

type outstandingDebtorItemsWithLanguageEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []DebtorItem `xml:"Item"`
				} `xml:"OutstandingDebtorItems"`
			} `xml:"OutstandingDebtorItemsWithLanguageResult"`
		} `xml:"OutstandingDebtorItemsWithLanguageResponse"`
	} `xml:"Body"`
}
