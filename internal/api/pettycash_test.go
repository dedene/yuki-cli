package api

import (
	"context"
	"strings"
	"testing"
)

func TestImportPettyCashStatementPostsDocumentedFields(t *testing.T) {
	statement := "Ledger account Petty Cash;Petty cash description\n570001;Kas Highpro\n"
	client := fixtureClientForService(t, "PettyCash", "ImportStatement", importStatementResponse, func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:administrationID>admin-1</they:administrationID>",
			"<they:statementText>Ledger account Petty Cash;Petty cash description\n570001;Kas Highpro\n</they:statementText>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.ImportPettyCashStatement(context.Background(), "session-1", PettyCashStatementImportOptions{
		AdministrationID: "admin-1",
		StatementText:    statement,
	})
	if err != nil {
		t.Fatalf("ImportPettyCashStatement: %v", err)
	}
	if result.Operation != "ImportStatement" ||
		result.AdministrationID != "admin-1" ||
		result.DocumentID != "doc-1" {
		t.Fatalf("result = %#v", result)
	}
}

func TestImportPettyCashLinePostsDocumentedFields(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "PettyCash", "ImportSingleStatementLine", importSingleStatementLineResponse, "sessionId", func(t *testing.T, body string) {
		t.Helper()
		for _, want := range []string{
			"<they:accountGlCode>570000</they:accountGlCode>",
			"<they:transactionCode>REV6</they:transactionCode>",
			"<they:offsetAccount>700000</they:offsetAccount>",
			"<they:offsetName>Revenue</they:offsetName>",
			"<they:transactionDate>2021-01-15</they:transactionDate>",
			"<they:transactionDescription>revenue 6%</they:transactionDescription>",
			"<they:amount>100</they:amount>",
		} {
			if !strings.Contains(body, want) {
				t.Fatalf("request body missing %q:\n%s", want, body)
			}
		}
	})

	result, err := client.ImportPettyCashLine(context.Background(), "session-1", PettyCashLineImportOptions{
		AccountGLCode:          "570000",
		TransactionCode:        "REV6",
		OffsetAccount:          "700000",
		OffsetName:             "Revenue",
		TransactionDate:        "2021-01-15",
		TransactionDescription: "revenue 6%",
		Amount:                 "100",
	})
	if err != nil {
		t.Fatalf("ImportPettyCashLine: %v", err)
	}
	if result.Operation != "ImportSingleStatementLine" ||
		result.AccountGLCode != "570000" ||
		result.DocumentID != "" {
		t.Fatalf("result = %#v", result)
	}
}

func TestImportPettyCashProjectLinePostsProjectFields(t *testing.T) {
	client := fixtureClientForServiceWithSessionElement(t, "PettyCash", "ImportSingleStatementProjectLine", importSingleStatementProjectLineResponse, "sessionId", func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:projectCode>PROJECTNEW</they:projectCode>") ||
			!strings.Contains(body, "<they:projectName>New Project</they:projectName>") {
			t.Fatalf("request body missing project fields:\n%s", body)
		}
	})

	result, err := client.ImportPettyCashProjectLine(context.Background(), "session-1", PettyCashLineImportOptions{
		AccountGLCode:          "570000",
		TransactionCode:        "REV6",
		OffsetAccount:          "700000",
		OffsetName:             "Revenue",
		TransactionDate:        "2021-01-15",
		TransactionDescription: "revenue 6%",
		Amount:                 "100",
		ProjectCode:            "PROJECTNEW",
		ProjectName:            "New Project",
	})
	if err != nil {
		t.Fatalf("ImportPettyCashProjectLine: %v", err)
	}
	if result.Operation != "ImportSingleStatementProjectLine" ||
		result.ProjectCode != "PROJECTNEW" ||
		result.ProjectName != "New Project" {
		t.Fatalf("result = %#v", result)
	}
}

const importStatementResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ImportStatementResponse xmlns="http://www.theyukicompany.com/">
      <ImportStatementResult>doc-1</ImportStatementResult>
    </ImportStatementResponse>
  </soap:Body>
</soap:Envelope>`

const importSingleStatementLineResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ImportSingleStatementLineResponse xmlns="http://www.theyukicompany.com/" />
  </soap:Body>
</soap:Envelope>`

const importSingleStatementProjectLineResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <ImportSingleStatementProjectLineResponse xmlns="http://www.theyukicompany.com/" />
  </soap:Body>
</soap:Envelope>`
