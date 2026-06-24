package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ArchiveDocumentsUploadBasicCmd struct {
	File           string `name:"file" required:"" help:"Document file to upload." type:"existingfile"`
	FileName       string `name:"file-name" help:"Override the uploaded filename. Defaults to the local basename."`
	FolderID       int    `name:"folder" required:"" help:"Archive folder ID."`
	Administration string `help:"Administration ID. Defaults to profile/global administration when set."`
	DryRun         bool   `name:"dry-run" help:"Preview the upload without authenticating or sending it."`
}

func (c *ArchiveDocumentsUploadBasicCmd) Run(rt *Runtime, globals *Globals) error {
	file, err := readArchiveUploadFile(c.File, c.FileName)
	if err != nil {
		return err
	}
	opts := api.ArchiveUploadOptions{
		FileName:         file.Name,
		DataBase64:       file.DataBase64,
		FolderID:         c.FolderID,
		AdministrationID: optionalAdministrationID(globals, c.Administration),
	}
	if c.DryRun {
		return renderArchiveUploadResult(rt, globals, archiveUploadDryRun("UploadDocument", opts, fmt.Sprintf("dry run; no archive document sent; bytes=%d", file.Bytes)))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: archive documents upload-basic")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.UploadDocument(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderArchiveUploadResult(rt, globals, result)
}

type ArchiveDocumentsUploadDataCmd struct {
	File           string `name:"file" required:"" help:"Document file to upload." type:"existingfile"`
	FileName       string `name:"file-name" help:"Override the uploaded filename. Defaults to the local basename."`
	FolderID       int    `name:"folder" required:"" help:"Archive folder ID."`
	Administration string `help:"Administration ID. Defaults to profile/global administration when set."`
	Currency       string `name:"currency" help:"Document currency, e.g. EUR."`
	Amount         string `name:"amount" default:"0" help:"Invoice or receipt amount. Defaults to 0 per Yuki examples."`
	CostCategory   string `name:"cost-category" help:"Archive cost category ID or code."`
	PaymentMethod  int    `name:"payment-method" default:"0" help:"Archive payment method ID. Defaults to 0 per Yuki examples."`
	Project        string `name:"project" help:"Project code or description."`
	Remarks        string `name:"remarks" help:"Free-form upload remarks."`
	DryRun         bool   `name:"dry-run" help:"Preview the upload without authenticating or sending it."`
}

func (c *ArchiveDocumentsUploadDataCmd) Run(rt *Runtime, globals *Globals) error {
	file, err := readArchiveUploadFile(c.File, c.FileName)
	if err != nil {
		return err
	}
	opts := c.options(file, globals)
	if c.DryRun {
		return renderArchiveUploadResult(rt, globals, archiveUploadDryRun("UploadDocumentWithData", opts, fmt.Sprintf("dry run; no archive document sent; bytes=%d", file.Bytes)))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: archive documents upload-data")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.UploadDocumentWithData(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderArchiveUploadResult(rt, globals, result)
}

func (c *ArchiveDocumentsUploadDataCmd) options(file archiveUploadFile, globals *Globals) api.ArchiveUploadOptions {
	return api.ArchiveUploadOptions{
		FileName:         file.Name,
		DataBase64:       file.DataBase64,
		FolderID:         c.FolderID,
		AdministrationID: optionalAdministrationID(globals, c.Administration),
		Currency:         c.Currency,
		Amount:           c.Amount,
		CostCategory:     c.CostCategory,
		PaymentMethod:    c.PaymentMethod,
		Project:          c.Project,
		Remarks:          c.Remarks,
	}
}

type ArchiveDocumentsUploadAttachmentCmd struct {
	File           string `name:"file" required:"" help:"Primary document file to upload." type:"existingfile"`
	FileName       string `name:"file-name" help:"Override the uploaded primary filename. Defaults to the local basename."`
	Attachment     string `name:"attachment" required:"" help:"Linked attachment file to upload." type:"existingfile"`
	AttachmentName string `name:"attachment-name" help:"Override the uploaded attachment filename. Defaults to the local basename."`
	FolderID       int    `name:"folder" required:"" help:"Archive folder ID."`
	Administration string `help:"Administration ID. Defaults to profile/global administration when set."`
	Currency       string `name:"currency" help:"Document currency, e.g. EUR."`
	Amount         string `name:"amount" default:"0" help:"Invoice or receipt amount. Defaults to 0 per Yuki examples."`
	CostCategory   string `name:"cost-category" help:"Archive cost category ID or code."`
	PaymentMethod  int    `name:"payment-method" default:"0" help:"Archive payment method ID. Defaults to 0 per Yuki examples."`
	Project        string `name:"project" help:"Project code or description."`
	Remarks        string `name:"remarks" help:"Free-form upload remarks."`
	DryRun         bool   `name:"dry-run" help:"Preview the upload without authenticating or sending it."`
}

func (c *ArchiveDocumentsUploadAttachmentCmd) Run(rt *Runtime, globals *Globals) error {
	file, err := readArchiveUploadFile(c.File, c.FileName)
	if err != nil {
		return err
	}
	attachment, err := readArchiveUploadFile(c.Attachment, c.AttachmentName)
	if err != nil {
		return err
	}
	opts := c.options(file, attachment, globals)
	if c.DryRun {
		message := fmt.Sprintf("dry run; no archive document attachment sent; bytes=%d+%d", file.Bytes, attachment.Bytes)
		return renderArchiveUploadResult(rt, globals, archiveAttachmentUploadDryRun(opts, message))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: archive documents upload-attachment")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.UploadDocumentWithAttachment(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderArchiveUploadResult(rt, globals, result)
}

func (c *ArchiveDocumentsUploadAttachmentCmd) options(file, attachment archiveUploadFile, globals *Globals) api.ArchiveAttachmentUploadOptions {
	return api.ArchiveAttachmentUploadOptions{
		FileName1:        file.Name,
		Data1Base64:      file.DataBase64,
		FileName2:        attachment.Name,
		Data2Base64:      attachment.DataBase64,
		FolderID:         c.FolderID,
		AdministrationID: optionalAdministrationID(globals, c.Administration),
		Currency:         c.Currency,
		Amount:           c.Amount,
		CostCategory:     c.CostCategory,
		PaymentMethod:    c.PaymentMethod,
		Project:          c.Project,
		Remarks:          c.Remarks,
	}
}

type archiveUploadFile struct {
	Path       string
	Name       string
	DataBase64 string
	Bytes      int
}

func readArchiveUploadFile(path string, nameOverride string) (archiveUploadFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return archiveUploadFile{}, fmt.Errorf("read %s: %w", path, err)
	}
	if len(data) == 0 {
		return archiveUploadFile{}, fmt.Errorf("validate %s: upload file is empty", path)
	}
	name := nameOverride
	if name == "" {
		name = filepath.Base(path)
	}
	return archiveUploadFile{
		Path:       path,
		Name:       name,
		DataBase64: base64.StdEncoding.EncodeToString(data),
		Bytes:      len(data),
	}, nil
}

func optionalAdministrationID(globals *Globals, explicit string) string {
	if explicit != "" {
		return explicit
	}
	return globals.Administration
}

func archiveUploadDryRun(operation string, opts api.ArchiveUploadOptions, message string) api.ArchiveUploadResult {
	return api.ArchiveUploadResult{
		Operation:        operation,
		FileName:         opts.FileName,
		FolderID:         opts.FolderID,
		AdministrationID: opts.AdministrationID,
		Currency:         opts.Currency,
		Amount:           opts.Amount,
		CostCategory:     opts.CostCategory,
		PaymentMethod:    opts.PaymentMethod,
		Project:          opts.Project,
		Remarks:          opts.Remarks,
		DryRun:           true,
		Message:          message,
	}
}

func archiveAttachmentUploadDryRun(opts api.ArchiveAttachmentUploadOptions, message string) api.ArchiveUploadResult {
	return api.ArchiveUploadResult{
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
		DryRun:           true,
		Message:          message,
	}
}

func renderArchiveUploadResult(rt *Runtime, globals *Globals, result api.ArchiveUploadResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"OPERATION", "FILE", "ATTACHMENT", "FOLDER", "ADMINISTRATION", "AMOUNT", "PAYMENT", "DOCUMENT", "DRY RUN", "MESSAGE"}, [][]string{{
		result.Operation,
		result.FileName,
		result.AttachmentName,
		fmt.Sprint(result.FolderID),
		result.AdministrationID,
		result.Amount,
		fmt.Sprint(result.PaymentMethod),
		result.DocumentID,
		output.Bool(result.DryRun),
		result.Message,
	}})
}
