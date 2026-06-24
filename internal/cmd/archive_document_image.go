package cmd

import (
	"errors"

	"github.com/dedene/yuki-cli/internal/output"
)

type ArchiveDocumentsImageCmd struct {
	Document  string `name:"document" required:"" help:"Document ID."`
	MaxWidth  int    `name:"max-width" required:"" help:"Maximum rendered image width in pixels."`
	MaxHeight int    `name:"max-height" required:"" help:"Maximum rendered image height in pixels."`
	Output    string `name:"output" short:"o" help:"Write decoded image bytes to this path."`
}

func (c *ArchiveDocumentsImageCmd) Run(rt *Runtime, globals *Globals) error {
	if c.MaxWidth <= 0 || c.MaxHeight <= 0 {
		return errors.New("invalid --max-width/--max-height; pass positive pixel dimensions")
	}
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	data, err := client.DocumentImage(rt.Context, sessionID, c.Document, c.MaxWidth, c.MaxHeight)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, data)
	}
	if c.Output == "" {
		return errors.New("missing --output; pass --output <path> or use --json to print the base64 image payload")
	}
	return writeBase64File(rt.Out, c.Output, data.DocumentID, data.ImageDataBase64)
}
