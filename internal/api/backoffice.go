package api

import (
	"context"
	"encoding/xml"
	"fmt"
)

func (c *Client) BackofficeWorkflow(ctx context.Context, sessionID, administrationID string) ([]BackofficeWorkflowDocument, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "Backoffice", "GetWorkflow", params)
	if err != nil {
		return nil, err
	}
	var env backofficeWorkflowEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetWorkflow response: %w", err)
	}
	return env.Body.Response.Result.Workflow.Documents, nil
}

func (c *Client) BackofficeOutstandingQuestions(ctx context.Context, sessionID, administrationID string) ([]BackofficeQuestion, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: administrationID},
	}
	data, err := c.call(ctx, "Backoffice", "GetOutstandingQuestions", params)
	if err != nil {
		return nil, err
	}
	var env backofficeQuestionsEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parse GetOutstandingQuestions response: %w", err)
	}
	return env.Body.Response.Result.Questions.Questions, nil
}

type backofficeWorkflowEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Workflow struct {
					Documents []BackofficeWorkflowDocument `xml:"Document"`
				} `xml:"Workflow"`
			} `xml:"GetWorkflowResult"`
		} `xml:"GetWorkflowResponse"`
	} `xml:"Body"`
}

type backofficeQuestionsEnvelope struct {
	Body struct {
		Response struct {
			Result struct {
				Questions struct {
					Questions []BackofficeQuestion `xml:"Question"`
				} `xml:"OutstandingQuestions"`
			} `xml:"GetOutstandingQuestionsResult"`
		} `xml:"GetOutstandingQuestionsResponse"`
	} `xml:"Body"`
}
