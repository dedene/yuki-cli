package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

type VATReturnListOptions struct {
	AdministrationID string
	Year             int
	ModifiedAfter    string
}

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

func (c *Client) VATReturns(ctx context.Context, sessionID string, opts VATReturnListOptions) ([]VATReturnInfo, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "year", Value: strconv.Itoa(opts.Year)},
		{Name: "modifiedAfter", Value: opts.ModifiedAfter},
	}
	data, err := c.call(ctx, "Vat", "VATReturnList", params)
	if err != nil {
		return nil, err
	}
	var env vatReturnListEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse VATReturnList response: %w", err)
	}
	return env.Body.Response.Result.Returns, nil
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

type vatReturnListEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Returns []VATReturnInfo `xml:"VATReturnInfo"`
			} `xml:"VATReturnListResult"`
		} `xml:"VATReturnListResponse"`
	} `xml:"Body"`
}
