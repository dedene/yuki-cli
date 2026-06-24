package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) CustomPaymentMethods(ctx context.Context, sessionID, administrationID string) ([]PaymentMethod, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetCustomPaymentMethods", params)
	if err != nil {
		return nil, err
	}
	var env customPaymentMethodsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetCustomPaymentMethods response: %w", err)
	}
	if len(env.Body.CustomResponse.Result.Methods) > 0 {
		return env.Body.CustomResponse.Result.Methods, nil
	}
	return env.Body.GetPaymentMethodsResponse.Result.Methods, nil
}

type customPaymentMethodsEnvelope struct {
	Body struct {
		CustomResponse struct {
			Result struct {
				Methods []PaymentMethod `xml:"PaymentMethod"`
			} `xml:"GetCustomPaymentMethodsResult"`
		} `xml:"GetCustomPaymentMethodsResponse"`
		GetPaymentMethodsResponse struct {
			Result struct {
				Methods []PaymentMethod `xml:"PaymentMethod"`
			} `xml:"GetPaymentMethodsResult"`
		} `xml:"GetPaymentMethodsResponse"`
	} `xml:"Body"`
}
