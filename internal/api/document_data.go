package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) DocumentXMLDataAsBinary(ctx context.Context, sessionID, documentID string) (DocumentXMLBinaryData, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
	}
	data, err := c.call(ctx, "Archive", "DocumentXMLDataAsBinary", params)
	if err != nil {
		return DocumentXMLBinaryData{}, err
	}
	var env documentXMLDataAsBinaryEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return DocumentXMLBinaryData{}, fmt.Errorf("parse DocumentXMLDataAsBinary response: %w", err)
	}
	return DocumentXMLBinaryData{
		DocumentID:    documentID,
		XMLDataBase64: env.Body.Response.Result,
	}, nil
}

type documentXMLDataAsBinaryEnvelope struct {
	Body struct {
		Response struct {
			Result string `xml:"DocumentXMLDataAsBinaryResult"`
		} `xml:"DocumentXMLDataAsBinaryResponse"`
	} `xml:"Body"`
}
