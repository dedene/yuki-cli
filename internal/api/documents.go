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
