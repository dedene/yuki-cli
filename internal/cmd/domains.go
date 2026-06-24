package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type DomainsCmd struct {
	List           DomainsListCmd           `cmd:"" help:"List accessible domains."`
	Current        DomainsCurrentCmd        `cmd:"" help:"Show the current session domain."`
	SetCurrent     DomainsSetCurrentCmd     `cmd:"" name:"set-current" help:"Set the current domain for this session."`
	Name           DomainsNameCmd           `cmd:"" help:"Resolve a domain name by administration name."`
	Users          DomainsUsersCmd          `cmd:"" help:"List users for a domain."`
	Create         DomainsCreateCmd         `cmd:"" help:"Create a domain."`
	CreateTrial    DomainsCreateTrialCmd    `cmd:"" name:"create-trial" help:"Create a trial domain."`
	AddUser        DomainsAddUserCmd        `cmd:"" name:"add-user" help:"Add a user to a domain."`
	Lyanthe        DomainsLyantheCmd        `cmd:"" help:"Enable or disable Lyanthe recognition for a domain."`
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

type DomainsNameCmd struct {
	AdministrationName string `name:"administration-name" required:"" help:"Administration name."`
}

func (c *DomainsNameCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.DomainName(rt.Context, sessionID, c.AdministrationName)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"ADMINISTRATION", "DOMAIN"}, [][]string{{
		result.AdministrationName,
		result.DomainName,
	}})
}

type DomainsUsersCmd struct {
	Domain string `name:"domain" required:"" help:"Domain ID."`
}

func (c *DomainsUsersCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	users, err := client.DomainUsers(rt.Context, sessionID, c.Domain)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, users)
	}
	rows := make([][]string, 0, len(users))
	for _, user := range users {
		rows = append(rows, []string{
			user.ID,
			firstNonEmptyString(user.FullName, user.Name),
			user.Login,
			user.Email,
			user.Roles,
			user.Active,
		})
	}
	return output.Table(rt.Out, []string{"ID", "NAME", "LOGIN", "EMAIL", "ROLES", "ACTIVE"}, rows)
}

type DomainsCreateCmd struct {
	AdministrationName string `name:"administration-name" required:"" help:"Administration name."`
	DomainName         string `name:"domain-name" required:"" help:"Domain name to create."`
	DefaultLanguage    string `name:"default-language" required:"" help:"Default language, e.g. en or nl-be."`
	DryRun             bool   `name:"dry-run" help:"Print the planned create without authenticating or sending it."`
}

