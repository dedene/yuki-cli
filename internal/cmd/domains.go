package cmd

import "github.com/dedene/yuki-cli/internal/output"

type DomainsCmd struct {
	List      DomainsListCmd      `cmd:"" help:"List accessible domains."`
	Current   DomainsCurrentCmd   `cmd:"" help:"Show the current session domain."`
	Functions DomainsFunctionsCmd `cmd:"" help:"Inspect backoffice role assignments for a domain."`
}

type DomainsListCmd struct{}

func (c *DomainsListCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	domains, err := client.Domains(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, domains)
	}
	rows := make([][]string, 0, len(domains))
	for _, domain := range domains {
		rows = append(rows, []string{domain.ID, domain.Name, domain.URL})
	}
	return output.Table(rt.Out, []string{"ID", "NAME", "URL"}, rows)
}

type DomainsCurrentCmd struct{}

func (c *DomainsCurrentCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	domain, err := client.CurrentDomain(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, domain)
	}
	return output.Table(rt.Out, []string{"ID", "NAME", "URL"}, [][]string{{domain.ID, domain.Name, domain.URL}})
}

type DomainsFunctionsCmd struct {
	Domain string `name:"domain" required:"" help:"Domain ID."`
}

func (c *DomainsFunctionsCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	assignments, err := client.DomainFunctions(rt.Context, sessionID, c.Domain)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, assignments)
	}
	rows := make([][]string, 0, len(assignments))
	for _, assignment := range assignments {
		rows = append(rows, []string{assignment.Function, assignment.FullName, assignment.Login})
	}
	return output.Table(rt.Out, []string{"FUNCTION", "FULL NAME", "LOGIN"}, rows)
}
