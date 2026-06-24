package api

import (
	"context"
	"fmt"
	"strings"
)

type PettyCashStatementImportOptions struct {
	AdministrationID string
	StatementText    string
}

type PettyCashLineImportOptions struct {
	AccountGLCode          string
	TransactionCode        string
	OffsetAccount          string
	OffsetName             string
	TransactionDate        string
	TransactionDescription string
	Amount                 string
	ProjectCode            string
	ProjectName            string
}

type PettyCashImportResult struct {
	Operation        string `json:"operation"`
	AdministrationID string `json:"administration_id,omitempty"`
	AccountGLCode    string `json:"account_gl_code,omitempty"`
	TransactionDate  string `json:"transaction_date,omitempty"`
	Amount           string `json:"amount,omitempty"`
	ProjectCode      string `json:"project_code,omitempty"`
	ProjectName      string `json:"project_name,omitempty"`
	DocumentID       string `json:"document_id,omitempty"`
	DryRun           bool   `json:"dry_run,omitempty"`
	Message          string `json:"message,omitempty"`
}

func (c *Client) ImportPettyCashStatement(ctx context.Context, sessionID string, opts PettyCashStatementImportOptions) (PettyCashImportResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "statementText", Value: opts.StatementText},
	}
	data, err := c.call(ctx, "PettyCash", "ImportStatement", params)
	if err != nil {
		return PettyCashImportResult{}, err
	}
	documentID, err := optionalTextAt(data, []string{"Envelope", "Body", "ImportStatementResponse", "ImportStatementResult"})
	if err != nil {
		return PettyCashImportResult{}, fmt.Errorf("parse ImportStatement response: %w", err)
	}
	return PettyCashImportResult{
		Operation:        "ImportStatement",
		AdministrationID: opts.AdministrationID,
		DocumentID:       documentID,
		Message:          "petty cash statement import accepted",
	}, nil
}

func (c *Client) ImportPettyCashLine(ctx context.Context, sessionID string, opts PettyCashLineImportOptions) (PettyCashImportResult, error) {
	return c.importPettyCashLine(ctx, "ImportSingleStatementLine", sessionID, opts)
}

func (c *Client) ImportPettyCashProjectLine(ctx context.Context, sessionID string, opts PettyCashLineImportOptions) (PettyCashImportResult, error) {
	return c.importPettyCashLine(ctx, "ImportSingleStatementProjectLine", sessionID, opts)
}

func (c *Client) importPettyCashLine(ctx context.Context, operation, sessionID string, opts PettyCashLineImportOptions) (PettyCashImportResult, error) {
	params := []Param{
		{Name: "sessionId", Value: sessionID},
		{Name: "accountGlCode", Value: opts.AccountGLCode},
		{Name: "transactionCode", Value: opts.TransactionCode},
		{Name: "offsetAccount", Value: opts.OffsetAccount},
		{Name: "offsetName", Value: opts.OffsetName},
		{Name: "transactionDate", Value: opts.TransactionDate},
		{Name: "transactionDescription", Value: opts.TransactionDescription},
		{Name: "amount", Value: opts.Amount},
	}
	if operation == "ImportSingleStatementProjectLine" {
		params = append(params,
			Param{Name: "projectCode", Value: opts.ProjectCode},
			Param{Name: "projectName", Value: opts.ProjectName},
		)
	}
	data, err := c.call(ctx, "PettyCash", operation, params)
	if err != nil {
		return PettyCashImportResult{}, err
	}
	resultElement := operation + "Result"
	documentID, err := optionalTextAt(data, []string{"Envelope", "Body", operation + "Response", resultElement})
	if err != nil {
		return PettyCashImportResult{}, fmt.Errorf("parse %s response: %w", operation, err)
	}
	return PettyCashImportResult{
		Operation:       operation,
		AccountGLCode:   opts.AccountGLCode,
		TransactionDate: opts.TransactionDate,
		Amount:          opts.Amount,
		ProjectCode:     opts.ProjectCode,
		ProjectName:     opts.ProjectName,
		DocumentID:      documentID,
		Message:         "petty cash line import accepted",
	}, nil
}

func optionalTextAt(data []byte, path []string) (string, error) {
	value, err := textAt(data, path)
	if err != nil {
		if strings.HasPrefix(err.Error(), "missing XML path ") {
			return "", nil
		}
		return "", err
	}
	return value, nil
}
