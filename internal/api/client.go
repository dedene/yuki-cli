package api

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
	userAgent  string
}

const generalService = "AccountingInfo"

func New(cfg Config) *Client {
	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = "yuki/dev"
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
		userAgent:  userAgent,
	}
}

func (c *Client) Authenticate(ctx context.Context, accessKey string) (string, error) {
	data, err := c.call(ctx, generalService, "Authenticate", []Param{{Name: "accessKey", Value: accessKey}})
	if err != nil {
		return "", err
	}
	return textAt(data, []string{"Envelope", "Body", "AuthenticateResponse", "AuthenticateResult"})
}

func (c *Client) Domains(ctx context.Context, sessionID string) ([]Domain, error) {
	data, err := c.call(ctx, generalService, "Domains", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env domainsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Domains response: %w", err)
	}
	return env.Body.Response.Result.Domains.Domains, nil
}

func (c *Client) Companies(ctx context.Context, sessionID string) ([]Company, error) {
	data, err := c.call(ctx, generalService, "Companies", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env companiesEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Companies response: %w", err)
	}
	return env.Body.Response.Result.Companies.Companies, nil
}

func (c *Client) Administrations(ctx context.Context, sessionID string) ([]Administration, error) {
	data, err := c.call(ctx, generalService, "Administrations", sessionParams(sessionID))
	if err != nil {
		return nil, err
	}
	var env administrationsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse Administrations response: %w", err)
	}
	return env.Body.Response.Result.Administrations.Administrations, nil
}

func (c *Client) CurrentDomain(ctx context.Context, sessionID string) (Domain, error) {
	data, err := c.call(ctx, generalService, "GetCurrentDomain", sessionParams(sessionID))
	if err != nil {
		return Domain{}, err
	}
	var env currentDomainEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return Domain{}, fmt.Errorf("parse GetCurrentDomain response: %w", err)
	}
	if len(env.Body.Response.Result.Domains.Domains) == 0 {
		return Domain{}, errors.New("current domain response did not contain a domain")
	}
	return env.Body.Response.Result.Domains.Domains[0], nil
}

func (c *Client) GLAccounts(ctx context.Context, sessionID, administrationID string) ([]GLAccount, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "AccountingInfo", "GetGLAccountScheme", params)
	if err != nil {
		return nil, err
	}
	var env glAccountsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetGLAccountScheme response: %w", err)
	}
	return env.Body.Response.Result.Accounts, nil
}

func (c *Client) call(ctx context.Context, service, operation string, params []Param) ([]byte, error) {
	body := []byte(Envelope(operation, params))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.serviceURL(service), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create %s request: %w", operation, err)
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", SOAPAction(operation))
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s request failed: %w", operation, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read %s response: %w", operation, err)
	}
	if fault := parseFault(data); fault != nil {
		return nil, fault
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s returned HTTP %d: %s", operation, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	return data, nil
}

func (c *Client) serviceURL(service string) string {
	return c.baseURL + "/" + service + ".asmx"
}

func sessionParams(sessionID string) []Param {
	return []Param{{Name: "sessionID", Value: sessionID}}
}

type SOAPFault struct {
	Code   string `xml:"faultcode"`
	String string `xml:"faultstring"`
	Detail string `xml:"detail"`
}

func (f *SOAPFault) Error() string {
	if f.String != "" {
		return f.String
	}
	if f.Detail != "" {
		return f.Detail
	}
	if f.Code != "" {
		return f.Code
	}
	return "SOAP fault"
}

func parseFault(data []byte) *SOAPFault {
	var env faultEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil
	}
	if env.Body.Fault == nil {
		return nil
	}
	return env.Body.Fault
}

func textAt(data []byte, path []string) (string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	var stack []string
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return "", fmt.Errorf("parse XML: %w", err)
		}
		switch t := token.(type) {
		case xml.StartElement:
			stack = append(stack, t.Name.Local)
		case xml.EndElement:
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
		case xml.CharData:
			if samePath(stack, path) {
				return strings.TrimSpace(string(t)), nil
			}
		}
	}
	return "", fmt.Errorf("missing XML path %s", strings.Join(path, ">"))
}

func samePath(got, want []string) bool {
	if len(got) != len(want) {
		return false
	}
	for i := range got {
		if got[i] != want[i] {
			return false
		}
	}
	return true
}

type domainsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Domains struct {
					Domains []Domain `xml:"Domain"`
				} `xml:"Domains"`
			} `xml:"DomainsResult"`
		} `xml:"DomainsResponse"`
	} `xml:"Body"`
}

type companiesEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Companies struct {
					Companies []Company `xml:"Company"`
				} `xml:"Companies"`
			} `xml:"CompaniesResult"`
		} `xml:"CompaniesResponse"`
	} `xml:"Body"`
}

type administrationsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Administrations struct {
					Administrations []Administration `xml:"Administration"`
				} `xml:"Administrations"`
			} `xml:"AdministrationsResult"`
		} `xml:"AdministrationsResponse"`
	} `xml:"Body"`
}

type currentDomainEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Domains struct {
					Domains []Domain `xml:"Domain"`
				} `xml:"Domains"`
			} `xml:"GetCurrentDomainResult"`
		} `xml:"GetCurrentDomainResponse"`
	} `xml:"Body"`
}

type glAccountsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Accounts []GLAccount `xml:"GlAccount"`
			} `xml:"GetGLAccountSchemeResult"`
		} `xml:"GetGLAccountSchemeResponse"`
	} `xml:"Body"`
}

type faultEnvelope struct {
	Body struct {
		Fault *SOAPFault `xml:"Fault"`
	} `xml:"Body"`
}
