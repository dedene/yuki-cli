package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ArchiveCmd struct {
	Documents      ArchiveDocumentsCmd      `cmd:"" help:"Inspect and download archive documents."`
	Folders        ArchiveFoldersCmd        `cmd:"" help:"Inspect archive folders."`
	PaymentMethods ArchivePaymentMethodsCmd `cmd:"" name:"payment-methods" help:"Inspect archive payment methods."`
}

type ArchiveFoldersCmd struct {
	List ArchiveFoldersListCmd `cmd:"" help:"List archive document folders."`
	Tabs ArchiveFoldersTabsCmd `cmd:"" help:"List tabs for an archive document folder."`
}

type ArchiveFoldersListCmd struct{}

func (c *ArchiveFoldersListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	folders, err := client.DocumentFolders(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, folders)
	}

	rows := make([][]string, 0, len(folders))
	for _, folder := range folders {
		rows = append(rows, []string{
			folder.ID,
			folder.Description,
			output.Bool(folder.ProcessedByYuki),
			folder.Icon,
		})
	}
	return output.Table(rt.Out, []string{"ID", "DESCRIPTION", "PROCESSED", "ICON"}, rows)
}

type ArchiveFoldersTabsCmd struct {
	Folder string `name:"folder" required:"" help:"Folder ID."`
}

func (c *ArchiveFoldersTabsCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	tabs, err := client.DocumentFolderTabs(rt.Context, sessionID, c.Folder)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, tabs)
	}

	rows := make([][]string, 0, len(tabs))
	for _, tab := range tabs {
		rows = append(rows, []string{
			tab.ID,
			tab.Description,
			output.Bool(tab.ProcessedByYuki),
		})
	}
	return output.Table(rt.Out, []string{"ID", "DESCRIPTION", "PROCESSED"}, rows)
}

type ArchiveDocumentsCmd struct {
	Search   ArchiveDocumentsSearchCmd   `cmd:"" help:"Search archive documents."`
	Find     ArchiveDocumentsFindCmd     `cmd:"" help:"Find document metadata by document ID."`
	Download ArchiveDocumentsDownloadCmd `cmd:"" help:"Download a document file by document ID."`
}

type ArchiveDocumentsSearchCmd struct {
	SearchOption string `name:"search-option" default:"All" help:"Yuki search option: All, Creator, Contact, Subject, Tag, or Type."`
	SearchText   string `name:"search-text" help:"Search text."`
	FolderID     int    `name:"folder" default:"-1" help:"Folder ID, or -1 for all folders."`
	TabID        int    `name:"tab" default:"-1" help:"Tab ID, or -1 for all tabs."`
	SortOrder    string `name:"sort-order" default:"CreatedDesc" help:"Yuki document sort order."`
	From         string `name:"from" default:"0001-01-01" help:"Start date, YYYY-MM-DD. Use 0001-01-01 for all years."`
	To           string `name:"to" default:"0001-01-01" help:"End date, YYYY-MM-DD. Use 0001-01-01 for all years."`
	Limit        int    `name:"limit" default:"25" help:"Number of records to request."`
	StartRecord  int    `name:"start-record" default:"1" help:"One-based start record."`
}

func (c *ArchiveDocumentsSearchCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	documents, err := client.SearchDocuments(rt.Context, sessionID, api.SearchDocumentsOptions{
		SearchOption:    c.SearchOption,
		SearchText:      c.SearchText,
		FolderID:        c.FolderID,
		TabID:           c.TabID,
		SortOrder:       c.SortOrder,
		StartDate:       c.From,
		EndDate:         c.To,
		NumberOfRecords: c.Limit,
		StartRecord:     c.StartRecord,
	})
	if err != nil {
		return err
	}
	return renderDocuments(rt, globals, documents)
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
	return renderDocuments(rt, globals, []api.Document{document})
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

type ArchivePaymentMethodsCmd struct {
	List ArchivePaymentMethodsListCmd `cmd:"" help:"List archive payment methods."`
}

type ArchivePaymentMethodsListCmd struct{}

func (c *ArchivePaymentMethodsListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	methods, err := client.PaymentMethods(rt.Context, sessionID)
	if err != nil {
		return err
	}
	return renderPaymentMethods(rt, globals, methods)
}

func renderDocuments(rt *Runtime, globals *Globals, documents []api.Document) error {
	if globals.JSON {
		return output.JSON(rt.Out, documents)
	}

	rows := make([][]string, 0, len(documents))
	for _, document := range documents {
		rows = append(rows, []string{
			document.ID,
			document.Subject,
			document.DocumentDate,
			document.Amount,
			document.TypeDescription,
			document.FileName,
			document.ContactName,
		})
	}
	return output.Table(rt.Out, []string{"ID", "SUBJECT", "DATE", "AMOUNT", "TYPE", "FILE", "CONTACT"}, rows)
}
