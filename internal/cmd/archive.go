package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/output"
)

type ArchiveCmd struct {
	Documents ArchiveDocumentsCmd `cmd:"" help:"Inspect and download archive documents."`
}

type ArchiveDocumentsCmd struct {
	Find     ArchiveDocumentsFindCmd     `cmd:"" help:"Find document metadata by document ID."`
	Download ArchiveDocumentsDownloadCmd `cmd:"" help:"Download a document file by document ID."`
}

type ArchiveDocumentsFindCmd struct {
	Document string `name:"document" required:"" help:"Document ID."`
}

func (c *ArchiveDocumentsFindCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	document, err := client.FindDocument(rt.Context, sessionID, c.Document)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, document)
	}
	return output.Table(rt.Out, []string{"ID", "SUBJECT", "DATE", "AMOUNT", "TYPE", "FILE", "CONTACT"}, [][]string{{
		document.ID,
		document.Subject,
		document.DocumentDate,
		document.Amount,
		document.TypeDescription,
		document.FileName,
		document.ContactName,
	}})
}

type ArchiveDocumentsDownloadCmd struct {
	Document string `name:"document" required:"" help:"Document ID."`
	Output   string `name:"output" short:"o" help:"Write decoded file bytes to this path."`
}

func (c *ArchiveDocumentsDownloadCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	file, err := client.DocumentFile(rt.Context, sessionID, c.Document)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, file)
	}
	if c.Output == "" {
		return errors.New("missing --output; pass --output <path> or use --json to print the base64 payload")
	}
	return writeBase64File(rt.Out, c.Output, file.FileName, file.FileData)
}
