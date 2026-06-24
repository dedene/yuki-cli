package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"strconv"
)

func (c *Client) DocumentImage(ctx context.Context, sessionID, documentID string, maxWidth, maxHeight int) (DocumentImageData, error) {
	params := []Param{
		{Name: "sessionID", Value: sessionID},
		{Name: "documentID", Value: documentID},
		{Name: "maxWidth", Value: strconv.Itoa(maxWidth)},
		{Name: "maxHeight", Value: strconv.Itoa(maxHeight)},
	}
	data, err := c.call(ctx, "Archive", "DocumentImage", params)
	if err != nil {
		return DocumentImageData{}, err
	}
	var env documentImageEnvelope
	if err := xml.Unmarshal(data, &env); err != nil {
		return DocumentImageData{}, fmt.Errorf("parse DocumentImage response: %w", err)
	}
	return DocumentImageData{
		DocumentID:      documentID,
		MaxWidth:        maxWidth,
		MaxHeight:       maxHeight,
		ImageDataBase64: env.Body.Response.Result,
	}, nil
}

type documentImageEnvelope struct {
	Body struct {
		Response struct {
			Result string `xml:"DocumentImageResult"`
		} `xml:"DocumentImageResponse"`
	} `xml:"Body"`
}
