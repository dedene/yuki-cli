package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) ActiveVATCodes(ctx context.Context, sessionID, administrationID string) ([]VATCode, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "Vat", "ActiveVATCodesList", params)
	if err != nil {
		return nil, err
	}
	var env activeVATCodesEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse ActiveVATCodesList response: %w", err)
	}
	return env.Body.Response.Result.Codes, nil
}

type activeVATCodesEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Codes []VATCode `xml:"VATCode"`
			} `xml:"ActiveVATCodesListResult"`
		} `xml:"ActiveVATCodesListResponse"`
	} `xml:"Body"`
}
