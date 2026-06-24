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
