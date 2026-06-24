package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type ArchiveCmd struct {
	Documents      ArchiveDocumentsCmd      `cmd:"" help:"Inspect and download archive documents."`
	Folders        ArchiveFoldersCmd        `cmd:"" help:"Inspect archive folders."`
	Currencies     ArchiveCurrenciesCmd     `cmd:"" help:"Inspect archive currencies."`
	CostCategories ArchiveCostCategoriesCmd `cmd:"" name:"cost-categories" help:"Inspect archive cost categories."`
	Menu           ArchiveMenuCmd           `cmd:"" help:"Inspect archive menu entries."`
	PaymentMethods ArchivePaymentMethodsCmd `cmd:"" name:"payment-methods" help:"Inspect archive payment methods."`
}

type ArchiveMenuCmd struct {
	List ArchiveMenuListCmd `cmd:"" help:"List archive menu entries."`
}

type ArchiveMenuListCmd struct{}

func (c *ArchiveMenuListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	entries, err := client.Menu(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, entries)
	}

	rows := make([][]string, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, []string{entry.ID, entry.Text, entry.Alert, entry.Link, entry.Icon})
	}
	return output.Table(rt.Out, []string{"ID", "TEXT", "ALERT", "LINK", "ICON"}, rows)
}

type ArchiveCostCategoriesCmd struct {
	List ArchiveCostCategoriesListCmd `cmd:"" help:"List archive cost categories."`
}

type ArchiveCostCategoriesListCmd struct{}

func (c *ArchiveCostCategoriesListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	categories, err := client.CostCategories(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, categories)
	}

	rows := make([][]string, 0, len(categories))
	for _, category := range categories {
		rows = append(rows, []string{category.ID, category.Description})
	}
	return output.Table(rt.Out, []string{"ID", "DESCRIPTION"}, rows)
}

type ArchiveCurrenciesCmd struct {
	List ArchiveCurrenciesListCmd `cmd:"" help:"List available archive currencies."`
}

type ArchiveCurrenciesListCmd struct{}

func (c *ArchiveCurrenciesListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	currencies, err := client.Currencies(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, currencies)
	}

	rows := make([][]string, 0, len(currencies))
	for _, currency := range currencies {
		rows = append(rows, []string{
			currency.ID,
			output.Bool(currency.Default),
			currency.Description,
		})
	}
	return output.Table(rt.Out, []string{"ID", "DEFAULT", "DESCRIPTION"}, rows)
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
	List     ArchiveDocumentsListCmd     `cmd:"" help:"List archive documents by document date range."`
	InFolder ArchiveDocumentsInFolderCmd `cmd:"" name:"in-folder" help:"List archive documents in a folder."`
	InTab    ArchiveDocumentsInTabCmd    `cmd:"" name:"in-tab" help:"List archive documents in a tab."`
	Search   ArchiveDocumentsSearchCmd   `cmd:"" help:"Search archive documents."`
	Find     ArchiveDocumentsFindCmd     `cmd:"" help:"Find document metadata by document ID."`
	Download ArchiveDocumentsDownloadCmd `cmd:"" help:"Download a document file by document ID."`
}

type ArchiveDocumentsListCmd struct {
	SortOrder   string `name:"sort-order" default:"CreatedAsc" help:"Yuki document sort order."`
	From        string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To          string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
	Limit       int    `name:"limit" default:"25" help:"Number of records to request."`
	StartRecord int    `name:"start-record" default:"0" help:"Start record to request."`
}

func (c *ArchiveDocumentsListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	documents, err := client.Documents(rt.Context, sessionID, api.DocumentsOptions{
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

type ArchiveDocumentsInFolderCmd struct {
	FolderID    int    `name:"folder" required:"" help:"Folder ID."`
	SortOrder   string `name:"sort-order" default:"CreatedAsc" help:"Yuki document sort order."`
	From        string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To          string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
	Limit       int    `name:"limit" default:"25" help:"Number of records to request."`
	StartRecord int    `name:"start-record" default:"0" help:"Start record to request."`
}

func (c *ArchiveDocumentsInFolderCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	documents, err := client.DocumentsInFolder(rt.Context, sessionID, api.DocumentsInFolderOptions{
		FolderID:        c.FolderID,
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

type ArchiveDocumentsInTabCmd struct {
	TabID       int    `name:"tab" required:"" help:"Tab ID."`
	SortOrder   string `name:"sort-order" default:"CreatedAsc" help:"Yuki document sort order."`
	From        string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To          string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
	Limit       int    `name:"limit" default:"25" help:"Number of records to request."`
	StartRecord int    `name:"start-record" default:"1" help:"Start record to request."`
}

func (c *ArchiveDocumentsInTabCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	documents, err := client.DocumentsInTab(rt.Context, sessionID, api.DocumentsInTabOptions{
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
