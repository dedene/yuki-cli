package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/output"
)

type ArchiveDocumentsXMLBinaryCmd struct {
	Document string `name:"document" required:"" help:"Document ID."`
	Output   string `name:"output" short:"o" help:"Write decoded XML bytes to this path."`
}

func (c *ArchiveDocumentsXMLBinaryCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	data, err := client.DocumentXMLDataAsBinary(rt.Context, sessionID, c.Document)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, data)
	}
	if c.Output == "" {
		return errors.New("missing --output; pass --output <path> or use --json to print the base64 XML payload")
	}
	return writeBase64File(rt.Out, c.Output, data.DocumentID, data.XMLDataBase64)
}
