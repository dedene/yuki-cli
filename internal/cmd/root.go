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
	AuthenticateClient(context.Context, string, string, string) (string, error)
	AuthenticateByUserName(context.Context, string, string) (string, error)
	Domains(context.Context, string) ([]api.Domain, error)
	CurrentDomain(context.Context, string) (api.Domain, error)
	SetCurrentDomain(context.Context, string, string) error
	DomainName(context.Context, string, string) (api.DomainNameResult, error)
	DomainUsers(context.Context, string, string) ([]api.DomainUser, error)
	CreateDomain(context.Context, string, api.DomainCreateOptions) (api.DomainAdminResult, error)
	CreateTrialDomain(context.Context, string, api.DomainCreateOptions) (api.DomainAdminResult, error)
	AddDomainUser(context.Context, string, api.DomainUserAddOptions) (api.DomainAdminResult, error)
	SetLyantheRecognitionEngine(context.Context, string, string, bool) (api.DomainAdminResult, error)
	DomainFunctions(context.Context, string, string) ([]api.DomainFunctionAssignment, error)
	UpdateDomainFunction(context.Context, string, api.UpdateDomainFunctionOptions) (api.DomainFunctionUpdateResult, error)
	SearchContacts(context.Context, string, api.ContactSearchOptions) ([]api.Contact, error)
	SuppliersAndCustomers(context.Context, string, api.ContactSearchOptions) ([]api.Contact, error)
	UpdateContact(context.Context, string, api.ContactUpdateOptions) (api.ContactUpdateResult, error)
	Administrations(context.Context, string) ([]api.Administration, error)
	AdministrationID(context.Context, string, string) (string, error)
	AdministrationsWithInternalCustomerCode(context.Context, string) ([]api.Administration, error)
	Companies(context.Context, string) ([]api.Company, error)
	Language(context.Context, string) (string, error)
	SupportedLanguages(context.Context, string) ([]api.SupportedLanguage, error)
	SetLanguage(context.Context, string, string) error
	GLAccounts(context.Context, string, string) ([]api.GLAccount, error)
	RGSScheme(context.Context, string, api.RGSSchemeOptions) ([]api.RGSEntry, error)
	StartBalanceByGLAccount(context.Context, string, api.StartBalanceByGLAccountOptions) ([]api.GLAccountStartBalance, error)
	GLAccountBalance(context.Context, string, api.GLAccountBalanceOptions) ([]api.GLAccountBalanceItem, error)
	GLAccountBalanceFiscal(context.Context, string, api.GLAccountBalanceOptions) ([]api.GLAccountBalanceItem, error)
	GLAccountBalanceYearEnd(context.Context, string, api.GLAccountBalanceOptions) ([]api.GLAccountBalanceItem, error)
	GLAccountTransactions(context.Context, string, api.GLAccountTransactionsOptions) ([]api.GLAccountTransaction, error)
	GLAccountTransactionsFiscal(context.Context, string, api.GLAccountTransactionsOptions) ([]api.GLAccountTransaction, error)
	GLAccountTransactionsAndContact(context.Context, string, api.GLAccountTransactionsOptions) ([]api.GLAccountTransaction, error)
	ProcessJournal(context.Context, string, api.JournalImportOptions) (api.JournalProcessResult, error)
	CheckOutstandingItem(context.Context, string, string) ([]api.OutstandingItem, error)
	CheckOutstandingItemAdmin(context.Context, string, string, string) ([]api.OutstandingItem, error)
	NetRevenue(context.Context, string, api.RevenueOptions) (api.RevenueReport, error)
	NetRevenueFiscal(context.Context, string, api.RevenueOptions) (api.RevenueReport, error)
	PeriodDateTable(context.Context, string, api.PeriodDateTableOptions) (api.AdministrationPeriod, error)
	FinancialYearModifiedDate(context.Context, string, api.PeriodDateTableOptions) (api.FinancialYearModifiedDate, error)
	ContactDefaultValues(context.Context, string, string, string) ([]api.ContactDefaultValues, error)
	OutstandingCreditorItems(context.Context, string, api.CreditorItemsOptions) ([]api.CreditorItem, error)
	OutstandingCreditorItemsByDate(context.Context, string, api.CreditorItemsOptions) ([]api.CreditorItem, error)
	OutstandingCreditorItemsByDateOutstanding(context.Context, string, api.CreditorItemsOptions) ([]api.CreditorItem, error)
	OutstandingCreditorWithPaymentReference(context.Context, string, api.CreditorItemsOptions) ([]api.CreditorItem, error)
	OutstandingDebtorItems(context.Context, string, api.DebtorItemsOptions) ([]api.DebtorItem, error)
	OutstandingDebtorItemsByDate(context.Context, string, api.DebtorItemsOptions) ([]api.DebtorItem, error)
	OutstandingDebtorItemsByDateOutstanding(context.Context, string, api.DebtorItemsOptions) ([]api.DebtorItem, error)
	OutstandingDebtorItemsWithLanguage(context.Context, string, api.DebtorItemsOptions) ([]api.DebtorItem, error)
	OutstandingDebtorWithPaymentReference(context.Context, string, api.DebtorItemsOptions) ([]api.DebtorItem, error)
	Transactions(context.Context, string, api.TransactionsOptions) ([]api.Transaction, error)
	TransactionDetails(context.Context, string, api.TransactionDetailsOptions) ([]api.TransactionInfo, error)
	TransactionDocument(context.Context, string, string, string) (api.TransactionDocument, error)
	UpdatedAndDeletedTransactions(context.Context, string, api.UpdatedAndDeletedTransactionsOptions) ([]api.UpdatedTransaction, error)
	ChangeDigestTransactionDetail(context.Context, string, string, string) (api.TransactionInfo, error)
	Projects(context.Context, string, api.ProjectsOptions) ([]api.AccountingProject, error)
	ProjectsAndID(context.Context, string, api.ProjectsOptions) ([]api.AccountingProject, error)
	ArchiveProjects(context.Context, string, string) ([]api.AccountingProject, error)
	UpdateProject(context.Context, string, api.ProjectUpdateOptions) (api.ProjectUpdateResult, error)
	ProjectBalance(context.Context, string, api.ProjectBalanceOptions) ([]api.ProjectBalance, error)
	CustomPaymentMethods(context.Context, string, string) ([]api.PaymentMethod, error)
	DocumentFolders(context.Context, string) ([]api.DocumentFolder, error)
	DocumentFolderCounts(context.Context, string, int) ([]api.DocumentFolderCount, error)
	DocumentFolderTabs(context.Context, string, string) ([]api.DocumentFolderTab, error)
	Documents(context.Context, string, api.DocumentsOptions) ([]api.Document, error)
	DocumentsInFolder(context.Context, string, api.DocumentsInFolderOptions) ([]api.Document, error)
	DocumentsInTab(context.Context, string, api.DocumentsInTabOptions) ([]api.Document, error)
	DocumentsByType(context.Context, string, api.DocumentsByTypeOptions) ([]api.Document, error)
	ModifiedDocumentsInFolder(context.Context, string, api.ModifiedDocumentsInFolderOptions) ([]api.Document, error)
	ModifiedDocumentsByType(context.Context, string, api.ModifiedDocumentsByTypeOptions) ([]api.Document, error)
	SearchDocuments(context.Context, string, api.SearchDocumentsOptions) ([]api.Document, error)
	FindDocument(context.Context, string, string) (api.Document, error)
	DocumentBundle(context.Context, string, string) ([]api.Document, error)
	DocumentFile(context.Context, string, string) (api.DocumentFile, error)
	DocumentDownloadURL(context.Context, string, string) (api.DocumentDownloadURL, error)
	DocumentBinaryData(context.Context, string, string) (api.DocumentBinaryData, error)
	DocumentImage(context.Context, string, string, int, int) (api.DocumentImageData, error)
	DocumentImageCount(context.Context, string, string) (api.DocumentImageCount, error)
	DocumentXMLData(context.Context, string, string) (api.DocumentXMLData, error)
	DocumentXMLDataAsBinary(context.Context, string, string) (api.DocumentXMLBinaryData, error)
	DocumentXMLDataAsString(context.Context, string, string) (api.DocumentXMLData, error)
	UploadDocument(context.Context, string, api.ArchiveUploadOptions) (api.ArchiveUploadResult, error)
	UploadDocumentWithData(context.Context, string, api.ArchiveUploadOptions) (api.ArchiveUploadResult, error)
	UploadDocumentWithAttachment(context.Context, string, api.ArchiveAttachmentUploadOptions) (api.ArchiveUploadResult, error)
	PaymentMethods(context.Context, string) ([]api.PaymentMethod, error)
	Currencies(context.Context, string) ([]api.Currency, error)
	CostCategories(context.Context, string) ([]api.CostCategory, error)
	Menu(context.Context, string) ([]api.MenuEntry, error)
	ActiveVATCodes(context.Context, string, string) ([]api.VATCode, error)
	VATReturns(context.Context, string, api.VATReturnListOptions) ([]api.VATReturnInfo, error)
	AdministrationData(context.Context, string, string) (api.AdministrationIntegrationData, error)
	FiscalTable(context.Context, string, string, int) (api.FiscalTableTotals, error)
	BackofficeWorkflow(context.Context, string, string) ([]api.BackofficeWorkflowDocument, error)
	BackofficeOutstandingQuestions(context.Context, string, string) ([]api.BackofficeQuestion, error)
	ImportPettyCashStatement(context.Context, string, api.PettyCashStatementImportOptions) (api.PettyCashImportResult, error)
	ImportPettyCashLine(context.Context, string, api.PettyCashLineImportOptions) (api.PettyCashImportResult, error)
	ImportPettyCashProjectLine(context.Context, string, api.PettyCashLineImportOptions) (api.PettyCashImportResult, error)
	SalesInvoiceSchemaPath(context.Context) (string, error)
	SalesItems(context.Context, string, string) ([]api.SalesItem, error)
	ProcessSalesInvoices(context.Context, string, api.SalesInvoiceImportOptions) (api.SalesInvoiceImportResponse, error)
	ProcessRecognizedSalesInvoices(context.Context, string, api.SalesInvoiceImportOptions) (api.SalesInvoiceImportResponse, error)
}

