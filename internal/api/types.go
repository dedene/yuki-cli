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

type CreditorItem struct {
	ID             string  `json:"id" xml:"ID,attr"`
	Date           string  `json:"date" xml:"Date"`
	Description    string  `json:"description" xml:"Description"`
	Contact        string  `json:"contact" xml:"Contact"`
	ContactID      string  `json:"contact_id,omitempty" xml:"ContactID"`
	OpenAmount     string  `json:"open_amount" xml:"OpenAmount"`
	OriginalAmount string  `json:"original_amount" xml:"OriginalAmount"`
	Type           XMLText `json:"type" xml:"Type"`
	Reference      string  `json:"reference,omitempty" xml:"Reference"`
	DueDate        string  `json:"due_date,omitempty" xml:"DueDate"`
	DocumentID     string  `json:"document_id,omitempty" xml:"DocumentID"`
	PaymentMethod  string  `json:"payment_method,omitempty" xml:"PaymentMethod"`
	ContactCode    string  `json:"contact_code,omitempty" xml:"ContactCode"`
	VATNumber      string  `json:"vat_number,omitempty" xml:"VATNumber"`
	Country        string  `json:"country,omitempty" xml:"Country"`
}

type TransactionInfo struct {
	ID                               string `json:"id" xml:"id"`
	HID                              string `json:"hid,omitempty" xml:"hID"`
	TransactionDate                  string `json:"transaction_date" xml:"transactionDate"`
	Description                      string `json:"description" xml:"description"`
	TransactionAmount                string `json:"transaction_amount" xml:"transactionAmount"`
	TransactionAmountForeignCurrency string `json:"transaction_amount_foreign_currency,omitempty" xml:"transactionAmountForeignCurrency"`
	CurrencyRate                     string `json:"currency_rate,omitempty" xml:"currencyRate"`
	Currency                         string `json:"currency,omitempty" xml:"currency"`
	TaxCodeDescription               string `json:"tax_code_description,omitempty" xml:"taxCodeDescription"`
	FullName                         string `json:"full_name,omitempty" xml:"fullName"`
	CoCNumber                        string `json:"coc_number,omitempty" xml:"CoCNumber"`
	VATNumber                        string `json:"vat_number,omitempty" xml:"VATNumber"`
	ContactID                        string `json:"contact_id,omitempty" xml:"contactID"`
	ContactCountry                   string `json:"contact_country,omitempty" xml:"contactCountry"`
	GLAccountCode                    string `json:"gl_account_code,omitempty" xml:"glAccountCode"`
	DocumentID                       string `json:"document_id,omitempty" xml:"documentID"`
	DocumentReference                string `json:"document_reference,omitempty" xml:"documentReference"`
	DocumentType                     string `json:"document_type,omitempty" xml:"documentType"`
	DocumentFolder                   string `json:"document_folder,omitempty" xml:"documentFolder"`
	DocumentFolderTab                string `json:"document_folder_tab,omitempty" xml:"documentFolderTab"`
	PeriodID                         string `json:"period_id,omitempty" xml:"periodId"`
	Company                          string `json:"company,omitempty" xml:"company"`
	MutationUser                     string `json:"mutation_user,omitempty" xml:"mutationUser"`
}

type TransactionDocument struct {
	FileName string `json:"file_name" xml:"fileName"`
	FileData string `json:"file_data" xml:"filedata"`
}

type Document struct {
	ID              string  `json:"id" xml:"ID,attr"`
	Subject         string  `json:"subject,omitempty" xml:"Subject"`
	DocumentDate    string  `json:"document_date,omitempty" xml:"DocumentDate"`
	Amount          string  `json:"amount,omitempty" xml:"Amount"`
	Folder          XMLText `json:"folder" xml:"Folder"`
	Tab             XMLText `json:"tab" xml:"Tab"`
	Type            string  `json:"type,omitempty" xml:"Type"`
	TypeDescription string  `json:"type_description,omitempty" xml:"TypeDescription"`
	FileName        string  `json:"file_name,omitempty" xml:"FileName"`
	ContentType     string  `json:"content_type,omitempty" xml:"ContentType"`
	FileSize        string  `json:"file_size,omitempty" xml:"FileSize"`
	ContactName     string  `json:"contact_name,omitempty" xml:"ContactName"`
	Created         string  `json:"created,omitempty" xml:"Created"`
	Creator         string  `json:"creator,omitempty" xml:"Creator"`
	Modified        string  `json:"modified,omitempty" xml:"Modified"`
	Modifier        string  `json:"modifier,omitempty" xml:"Modifier"`
}

type DocumentFile struct {
	ID       string `json:"id" xml:"ID,attr"`
	FileName string `json:"file_name" xml:"FileName"`
	FileSize string `json:"file_size,omitempty" xml:"FileSize"`
	FileData string `json:"file_data" xml:"FileData"`
}

type XMLText struct {
	ID   string `json:"id,omitempty" xml:"ID,attr"`
	Text string `json:"text" xml:",chardata"`
}
