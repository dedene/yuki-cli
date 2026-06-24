package cmd

import "fmt"

type VersionCmd struct{}

func (c *VersionCmd) Run(rt *Runtime) error {
	if Commit != "none" || Date != "unknown" {
		_, err := fmt.Fprintf(rt.Out, "yuki %s (%s, %s)\n", Version, Commit, Date)
		return err
	}
	_, err := fmt.Fprintf(rt.Out, "yuki %s\n", Version)
	return err
}
