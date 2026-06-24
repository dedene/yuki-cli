package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

type UpdatedAndDeletedTransactionsOptions struct {
	AdministrationID string
	StartDate        string
	EndDate          string
	NumberOfRecords  int
	StartRecord      int
}

func (c *Client) UpdatedAndDeletedTransactions(ctx context.Context, sessionID string, opts UpdatedAndDeletedTransactionsOptions) ([]UpdatedTransaction, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "ChangeDigest", "GetUpdatedAndDeletedTransactions", params)
	if err != nil {
		return nil, err
	}
	var env updatedAndDeletedTransactionsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetUpdatedAndDeletedTransactions response: %w", err)
	}
	return env.Body.Response.Result.Transactions, nil
}

func (c *Client) ChangeDigestTransactionDetail(ctx context.Context, sessionID, administrationID, transactionID string) (TransactionInfo, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
		{Name: "transactionID", Value: transactionID},
	}
	data, err := c.call(ctx, "ChangeDigest", "GetTransactionDetail", params)
	if err != nil {
		return TransactionInfo{}, err
	}
	var env changeDigestTransactionDetailEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return TransactionInfo{}, fmt.Errorf("parse GetTransactionDetail response: %w", err)
	}
	return env.Body.Response.Result, nil
}

type updatedAndDeletedTransactionsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Transactions []UpdatedTransaction `xml:"UpdatedTransaction"`
			} `xml:"GetUpdatedAndDeletedTransactionsResult"`
		} `xml:"GetUpdatedAndDeletedTransactionsResponse"`
	} `xml:"Body"`
}

type changeDigestTransactionDetailEnvelope struct {
	Body struct {
		Response struct {
			Result TransactionInfo `xml:"GetTransactionDetailResult"`
		} `xml:"GetTransactionDetailResponse"`
	} `xml:"Body"`
}