type Runtime struct {
	Context   context.Context
	Out       io.Writer
	Err       io.Writer
	Store     auth.Store
	NewClient func(api.Config) Client
}

type Globals struct {
	JSON            bool   `help:"Output JSON to stdout."`
	NoInput         bool   `name:"no-input" help:"Fail instead of prompting."`
	Readonly        bool   `help:"Block mutating commands before network calls."`
	Profile         string `help:"Config/auth profile." default:"default"`
	BaseURL         string `name:"base-url" help:"Override Yuki SOAP base URL, e.g. https://api.yukiworks.nl/ws."`
	Administration  string `name:"default-administration" help:"Default administration ID for commands that need one."`
	Domain          string `name:"default-domain" help:"Set this domain ID on each authenticated session."`
	SessionLanguage string `name:"session-language" help:"Set this session language after authentication, e.g. en or nl-be."`
}

type CLI struct {
	Globals `embed:""`

	VersionCmd      VersionCmd         `cmd:"" name:"version" help:"Print version information."`
	Auth            AuthCmd            `cmd:"" help:"Manage authentication."`
	Domains         DomainsCmd         `cmd:"" help:"Manage accessible domains."`
	Administrations AdministrationsCmd `cmd:"" help:"Inspect accessible administrations."`
	Companies       CompaniesCmd       `cmd:"" help:"Inspect accessible companies."`
	Language        LanguageCmd        `cmd:"" help:"Inspect session language."`
	Accounting      AccountingCmd      `cmd:"" help:"Read accounting information."`
	Archive         ArchiveCmd         `cmd:"" help:"Read archive document information."`
	Contacts        ContactsCmd        `cmd:"" help:"Manage contact information."`
	VAT             VATCmd             `cmd:"" name:"vat" help:"Read VAT information."`
	Integration     IntegrationCmd     `cmd:"" help:"Read integration information."`
	FiscalTable     FiscalTableCmd     `cmd:"" name:"fiscal-table" help:"Read fiscal table information."`
	Backoffice      BackofficeCmd      `cmd:"" help:"Read backoffice information."`
	PettyCash       PettyCashCmd       `cmd:"" name:"petty-cash" help:"Import petty cash statements."`
	Sales           SalesCmd           `cmd:"" help:"Manage sales invoices and sales items."`
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
	if globals.Domain != "" {
		profile.DomainID = globals.Domain
	}
	return profile, nil
}

