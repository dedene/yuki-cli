package api

import (
	"encoding/xml"
	"net/http"
)

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

type RGSEntry struct {
	YukiCode         string `json:"yuki_code" xml:"YukiCode"`
	YukiIsEnabled    string `json:"yuki_is_enabled" xml:"YukiIsEnabled"`
	YukiDescription  string `json:"yuki_description" xml:"YukiDescription"`
	RGSReferenceCode string `json:"rgs_reference_code,omitempty" xml:"RgsReferenceCode"`
	RGSDescription   string `json:"rgs_description,omitempty" xml:"RgsDescription"`
	RGSFlipCode      string `json:"rgs_flip_code,omitempty" xml:"RgsFlipCode"`
	AdministrationID string `json:"administration_id,omitempty"`
	RGSVersion       string `json:"rgs_version,omitempty"`
}

type GLAccountStartBalance struct {
	AccountID          string `json:"account_id" xml:"accountID"`
	StartBalance       string `json:"start_balance" xml:"startBalance"`
	AccountDescription string `json:"account_description" xml:"accountDescription"`
	AdministrationID   string `json:"administration_id,omitempty"`
	Bookyear           int    `json:"bookyear,omitempty"`
	FinancialMode      int    `json:"financial_mode,omitempty"`
}

type GLAccountBalanceItem struct {
	Code        string `json:"code" xml:"Code,attr"`
	BalanceType string `json:"balance_type" xml:"BalanceType,attr"`
	Description string `json:"description" xml:"Description"`
	Amount      string `json:"amount" xml:"Amount"`
}

type GLAccountTransaction struct {
	ID            string               `json:"id" xml:"ID,attr"`
	Date          string               `json:"date" xml:"Date"`
	Description   string               `json:"description" xml:"Description"`
	Amount        string               `json:"amount" xml:"Amount"`
	SalesItem     string               `json:"sales_item,omitempty" xml:"SalesItem"`
	Contact       string               `json:"contact,omitempty" xml:"Contact"`
	ContactID     string               `json:"contact_id,omitempty" xml:"ContactID"`
	Project       GLTransactionProject `json:"project" xml:"Project"`
	GLAccountCode string               `json:"gl_account_code" xml:"GLAccountCode"`
	FileName      string               `json:"file_name,omitempty" xml:"FileName"`
}

type GLTransactionProject struct {
	Code string `json:"code,omitempty" xml:"Code,attr"`
	Text string `json:"text,omitempty" xml:",chardata"`
}

