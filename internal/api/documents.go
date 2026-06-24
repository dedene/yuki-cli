package api

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
)

type SearchDocumentsOptions struct {
	SearchOption    string
	SearchText      string
	FolderID        int
	TabID           int
	SortOrder       string
	StartDate       string
	EndDate         string
	NumberOfRecords int
	StartRecord     int
}

type DocumentsOptions struct {
	SortOrder       string
	StartDate       string
	EndDate         string
	NumberOfRecords int
	StartRecord     int
}

type DocumentsInFolderOptions struct {
	FolderID        int
	SortOrder       string
	StartDate       string
	EndDate         string
	NumberOfRecords int
	StartRecord     int
}

type DocumentsInTabOptions struct {
	TabID           int
	SortOrder       string
	StartDate       string
	EndDate         string
	NumberOfRecords int
	StartRecord     int
}

type DocumentsByTypeOptions struct {
	DocumentType    int
	SortOrder       string
	StartDate       string
	EndDate         string
	NumberOfRecords int
	StartRecord     int
}

type ModifiedDocumentsInFolderOptions struct {
	FolderID        int
	SortOrder       string
	ModifiedSince   string
	NumberOfRecords int
	StartRecord     int
}

type ModifiedDocumentsByTypeOptions struct {
	DocumentType    int
	SortOrder       string
	ModifiedSince   string
	NumberOfRecords int
	StartRecord     int
}

func (c *Client) DocumentFolders(ctx context.Context, sessionID string) ([]DocumentFolder, error) {
	data, err := c.call(ctx, "Archive", "DocumentFolders", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env documentFoldersEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse DocumentFolders response: %w", err)
	}
	return env.Body.Response.Result.DocumentFolders.Folders, nil
}

func (c *Client) DocumentFolderTabs(ctx context.Context, sessionID, folderID string) ([]DocumentFolderTab, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "folderID", Value: folderID},
	}
	data, err := c.call(ctx, "Archive", "DocumentFolderTabs", params)
	if err != nil {
		return nil, err
	}
	var env documentFolderTabsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse DocumentFolderTabs response: %w", err)
	}
	return env.Body.Response.Result.DocumentFolderTabs.Tabs, nil
}

func (c *Client) Documents(ctx context.Context, sessionID string, opts DocumentsOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "Documents", params)
	if err != nil {
		return nil, err
	}
	var env documentsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Documents response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) DocumentsInFolder(ctx context.Context, sessionID string, opts DocumentsInFolderOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "folderID", Value: strconv.Itoa(opts.FolderID)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "DocumentsInFolder", params)
	if err != nil {
		return nil, err
	}
	var env documentsInFolderEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse DocumentsInFolder response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) DocumentsInTab(ctx context.Context, sessionID string, opts DocumentsInTabOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "tabID", Value: strconv.Itoa(opts.TabID)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "DocumentsInTab", params)
	if err != nil {
		return nil, err
	}
	var env documentsInTabEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse DocumentsInTab response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) DocumentsByType(ctx context.Context, sessionID string, opts DocumentsByTypeOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentType", Value: strconv.Itoa(opts.DocumentType)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "DocumentsByType", params)
	if err != nil {
		return nil, err
	}
	var env documentsByTypeEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse DocumentsByType response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) ModifiedDocumentsInFolder(ctx context.Context, sessionID string, opts ModifiedDocumentsInFolderOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "folderID", Value: strconv.Itoa(opts.FolderID)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "modifiedSince", Value: opts.ModifiedSince},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "ModifiedDocumentsInFolder", params)
	if err != nil {
		return nil, err
	}
	var env modifiedDocumentsInFolderEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse ModifiedDocumentsInFolder response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) ModifiedDocumentsByType(ctx context.Context, sessionID string, opts ModifiedDocumentsByTypeOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentType", Value: strconv.Itoa(opts.DocumentType)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "modifiedSince", Value: opts.ModifiedSince},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "ModifiedDocumentsByType", params)
	if err != nil {
		return nil, err
	}
	var env modifiedDocumentsByTypeEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse ModifiedDocumentsByType response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) SearchDocuments(ctx context.Context, sessionID string, opts SearchDocumentsOptions) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "searchOption", Value: opts.SearchOption},
		{Name: "searchText", Value: opts.SearchText},
		{Name: "folderID", Value: strconv.Itoa(opts.FolderID)},
		{Name: "tabID", Value: strconv.Itoa(opts.TabID)},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "startDate", Value: opts.StartDate},
		{Name: "endDate", Value: opts.EndDate},
		{Name: "numberOfRecords", Value: strconv.Itoa(opts.NumberOfRecords)},
		{Name: "startRecord", Value: strconv.Itoa(opts.StartRecord)},
	}
	data, err := c.call(ctx, "Archive", "SearchDocuments", params)
	if err != nil {
		return nil, err
	}
	var env searchDocumentsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse SearchDocuments response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

func (c *Client) FindDocument(ctx context.Context, sessionID, documentID string) (Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
	}
	data, err := c.call(ctx, "Archive", "FindDocument", params)
	if err != nil {
		return Document{}, err
	}
	var env findDocumentEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return Document{}, fmt.Errorf("parse FindDocument response: %w", err)
	}
	if len(env.Body.Response.Result.Documents.Documents) == 0 {
		return Document{}, errors.New("FindDocument response did not contain a document")
	}
	return env.Body.Response.Result.Documents.Documents[0], nil
}

