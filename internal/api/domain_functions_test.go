package api

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

func TestDomainNameParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "GetDomainName", domainNameResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationName>Highpro BV</they:administrationName>") {
			t.Fatalf("request body missing administration name:\n%s", body)
		}
	})

	result, err := client.DomainName(context.Background(), "session-1", "Highpro BV")
	if err != nil {
		t.Fatalf("DomainName: %v", err)
	}
	if result.AdministrationName != "Highpro BV" || result.DomainName != "highpro.yukiworks.be" {
		t.Fatalf("result = %#v", result)
	}
}

func TestDomainUsersParsesWSDLAnyResponse(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "GetDomainUsers", domainUsersResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:domain>domain-1</they:domain>") {
			t.Fatalf("request body missing domain:\n%s", body)
		}
	})

	users, err := client.DomainUsers(context.Background(), "session-1", "domain-1")
	if err != nil {
		t.Fatalf("DomainUsers: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("len(users) = %d, want 1", len(users))
	}
	if users[0].ID != "user-1" ||
		users[0].FullName != "Peter Dedene" ||
		users[0].Login != "peter@example.com" ||
		users[0].Email != "peter@example.com" ||
		users[0].Roles != "Backoffice" ||
		users[0].Active != "true" {
		t.Fatalf("users[0] = %#v", users[0])
	}
	data, err := json.Marshal(users[0])
	if err != nil {
		t.Fatalf("json: %v", err)
	}
	if !strings.Contains(string(data), `"name":"Department"`) || !strings.Contains(string(data), `"value":"Finance"`) {
		t.Fatalf("unknown fields not preserved in JSON: %s", data)
	}
}

func TestCreateDomainPostsWSDLFields(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "CreateDomain", createDomainResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationName>Highpro BV</they:administrationName>",
			"<they:domainName>Highpro</they:domainName>",
			"<they:defaultLanguage>nl-be</they:defaultLanguage>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.CreateDomain(context.Background(), "session-1", DomainCreateOptions{
		AdministrationName: "Highpro BV",
		DomainName:         "Highpro",
		DefaultLanguage:    "nl-be",
	})
	if err != nil {
		t.Fatalf("CreateDomain: %v", err)
	}
	if result.Operation != "CreateDomain" || result.Message != "domain-created" {
		t.Fatalf("result = %#v", result)
	}
}

func TestCreateTrialDomainPostsWSDLFields(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "CreateTrialDomain", createTrialDomainResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:domainName>Trial Highpro</they:domainName>") {
			t.Fatalf("request body missing trial domain name:\n%s", body)
		}
	})

	result, err := client.CreateTrialDomain(context.Background(), "session-1", DomainCreateOptions{
		AdministrationName: "Highpro BV",
		DomainName:         "Trial Highpro",
		DefaultLanguage:    "en",
	})
	if err != nil {
		t.Fatalf("CreateTrialDomain: %v", err)
	}
	if result.Operation != "CreateTrialDomain" || result.Message != "trial-domain-created" {
		t.Fatalf("result = %#v", result)
	}
}

func TestAddDomainUserPostsWSDLFields(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "AddDomainUser", addDomainUserResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:domain>domain-1</they:domain>",
			"<they:email>peter@example.com</they:email>",
			"<they:name>Peter Dedene</they:name>",
			"<they:roles>Backoffice</they:roles>",
			"<they:administrations>admin-1</they:administrations>",
			"<they:sendMessage>true</they:sendMessage>",
			"<they:customMessage>Welcome</they:customMessage>",
			"<they:language>en</they:language>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.AddDomainUser(context.Background(), "session-1", DomainUserAddOptions{
		DomainID:        "domain-1",
		Email:           "peter@example.com",
		Name:            "Peter Dedene",
		Roles:           "Backoffice",
		Administrations: "admin-1",
		SendMessage:     true,
		CustomMessage:   "Welcome",
		Language:        "en",
	})
	if err != nil {
		t.Fatalf("AddDomainUser: %v", err)
	}
	if result.Operation != "AddDomainUser" || result.Email != "peter@example.com" || result.Message != "user-added" {
		t.Fatalf("result = %#v", result)
	}
}

func TestSetLyantheRecognitionEnginePostsWSDLFields(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "LyantheRecognitionEngine", lyantheRecognitionEngineResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:domain>domain-1</they:domain>",
			"<they:enableLyanthe>true</they:enableLyanthe>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.SetLyantheRecognitionEngine(context.Background(), "session-1", "domain-1", true)
	if err != nil {
		t.Fatalf("SetLyantheRecognitionEngine: %v", err)
	}
	if result.Operation != "LyantheRecognitionEngine" || result.Enabled == nil || !*result.Enabled || result.Message != "lyanthe-enabled" {
		t.Fatalf("result = %#v", result)
	}
}

func TestDomainFunctionsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "GetDomainFunctions", domainFunctionsResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:domain>domain-1</they:domain>") {
			t.Fatalf("request body missing domain:\n%s", body)
		}
	})

	assignments, err := client.DomainFunctions(context.Background(), "session-1", "domain-1")
	if err != nil {
		t.Fatalf("DomainFunctions: %v", err)
	}
	if len(assignments) != 4 {
		t.Fatalf("len(assignments) = %d, want 4", len(assignments))
	}
	if assignments[0].DomainID != "domain-1" ||
		assignments[0].Function != "BOResponsible" ||
		assignments[0].FullName != "Oliver Test" ||
		assignments[0].Login != "oliver.test@test.com" {
		t.Fatalf("assignments[0] = %#v", assignments[0])
	}
	if assignments[1].Function != "BOBackup" ||
		assignments[1].FullName != "katrien portaluser" ||
		assignments[1].Login != "katrien.portaltestuser@test.be" {
		t.Fatalf("assignments[1] = %#v", assignments[1])
	}
	if assignments[2].Function != "BOController" || assignments[2].FullName != "" || assignments[2].Login != "" {
		t.Fatalf("assignments[2] = %#v", assignments[2])
	}
}

func TestUpdateDomainFunctionParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "Domains", "UpdateDomainFunctions", updateDomainFunctionsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:domain>domain-1</they:domain>",
			"<they:domainFunction>BOAccountManager</they:domainFunction>",
			"<they:login>test@test.be</they:login>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.UpdateDomainFunction(context.Background(), "session-1", UpdateDomainFunctionOptions{
		DomainID: "domain-1",
		Function: "BOAccountManager",
		Login:    "test@test.be",
	})
	if err != nil {
		t.Fatalf("UpdateDomainFunction: %v", err)
	}
	if result.DomainID != "domain-1" ||
		result.Function != "BOAccountManager" ||
		result.Login != "test@test.be" ||
		result.Message != "Domain function successfully updated" {
		t.Fatalf("result = %#v", result)
	}
}

const domainNameResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetDomainNameResponse xmlns="http://www.theyukicompany.com/">
      <GetDomainNameResult>highpro.yukiworks.be</GetDomainNameResult>
    </GetDomainNameResponse>
  </soap:Body>
</soap:Envelope>`

const domainUsersResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetDomainUsersResponse xmlns="http://www.theyukicompany.com/">
      <GetDomainUsersResult>
        <DomainUsers>
          <DomainUser ID="user-1">
            <FullName>Peter Dedene</FullName>
            <Login>peter@example.com</Login>
            <Email>peter@example.com</Email>
            <Roles>Backoffice</Roles>
            <Active>true</Active>
            <Department>Finance</Department>
          </DomainUser>
        </DomainUsers>
      </GetDomainUsersResult>
    </GetDomainUsersResponse>
  </soap:Body>
</soap:Envelope>`

const createDomainResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateDomainResponse xmlns="http://www.theyukicompany.com/">
      <CreateDomainResult>domain-created</CreateDomainResult>
    </CreateDomainResponse>
  </soap:Body>
</soap:Envelope>`

const createTrialDomainResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateTrialDomainResponse xmlns="http://www.theyukicompany.com/">
      <CreateTrialDomainResult>trial-domain-created</CreateTrialDomainResult>
    </CreateTrialDomainResponse>
  </soap:Body>
</soap:Envelope>`

const addDomainUserResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <AddDomainUserResponse xmlns="http://www.theyukicompany.com/">
      <AddDomainUserResult>user-added</AddDomainUserResult>
    </AddDomainUserResponse>
  </soap:Body>
</soap:Envelope>`

const lyantheRecognitionEngineResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <LyantheRecognitionEngineResponse xmlns="http://www.theyukicompany.com/">
      <LyantheRecognitionEngineResult>lyanthe-enabled</LyantheRecognitionEngineResult>
    </LyantheRecognitionEngineResponse>
  </soap:Body>
</soap:Envelope>`

const domainFunctionsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetDomainFunctionsResponse xmlns="http://www.theyukicompany.com/">
      <GetDomainFunctionsResult>
        <DomainFunctions xmlns="">
          <BOResponsible>
            <FullName>Oliver Test</FullName>
            <Login>oliver.test@test.com</Login>
          </BOResponsible>
          <BOBackup>
            <FullName>katrien portaluser</FullName>
            <Login>katrien.portaltestuser@test.be</Login>
          </BOBackup>
          <BOController>
            <FullName></FullName>
            <Login></Login>
          </BOController>
          <BOAccountManager>
            <FullName></FullName>
            <Login></Login>
          </BOAccountManager>
        </DomainFunctions>
      </GetDomainFunctionsResult>
    </GetDomainFunctionsResponse>
  </soap:Body>
</soap:Envelope>`

const updateDomainFunctionsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <UpdateDomainFunctionsResponse xmlns="http://www.theyukicompany.com/">
      <UpdateDomainFunctionsResult>Domain function successfully updated</UpdateDomainFunctionsResult>
    </UpdateDomainFunctionsResponse>
  </soap:Body>
</soap:Envelope>`