func authenticatedClient(ctx context.Context, rt *Runtime, globals *Globals) (Client, string, error) {
	client, sessionID, profile, err := authenticatedSession(ctx, rt, globals)
	if err != nil {
		return nil, "", err
	}
	if profile.DomainID != "" {
		if err := client.SetCurrentDomain(ctx, sessionID, profile.DomainID); err != nil {
			return nil, "", fmt.Errorf("set current domain %s: %w", profile.DomainID, err)
		}
	}
	if globals.SessionLanguage != "" {
		if err := client.SetLanguage(ctx, sessionID, globals.SessionLanguage); err != nil {
			return nil, "", fmt.Errorf("set session language %s: %w", globals.SessionLanguage, err)
		}
	}
	return client, sessionID, nil
}

func authenticatedSession(ctx context.Context, rt *Runtime, globals *Globals) (Client, string, config.Profile, error) {
	profile, err := loadProfile(globals)
	if err != nil {
		return nil, "", config.Profile{}, err
	}
	accessKey, _, err := resolveAccessKey(ctx, rt, globals.Profile)
	if err != nil {
		if errors.Is(err, auth.ErrAccessKeyNotFound) {
			return nil, "", config.Profile{}, fmt.Errorf("%w; run 'yuki auth login' to enter it securely, or set %s for agent/CI runs", err, auth.AccessKeyEnv)
		}
		return nil, "", config.Profile{}, err
	}
	client := rt.client(globals, profile)
	sessionID, err := client.Authenticate(ctx, accessKey)
	if err != nil {
		return nil, "", config.Profile{}, err
	}
	return client, sessionID, profile, nil
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
