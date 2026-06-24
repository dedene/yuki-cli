package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

func writeBase64File(w io.Writer, path string, fileName string, encoded string) error {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return fmt.Errorf("decode %s: %w", fileName, err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	_, err = fmt.Fprintf(w, "Wrote %s (%d bytes)\n", path, len(data))
	return err
}

func writeTextFile(w io.Writer, path string, content string) error {
	data := []byte(content)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	_, err := fmt.Fprintf(w, "Wrote %s (%d bytes)\n", path, len(data))
	return err
}
