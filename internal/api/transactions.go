package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

type TransactionDetailsOptions struct {
	AdministrationID string
	GLAccountCode    string
	StartDate        string
	EndDate          string
	FinancialMode    string
}

type TransactionsOptions struct {
	AdministrationID string
	GLAccountCode    string
	StartDate        string
	EndDate          string
	FinancialMode    string
	DataGroups       string
	NumberOfRecords  int
	StartRecord      int
}

func (c *Client) Transactions(ctx context.Context, sessionID string, opts TransactionsOptions) ([]Transaction, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "glAccountCode", Value: opts.GLAccountCode},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "financialMode", Value: opts.FinancialMode},
		{Name: "dataGroups", Value: opts.DataGroups},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetTransactions", params)
	if err != nil {
		return nil, err
	}
	var env transactionsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetTransactions response: %w", err)
	}
	return env.Body.Response.Result.Transactions, nil
}

func (c *Client) TransactionDetails(ctx context.Context, sessionID string, opts TransactionDetailsOptions) ([]TransactionInfo, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "GLAccountCode", Value: opts.GLAccountCode},
		{Name: "StartDate", Value: opts.StartDate},
		{Name: "EndDate", Value: opts.EndDate},
		{Name: "financialMode", Value: opts.FinancialMode},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetTransactionDetails", params)
	if err != nil {
		return nil, err
	}
	var env transactionDetailsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetTransactionDetails response: %w", err)
	}
	return env.Body.Response.Result.Transactions, nil
}

func (c *Client) TransactionDocument(ctx context.Context, sessionID, administrationID, transactionID string) (TransactionDocument, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
		{Name: "transactionID", Value: transactionID},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetTransactionDocument", params)
	if err != nil {
		return TransactionDocument{}, err
	}
	var env transactionDocumentEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return TransactionDocument{}, fmt.Errorf("parse GetTransactionDocument response: %w", err)
	}
	return env.Body.Response.Result, nil
}

type transactionDetailsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Transactions []TransactionInfo `xml:"TransactionInfo"`
			} `xml:"GetTransactionDetailsResult"`
		} `xml:"GetTransactionDetailsResponse"`
	} `xml:"Body"`
}

type transactionsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Transactions []Transaction `xml:"Transaction"`
			} `xml:"GetTransactionsResult"`
		} `xml:"GetTransactionsResponse"`
	} `xml:"Body"`
}

type transactionDocumentEnvelope struct {
	Body struct {
		Response struct {
			Result TransactionDocument `xml:"GetTransactionDocumentResult"`
		} `xml:"GetTransactionDocumentResponse"`
	} `xml:"Body"`
}
