package api

import (
	"context"
	"strings"
	"testing"
)

func TestBackofficeWorkflowParsesPostmanResponse(t *testing.T) {
	client := fixtureClientForService(t, "Backoffice", "GetWorkflow", backofficeWorkflowResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationID>admin-1</they:administrationID>") {
			t.Fatalf("request body missing administrationID:\n%s", body)
		}
	})

	documents, err := client.BackofficeWorkflow(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("BackofficeWorkflow: %v", err)
	}
	if len(documents) != 1 {
		t.Fatalf("len(documents) = %d, want 1", len(documents))
	}
	if documents[0].SubmitDate != "2020-08-26T14:10:05" ||
		documents[0].DocumentType.ID != "2" ||
		documents[0].DocumentType.Text != "Purchase invoice" ||
		documents[0].FileName != "ININV-00004.pdf" {
		t.Fatalf("documents = %#v", documents)
	}
}

func TestBackofficeOutstandingQuestionsParsesPostmanResponse(t *testing.T) {
	client := fixtureClientForService(t, "Backoffice", "GetOutstandingQuestions", backofficeQuestionsResponse, func(t *testing.T, body string) {
		t.Helper()
		if !strings.Contains(body, "<they:administrationID>admin-1</they:administrationID>") {
			t.Fatalf("request body missing administrationID:\n%s", body)
		}
	})

	questions, err := client.BackofficeOutstandingQuestions(context.Background(), "session-1", "admin-1")
	if err != nil {
		t.Fatalf("BackofficeOutstandingQuestions: %v", err)
	}
	if len(questions) != 2 {
		t.Fatalf("len(questions) = %d, want 2", len(questions))
	}
	if questions[0].Date != "2022-03-09T16:57:47" ||
		questions[0].Type.ID != "29" ||
		questions[0].Type.Text != "Question" ||
		!strings.Contains(questions[0].Description, "Kind regards") ||
		questions[0].From != "Katrien" ||
		questions[1].Type.Text != "Incorrectly processed" {
		t.Fatalf("questions = %#v", questions)
	}
}

const backofficeWorkflowResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetWorkflowResponse xmlns="http://www.theyukicompany.com/">
      <GetWorkflowResult>
        <Workflow xmlns="">
          <Document>
            <SubmitDate>2020-08-26T14:10:05</SubmitDate>
            <DocumentType ID="2">Purchase invoice</DocumentType>
            <Filename>ININV-00004.pdf</Filename>
          </Document>
        </Workflow>
      </GetWorkflowResult>
    </GetWorkflowResponse>
  </soap:Body>
</soap:Envelope>`

const backofficeQuestionsResponse = `<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    <GetOutstandingQuestionsResponse xmlns="http://www.theyukicompany.com/">
      <GetOutstandingQuestionsResult>
        <OutstandingQuestions xmlns="">
          <Question>
            <Date>2022-03-09T16:57:47</Date>
            <Type ID="29">Question</Type>
            <Description>Stijn,

vraag

Kind regards,
Katrien of Belgian portal</Description>
            <From>Katrien</From>
          </Question>
          <Question>
            <Date>2022-07-06T09:41:59</Date>
            <Type ID="29">Incorrectly processed</Type>
            <Description>Hey Katrien,

Can you please review this invoice. I think something is incorrect.</Description>
            <From>Matthijs</From>
          </Question>
        </OutstandingQuestions>
      </GetOutstandingQuestionsResult>
    </GetOutstandingQuestionsResponse>
  </soap:Body>
</soap:Envelope>`
