package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) CheckOutstandingItem(ctx context.Context, sessionID, reference string) ([]OutstandingItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "Reference", Value: reference},
	}
	data, err := c.call(ctx, "Accounting", "CheckOutstandingItem", params)
	if err != nil {
		return nil, err
	}
	var env checkOutstandingItemEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse CheckOutstandingItem response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

func (c *Client) CheckOutstandingItemAdmin(ctx context.Context, sessionID, administrationID, reference string) ([]OutstandingItem, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
		{Name: "Reference", Value: reference},
	}
	data, err := c.call(ctx, "Accounting", "CheckOutstandingItemAdmin", params)
	if err != nil {
		return nil, err
	}
	var env checkOutstandingItemAdminEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse CheckOutstandingItemAdmin response: %w", err)
	}
	return env.Body.Response.Result.Items.Items, nil
}

type checkOutstandingItemEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []OutstandingItem `xml:"Item"`
				} `xml:"OutstandingItems"`
			} `xml:"CheckOutstandingItemResult"`
		} `xml:"CheckOutstandingItemResponse"`
	} `xml:"Body"`
}

type checkOutstandingItemAdminEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Items struct {
					Items []OutstandingItem `xml:"Item"`
				} `xml:"OutstandingItems"`
			} `xml:"CheckOutstandingItemAdminResult"`
		} `xml:"CheckOutstandingItemAdminResponse"`
	} `xml:"Body"`
}
