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

type RGSSchemeOptions struct {
	AdministrationID string
	RGSVersion       string
}

type StartBalanceByGLAccountOptions struct {
	AdministrationID string
	Bookyear         int
	FinancialMode    int
}

func (c *Client) RGSScheme(ctx context.Context, sessionID string, opts RGSSchemeOptions) ([]RGSEntry, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "rgsVersion", Value: opts.RGSVersion},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetRGSScheme", params)
	if err != nil {
		return nil, err
	}
	var env rgsSchemeEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetRGSScheme response: %w", err)
	}
	entries := env.Body.Response.Result.Entries
	for i := range entries {
		entries[i].AdministrationID = opts.AdministrationID
		entries[i].RGSVersion = opts.RGSVersion
	}
	return entries, nil
}

func (c *Client) StartBalanceByGLAccount(ctx context.Context, sessionID string, opts StartBalanceByGLAccountOptions) ([]GLAccountStartBalance, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "bookyear", Value: strconv.Itoa(opts.Bookyear)},
		{Name: "financialMode", Value: strconv.Itoa(opts.FinancialMode)},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetStartBalanceByGlAccount", params)
	if err != nil {
		return nil, err
	}
	var env startBalanceByGLAccountEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetStartBalanceByGlAccount response: %w", err)
	}
	balances := env.Body.Response.Result.Balances
	for i := range balances {
		balances[i].AdministrationID = opts.AdministrationID
		balances[i].Bookyear = opts.Bookyear
		balances[i].FinancialMode = opts.FinancialMode
	}
	return balances, nil
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

func (c *Client) FinancialYearModifiedDate(ctx context.Context, sessionID string, opts PeriodDateTableOptions) (FinancialYearModifiedDate, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "yearID", Value: strconv.Itoa(opts.YearID)},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetFinancialYearModifiedDate", params)
	if err != nil {
		return FinancialYearModifiedDate{}, err
	}
	modifiedDate, err := textAt(data, []string{"Envelope", "Body", "GetFinancialYearModifiedDateResponse", "GetFinancialYearModifiedDateResult"})
	if err != nil {
		return FinancialYearModifiedDate{}, err
	}
	return FinancialYearModifiedDate{
		AdministrationID: opts.AdministrationID,
		YearID:           opts.YearID,
		ModifiedDate:     modifiedDate,
	}, nil
}

func (c *Client) ContactDefaultValues(ctx context.Context, sessionID, administrationID, contactID string) ([]ContactDefaultValues, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
		{Name: "contactID", Value: contactID},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetContactDefaultValues", params)
	if err != nil {
		return nil, err
	}
	var env contactDefaultValuesEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetContactDefaultValues response: %w", err)
	}
	return env.Body.Response.Result.Values, nil
}

type periodDateTableEnvelope struct {
	Body struct {
		Response struct {
			Result AdministrationPeriod `xml:"GetPeriodDateTableResult"`
		} `xml:"GetPeriodDateTableResponse"`
	} `xml:"Body"`
}

type contactDefaultValuesEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Values []ContactDefaultValues `xml:"ContactDefaultValues"`
			} `xml:"GetContactDefaultValuesResult"`
		} `xml:"GetContactDefaultValuesResponse"`
	} `xml:"Body"`
}

type rgsSchemeEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Entries []RGSEntry `xml:"RGSEntry"`
			} `xml:"GetRGSSchemeResult"`
		} `xml:"GetRGSSchemeResponse"`
	} `xml:"Body"`
}

type startBalanceByGLAccountEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Balances []GLAccountStartBalance `xml:"AccountStartBalance"`
			} `xml:"GetStartBalanceByGlAccountResult"`
		} `xml:"GetStartBalanceByGlAccountResponse"`
	} `xml:"Body"`
}
