package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type DomainsCmd struct {
	List           DomainsListCmd           `cmd:"" help:"List accessible domains."`
	Current        DomainsCurrentCmd        `cmd:"" help:"Show the current session domain."`
	SetCurrent     DomainsSetCurrentCmd     `cmd:"" name:"set-current" help:"Set the current domain for this session."`
	Functions      DomainsFunctionsCmd      `cmd:"" help:"Inspect backoffice role assignments for a domain."`
	UpdateFunction DomainsUpdateFunctionCmd `cmd:"" name:"update-function" help:"Update one backoffice role assignment for a domain."`
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

type DomainsSetCurrentCmd struct {
	Domain string `name:"domain" required:"" help:"Domain ID."`
	DryRun bool   `name:"dry-run" help:"Print the planned session update without authenticating or sending it."`
}

func (c *DomainsSetCurrentCmd) Run(rt *Runtime, globals *Globals) error {
	result := sessionSettingResult{
		DomainID: c.Domain,
		Message:  "current domain set for this session",
	}
	if c.DryRun {
		result.DryRun = true
		result.Message = "dry run; current domain not sent"
		return renderSessionSetting(rt, globals, result)
	}
	client, sessionID, _, err := authenticatedSession(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	if err := client.SetCurrentDomain(rt.Context, sessionID, c.Domain); err != nil {
		return err
	}
	return renderSessionSetting(rt, globals, result)
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

type DomainsUpdateFunctionCmd struct {
	Domain   string `name:"domain" required:"" help:"Domain ID."`
	Function string `name:"function" required:"" help:"Backoffice role: BOResponsible, BOBackup, BOController, or BOAccountManager."`
	Login    string `name:"login" required:"" help:"Yuki login to assign to the role."`
	DryRun   bool   `name:"dry-run" help:"Print the planned update without authenticating or sending it."`
}

func (c *DomainsUpdateFunctionCmd) Run(rt *Runtime, globals *Globals) error {
	opts := api.UpdateDomainFunctionOptions{
		DomainID: c.Domain,
		Function: c.Function,
		Login:    c.Login,
	}
	if !validDomainFunction(opts.Function) {
		return fmt.Errorf("invalid --function %q; expected one of %s", opts.Function, strings.Join(domainFunctionValues, ", "))
	}
	if c.DryRun {
		return renderDomainFunctionUpdate(rt, globals, api.DomainFunctionUpdateResult{
			DomainID: opts.DomainID,
			Function: opts.Function,
			Login:    opts.Login,
			Message:  "dry run; no update sent",
			DryRun:   true,
		})
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: domains update-function")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.UpdateDomainFunction(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderDomainFunctionUpdate(rt, globals, result)
}

func renderDomainFunctionUpdate(rt *Runtime, globals *Globals, result api.DomainFunctionUpdateResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	dryRun := "false"
	if result.DryRun {
		dryRun = "true"
	}
	return output.Table(rt.Out, []string{"DOMAIN", "FUNCTION", "LOGIN", "DRY RUN", "MESSAGE"}, [][]string{{
		result.DomainID,
		result.Function,
		result.Login,
		dryRun,
		result.Message,
	}})
}

var domainFunctionValues = []string{"BOResponsible", "BOBackup", "BOController", "BOAccountManager"}

func validDomainFunction(value string) bool {
	for _, allowed := range domainFunctionValues {
		if value == allowed {
			return true
		}
	}
	return false
}

type sessionSettingResult struct {
	DomainID string `json:"domain_id,omitempty"`
	Language string `json:"language,omitempty"`
	DryRun   bool   `json:"dry_run,omitempty"`
	Message  string `json:"message"`
}

func renderSessionSetting(rt *Runtime, globals *Globals, result sessionSettingResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"DOMAIN", "LANGUAGE", "DRY RUN", "MESSAGE"}, [][]string{{
		result.DomainID,
		result.Language,
		output.Bool(result.DryRun),
		result.Message,
	}})
}
