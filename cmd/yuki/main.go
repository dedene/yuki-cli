package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/dedene/yuki-cli/internal/cmd"
)

func main() {
	err := cmd.Execute(context.Background(), os.Args[1:], cmd.Runtime{})
	if err == nil {
		return
	}
	var exitErr *cmd.ExitError
	if errors.As(err, &exitErr) {
		fmt.Fprintln(os.Stderr, exitErr)
		os.Exit(exitErr.Code)
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(cmd.ExitFailure)
}
