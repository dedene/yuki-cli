package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

type ContactUpdateOptions struct {
	DomainID string
	XMLDoc   string
}

type ContactUpdateResult struct {
	DomainID  string `json:"domain_id,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Succeeded string `json:"succeeded,omitempty"`
	DryRun    bool   `json:"dry_run,omitempty"`
	Message   string `json:"message,omitempty"`
}

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

func (c *Client) UpdateContact(ctx context.Context, sessionID string, opts ContactUpdateOptions) (ContactUpdateResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domainID", Value: opts.DomainID},
		{Name: "xmlDoc", Value: opts.XMLDoc, Raw: true},
	}
	data, err := c.call(ctx, "Contact", "UpdateContact", params)
	if err != nil {
		return ContactUpdateResult{}, err
	}
	var env updateContactEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return ContactUpdateResult{}, fmt.Errorf("parse UpdateContact response: %w", err)
	}
	result := env.Body.Response.Result.result()
	result.DomainID = opts.DomainID
	if result.Message == "" {
		result.Message = result.Succeeded
	}
	return result, nil
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

type updateContactEnvelope struct {
	Body struct {
		Response struct {
			Result updateContactResultXML `xml:"UpdateContactResult"`
		} `xml:"UpdateContactResponse"`
	} `xml:"Body"`
}

type updateContactResultXML struct {
	ContactResponse contactUpdateResponseXML `xml:"ContactResponse"`
	Timestamp       string                   `xml:"TimeStamp"`
	Succeeded       string                   `xml:"Succeeded"`
}

func (r updateContactResultXML) result() ContactUpdateResult {
	result := ContactUpdateResult{
		Timestamp: r.Timestamp,
		Succeeded: r.Succeeded,
	}
	if result.Timestamp == "" {
		result.Timestamp = r.ContactResponse.Timestamp
	}
	if result.Succeeded == "" {
		result.Succeeded = r.ContactResponse.Succeeded
	}
	return result
}

type contactUpdateResponseXML struct {
	Timestamp string `xml:"TimeStamp"`
	Succeeded string `xml:"Succeeded"`
}
