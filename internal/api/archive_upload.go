package api

import (
	"context"
	"fmt"
	"strconv"
)

type ArchiveUploadOptions struct {
	FileName         string
	DataBase64       string
	FolderID         int
	AdministrationID string
	Currency         string
	Amount           string
	CostCategory     string
	PaymentMethod    int
	Project          string
	Remarks          string
}

type ArchiveAttachmentUploadOptions struct {
	FileName1        string
	Data1Base64      string
	FileName2        string
	Data2Base64      string
	FolderID         int
	AdministrationID string
	Currency         string
	Amount           string
	CostCategory     string
	PaymentMethod    int
	Project          string
	Remarks          string
}

type ArchiveUploadResult struct {
	Operation        string `json:"operation"`
	DocumentID       string `json:"document_id,omitempty"`
	FileName         string `json:"file_name,omitempty"`
	AttachmentName   string `json:"attachment_name,omitempty"`
	FolderID         int    `json:"folder_id,omitempty"`
	AdministrationID string `json:"administration_id,omitempty"`
	Currency         string `json:"currency,omitempty"`
	Amount           string `json:"amount,omitempty"`
	CostCategory     string `json:"cost_category,omitempty"`
	PaymentMethod    int    `json:"payment_method"`
	Project          string `json:"project,omitempty"`
	Remarks          string `json:"remarks,omitempty"`
	DryRun           bool   `json:"dry_run,omitempty"`
	Message          string `json:"message,omitempty"`
}

func (c *Client) UploadDocument(ctx context.Context, sessionID string, opts ArchiveUploadOptions) (ArchiveUploadResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "fileName", Value: opts.FileName},
		{Name: "data", Value: opts.DataBase64},
		{Name: "folder", Value: strconv.Itoa(opts.FolderID)},
		{Name: "administrationID", Value: opts.AdministrationID},
	}
	return c.uploadArchiveDocument(ctx, "UploadDocument", params, ArchiveUploadResult{
		Operation:        "UploadDocument",
		FileName:         opts.FileName,
		FolderID:         opts.FolderID,
		AdministrationID: opts.AdministrationID,
		Message:          "archive document upload accepted",
	})
}

func (c *Client) UploadDocumentWithData(ctx context.Context, sessionID string, opts ArchiveUploadOptions) (ArchiveUploadResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "fileName", Value: opts.FileName},
		{Name: "data", Value: opts.DataBase64},
		{Name: "folder", Value: strconv.Itoa(opts.FolderID)},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "currency", Value: opts.Currency},
		{Name: "amount", Value: opts.Amount},
		{Name: "costCategory", Value: opts.CostCategory},
		{Name: "paymentMethod", Value: strconv.Itoa(opts.PaymentMethod)},
		{Name: "project", Value: opts.Project},
		{Name: "remarks", Value: opts.Remarks},
	}
	return c.uploadArchiveDocument(ctx, "UploadDocumentWithData", params, ArchiveUploadResult{
		Operation:        "UploadDocumentWithData",
		FileName:         opts.FileName,
		FolderID:         opts.FolderID,
		AdministrationID: opts.AdministrationID,
		Currency:         opts.Currency,
		Amount:           opts.Amount,
		CostCategory:     opts.CostCategory,
		PaymentMethod:    opts.PaymentMethod,
		Project:          opts.Project,
		Remarks:          opts.Remarks,
		Message:          "archive document upload accepted",
	})
}

func (c *Client) UploadDocumentWithAttachment(ctx context.Context, sessionID string, opts ArchiveAttachmentUploadOptions) (ArchiveUploadResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "fileName1", Value: opts.FileName1},
		{Name: "data1", Value: opts.Data1Base64},
		{Name: "fileName2", Value: opts.FileName2},
		{Name: "data2", Value: opts.Data2Base64},
		{Name: "folder", Value: strconv.Itoa(opts.FolderID)},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "currency", Value: opts.Currency},
		{Name: "amount", Value: opts.Amount},
		{Name: "costCategory", Value: opts.CostCategory},
		{Name: "paymentMethod", Value: strconv.Itoa(opts.PaymentMethod)},
		{Name: "project", Value: opts.Project},
		{Name: "remarks", Value: opts.Remarks},
	}
	return c.uploadArchiveDocument(ctx, "UploadDocumentWithAttachment", params, ArchiveUploadResult{
		Operation:        "UploadDocumentWithAttachment",
		FileName:         opts.FileName1,
		AttachmentName:   opts.FileName2,
		FolderID:         opts.FolderID,
		AdministrationID: opts.AdministrationID,
		Currency:         opts.Currency,
		Amount:           opts.Amount,
		CostCategory:     opts.CostCategory,
		PaymentMethod:    opts.PaymentMethod,
		Project:          opts.Project,
		Remarks:          opts.Remarks,
		Message:          "archive document attachment upload accepted",
	})
}

func (c *Client) uploadArchiveDocument(ctx context.Context, operation string, params []Param, result ArchiveUploadResult) (ArchiveUploadResult, error) {
	data, err := c.call(ctx, "Archive", operation, params)
	if err != nil {
		return ArchiveUploadResult{}, err
	}
	documentID, err := optionalTextAt(data, []string{"Envelope", "Body", operation + "Response", operation + "Result"})
	if err != nil {
		return ArchiveUploadResult{}, fmt.Errorf("parse %s response: %w", operation, err)
	}
	result.DocumentID = documentID
	return result, nil
}