func (c *Client) DocumentFile(ctx context.Context, sessionID, documentID string) (DocumentFile, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
	}
	data, err := c.call(ctx, "Archive", "DocumentFile", params)
	if err != nil {
		return DocumentFile{}, err
	}
	var env documentFileEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return DocumentFile{}, fmt.Errorf("parse DocumentFile response: %w", err)
	}
	return env.Body.Response.Result.File, nil
}

func (c *Client) DocumentImageCount(ctx context.Context, sessionID, documentID string) (DocumentImageCount, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
	}
	data, err := c.call(ctx, "Archive", "DocumentImageCount", params)
	if err != nil {
		return DocumentImageCount{}, err
	}
	var env documentImageCountEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return DocumentImageCount{}, fmt.Errorf("parse DocumentImageCount response: %w", err)
	}
	return DocumentImageCount{
		DocumentID: documentID,
		ImageCount: env.Body.Response.Result,
	}, nil
}

func (c *Client) PaymentMethods(ctx context.Context, sessionID string) ([]PaymentMethod, error) {
	data, err := c.call(ctx, "Archive", "PaymentMethods", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env archivePaymentMethodsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse PaymentMethods response: %w", err)
	}
	methods := make([]PaymentMethod, 0, len(env.Body.Response.Result.PaymentMethods.Methods))
	for _, method := range env.Body.Response.Result.PaymentMethods.Methods {
		methods = append(methods, PaymentMethod(method))
	}
	return methods, nil
}

func (c *Client) Currencies(ctx context.Context, sessionID string) ([]Currency, error) {
	data, err := c.call(ctx, "Archive", "Currencies", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env currenciesEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Currencies response: %w", err)
	}
	return env.Body.Response.Result.Currencies.Currencies, nil
}

func (c *Client) CostCategories(ctx context.Context, sessionID string) ([]CostCategory, error) {
	data, err := c.call(ctx, "Archive", "CostCategories", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env costCategoriesEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse CostCategories response: %w", err)
	}
	return env.Body.Response.Result.CostCategories.CostCategories, nil
}

func (c *Client) Menu(ctx context.Context, sessionID string) ([]MenuEntry, error) {
	data, err := c.call(ctx, "Archive", "Menu", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env menuEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Menu response: %w", err)
	}
	return env.Body.Response.Result.Menu.Entries, nil
}

type documentFoldersEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				DocumentFolders struct {
					Folders []DocumentFolder `xml:"DocumentFolder"`
				} `xml:"DocumentFolders"`
			} `xml:"DocumentFoldersResult"`
		} `xml:"DocumentFoldersResponse"`
	} `xml:"Body"`
}

type documentFolderTabsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				DocumentFolderTabs struct {
					Tabs []DocumentFolderTab `xml:"DocumentFolderTab"`
				} `xml:"DocumentFolderTabs"`
			} `xml:"DocumentFolderTabsResult"`
		} `xml:"DocumentFolderTabsResponse"`
	} `xml:"Body"`
}

type currenciesEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Currencies struct {
					Currencies []Currency `xml:"Currency"`
				} `xml:"Currencies"`
			} `xml:"CurrenciesResult"`
		} `xml:"CurrenciesResponse"`
	} `xml:"Body"`
}

type costCategoriesEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				CostCategories struct {
					CostCategories []CostCategory `xml:"CostCategory"`
				} `xml:"CostCategories"`
			} `xml:"CostCategoriesResult"`
		} `xml:"CostCategoriesResponse"`
	} `xml:"Body"`
}

type menuEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Menu struct {
					Entries []MenuEntry `xml:"MenuEntry"`
				} `xml:"Menu"`
			} `xml:"MenuResult"`
		} `xml:"MenuResponse"`
	} `xml:"Body"`
}

type documentsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"DocumentsResult"`
		} `xml:"DocumentsResponse"`
	} `xml:"Body"`
}

type documentsInFolderEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"DocumentsInFolderResult"`
		} `xml:"DocumentsInFolderResponse"`
	} `xml:"Body"`
}

type documentsInTabEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"DocumentsInTabResult"`
		} `xml:"DocumentsInTabResponse"`
	} `xml:"Body"`
}

type documentsByTypeEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"DocumentsByTypeResult"`
		} `xml:"DocumentsByTypeResponse"`
	} `xml:"Body"`
}

type modifiedDocumentsInFolderEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"ModifiedDocumentsInFolderResult"`
		} `xml:"ModifiedDocumentsInFolderResponse"`
	} `xml:"Body"`
}

type modifiedDocumentsByTypeEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"ModifiedDocumentsByTypeResult"`
		} `xml:"ModifiedDocumentsByTypeResponse"`
	} `xml:"Body"`
}

type searchDocumentsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"SearchDocumentsResult"`
		} `xml:"SearchDocumentsResponse"`
	} `xml:"Body"`
}

type findDocumentEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"FindDocumentResult"`
		} `xml:"FindDocumentResponse"`
	} `xml:"Body"`
}

type archivePaymentMethodsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				PaymentMethods struct {
					Methods []archivePaymentMethod `xml:"PaymentMethod"`
				} `xml:"PaymentMethods"`
			} `xml:"PaymentMethodsResult"`
		} `xml:"PaymentMethodsResponse"`
	} `xml:"Body"`
}

type archivePaymentMethod struct {
	ID          string `xml:"ID,attr"`
	Description string `xml:"Description"`
}

type documentFileEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				File DocumentFile `xml:"Document"`
			} `xml:"DocumentFileResult"`
		} `xml:"DocumentFileResponse"`
	} `xml:"Body"`
}

type documentImageCountEnvelope struct {
	Body struct {
		Response struct {
			Result int `xml:"DocumentImageCountResult"`
		} `xml:"DocumentImageCountResponse"`
	} `xml:"Body"`
}