func (c *DomainsCreateCmd) Run(rt *Runtime, globals *Globals) error {
	opts := api.DomainCreateOptions{
		AdministrationName: c.AdministrationName,
		DomainName:         c.DomainName,
		DefaultLanguage:    c.DefaultLanguage,
	}
	if c.DryRun {
		return renderDomainAdminResult(rt, globals, domainCreateDryRun("CreateDomain", opts))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: domains create")
	}
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.CreateDomain(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderDomainAdminResult(rt, globals, result)
}

type DomainsCreateTrialCmd struct {
	AdministrationName string `name:"administration-name" required:"" help:"Administration name."`
	DomainName         string `name:"domain-name" required:"" help:"Trial domain name to create."`
	DefaultLanguage    string `name:"default-language" required:"" help:"Default language, e.g. en or nl-be."`
	DryRun             bool   `name:"dry-run" help:"Print the planned trial create without authenticating or sending it."`
}

func (c *DomainsCreateTrialCmd) Run(rt *Runtime, globals *Globals) error {
	opts := api.DomainCreateOptions{
		AdministrationName: c.AdministrationName,
		DomainName:         c.DomainName,
		DefaultLanguage:    c.DefaultLanguage,
	}
	if c.DryRun {
		return renderDomainAdminResult(rt, globals, domainCreateDryRun("CreateTrialDomain", opts))
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: domains create-trial")
	}
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.CreateTrialDomain(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderDomainAdminResult(rt, globals, result)
}

func domainCreateDryRun(operation string, opts api.DomainCreateOptions) api.DomainAdminResult {
	return api.DomainAdminResult{
		Operation:          operation,
		AdministrationName: opts.AdministrationName,
		DomainName:         opts.DomainName,
		DefaultLanguage:    opts.DefaultLanguage,
		DryRun:             true,
		Message:            "dry run; no domain create sent",
	}
}

type DomainsAddUserCmd struct {
	Domain          string `name:"domain" required:"" help:"Domain ID."`
	Email           string `name:"email" required:"" help:"User email address."`
	Name            string `name:"name" required:"" help:"User name."`
	Roles           string `name:"roles" help:"Comma-separated Yuki roles."`
	Administrations string `name:"administrations" help:"Comma-separated administration IDs or names."`
	SendMessage     bool   `name:"send-message" help:"Send Yuki's invitation message."`
	CustomMessage   string `name:"custom-message" help:"Custom invitation message."`
	Language        string `name:"language" help:"Invitation/user language, e.g. en."`
	DryRun          bool   `name:"dry-run" help:"Print the planned user add without authenticating or sending it."`
}

func (c *DomainsAddUserCmd) Run(rt *Runtime, globals *Globals) error {
	opts := api.DomainUserAddOptions{
		DomainID:        c.Domain,
		Email:           c.Email,
		Name:            c.Name,
		Roles:           c.Roles,
		Administrations: c.Administrations,
		SendMessage:     c.SendMessage,
		CustomMessage:   c.CustomMessage,
		Language:        c.Language,
	}
	if c.DryRun {
		return renderDomainAdminResult(rt, globals, api.DomainAdminResult{
			Operation:       "AddDomainUser",
			DomainID:        opts.DomainID,
			Email:           opts.Email,
			Name:            opts.Name,
			Roles:           opts.Roles,
			Administrations: opts.Administrations,
			SendMessage:     &opts.SendMessage,
			CustomMessage:   opts.CustomMessage,
			Language:        opts.Language,
			DryRun:          true,
			Message:         "dry run; no domain user add sent",
		})
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: domains add-user")
	}
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.AddDomainUser(rt.Context, sessionID, opts)
	if err != nil {
		return err
	}
	return renderDomainAdminResult(rt, globals, result)
}

type DomainsLyantheCmd struct {
	Domain  string `name:"domain" required:"" help:"Domain ID."`
	Enabled string `name:"enabled" required:"" help:"Whether Lyanthe recognition should be enabled: true or false."`
	DryRun  bool   `name:"dry-run" help:"Print the planned Lyanthe update without authenticating or sending it."`
}

func (c *DomainsLyantheCmd) Run(rt *Runtime, globals *Globals) error {
	enabled, err := strconv.ParseBool(c.Enabled)
	if err != nil {
		return fmt.Errorf("invalid --enabled %q; expected true or false", c.Enabled)
	}
	if c.DryRun {
		return renderDomainAdminResult(rt, globals, api.DomainAdminResult{
			Operation: "LyantheRecognitionEngine",
			DomainID:  c.Domain,
			Enabled:   &enabled,
			DryRun:    true,
			Message:   "dry run; no Lyanthe update sent",
		})
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: domains lyanthe")
	}
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.SetLyantheRecognitionEngine(rt.Context, sessionID, c.Domain, enabled)
	if err != nil {
		return err
	}
	return renderDomainAdminResult(rt, globals, result)
}

func renderDomainAdminResult(rt *Runtime, globals *Globals, result api.DomainAdminResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	enabled := ""
	if result.Enabled != nil {
		enabled = output.Bool(*result.Enabled)
	}
	sendMessage := ""
	if result.SendMessage != nil {
		sendMessage = output.Bool(*result.SendMessage)
	}
	language := firstNonEmptyString(result.DefaultLanguage, result.Language)
	return output.Table(rt.Out, []string{"OPERATION", "DOMAIN", "ADMINISTRATION", "DOMAIN NAME", "EMAIL", "NAME", "ROLES", "ADMINISTRATIONS", "LANGUAGE", "SEND MESSAGE", "ENABLED", "DRY RUN", "MESSAGE"}, [][]string{{
		result.Operation,
		result.DomainID,
		result.AdministrationName,
		result.DomainName,
		result.Email,
		result.Name,
		result.Roles,
		result.Administrations,
		language,
		sendMessage,
		enabled,
		output.Bool(result.DryRun),
		result.Message,
	}})
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
