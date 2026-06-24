package cmd

type ArchiveDocumentsBundleCmd struct {
	Document string `name:"document" required:"" help:"Main document ID."`
}

func (c *ArchiveDocumentsBundleCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	documents, err := client.DocumentBundle(rt.Context, sessionID, c.Document)
	if err != nil {
		return err
	}
	return renderDocuments(rt, globals, documents)
}
