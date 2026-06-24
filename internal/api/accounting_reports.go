package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

type GLAccountBalanceOptions struct {
	AdministrationID string
	TransactionDate  string
}

type GLAccountTransactionsOptions struct {
	AdministrationID string
	GLAccountCode    string
	StartDate        string
	EndDate          string
}

type RevenueOptions struct {
	AdministrationID string
	StartDate        string
	EndDate          string
}

func (c *Client) GLAccountBalance(ctx context.Context, sessionID string, opts GLAccountBalanceOptions) ([]GLAccountBalanceItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "transactionDate", Value: opts.TransactionDate},
	}
	data, err := c.call(ctx, "Accounting", "GLAccountBalance", params)
	if err != nil {
		return nil, err
	}
	var env glAccountBalanceEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GLAccountBalance response: %w", err)
	}
	return env.Body.Response.Result.Balance.Items, nil
}

func (c *Client) GLAccountBalanceFiscal(ctx context.Context, sessionID string, opts GLAccountBalanceOptions) ([]GLAccountBalanceItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "transactionDate", Value: opts.TransactionDate},
	}
	data, err := c.call(ctx, "Accounting", "GLAccountBalanceFiscal", params)
	if err != nil {
		return nil, err
	}
	var env glAccountBalanceFiscalEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GLAccountBalanceFiscal response: %w", err)
	}
	return env.Body.Response.Result.Balance.Items, nil
}

func (c *Client) GLAccountBalanceYearEnd(ctx context.Context, sessionID string, opts GLAccountBalanceOptions) ([]GLAccountBalanceItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "transactionDate", Value: opts.TransactionDate},
	}
	data, err := c.call(ctx, "Accounting", "GLAccountBalanceYearEnd", params)
	if err != nil {
		return nil, err
	}
	var env glAccountBalanceYearEndEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GLAccountBalanceYearEnd response: %w", err)
	}
	return env.Body.Response.Result.Balance.Items, nil
}

func (c *Client) GLAccountTransactions(ctx context.Context, sessionID string, opts GLAccountTransactionsOptions) ([]GLAccountTransaction, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "GLAccountCode", Value: opts.GLAccountCode},
		{Name: "StartDate", Value: opts.StartDate},
		{Name: "EndDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "Accounting", "GLAccountTransactions", params)
	if err != nil {
		return nil, err
	}
	var env glAccountTransactionsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GLAccountTransactions response: %w", err)
	}
	return env.Body.Response.Result.Transactions.Items, nil
}

func (c *Client) NetRevenue(ctx context.Context, sessionID string, opts RevenueOptions) (RevenueReport, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "StartDate", Value: opts.StartDate},
		{Name: "EndDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "Accounting", "NetRevenue", params)
	if err != nil {
		return RevenueReport{}, err
	}
	var env netRevenueEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return RevenueReport{}, fmt.Errorf("parse NetRevenue response: %w", err)
	}
	return RevenueReport{
		AdministrationID: opts.AdministrationID,
		StartDate:        opts.StartDate,
		EndDate:          opts.EndDate,
		Amount:           env.Body.Response.Result,
	}, nil
}

func (c *Client) NetRevenueFiscal(ctx context.Context, sessionID string, opts RevenueOptions) (RevenueReport, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "StartDate", Value: opts.StartDate},
		{Name: "EndDate", Value: opts.EndDate},
	}
	data, err := c.call(ctx, "Accounting", "NetRevenueFiscal", params)
	if err != nil {
		return RevenueReport{}, err
	}
	var env netRevenueFiscalEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return RevenueReport{}, fmt.Errorf("parse NetRevenueFiscal response: %w", err)
	}
	return RevenueReport{
		AdministrationID: opts.AdministrationID,
		StartDate:        opts.StartDate,
		EndDate:          opts.EndDate,
		Amount:           env.Body.Response.Result,
	}, nil
}

type glAccountBalanceEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Balance struct {
					Items []GLAccountBalanceItem `xml:"GLAccount"`
				} `xml:"GLAccountBalance"`
			} `xml:"GLAccountBalanceResult"`
		} `xml:"GLAccountBalanceResponse"`
	} `xml:"Body"`
}

type glAccountTransactionsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Transactions struct {
					Items []GLAccountTransaction `xml:"GLAccountTransaction"`
				} `xml:"GLAccountTransactions"`
			} `xml:"GLAccountTransactionsResult"`
		} `xml:"GLAccountTransactionsResponse"`
	} `xml:"Body"`
}

type netRevenueEnvelope struct {
	Body struct {
		Response struct {
			Result string `xml:"NetRevenueResult"`
		} `xml:"NetRevenueResponse"`
	} `xml:"Body"`
}

type netRevenueFiscalEnvelope struct {
	Body struct {
		Response struct {
			Result string `xml:"NetRevenueFiscalResult"`
		} `xml:"NetRevenueFiscalResponse"`
	} `xml:"Body"`
}

type glAccountBalanceYearEndEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Balance struct {
					Items []GLAccountBalanceItem `xml:"GLAccount"`
				} `xml:"GLAccountBalance"`
			} `xml:"GLAccountBalanceYearEndResult"`
		} `xml:"GLAccountBalanceYearEndResponse"`
	} `xml:"Body"`
}

type glAccountBalanceFiscalEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Balance struct {
					Items []GLAccountBalanceItem `xml:"GLAccount"`
				} `xml:"GLAccountBalance"`
			} `xml:"GLAccountBalanceFiscalResult"`
		} `xml:"GLAccountBalanceFiscalResponse"`
	} `xml:"Body"`
}
