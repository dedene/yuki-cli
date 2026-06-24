package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

func (c *Client) FiscalTable(ctx context.Context, sessionID, companyID string, year int) (FiscalTableTotals, error) {
	params := []Param{
		{Name: "sessionId", Value: sessionID},
		{Name: "companyId", Value: companyID},
		{Name: "year", Value: strconv.Itoa(year)},
	}
	data, err := c.call(ctx, "FiscalTable", "GetFiscalTable", params)
	if err != nil {
		return FiscalTableTotals{}, err
	}
	var env fiscalTableEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return FiscalTableTotals{}, fmt.Errorf("parse GetFiscalTable response: %w", err)
	}
	result := env.Body.Response.Result
	result.CompanyID = companyID
	result.Year = year
	return result, nil
}

type fiscalTableEnvelope struct {
	Body struct {
		Response struct {
			Result FiscalTableTotals `xml:"GetFiscalTableResult"`
		} `xml:"GetFiscalTableResponse"`
	} `xml:"Body"`
}
