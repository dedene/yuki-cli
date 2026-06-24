package api

import "net/http"

const (
	DefaultBaseURL = "https://api.yukiworks.be/ws"
	Namespace      = "http://www.theyukicompany.com/"
)

type Config struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

type Param struct {
	Name  string
	Value string
}

type Domain struct {
	ID   string `json:"id" xml:"ID,attr"`
	Name string `json:"name" xml:"Name"`
	URL  string `json:"url,omitempty" xml:"URL"`
}

type Company struct {
	ID     string `json:"id" xml:"ID,attr"`
	Name   string `json:"name" xml:"Name"`
	Active bool   `json:"active" xml:"Active"`
}

type Administration struct {
	ID          string `json:"id" xml:"ID,attr"`
	Name        string `json:"name" xml:"Name"`
	AddressLine string `json:"address_line,omitempty" xml:"AddressLine"`
	Postcode    string `json:"postcode,omitempty" xml:"Postcode"`
	City        string `json:"city,omitempty" xml:"City"`
	Country     string `json:"country,omitempty" xml:"Country"`
	CoCNumber   string `json:"coc_number,omitempty" xml:"CoCNumber"`
	VATNumber   string `json:"vat_number,omitempty" xml:"VATNumber"`
	StartDate   string `json:"start_date,omitempty" xml:"StartDate"`
	DomainID    string `json:"domain_id,omitempty" xml:"DomainID"`
	Active      bool   `json:"active" xml:"Active"`
}

type GLAccount struct {
	Code        string `json:"code" xml:"code"`
	Type        string `json:"type" xml:"type"`
	Subtype     string `json:"subtype" xml:"subtype"`
	Enabled     bool   `json:"enabled" xml:"isEnabled"`
	Description string `json:"description" xml:"descripton"`
}
