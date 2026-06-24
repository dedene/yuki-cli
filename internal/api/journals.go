package api

import "context"

type JournalImportOptions struct {
	AdministrationID string
	XMLDoc           string
}

type JournalProcessResult struct {
	AdministrationID string `json:"administration_id,omitempty"`
	DocumentID       string `json:"document_id"`
}

func (c *Client) ProcessJournal(ctx context.Context, sessionID string, opts JournalImportOptions) (JournalProcessResult, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "administrationID", Value: opts.AdministrationID},
		{Name: "xmlDoc", Value: opts.XMLDoc, Raw: true},
	}
	data, err := c.call(ctx, "Accounting", "ProcessJournal", params)
	if err != nil {
		return JournalProcessResult{}, err
	}
	documentID, err := textAt(data, []string{"Envelope", "Body", "ProcessJournalResponse", "ProcessJournalResult"})
	if err != nil {
		return JournalProcessResult{}, err
	}
	return JournalProcessResult{
		AdministrationID: opts.AdministrationID,
		DocumentID:       documentID,
	}, nil
}
