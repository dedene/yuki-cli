package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

func (c *Client) SearchContacts(ctx context.Context, sessionID string, opts ContactSearchOptions) ([]Contact, error) {
	data, err := c.call(ctx, "Contact", "SearchContacts", contactSearchParams(sessionID, opts, false))
	if err != nil {
		return nil, err
	}
	var env searchContactsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse SearchContacts response: %w", err)
	}
	return env.Body.Response.Result.Contacts.Contacts, nil
}

func (c *Client) SuppliersAndCustomers(ctx context.Context, sessionID string, opts ContactSearchOptions) ([]Contact, error) {
	data, err := c.call(ctx, "Contact", "GetSuppliersAndCustomers", contactSearchParams(sessionID, opts, true))
	if err != nil {
		return nil, err
	}
	var env suppliersAndCustomersEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetSuppliersAndCustomers response: %w", err)
	}
	return env.Body.Response.Result.Contacts.Contacts, nil
}

func contactSearchParams(sessionID string, opts ContactSearchOptions, includeContactType bool) []Param {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domainID", Value: opts.DomainID},
		{Name: "searchOption", Value: opts.SearchOption},
		{Name: "searchValue", Value: opts.SearchValue},
		{Name: "sortOrder", Value: opts.SortOrder},
		{Name: "modifiedAfter", Value: opts.ModifiedAfter},
		{Name: "active", Value: opts.Active},
		{Name: "pageNumber", Value: strconv.Itoa(opts.PageNumber)},
	}
	if includeContactType {
		params = append(params, Param{Name: "contactType", Value: opts.ContactType})
	}
	return params
}

type searchContactsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Contacts struct {
					Contacts []Contact `xml:"Contact"`
				} `xml:"Contacts"`
			} `xml:"SearchContactsResult"`
		} `xml:"SearchContactsResponse"`
	} `xml:"Body"`
}

type suppliersAndCustomersEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Contacts struct {
					Contacts []Contact `xml:"Contact"`
				} `xml:"Contacts"`
			} `xml:"GetSuppliersAndCustomersResult"`
		} `xml:"GetSuppliersAndCustomersResponse"`
	} `xml:"Body"`
}
