package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strings"
)

func (c *Client) DomainName(ctx context.Context, sessionID, administrationName string) (DomainNameResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationName", Value: administrationName},
	}
	data, err := c.call(ctx, "Domains", "GetDomainName", params)
	if err != nil {
		return DomainNameResult{}, err
	}
	domainName, err := textAt(data, []string{"Envelope", "Body", "GetDomainNameResponse", "GetDomainNameResult"})
	if err != nil {
		return DomainNameResult{}, fmt.Errorf("parse GetDomainName response: %w", err)
	}
	return DomainNameResult{
		AdministrationName: administrationName,
		DomainName:         domainName,
	}, nil
}

func (c *Client) DomainUsers(ctx context.Context, sessionID, domainID string) ([]DomainUser, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domain", Value: domainID},
	}
	data, err := c.call(ctx, "Domains", "GetDomainUsers", params)
	if err != nil {
		return nil, err
	}
	var env domainUsersEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetDomainUsers response: %w", err)
	}
	return env.Body.Response.Result.users(), nil
}

func (c *Client) CreateDomain(ctx context.Context, sessionID string, opts DomainCreateOptions) (DomainAdminResult, error) {
	return c.createDomain(ctx, "CreateDomain", sessionID, opts)
}

func (c *Client) CreateTrialDomain(ctx context.Context, sessionID string, opts DomainCreateOptions) (DomainAdminResult, error) {
	return c.createDomain(ctx, "CreateTrialDomain", sessionID, opts)
}

func (c *Client) createDomain(ctx context.Context, operation, sessionID string, opts DomainCreateOptions) (DomainAdminResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationName", Value: opts.AdministrationName},
		{Name: "domainName", Value: opts.DomainName},
		{Name: "defaultLanguage", Value: opts.DefaultLanguage},
	}
	data, err := c.call(ctx, "Domains", operation, params)
	if err != nil {
		return DomainAdminResult{}, err
	}
	message, err := optionalTextAt(data, []string{"Envelope", "Body", operation + "Response", operation + "Result"})
	if err != nil {
		return DomainAdminResult{}, fmt.Errorf("parse %s response: %w", operation, err)
	}
	return DomainAdminResult{
		Operation:          operation,
		AdministrationName: opts.AdministrationName,
		DomainName:         opts.DomainName,
		DefaultLanguage:    opts.DefaultLanguage,
		Message:            message,
	}, nil
}

func (c *Client) AddDomainUser(ctx context.Context, sessionID string, opts DomainUserAddOptions) (DomainAdminResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domain", Value: opts.DomainID},
		{Name: "email", Value: opts.Email},
		{Name: "name", Value: opts.Name},
		{Name: "roles", Value: opts.Roles},
		{Name: "administrations", Value: opts.Administrations},
		{Name: "sendMessage", Value: boolString(opts.SendMessage)},
		{Name: "customMessage", Value: opts.CustomMessage},
		{Name: "language", Value: opts.Language},
	}
	data, err := c.call(ctx, "Domains", "AddDomainUser", params)
	if err != nil {
		return DomainAdminResult{}, err
	}
	message, err := optionalTextAt(data, []string{"Envelope", "Body", "AddDomainUserResponse", "AddDomainUserResult"})
	if err != nil {
		return DomainAdminResult{}, fmt.Errorf("parse AddDomainUser response: %w", err)
	}
	return DomainAdminResult{
		Operation:       "AddDomainUser",
		DomainID:        opts.DomainID,
		Email:           opts.Email,
		Name:            opts.Name,
		Roles:           opts.Roles,
		Administrations: opts.Administrations,
		SendMessage:     &opts.SendMessage,
		CustomMessage:   opts.CustomMessage,
		Language:        opts.Language,
		Message:         message,
	}, nil
}

func (c *Client) SetLyantheRecognitionEngine(ctx context.Context, sessionID, domainID string, enabled bool) (DomainAdminResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domain", Value: domainID},
		{Name: "enableLyanthe", Value: boolString(enabled)},
	}
	data, err := c.call(ctx, "Domains", "LyantheRecognitionEngine", params)
	if err != nil {
		return DomainAdminResult{}, err
	}
	message, err := textAt(data, []string{"Envelope", "Body", "LyantheRecognitionEngineResponse", "LyantheRecognitionEngineResult"})
	if err != nil {
		return DomainAdminResult{}, fmt.Errorf("parse LyantheRecognitionEngine response: %w", err)
	}
	return DomainAdminResult{
		Operation: "LyantheRecognitionEngine",
		DomainID:  domainID,
		Enabled:   &enabled,
		Message:   message,
	}, nil
}

func (c *Client) DomainFunctions(ctx context.Context, sessionID, domainID string) ([]DomainFunctionAssignment, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domain", Value: domainID},
	}
	data, err := c.call(ctx, "Domains", "GetDomainFunctions", params)
	if err != nil {
		return nil, err
	}
	var env domainFunctionsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetDomainFunctions response: %w", err)
	}
	assignments := env.Body.Response.Result.Functions.assignments()
	for i := range assignments {
		assignments[i].DomainID = domainID
	}
	return assignments, nil
}

