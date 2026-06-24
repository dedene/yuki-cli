package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

type SalesInvoiceImportOptions struct {
	AdministrationID string
	XMLDoc           string
}

type SalesItem struct {
	ID          string `json:"id" xml:"id"`
	Description string `json:"description" xml:"description"`
}

type SalesInvoiceImportResponse struct {
	TimeStamp        string                      `json:"timestamp,omitempty" xml:"TimeStamp"`
	AdministrationID string                      `json:"administration_id,omitempty" xml:"AdministrationId"`
	TotalSucceeded   int                         `json:"total_succeeded" xml:"TotalSucceeded"`
	TotalFailed      int                         `json:"total_failed" xml:"TotalFailed"`
	TotalSkipped     int                         `json:"total_skipped" xml:"TotalSkipped"`
	Invoices         []SalesInvoiceImportInvoice `json:"invoices,omitempty" xml:"Invoice"`
}

type SalesInvoiceImportInvoice struct {
	Succeeded  bool   `json:"succeeded" xml:"Succeeded"`
	Processed  bool   `json:"processed" xml:"Processed"`
	EmailSent  bool   `json:"email_sent" xml:"EmailSent"`
	Reference  string `json:"reference,omitempty" xml:"Reference"`
	Subject    string `json:"subject,omitempty" xml:"Subject"`
	Contact    string `json:"contact,omitempty" xml:"Contact"`
	Message    string `json:"message,omitempty" xml:"Message"`
	DocumentID string `json:"document_id,omitempty" xml:"DocumentId"`
}

func (c *Client) SalesInvoiceSchemaPath(ctx context.Context) (string, error) {
	data, err := c.call(ctx, "Sales", "SalesInvoiceSchemaPath", nil)
	if err != nil {
		return "", err
	}
	return textAt(data, []string{"Envelope", "Body", "SalesInvoiceSchemaPathResponse", "SalesInvoiceSchemaPathResult"})
}

func (c *Client) SalesItems(ctx context.Context, sessionID, administrationID string) ([]SalesItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "Sales", "GetSalesItems", params)
	if err != nil {
		return nil, err
	}
	var env salesItemsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetSalesItems response: %w", err)
	}
	return env.Body.Response.Result.Items, nil
}

func (c *Client) ProcessSalesInvoices(ctx context.Context, sessionID string, opts SalesInvoiceImportOptions) (SalesInvoiceImportResponse, error) {
	return c.processSalesInvoices(ctx, "ProcessSalesInvoices", sessionID, opts)
}

func (c *Client) ProcessRecognizedSalesInvoices(ctx context.Context, sessionID string, opts SalesInvoiceImportOptions) (SalesInvoiceImportResponse, error) {
	return c.processSalesInvoices(ctx, "ProcessRecognizedSalesInvoices", sessionID, opts)
}

func (c *Client) processSalesInvoices(ctx context.Context, operation, sessionID string, opts SalesInvoiceImportOptions) (SalesInvoiceImportResponse, error) {
	params := []Param{
		{Name: "sessionId", Value: sessionID},
		{Name: "administrationId", Value: opts.AdministrationID},
		{Name: "xmlDoc", Value: opts.XMLDoc, Raw: true},
	}
	data, err := c.call(ctx, "Sales", operation, params)
	if err != nil {
		return SalesInvoiceImportResponse{}, err
	}
	response, err := parseSalesInvoiceImportResponse(data, operation)
	if err != nil {
		return SalesInvoiceImportResponse{}, err
	}
	return response, nil
}

func parseSalesInvoiceImportResponse(data []byte, operation string) (SalesInvoiceImportResponse, error) {
	var env salesInvoiceImportEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return SalesInvoiceImportResponse{}, fmt.Errorf("parse %s response: %w", operation, err)
	}
	if operation == "ProcessRecognizedSalesInvoices" {
		return env.Body.ProcessRecognizedSalesInvoicesResponse.Result.Response, nil
	}
	return env.Body.ProcessSalesInvoicesResponse.Result.Response, nil
}

type salesItemsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items []SalesItem `xml:"SalesItem"`
			} `xml:"GetSalesItemsResult"`
		} `xml:"GetSalesItemsResponse"`
	} `xml:"Body"`
}

type salesInvoiceImportEnvelope struct {
	Body struct {
		ProcessSalesInvoicesResponse struct {
			Result struct {
				Response SalesInvoiceImportResponse `xml:"SalesInvoicesImportResponse"`
			} `xml:"ProcessSalesInvoicesResult"`
		} `xml:"ProcessSalesInvoicesResponse"`
		ProcessRecognizedSalesInvoicesResponse struct {
			Result struct {
				Response SalesInvoiceImportResponse `xml:"SalesInvoicesImportResponse"`
			} `xml:"ProcessRecognizedSalesInvoicesResult"`
		} `xml:"ProcessRecognizedSalesInvoicesResponse"`
	} `xml:"Body"`
}