type RevenueReport struct {
	AdministrationID string `json:"administration_id,omitempty"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	Amount           string `json:"amount"`
}

type AdministrationPeriod struct {
	AdministrationID string `json:"administration_id,omitempty"`
	YearID           int    `json:"year_id,omitempty"`
	Name             string `json:"name" xml:"name"`
	Period           string `json:"period" xml:"period"`
	WholePeriod      string `json:"whole_period,omitempty" xml:"wholePeriod"`
	ISO8601Period    bool   `json:"iso8601_period" xml:"ISO8601Period"`
}

type OutstandingItem struct {
	ID               string  `json:"id" xml:"ID,attr"`
	Date             string  `json:"date" xml:"Date"`
	Description      string  `json:"description" xml:"Description"`
	Contact          string  `json:"contact" xml:"Contact"`
	ContactID        string  `json:"contact_id,omitempty" xml:"ContactID"`
	OpenAmount       string  `json:"open_amount" xml:"OpenAmount"`
	OriginalAmount   string  `json:"original_amount" xml:"OriginalAmount"`
	Type             XMLText `json:"type" xml:"Type"`
	Reference        string  `json:"reference,omitempty" xml:"Reference"`
	PaymentReference string  `json:"payment_reference,omitempty" xml:"PaymentReference"`
	DueDate          string  `json:"due_date,omitempty" xml:"DueDate"`
	DocumentID       string  `json:"document_id,omitempty" xml:"DocumentID"`
	PaymentMethod    string  `json:"payment_method,omitempty" xml:"PaymentMethod"`
	ContactCode      string  `json:"contact_code,omitempty" xml:"ContactCode"`
	CoCNumber        string  `json:"coc_number,omitempty" xml:"CoCNumber"`
	VATNumber        string  `json:"vat_number,omitempty" xml:"VATNumber"`
	AddressLine1     string  `json:"address_line_1,omitempty" xml:"AddressLine_1"`
	AddressLine2     string  `json:"address_line_2,omitempty" xml:"AddressLine_2"`
	Postcode         string  `json:"postcode,omitempty" xml:"Postcode"`
	City             string  `json:"city,omitempty" xml:"City"`
	MailAddressLine1 string  `json:"mail_address_line_1,omitempty" xml:"MailAddressLine_1"`
	MailAddressLine2 string  `json:"mail_address_line_2,omitempty" xml:"MailAddressLine_2"`
	MailPostcode     string  `json:"mail_postcode,omitempty" xml:"MailPostcode"`
	MailCity         string  `json:"mail_city,omitempty" xml:"MailCity"`
	Country          string  `json:"country,omitempty" xml:"Country"`
	RecipientEmail   string  `json:"recipient_email,omitempty" xml:"RecipientEmail"`
	PhoneHome        string  `json:"phone_home,omitempty" xml:"PhoneHome"`
	PhoneWork        string  `json:"phone_work,omitempty" xml:"PhoneWork"`
	EmailHome        string  `json:"email_home,omitempty" xml:"EmailHome"`
	EmailWork        string  `json:"email_work,omitempty" xml:"EmailWork"`
}

type (
	CreditorItem = OutstandingItem
	DebtorItem   = OutstandingItem
)

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
	ContactHID                       string `json:"contact_hid,omitempty" xml:"contactHID"`
	ContactID                        string `json:"contact_id,omitempty" xml:"contactID"`
	ContactCode                      string `json:"contact_code,omitempty" xml:"contactCode"`
	ContactCountry                   string `json:"contact_country,omitempty" xml:"contactCountry"`
	ContactZipCode                   string `json:"contact_zip_code,omitempty" xml:"contactZipCode"`
	Project                          string `json:"project,omitempty" xml:"project"`
	ProjectCode                      string `json:"project_code,omitempty" xml:"projectCode"`
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

type UpdatedTransaction struct {
	ID                               string `json:"id" xml:"id"`
	TransactionDate                  string `json:"transaction_date" xml:"transactionDate"`
	Description                      string `json:"description,omitempty" xml:"description"`
	TransactionAmount                string `json:"transaction_amount" xml:"transactionAmount"`
	TransactionAmountForeignCurrency string `json:"transaction_amount_foreign_currency,omitempty" xml:"transactionAmountForeignCurrency"`
	CurrencyRate                     string `json:"currency_rate,omitempty" xml:"currencyRate"`
	Currency                         string `json:"currency,omitempty" xml:"currency"`
	ContactID                        string `json:"contact_id,omitempty" xml:"contactID"`
	FullName                         string `json:"full_name,omitempty" xml:"fullName"`
	GLAccountCode                    string `json:"gl_account_code,omitempty" xml:"glAccountCode"`
	DocumentID                       string `json:"document_id,omitempty" xml:"documentID"`
	Project                          string `json:"project,omitempty" xml:"project"`
	ProjectCode                      string `json:"project_code,omitempty" xml:"projectCode"`
	Created                          string `json:"created" xml:"created"`
	Updated                          string `json:"updated" xml:"updated"`
	Deleted                          string `json:"deleted,omitempty" xml:"deleted"`
}

type Transaction struct {
	ID                string                        `json:"id" xml:"id"`
	HID               string                        `json:"hid,omitempty" xml:"hID"`
	TransactionDate   string                        `json:"transaction_date" xml:"transactionDate"`
	Description       string                        `json:"description,omitempty" xml:"description"`
	Amount            string                        `json:"amount" xml:"amount"`
	GLAccountCode     string                        `json:"gl_account_code,omitempty" xml:"glAccountCode"`
	Contact           *TransactionContact           `json:"contact,omitempty" xml:"contact"`
	Document          *TransactionDocumentReference `json:"document,omitempty" xml:"document"`
	DocumentProcessed *TransactionDocumentProcessed `json:"document_processed,omitempty" xml:"documentProcessed"`
	DocumentMatched   *TransactionDocumentMatched   `json:"document_matched,omitempty" xml:"documentMatched"`
	ForeignCurrency   *TransactionForeignCurrency   `json:"foreign_currency,omitempty" xml:"foreignCurrency"`
	VAT               *TransactionVAT               `json:"vat,omitempty" xml:"vat"`
	Project           *ProjectInfo                  `json:"project,omitempty" xml:"project"`
}

type TransactionContact struct {
	ID          string `json:"id,omitempty" xml:"id,attr"`
	HID         string `json:"hid,omitempty" xml:"HID"`
	FullName    string `json:"full_name,omitempty" xml:"fullName"`
	ZipCode     string `json:"zip_code,omitempty" xml:"zipCode"`
	Country     string `json:"country,omitempty" xml:"country"`
	ContactCode string `json:"contact_code,omitempty" xml:"contactCode"`
	CoCNumber   string `json:"coc_number,omitempty" xml:"CoCNumber"`
	VATNumber   string `json:"vat_number,omitempty" xml:"VATNumber"`
}

type TransactionDocumentReference struct {
	ID              string `json:"id,omitempty" xml:"id,attr"`
	HID             string `json:"hid,omitempty" xml:"HID"`
	Reference       string `json:"reference,omitempty" xml:"reference"`
	Type            string `json:"type,omitempty" xml:"type"`
	TypeDescription string `json:"type_description,omitempty" xml:"typeDescription"`
	FolderID        string `json:"folder_id,omitempty" xml:"folderId"`
	Folder          string `json:"folder,omitempty" xml:"folder"`
	FolderTabID     string `json:"folder_tab_id,omitempty" xml:"folderTabId"`
	FolderTab       string `json:"folder_tab,omitempty" xml:"folderTab"`
	Created         string `json:"created,omitempty" xml:"created"`
	Modified        string `json:"modified,omitempty" xml:"modified"`
	UploadMethod    string `json:"upload_method,omitempty" xml:"uploadMethod"`
}

type TransactionDocumentProcessed struct {
	ProcessedDate string `json:"processed_date,omitempty" xml:"processedDate"`
	ProcessedBy   string `json:"processed_by,omitempty" xml:"processedBy"`
}

type TransactionDocumentMatched struct {
	MatchDate string `json:"match_date,omitempty" xml:"matchDate"`
	MatchedBy string `json:"matched_by,omitempty" xml:"matchedBy"`
}

type TransactionForeignCurrency struct {
	AmountFC string `json:"amount_fc,omitempty" xml:"amountFC"`
	Rate     string `json:"rate,omitempty" xml:"rate"`
	Currency string `json:"currency,omitempty" xml:"currency"`
}

type TransactionVAT struct {
	CodeType        string `json:"code_type,omitempty" xml:"codeType"`
	CodeDescription string `json:"code_description,omitempty" xml:"codeDescription"`
	CodePercentage  string `json:"code_percentage,omitempty" xml:"codePercentage"`
}

type VATCode struct {
	Description     string `json:"description,omitempty" xml:"description"`
	Type            string `json:"type" xml:"type"`
	TypeDescription string `json:"type_description,omitempty" xml:"typeDescription"`
	Percentage      string `json:"percentage" xml:"percentage"`
	Country         string `json:"country,omitempty" xml:"country"`
	StartDate       string `json:"start_date,omitempty" xml:"startDate"`
	EndDate         string `json:"end_date,omitempty" xml:"endDate"`
}

type VATReturnInfo struct {
	DocumentID      string `json:"document_id,omitempty" xml:"documentID"`
	StartDate       string `json:"start_date" xml:"startDate"`
	EndDate         string `json:"end_date" xml:"endDate"`
	Status          string `json:"status,omitempty" xml:"status"`
	SendDate        string `json:"send_date,omitempty" xml:"sendDate"`
	AcknowledgeDate string `json:"acknowledge_date,omitempty" xml:"acknowledgeDate"`
	Modified        string `json:"modified" xml:"modified"`
}

type AdministrationIntegrationData struct {
	CompanyName           string `json:"company_name,omitempty" xml:"CompanyName"`
	Description           string `json:"description,omitempty" xml:"Description"`
	MainContactName       string `json:"main_contact_name,omitempty" xml:"MainContactName"`
	MainContactEmail      string `json:"main_contact_email,omitempty" xml:"MainContactEmail"`
	AddressLine1          string `json:"address_line_1,omitempty" xml:"AddressLine_1"`
	AddressLine2          string `json:"address_line_2,omitempty" xml:"AddressLine_2"`
	Postcode              string `json:"postcode,omitempty" xml:"Postcode"`
	City                  string `json:"city,omitempty" xml:"City"`
	Country               string `json:"country,omitempty" xml:"Country"`
	EmailOutgoingInvoices string `json:"email_outgoing_invoices,omitempty" xml:"EmailOutgoingInvoices"`
	PhoneWork             string `json:"phone_work,omitempty" xml:"PhoneWork"`
	MobileWork            string `json:"mobile_work,omitempty" xml:"MobileWork"`
	FaxWork               string `json:"fax_work,omitempty" xml:"FaxWork"`
	CompanyLogoB64        string `json:"company_logo_b64,omitempty" xml:"CompanyLogoB64"`
	NavigationLogoB64     string `json:"navigation_logo_b64,omitempty" xml:"NavigationLogoB64"`
	IBAN                  string `json:"iban,omitempty" xml:"IBAN"`
	BankAccountName       string `json:"bank_account_name,omitempty" xml:"BankAccountName"`
	CoCNumber             string `json:"coc_number,omitempty" xml:"CoCNumber"`
	VATNumber             string `json:"vat_number,omitempty" xml:"VATNumber"`
}

type FiscalTableTotals struct {
	CompanyID                string `json:"company_id,omitempty"`
	Year                     int    `json:"year,omitempty"`
	RevenueTotal             string `json:"revenue_total,omitempty" xml:"RevenueTotal"`
	GrossMarginTotal         string `json:"gross_margin_total,omitempty" xml:"GrossMarginTotal"`
	ProfessionalCostsTotal   string `json:"professional_costs_total,omitempty" xml:"ProfessionalCostsTotal"`
	SocialContributionsTotal string `json:"social_contributions_total,omitempty" xml:"SocialContributionsTotal"`
}

type BackofficeWorkflowDocument struct {
	SubmitDate   string  `json:"submit_date,omitempty" xml:"SubmitDate"`
	DocumentType XMLText `json:"document_type" xml:"DocumentType"`
	FileName     string  `json:"file_name,omitempty" xml:"Filename"`
}

type BackofficeQuestion struct {
	Date        string  `json:"date,omitempty" xml:"Date"`
	Type        XMLText `json:"type" xml:"Type"`
	Description string  `json:"description,omitempty" xml:"Description"`
	From        string  `json:"from,omitempty" xml:"From"`
}

type ProjectInfo struct {
	Code        string `json:"code,omitempty" xml:"code"`
	Description string `json:"description,omitempty" xml:"description"`
}

type AccountingProject struct {
	ID             string `json:"id,omitempty" xml:"id"`
	HID            string `json:"hid" xml:"HID"`
	Code           string `json:"code,omitempty" xml:"code"`
	Description    string `json:"description,omitempty" xml:"description"`
	StartDate      string `json:"start_date" xml:"startDate"`
	EndDate        string `json:"end_date" xml:"endDate"`
	Company        string `json:"company,omitempty" xml:"company"`
	Contact        string `json:"contact,omitempty" xml:"contact"`
	Tags           string `json:"tags,omitempty" xml:"tags"`
	BudgetSales    string `json:"budget_sales" xml:"budgetSales"`
	BudgetPurchase string `json:"budget_purchase" xml:"budgetPurchase"`
	ContactID      string `json:"contact_id,omitempty" xml:"contactID"`
}

type ProjectBalance struct {
	GLAccountCode string `json:"gl_account_code,omitempty" xml:"glAccountCode"`
	Project       string `json:"project,omitempty" xml:"project"`
	ProjectCode   string `json:"project_code,omitempty" xml:"projectCode"`
	Amount        string `json:"amount" xml:"amount"`
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

type DocumentImageCount struct {
	DocumentID string `json:"document_id"`
	ImageCount int    `json:"image_count"`
}

type DocumentImageData struct {
	DocumentID      string `json:"document_id"`
	MaxWidth        int    `json:"max_width"`
	MaxHeight       int    `json:"max_height"`
	ImageDataBase64 string `json:"image_data_base64"`
}

type DocumentXMLData struct {
	DocumentID string `json:"document_id"`
	XML        string `json:"xml"`
}

type DocumentXMLBinaryData struct {
	DocumentID    string `json:"document_id"`
	XMLDataBase64 string `json:"xml_data_base64"`
}

type DocumentBinaryData struct {
	DocumentID string `json:"document_id"`
	FileData   string `json:"file_data"`
}

type DocumentFolder struct {
	ID              string `json:"id" xml:"ID,attr"`
	Description     string `json:"description" xml:"Description"`
	Icon            string `json:"icon,omitempty" xml:"Icon"`
	ProcessedByYuki bool   `json:"processed_by_yuki" xml:"ProcessedByYuki"`
}

type DocumentFolderCount struct {
	ID          string `json:"id"`
	Description string `json:"description,omitempty"`
	Count       string `json:"count,omitempty"`
}

func (c *DocumentFolderCount) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw struct {
		IDAttr            string `xml:"ID,attr"`
		ID                string `xml:"ID"`
		FolderID          string `xml:"FolderID"`
		Description       string `xml:"Description"`
		Count             string `xml:"Count"`
		DocumentCount     string `xml:"DocumentCount"`
		Documents         string `xml:"Documents"`
		FolderCount       string `xml:"FolderCount"`
		NumberOfDocuments string `xml:"NumberOfDocuments"`
		Total             string `xml:"Total"`
	}
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}
	c.ID = firstNonEmpty(raw.IDAttr, raw.ID, raw.FolderID)
	c.Description = raw.Description
	c.Count = firstNonEmpty(raw.Count, raw.DocumentCount, raw.Documents, raw.FolderCount, raw.NumberOfDocuments, raw.Total)
	return nil
}

type DocumentFolderTab struct {
	ID              string `json:"id" xml:"ID,attr"`
	Description     string `json:"description" xml:"Description"`
	ProcessedByYuki bool   `json:"processed_by_yuki" xml:"ProcessedByYuki"`
}

type Currency struct {
	ID          string `json:"id" xml:"ID,attr"`
	Default     bool   `json:"default" xml:"Default,attr"`
	Description string `json:"description" xml:"Description"`
}

type CostCategory struct {
	ID          string `json:"id" xml:"ID,attr"`
	Description string `json:"description" xml:"Description"`
}

type MenuEntry struct {
	ID    string `json:"id" xml:"ID,attr"`
	Text  string `json:"text" xml:"Text"`
	Icon  string `json:"icon,omitempty" xml:"Icon"`
	Link  string `json:"link,omitempty" xml:"Link"`
	Alert string `json:"alert,omitempty" xml:"Alert"`
}

type PaymentMethod struct {
	ID          string `json:"id" xml:"ID"`
	Description string `json:"description" xml:"Description"`
}

type XMLText struct {
	ID   string `json:"id,omitempty" xml:"ID,attr"`
	Text string `json:"text" xml:",chardata"`
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