func (c *Client) UpdateDomainFunction(ctx context.Context, sessionID string, opts UpdateDomainFunctionOptions) (DomainFunctionUpdateResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "domain", Value: opts.DomainID},
		{Name: "domainFunction", Value: opts.Function},
		{Name: "login", Value: opts.Login},
	}
	data, err := c.call(ctx, "Domains", "UpdateDomainFunctions", params)
	if err != nil {
		return DomainFunctionUpdateResult{}, err
	}
	message, err := textAt(data, []string{"Envelope", "Body", "UpdateDomainFunctionsResponse", "UpdateDomainFunctionsResult"})
	if err != nil {
		return DomainFunctionUpdateResult{}, fmt.Errorf("parse UpdateDomainFunctions response: %w", err)
	}
	return DomainFunctionUpdateResult{
		DomainID: opts.DomainID,
		Function: opts.Function,
		Login:    opts.Login,
		Message:  message,
	}, nil
}

type domainFunctionsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Functions domainFunctions `xml:"DomainFunctions"`
			} `xml:"GetDomainFunctionsResult"`
		} `xml:"GetDomainFunctionsResponse"`
	} `xml:"Body"`
}

type domainUsersEnvelope struct {
	Body struct {
		Response struct {
			Result domainUsersResult `xml:"GetDomainUsersResult"`
		} `xml:"GetDomainUsersResponse"`
	} `xml:"Body"`
}

type domainUsersResult struct {
	DirectUsers       []DomainUser `xml:"User"`
	DirectDomainUsers []DomainUser `xml:"DomainUser"`
	DomainUsers       struct {
		Users       []DomainUser `xml:"User"`
		DomainUsers []DomainUser `xml:"DomainUser"`
	} `xml:"DomainUsers"`
	Users struct {
		Users       []DomainUser `xml:"User"`
		DomainUsers []DomainUser `xml:"DomainUser"`
	} `xml:"Users"`
}

func (r domainUsersResult) users() []DomainUser {
	users := make([]DomainUser, 0, len(r.DirectUsers)+len(r.DirectDomainUsers)+len(r.DomainUsers.Users)+len(r.DomainUsers.DomainUsers)+len(r.Users.Users)+len(r.Users.DomainUsers))
	users = append(users, r.DirectUsers...)
	users = append(users, r.DirectDomainUsers...)
	users = append(users, r.DomainUsers.Users...)
	users = append(users, r.DomainUsers.DomainUsers...)
	users = append(users, r.Users.Users...)
	users = append(users, r.Users.DomainUsers...)
	return users
}

func (u *DomainUser) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	*u = DomainUser{}
	for _, attr := range start.Attr {
		if strings.EqualFold(attr.Name.Local, "ID") || strings.EqualFold(attr.Name.Local, "id") {
			u.ID = attr.Value
		}
	}
	for {
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		switch token := token.(type) {
		case xml.StartElement:
			var value string
			if err := decoder.DecodeElement(&value, &token); err != nil {
				return err
			}
			value = strings.TrimSpace(value)
			name := token.Name.Local
			u.Fields = append(u.Fields, XMLField{Name: name, Value: value})
			u.assignField(name, value)
		case xml.EndElement:
			if token.Name.Local == start.Name.Local {
				return nil
			}
		}
	}
}

func (u *DomainUser) assignField(name, value string) {
	switch strings.ToLower(name) {
	case "id":
		if u.ID == "" {
			u.ID = value
		}
	case "name":
		u.Name = value
	case "fullname", "full_name", "full name":
		u.FullName = value
	case "login":
		u.Login = value
	case "email", "e-mail", "emailaddress", "emailaddresswork":
		u.Email = value
	case "language":
		u.Language = value
	case "roles":
		u.Roles = value
	case "administrations":
		u.Administrations = value
	case "active", "isactive":
		u.Active = value
	}
}

type domainFunctions struct {
	BOResponsible    domainFunctionUser `xml:"BOResponsible"`
	BOBackup         domainFunctionUser `xml:"BOBackup"`
	BOController     domainFunctionUser `xml:"BOController"`
	BOAccountManager domainFunctionUser `xml:"BOAccountManager"`
}

func (f domainFunctions) assignments() []DomainFunctionAssignment {
	return []DomainFunctionAssignment{
		f.BOResponsible.assignment("BOResponsible"),
		f.BOBackup.assignment("BOBackup"),
		f.BOController.assignment("BOController"),
		f.BOAccountManager.assignment("BOAccountManager"),
	}
}

type domainFunctionUser struct {
	FullName string `xml:"FullName"`
	Login    string `xml:"Login"`
}

func (u domainFunctionUser) assignment(function string) DomainFunctionAssignment {
	return DomainFunctionAssignment{
		Function: function,
		FullName: u.FullName,
		Login:    u.Login,
	}
}
