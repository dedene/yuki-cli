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

func TestArchiveProjectsParsesWSDLResponse(t *testing.T) {
	client := fixtureClientForService(t, "Archive", "Projects", archiveProjectsResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationID>admin-1</they:administrationID>") {
			t.Fatalf("request body missing administration ID:\n%s", body)
		}
	})

	projects, err := client.ArchiveProjects(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("ArchiveProjects: %v", err)
	}
	if len(projects) != 1 {
		t.Fatalf("len(projects) = %d, want 1", len(projects))
	}
	if projects[0].HID != "5" ||
		projects[0].Code != "ARCHIVE" ||
		projects[0].Description != "Archive Project" {
		t.Fatalf("projects = %#v", projects)
	}
}

func TestProjectBalanceParsesDocumentedResponse(t *testing.T) {
	client := fixtureClientForService(t, "AccountingInfo", "GetProjectBalance", projectBalanceResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:GLAccountCode>700000</they:GLAccountCode>",
			"<they:projectCode>DOS1</they:projectCode>",
			"<they:StartDate>2018-01-01</they:StartDate>",
			"<they:EndDate>2020-12-31</they:EndDate>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	balances, err := client.ProjectBalance(context.Background(), "session-1", ProjectBalanceOptions{
		AdministrationID: "admin-1",
		GLAccountCode:    "700000",
		ProjectCode:      "DOS1",
		StartDate:        "2018-01-01",
		EndDate:          "2020-12-31",
	})
	if err != nil {
		t.Fatalf("ProjectBalance: %v", err)
	}
	if len(balances) != 3 {
		t.Fatalf("len(balances) = %d, want 3", len(balances))
	}
	if balances[0].GLAccountCode != "400000" ||
		balances[0].Project != "Dossier1" ||
		balances[0].ProjectCode != "DOS1" ||
		balances[0].Amount != "542.00" ||
		balances[2].GLAccountCode != "494100" ||
		balances[2].Amount != "-178.50" {
		t.Fatalf("balances = %#v", balances)
	}
}

func TestUpdateProjectPostsDocumentedProjectFields(t *testing.T) {
	client := fixtureClientForService(t, "Projects", "UpdateProject", updateProjectResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:project><they:Description>New Project</they:Description>",
			"<they:Code>PROJECTNEW</they:Code>",
			"<they:Manager>manager@example.com</they:Manager>",
			"<they:Notes>R&amp;D &lt;internal&gt;</they:Notes>",
			"<they:SecurityLevel>1</they:SecurityLevel>",
			"<they:AllowOCRMatching>true</they:AllowOCRMatching>",
			"<they:BudgetCosts>1000</they:BudgetCosts>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.UpdateProject(context.Background(), "session-1", ProjectUpdateOptions{
		AdministrationID: "admin-1",
		Project: ProjectUpdate{
			Description:      "New Project",
			Code:             "PROJECTNEW",
			Company:          "admin-1",
			Manager:          "manager@example.com",
			Contact:          "contact-1",
			Notes:            "R&D <internal>",
			SecurityLevel:    "1",
			AllowOCRMatching: "true",
			StartDate:        "2020-01-20",
			EndDate:          "2022-12-31",
			BudgetRevenue:    "3000",
			BudgetCosts:      "1000",
		},
	})
	if err != nil {
		t.Fatalf("UpdateProject: %v", err)
	}
	if result.AdministrationID != "admin-1" ||
		result.Project.Description != "New Project" ||
		result.Message != "project upserted" {
		t.Fatalf("result = %#v", result)
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

const projectBalanceResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetProjectBalanceResponse xmlns="http://www.theyukicompany.com/">
      <GetProjectBalanceResult>
        <ProjectBalance>
          <glAccountCode>400000</glAccountCode>
          <project>Dossier1</project>
          <projectCode>DOS1</projectCode>
          <amount>542.00</amount>
        </ProjectBalance>
        <ProjectBalance>
          <glAccountCode>451020</glAccountCode>
          <project>Dossier1</project>
          <projectCode>DOS1</projectCode>
          <amount>0.00</amount>
        </ProjectBalance>
        <ProjectBalance>
          <glAccountCode>494100</glAccountCode>
          <project>Dossier1</project>
          <projectCode>DOS1</projectCode>
          <amount>-178.50</amount>
        </ProjectBalance>
      </GetProjectBalanceResult>
    </GetProjectBalanceResponse>
  </soap:Body>
</soap:Envelope>`

const updateProjectResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <UpdateProjectResponse xmlns="http://www.theyukicompany.com/" />
  </soap:Body>
</soap:Envelope>`

const archiveProjectsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ProjectsResponse xmlns="http://www.theyukicompany.com/">
      <ProjectsResult>
        <Projects>
          <Project>
            <description>Archive Project</description>
            <HID>5</HID>
            <code>ARCHIVE</code>
            <startDate>2020-01-01T00:00:00</startDate>
            <endDate>0001-01-01T00:00:00</endDate>
            <company>Highpro BV</company>
            <budgetSales>0</budgetSales>
            <budgetPurchase>0</budgetPurchase>
          </Project>
        </Projects>
      </ProjectsResult>
    </ProjectsResponse>
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
