package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAuthenticatePostsSOAPEnvelopeAndParsesSession(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/AccountingInfo.asmx" {
			t.Fatalf("path = %s, want /AccountingInfo.asmx", r.URL.Path)
		}
		if got := r.Header.Get("SOAPAction"); got != SOAPAction("Authenticate") {
			t.Fatalf("SOAPAction = %q", got)
		}
		if got := r.Header.Get("Content-Type"); !strings.Contains(got, "text/xml") {
			t.Fatalf("Content-Type = %q, want text/xml", got)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !strings.Contains(string(body), "<they:accessKey>test-key</they:accessKey>") {
			t.Fatalf("request body missing access key:\n%s", body)
		}
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		_, _ = w.Write([]byte(authenticateResponse))
	}))
	defer srv.Close()

	client := New(Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
	sessionID, err := client.Authenticate(context.Background(), "test-key")
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if sessionID != "a00912fb-558c-4165-a521-d3a095f88cc7" {
		t.Fatalf("sessionID = %q", sessionID)
	}
}

func TestAuthenticateClientPostsDocumentedFieldsAndParsesPostmanResponse(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "AccountingInfo", "AuthenticateClient", authenticateClientPostmanResponse, "", func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:clientID>client-1</they:clientID>",
			"<they:clientSecret>secret-1</they:clientSecret>",
			"<they:accessKey>access-1</they:accessKey>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	sessionID, err := client.AuthenticateClient(context.Background(), "client-1", "secret-1", "access-1")
	if err != nil {
		t.Fatalf("AuthenticateClient: %v", err)
	}
	if sessionID != "client-session-1" {
		t.Fatalf("sessionID = %q", sessionID)
	}
}

func TestAuthenticateClientParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "AccountingInfo", "AuthenticateClient", authenticateClientWSDLResponse, "", nil)

	sessionID, err := client.AuthenticateClient(context.Background(), "client-1", "secret-1", "access-1")
	if err != nil {
		t.Fatalf("AuthenticateClient: %v", err)
	}
	if sessionID != "client-session-wsdl" {
		t.Fatalf("sessionID = %q", sessionID)
	}
}

func TestAuthenticateByUserNamePostsDocumentedFieldsAndParsesResponse(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "AccountingInfo", "AuthenticateByUserName", authenticateByUserNameResponse, "", func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:userName>peter@example.com</they:userName>",
			"<they:password>secret-1</they:password>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	sessionID, err := client.AuthenticateByUserName(context.Background(), "peter@example.com", "secret-1")
	if err != nil {
		t.Fatalf("AuthenticateByUserName: %v", err)
	}
	if sessionID != "username-session-1" {
		t.Fatalf("sessionID = %q", sessionID)
	}
}

func TestSetCurrentDomainPostsDocumentedFields(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "SetCurrentDomain", setCurrentDomainResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:sessionID>session-1</they:sessionID>",
			"<they:domainID>domain-1</they:domainID>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	if err := client.SetCurrentDomain(context.Background(), "session-1", "domain-1"); err != nil {
		t.Fatalf("SetCurrentDomain: %v", err)
	}
}

func TestSetLanguagePostsDocumentedFields(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "SetLanguage", setLanguageResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:sessionID>session-1</they:sessionID>",
			"<they:language>en</they:language>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	if err := client.SetLanguage(context.Background(), "session-1", "en"); err != nil {
		t.Fatalf("SetLanguage: %v", err)
	}
}

func TestListDomainsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "Domains", domainsResponse)

	domains, err := client.Domains(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("Domains: %v", err)
	}
	if len(domains) != 1 {
		t.Fatalf("len(domains) = %d, want 1", len(domains))
	}
	if domains[0].ID != "e9570a5f-xxxx-xxxx-xxxx-144dd8574468" ||
		domains[0].Name != "katrien-highpro" ||
		domains[0].URL != "katrien-highpro.yukiworks.be" {
		t.Fatalf("domain = %#v", domains[0])
	}
}

func TestListAdministrationsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "Administrations", administrationsResponse)

	admins, err := client.Administrations(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("Administrations: %v", err)
	}
	if len(admins) != 1 {
		t.Fatalf("len(admins) = %d, want 1", len(admins))
	}
	admin := admins[0]
	if admin.ID != "38a97d9b-xxxx-xxxxx-xxxx-b9d5793723ee" ||
		admin.Name != "Highpro BV" ||
		admin.Country != "BE" ||
		admin.VATNumber != "BE0123.456.749" ||
		admin.DomainID != "e9570a5f-9339-452b-b621-144dd8574468" ||
		!admin.Active {
		t.Fatalf("administration = %#v", admin)
	}
}

func TestAdministrationIDPostsNameAndParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "AdministrationID", administrationIDResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationName>Highpro NV</they:administrationName>") {
			t.Fatalf("request body missing administration name:\n%s", body)
		}
	})

	id, err := client.AdministrationID(context.Background(), "session-1", "Highpro NV")
	if err != nil {
		t.Fatalf("AdministrationID: %v", err)
	}
	if id != "38axxxxx-xxxx-xxxx-xxxx-xxxxxxxx23ee" {
		t.Fatalf("id = %q", id)
	}
}

func TestAdministrationsWithInternalCustomerCodeParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "AdministrationsWithInternalCustomerCode", administrationsWithInternalCustomerCodeResponse)

	admins, err := client.AdministrationsWithInternalCustomerCode(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("AdministrationsWithInternalCustomerCode: %v", err)
	}
	if len(admins) != 1 {
		t.Fatalf("len(admins) = %d, want 1", len(admins))
	}
	if admins[0].ID != "a484xxxxxxxxxxxxxxxxxd239ed" ||
		admins[0].Name != "Yuki test account" ||
		admins[0].InternalCustomerCode != "MI6-005" {
		t.Fatalf("administration = %#v", admins[0])
	}
}

func TestListCompaniesParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "Companies", companiesResponse)

	companies, err := client.Companies(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("Companies: %v", err)
	}
	if len(companies) != 1 {
		t.Fatalf("len(companies) = %d, want 1", len(companies))
	}
	if companies[0].ID != "38a97d9b-xxxx-xxxx-xxxx-b9d5793723ee" ||
		companies[0].Name != "Highpro BV" ||
		!companies[0].Active {
		t.Fatalf("company = %#v", companies[0])
	}
}

func TestLanguageParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "Language", languageResponse)

	language, err := client.Language(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("Language: %v", err)
	}
	if language != "nl-be" {
		t.Fatalf("language = %q", language)
	}
}

func TestSupportedLanguagesParsesWSDLAnyResponse(t *testing.T) {
	client := fixtureClient(t, "SupportedLanguages", supportedLanguagesResponse)

	languages, err := client.SupportedLanguages(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("SupportedLanguages: %v", err)
	}
	if len(languages) != 2 {
		t.Fatalf("len(languages) = %d, want 2", len(languages))
	}
	if languages[0].Code != "nl-be" ||
		languages[0].Description != "Néerlandais" ||
		languages[0].NativeDescription != "Nederlands" ||
		languages[1].Code != "en" ||
		languages[1].Description != "Anglais" ||
		languages[1].NativeDescription != "English" {
		t.Fatalf("languages = %#v", languages)
	}
}

func TestGetCurrentDomainParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "GetCurrentDomain", currentDomainResponse)

	domain, err := client.CurrentDomain(context.Background(), "session-1")
	if err != nil {
		t.Fatalf("CurrentDomain: %v", err)
	}
	if domain.ID != "e9570a5f-xxxx-xxxx-xxxx-144dd8574468" || domain.Name != "katrien-highpro" {
		t.Fatalf("domain = %#v", domain)
	}
}

func TestListGLAccountsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClient(t, "GetGLAccountScheme", glAccountsResponse)

	accounts, err := client.GLAccounts(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("GLAccounts: %v", err)
	}
	if len(accounts) != 2 {
		t.Fatalf("len(accounts) = %d, want 2", len(accounts))
	}
	if accounts[0].Code != "100000" ||
		accounts[0].Type != "2" ||
		accounts[0].Subtype != "0" ||
		!accounts[0].Enabled ||
		accounts[0].Description != "Geplaatst kapitaal" {
		t.Fatalf("account[0] = %#v", accounts[0])
	}
}

func TestSOAPFaultParsesFaultString(t *testing.T) {
	client := fixtureClient(t, "Domains", soapFaultResponse)

	_, err := client.Domains(context.Background(), "session-1")
	if err == nil {
		t.Fatal("expected SOAP fault")
	}
	fault, ok := err.(*SOAPFault)
	if !ok {
		t.Fatalf("error = %T %v, want *SOAPFault", err, err)
	}
	if fault.Code != "soap:Client" || fault.String != "Daily limit exceeded" || fault.Error() != "Daily limit exceeded" {
		t.Fatalf("fault = %#v, error = %q", fault, fault.Error())
	}
}

