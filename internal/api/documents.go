package api

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
)

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

type documentFileEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				File DocumentFile `xml:"Document"`
			} `xml:"DocumentFileResult"`
		} `xml:"DocumentFileResponse"`
	} `xml:"Body"`
}
