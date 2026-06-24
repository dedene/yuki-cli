package api

import (
	"context"
	"strings"
	"testing"
)

func TestProjectsParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetProjects", projectsResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:searchOption>All</they:searchOption>",
			"<they:searchValue></they:searchValue>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	projects, err := client.Projects(context.Background(), "session-1", ProjectsOptions{
		AdministrationID: "admin-1",
		SearchOption:     "All",
	})
	if err != nil {
		t.Fatalf("Projects: %v", err)
	}
	if len(projects) != 3 {
		t.Fatalf("len(projects) = %d, want 3", len(projects))
	}
	if projects[0].HID != "1" ||
		projects[0].Code != "WELLNESS" ||
		projects[0].Description != "Wellness Event" ||
		projects[1].Code != "DOS1" ||
		projects[2].Contact != "AD Delhaize" ||
		projects[2].ContactID != "d390f5bb-3c2e-41a8-9dff-5023595ded16" {
		t.Fatalf("projects = %#v", projects)
	}
}

func TestProjectsAndIDParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetProjectsAndID", projectsAndIDResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:searchOption>Code</they:searchOption>",
			"<they:searchValue>WELLNESS</they:searchValue>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	projects, err := client.ProjectsAndID(context.Background(), "session-1", ProjectsOptions{
		AdministrationID: "admin-1",
		SearchOption:     "Code",
		SearchValue:      "WELLNESS",
	})
	if err != nil {
		t.Fatalf("ProjectsAndID: %v", err)
	}
	if len(projects) != 2 {
		t.Fatalf("len(projects) = %d, want 2", len(projects))
	}
	if projects[0].ID != "f2e749e4-af93-4259-9351-930ab20a2991" ||
		projects[0].Code != "WELLNESS" ||
		projects[1].ID != "b37d5de3-05f1-4061-a495-5591d8baf745" ||
		projects[1].Code != "18-OH-ALG" {
		t.Fatalf("projects = %#v", projects)
	}
}

const projectsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetProjectsResponse xmlns="http://www.theyukicompany.com/">
      <GetProjectsResult>
        <Project>
          <description>Wellness Event</description>
          <HID>1</HID>
          <code>WELLNESS</code>
          <startDate>2018-09-10T00:00:00</startDate>
          <endDate>0001-01-01T00:00:00</endDate>
          <company>Highpro BV</company>
          <budgetSales>0</budgetSales>
          <budgetPurchase>0</budgetPurchase>
        </Project>
        <Project>
          <description>Dossier1</description>
          <HID>3</HID>
          <code>DOS1</code>
          <startDate>2019-02-27T00:00:00</startDate>
          <endDate>0001-01-01T00:00:00</endDate>
          <company>Highpro BV</company>
          <budgetSales>0</budgetSales>
          <budgetPurchase>0</budgetPurchase>
        </Project>
        <Project>
          <description>Project 1</description>
          <HID>9</HID>
          <code>PROJECT1</code>
          <startDate>2020-01-20T00:00:00</startDate>
          <endDate>2022-12-31T00:00:00</endDate>
          <company>Highpro BV</company>
          <contact>AD Delhaize</contact>
          <budgetSales>0</budgetSales>
          <budgetPurchase>0</budgetPurchase>
          <contactID>d390f5bb-3c2e-41a8-9dff-5023595ded16</contactID>
        </Project>
      </GetProjectsResult>
    </GetProjectsResponse>
  </soap:Body>
</soap:Envelope>`

const projectsAndIDResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetProjectsAndIDResponse xmlns="http://www.theyukicompany.com/">
      <GetProjectsAndIDResult>
        <Project>
          <description>Wellness Event</description>
          <HID>1</HID>
          <code>WELLNESS</code>
          <startDate>2018-09-10T00:00:00</startDate>
          <endDate>0001-01-01T00:00:00</endDate>
          <company>Highpro BV</company>
          <budgetSales>0</budgetSales>
          <budgetPurchase>0</budgetPurchase>
          <id>f2e749e4-af93-4259-9351-930ab20a2991</id>
        </Project>
        <Project>
          <description>test ttpe</description>
          <HID>2</HID>
          <code>18-OH-ALG</code>
          <startDate>2018-10-15T00:00:00</startDate>
          <endDate>0001-01-01T00:00:00</endDate>
          <company>Highpro BV</company>
          <budgetSales>0</budgetSales>
          <budgetPurchase>0</budgetPurchase>
          <id>b37d5de3-05f1-4061-a495-5591d8baf745</id>
        </Project>
      </GetProjectsAndIDResult>
    </GetProjectsAndIDResponse>
  </soap:Body>
</soap:Envelope>`
