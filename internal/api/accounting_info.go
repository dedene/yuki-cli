package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

type PeriodDateTableOptions struct {
	AdministrationID string
	YearID           int
}

func (c *Client) PeriodDateTable(ctx context.Context, sessionID string, opts PeriodDateTableOptions) (AdministrationPeriod, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "yearID", Value: strconv.Itoa(opts.YearID)},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetPeriodDateTable", params)
	if err != nil {
		return AdministrationPeriod{}, err
	}
	var env periodDateTableEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return AdministrationPeriod{}, fmt.Errorf("parse GetPeriodDateTable response: %w", err)
	}
	period := env.Body.Response.Result
	period.AdministrationID = opts.AdministrationID
	period.YearID = opts.YearID
	return period, nil
}

type periodDateTableEnvelope struct {
	Body struct {
		Response struct {
			Result AdministrationPeriod `xml:"GetPeriodDateTableResult"`
		} `xml:"GetPeriodDateTableResponse"`
	} `xml:"Body"`
}
