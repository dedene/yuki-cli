package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) AdministrationData(ctx context.Context, sessionID, administrationID string) (AdministrationIntegrationData, error) {
	params := []Param{
		{Name: "sessionId", Value: sessionID},
		{Name: "administrationId", Value: administrationID},
	}
	data, err := c.call(ctx, "Integration", "GetAdministrationData", params)
	if err != nil {
		return AdministrationIntegrationData{}, err
	}
	var env administrationDataEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return AdministrationIntegrationData{}, fmt.Errorf("parse GetAdministrationData response: %w", err)
	}
	return env.Body.Response.Result, nil
}

type administrationDataEnvelope struct {
	Body struct {
		Response struct {
			Result AdministrationIntegrationData `xml:"GetAdministrationDataResult"`
		} `xml:"GetAdministrationDataResponse"`
	} `xml:"Body"`
}
