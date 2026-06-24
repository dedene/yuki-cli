package api

import (
	"context"
	"strings"
	"testing"
)

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
