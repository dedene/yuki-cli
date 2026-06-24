package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alecthomas/kong"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/auth"
	"github.com/dedene/yuki-cli/internal/config"
)

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

type Client interface {
	Authenticate(context.Context, string) (string, error)
	Domains(context.Context, string) ([]api.Domain, error)
	CurrentDomain(context.Context, string) (api.Domain, error)
	Administrations(context.Context, string) ([]api.Administration, error)
	Companies(context.Context, string) ([]api.Company, error)
	GLAccounts(context.Context, string, string) ([]api.GLAccount, error)
	OutstandingCreditorItemsByDate(context.Context, string, api.CreditorItemsOptions) ([]api.CreditorItem, error)
	Transactions(context.Context, string, api.TransactionsOptions) ([]api.Transaction, error)
	TransactionDetails(context.Context, string, api.TransactionDetailsOptions) ([]api.TransactionInfo, error)
	TransactionDocument(context.Context, string, string, string) (api.TransactionDocument, error)
	CustomPaymentMethods(context.Context, string, string) ([]api.PaymentMethod, error)
	DocumentFolders(context.Context, string) ([]api.DocumentFolder, error)
	DocumentFolderTabs(context.Context, string, string) ([]api.DocumentFolderTab, error)
	Documents(context.Context, string, api.DocumentsOptions) ([]api.Document, error)
	DocumentsInFolder(context.Context, string, api.DocumentsInFolderOptions) ([]api.Document, error)
	DocumentsInTab(context.Context, string, api.DocumentsInTabOptions) ([]api.Document, error)
	DocumentsByType(context.Context, string, api.DocumentsByTypeOptions) ([]api.Document, error)
	SearchDocuments(context.Context, string, api.SearchDocumentsOptions) ([]api.Document, error)
	FindDocument(context.Context, string, string) (api.Document, error)
	DocumentFile(context.Context, string, string) (api.DocumentFile, error)
	PaymentMethods(context.Context, string) ([]api.PaymentMethod, error)
	Currencies(context.Context, string) ([]api.Currency, error)
	CostCategories(context.Context, string) ([]api.CostCategory, error)
	Menu(context.Context, string) ([]api.MenuEntry, error)
}

type Runtime struct {
	Context   context.Context
	Out       io.Writer
	Err       io.Writer
	Store     auth.Store
	NewClient func(api.Config) Client
}

type Globals struct {
	JSON           bool   `help:"Output JSON to stdout."`
	NoInput        bool   `name:"no-input" help:"Fail instead of prompting."`
	Readonly       bool   `help:"Block mutating commands before network calls."`
	Profile        string `help:"Config/auth profile." default:"default"`
	BaseURL        string `name:"base-url" help:"Override Yuki SOAP base URL, e.g. https://api.yukiworks.nl/ws."`
	Administration string `name:"default-administration" help:"Default administration ID for commands that need one."`
}

type CLI struct {
	Globals `embed:""`

	VersionCmd      VersionCmd         `cmd:"" name:"version" help:"Print version information."`
	Auth            AuthCmd            `cmd:"" help:"Manage authentication."`
	Domains         DomainsCmd         `cmd:"" help:"Inspect accessible domains."`
	Administrations AdministrationsCmd `cmd:"" help:"Inspect accessible administrations."`
	Accounting      AccountingCmd      `cmd:"" help:"Read accounting information."`
	Archive         ArchiveCmd         `cmd:"" help:"Read archive document information."`
}

func Execute(ctx context.Context, args []string, rt Runtime) (err error) {
	rt.Context = ctx
	rt.setDefaults()
	if len(args) == 0 {
		args = []string{"--help"}
	}

	cli := &CLI{}
	parser, err := kong.New(
		cli,
		kong.Name("yuki"),
		kong.Description("CLI for Yuki accounting SOAP webservices."),
		kong.Writers(rt.Out, rt.Err),
		kong.ConfigureHelp(helpOptions()),
		kong.Help(helpPrinter()),
		kong.Exit(func(code int) {
			panic(exitPanic{code: code})
		}),
	)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			if ep, ok := r.(exitPanic); ok {
				if ep.code == 0 {
					err = nil
					return
				}
				err = &ExitError{Code: ep.code, Err: errors.New("exited")}
				return
			}
			panic(r)
		}
	}()

	kctx, err := parser.Parse(args)
	if err != nil {
		return &ExitError{Code: ExitUsage, Err: err}
	}
	return kctx.Run(rt.Context, &rt, &cli.Globals)
}

func (rt *Runtime) setDefaults() {
	if rt.Context == nil {
		rt.Context = context.Background()
	}
	if rt.Out == nil {
		rt.Out = os.Stdout
	}
	if rt.Err == nil {
		rt.Err = os.Stderr
	}
	if rt.NewClient == nil {
		rt.NewClient = func(cfg api.Config) Client {
			return api.New(cfg)
		}
	}
}

func (rt *Runtime) store() (auth.Store, error) {
	if rt.Store != nil {
		return rt.Store, nil
	}
	store, err := auth.OpenDefault()
	if err != nil {
		return nil, fmt.Errorf("open keyring: %w", err)
	}
	rt.Store = store
	return store, nil
}

func (rt *Runtime) client(globals *Globals, profile config.Profile) Client {
	baseURL := globals.BaseURL
	if baseURL == "" {
		baseURL = profile.BaseURL
	}
	return rt.NewClient(api.Config{
		BaseURL:   baseURL,
		UserAgent: "yuki/" + Version,
	})
}

func loadProfile(globals *Globals) (config.Profile, error) {
	cfg, err := config.Load()
	if err != nil {
		return config.Profile{}, err
	}
	profile := cfg.Profile(globals.Profile)
	if globals.BaseURL != "" {
		profile.BaseURL = globals.BaseURL
	}
	if globals.Administration != "" {
		profile.AdministrationID = globals.Administration
	}
	return profile, nil
}

func authenticatedClient(ctx context.Context, rt *Runtime, globals *Globals) (Client, string, error) {
	profile, err := loadProfile(globals)
	if err != nil {
		return nil, "", err
	}
	accessKey, _, err := resolveAccessKey(ctx, rt, globals.Profile)
	if err != nil {
		if errors.Is(err, auth.ErrAccessKeyNotFound) {
			return nil, "", fmt.Errorf("%w; run 'yuki auth login --access-key <key>' or set %s", err, auth.AccessKeyEnv)
		}
		return nil, "", err
	}
	client := rt.client(globals, profile)
	sessionID, err := client.Authenticate(ctx, accessKey)
	if err != nil {
		return nil, "", err
	}
	return client, sessionID, nil
}

func resolveAccessKey(ctx context.Context, rt *Runtime, profile string) (string, auth.Source, error) {
	if key, ok := auth.EnvAccessKey(); ok {
		return key, auth.SourceEnv, nil
	}
	store, err := rt.store()
	if err != nil {
		return "", "", err
	}
	return auth.ResolveAccessKey(ctx, store, profile)
}

type exitPanic struct{ code int }

const (
	ExitOK      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err == nil {
		return "exit"
	}
	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	return e.Err
}
