package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) DocumentBundle(ctx context.Context, sessionID, documentID string) ([]Document, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
	}
	data, err := c.call(ctx, "Archive", "DocumentBundle", params)
	if err != nil {
		return nil, err
	}
	var env documentBundleEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse DocumentBundle response: %w", err)
	}
	return env.Body.Response.Result.Documents.Documents, nil
}

type documentBundleEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Documents struct {
					Documents []Document `xml:"Document"`
				} `xml:"Documents"`
			} `xml:"DocumentBundleResult"`
		} `xml:"DocumentBundleResponse"`
	} `xml:"Body"`
}
