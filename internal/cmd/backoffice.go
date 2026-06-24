package cmd

import (
	"github.com/dedene/yuki-cli/internal/output"
)

type BackofficeCmd struct {
	Workflow             BackofficeWorkflowCmd             `cmd:"" help:"List documents waiting in the backoffice workflow."`
	OutstandingQuestions BackofficeOutstandingQuestionsCmd `cmd:"" name:"outstanding-questions" help:"List outstanding backoffice questions."`
}

type BackofficeWorkflowCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *BackofficeWorkflowCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	documents, err := client.BackofficeWorkflow(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, documents)
	}

	rows := make([][]string, 0, len(documents))
	for _, document := range documents {
		rows = append(rows, []string{
			document.SubmitDate,
			document.DocumentType.ID,
			document.DocumentType.Text,
			document.FileName,
		})
	}
	return output.Table(rt.Out, []string{"SUBMITTED", "TYPE ID", "TYPE", "FILE"}, rows)
}

type BackofficeOutstandingQuestionsCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *BackofficeOutstandingQuestionsCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	questions, err := client.BackofficeOutstandingQuestions(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, questions)
	}

	rows := make([][]string, 0, len(questions))
	for _, question := range questions {
		rows = append(rows, []string{
			question.Date,
			question.Type.ID,
			question.Type.Text,
			question.From,
			question.Description,
		})
	}
	return output.Table(rt.Out, []string{"DATE", "TYPE ID", "TYPE", "FROM", "DESCRIPTION"}, rows)
}
