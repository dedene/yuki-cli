package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
)

func (c *Client) DocumentXMLData(ctx context.Context, sessionID, documentID string) (DocumentXMLData, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
	}
	data, err := c.call(ctx, "Archive", "DocumentXMLData", params)
	if err != nil {
		return DocumentXMLData{}, err
	}
	var env documentXMLDataEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return DocumentXMLData{}, fmt.Errorf("parse DocumentXMLData response: %w", err)
	}
	return DocumentXMLData{
		DocumentID: documentID,
		XML:        env.Body.Response.Result.xml(),
	}, nil
}

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

type documentXMLDataEnvelope struct {
	Body struct {
		Response struct {
			Result documentXMLDataResult `xml:"DocumentXMLDataResult"`
		} `xml:"DocumentXMLDataResponse"`
	} `xml:"Body"`
}

type documentXMLDataResult struct {
	InnerXML string `xml:",innerxml"`
	Text     string `xml:",chardata"`
}

func (r documentXMLDataResult) xml() string {
	innerXML := strings.TrimSpace(r.InnerXML)
	if innerXML != "" {
		return innerXML
	}
	return strings.TrimSpace(r.Text)
}

type documentXMLDataAsBinaryEnvelope struct {
	Body struct {
		Response struct {
			Result string `xml:"DocumentXMLDataAsBinaryResult"`
		} `xml:"DocumentXMLDataAsBinaryResponse"`
	} `xml:"Body"`
}