func fixtureClient(t *testing.T, operation string, response string) *Client {
	t.Helper()
	return fixtureClientForService(t, "AccountingInfo", operation, response, nil)
}

func fixtureClientForService(t *testing.T, service string, operation string, response string, assertBody func(*testing.T, string)) *Client {
	t.Helper()
	return fixtureClientForServiceWithSessionElement(t, service, operation, response, "sessionID", assertBody)
}

func fixtureClientForServiceWithSessionElement(t *testing.T, service string, operation string, response string, sessionElement string, assertBody func(*testing.T, string)) *Client {
	t.Helper()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Cleanup(srvCloseOnce(t, r))
		wantPath := "/" + service + ".asmx"
		if r.URL.Path != wantPath {
			t.Fatalf("path = %s, want %s", r.URL.Path, wantPath)
		}
		if got := r.Header.Get("SOAPAction"); got != SOAPAction(operation) {
			t.Fatalf("SOAPAction = %q, want %q", got, SOAPAction(operation))
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if sessionElement != "" {
			wantSession := "<they:" + sessionElement + ">session-1</they:" + sessionElement + ">"
			if !strings.Contains(string(body), wantSession) {
				t.Fatalf("request body missing session ID:\n%s", body)
			}
		}
		if assertBody != nil {
			assertBody(t, string(body))
		}
		w.Header().Set("Content-Type", "text/xml; charset=utf-8")
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(srv.Close)

	return New(Config{BaseURL: srv.URL, HTTPClient: srv.Client()})
}

func srvCloseOnce(t *testing.T, r *http.Request) func() {
	t.Helper()
	return func() {
		_ = r.Body.Close()
	}
}

const authenticateResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AuthenticateResponse xmlns="http://www.theyukicompany.com/">
      <AuthenticateResult>a00912fb-558c-4165-a521-d3a095f88cc7</AuthenticateResult>
    </AuthenticateResponse>
  </soap:Body>
</soap:Envelope>`

const authenticateClientPostmanResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AuthenticateResponse xmlns="http://www.theyukicompany.com/">
      <AuthenticateClientResult>client-session-1</AuthenticateClientResult>
    </AuthenticateResponse>
  </soap:Body>
</soap:Envelope>`

const authenticateClientWSDLResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AuthenticateClientResponse xmlns="http://www.theyukicompany.com/">
      <AuthenticateClientResult>client-session-wsdl</AuthenticateClientResult>
    </AuthenticateClientResponse>
  </soap:Body>
</soap:Envelope>`

const authenticateByUserNameResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AuthenticateByUserNameResponse xmlns="http://www.theyukicompany.com/">
      <AuthenticateByUserNameResult>username-session-1</AuthenticateByUserNameResult>
    </AuthenticateByUserNameResponse>
  </soap:Body>
</soap:Envelope>`

const setCurrentDomainResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <SetCurrentDomainResponse xmlns="http://www.theyukicompany.com/" />
  </soap:Body>
</soap:Envelope>`

const setLanguageResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <SetLanguageResponse xmlns="http://www.theyukicompany.com/" />
  </soap:Body>
</soap:Envelope>`

const domainsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <DomainsResponse xmlns="http://www.theyukicompany.com/">
      <DomainsResult>
        <Domains xmlns="">
          <Domain ID="e9570a5f-xxxx-xxxx-xxxx-144dd8574468">
            <Name>katrien-highpro</Name>
            <URL>katrien-highpro.yukiworks.be</URL>
          </Domain>
        </Domains>
      </DomainsResult>
    </DomainsResponse>
  </soap:Body>
</soap:Envelope>`

const currentDomainResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetCurrentDomainResponse xmlns="http://www.theyukicompany.com/">
      <GetCurrentDomainResult>
        <Domains xmlns="">
          <Domain ID="e9570a5f-xxxx-xxxx-xxxx-144dd8574468">
            <Name>katrien-highpro</Name>
          </Domain>
        </Domains>
      </GetCurrentDomainResult>
    </GetCurrentDomainResponse>
  </soap:Body>
</soap:Envelope>`

const administrationsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AdministrationsResponse xmlns="http://www.theyukicompany.com/">
      <AdministrationsResult>
        <Administrations xmlns="">
          <Administration ID="38a97d9b-xxxx-xxxxx-xxxx-b9d5793723ee">
            <Name>Highpro BV</Name>
            <AddressLine>Rijnkaai 37</AddressLine>
            <Postcode>2000</Postcode>
            <City>Antwerpen</City>
            <Country>BE</Country>
            <CoCNumber>0123456749</CoCNumber>
            <VATNumber>BE0123.456.749</VATNumber>
            <StartDate>2016-02-05</StartDate>
            <DomainID>e9570a5f-9339-452b-b621-144dd8574468</DomainID>
            <Active>true</Active>
          </Administration>
        </Administrations>
      </AdministrationsResult>
    </AdministrationsResponse>
  </soap:Body>
</soap:Envelope>`

const administrationIDResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AdministrationIDResponse xmlns="http://www.theyukicompany.com/">
      <AdministrationIDResult>38axxxxx-xxxx-xxxx-xxxx-xxxxxxxx23ee</AdministrationIDResult>
    </AdministrationIDResponse>
  </soap:Body>
</soap:Envelope>`

const administrationsWithInternalCustomerCodeResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AdministrationsWithInternalCustomerCodeResponse xmlns="http://www.theyukicompany.com/">
      <AdministrationsWithInternalCustomerCodeResult>
        <Administrations xmlns="">
          <Administration ID="a484xxxxxxxxxxxxxxxxxd239ed">
            <Name>Yuki test account</Name>
            <AddressLine>Orteliuskaai 4</AddressLine>
            <Postcode>2000</Postcode>
            <City>Antwerpen</City>
            <Country>BE</Country>
            <CoCNumber>0831.xxx.xxx</CoCNumber>
            <VATNumber>BE0831.xxx.xxx</VATNumber>
            <StartDate>2018-01-01</StartDate>
            <DomainID>ff992xxxxxxxxxxxxxxx864858</DomainID>
            <Active>true</Active>
            <InternalCustomerCode>MI6-005</InternalCustomerCode>
          </Administration>
        </Administrations>
      </AdministrationsWithInternalCustomerCodeResult>
    </AdministrationsWithInternalCustomerCodeResponse>
  </soap:Body>
</soap:Envelope>`

const companiesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CompaniesResponse xmlns="http://www.theyukicompany.com/">
      <CompaniesResult>
        <Companies xmlns="">
          <Company ID="38a97d9b-xxxx-xxxx-xxxx-b9d5793723ee">
            <Name>Highpro BV</Name>
            <Active>true</Active>
          </Company>
        </Companies>
      </CompaniesResult>
    </CompaniesResponse>
  </soap:Body>
</soap:Envelope>`

const languageResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <LanguageResponse xmlns="http://www.theyukicompany.com/">
      <LanguageResult>nl-be</LanguageResult>
    </LanguageResponse>
  </soap:Body>
</soap:Envelope>`

const supportedLanguagesResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <SupportedLanguagesResponse xmlns="http://www.theyukicompany.com/">
      <SupportedLanguagesResult>
        <Languages xmlns="">
          <Language Code="nl-be">
            <Description>Néerlandais</Description>
            <NativeDescription>Nederlands</NativeDescription>
          </Language>
          <Language Code="en">
            <Description>Anglais</Description>
            <NativeDescription>English</NativeDescription>
          </Language>
        </Languages>
      </SupportedLanguagesResult>
    </SupportedLanguagesResponse>
  </soap:Body>
</soap:Envelope>`

const glAccountsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetGLAccountSchemeResponse xmlns="http://www.theyukicompany.com/">
      <GetGLAccountSchemeResult>
        <GlAccount>
          <code>100000</code>
          <type>2</type>
          <subtype>0</subtype>
          <isEnabled>true</isEnabled>
          <descripton>Geplaatst kapitaal</descripton>
        </GlAccount>
        <GlAccount>
          <code>100100</code>
          <type>2</type>
          <subtype>16</subtype>
          <isEnabled>true</isEnabled>
          <descripton>Kapitaal Katrien</descripton>
        </GlAccount>
      </GetGLAccountSchemeResult>
    </GetGLAccountSchemeResponse>
  </soap:Body>
</soap:Envelope>`

const soapFaultResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <soap:Fault>
      <faultcode>soap:Client</faultcode>
      <faultstring>Daily limit exceeded</faultstring>
      <detail>Every domain has 1000 free webservice calls a day.</detail>
    </soap:Fault>
  </soap:Body>
</soap:Envelope>`
