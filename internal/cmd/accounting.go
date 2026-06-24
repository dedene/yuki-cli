package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/dedene/yuki-cli/internal/api"
	"github.com/dedene/yuki-cli/internal/output"
)

type AccountingCmd struct {
	GLAccounts      GLAccountsCmd               `cmd:"" name:"gl-accounts" help:"Inspect GL accounts."`
	Revenue         RevenueCmd                  `cmd:"" help:"Inspect revenue reports."`
	CreditorItems   CreditorItemsCmd            `cmd:"" name:"creditor-items" help:"Inspect outstanding creditor purchase invoices."`
	DebtorItems     DebtorItemsCmd              `cmd:"" name:"debtor-items" help:"Inspect outstanding debtor sales invoices."`
	Outstanding     OutstandingCmd              `cmd:"" help:"Inspect outstanding items."`
	Journals        JournalsCmd                 `cmd:"" help:"Process general journal entries."`
	Transactions    TransactionsCmd             `cmd:"" help:"Inspect accounting transactions."`
	ChangeDigest    ChangeDigestCmd             `cmd:"" name:"change-digest" help:"Inspect change digest feeds."`
	Projects        ProjectsCmd                 `cmd:"" help:"Inspect accounting projects."`
	PaymentMethods  AccountingPaymentMethodsCmd `cmd:"" name:"payment-methods" help:"Inspect accounting payment methods."`
	Periods         PeriodsCmd                  `cmd:"" help:"Inspect administration periods."`
	ContactDefaults ContactDefaultsCmd          `cmd:"" name:"contact-defaults" help:"Inspect contact accounting defaults."`
}

type GLAccountsCmd struct {
	List                    GLAccountsListCmd                    `cmd:"" help:"List GL accounts for an administration."`
	RGSScheme               GLAccountsRGSSchemeCmd               `cmd:"" name:"rgs-scheme" help:"List GL accounts with RGS codes."`
	StartBalance            GLAccountsStartBalanceCmd            `cmd:"" name:"start-balance" help:"List GL account start balances for a bookyear."`
	Balance                 GLAccountsBalanceCmd                 `cmd:"" help:"List GL account balances at a transaction date."`
	BalanceFiscal           GLAccountsBalanceFiscalCmd           `cmd:"" name:"balance-fiscal" help:"List fiscal GL account balances at a transaction date."`
	BalanceYearEnd          GLAccountsBalanceYearEndCmd          `cmd:"" name:"balance-year-end" help:"List year-end GL account balances at a transaction date."`
	Transactions            GLAccountsTransactionsCmd            `cmd:"" help:"List transactions for a GL account."`
	TransactionsFiscal      GLAccountsTransactionsFiscalCmd      `cmd:"" name:"transactions-fiscal" help:"List fiscal transactions for a GL account."`
	TransactionsWithContact GLAccountsTransactionsWithContactCmd `cmd:"" name:"transactions-with-contact" help:"List GL account transactions with contact and document filename fields."`
}

type GLAccountsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *GLAccountsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	accounts, err := client.GLAccounts(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, accounts)
	}
	rows := make([][]string, 0, len(accounts))
	for _, account := range accounts {
		rows = append(rows, []string{account.Code, account.Type, account.Subtype, output.Bool(account.Enabled), account.Description})
	}
	return output.Table(rt.Out, []string{"CODE", "TYPE", "SUBTYPE", "ENABLED", "DESCRIPTION"}, rows)
}

type GLAccountsRGSSchemeCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	RGSVersion     string `name:"rgs-version" default:"2.0" help:"RGS version to request."`
}

func (c *GLAccountsRGSSchemeCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	entries, err := client.RGSScheme(rt.Context, sessionID, api.RGSSchemeOptions{
		AdministrationID: administrationID,
		RGSVersion:       c.RGSVersion,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, entries)
	}
	rows := make([][]string, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, []string{
			entry.YukiCode,
			entry.YukiIsEnabled,
			entry.YukiDescription,
			entry.RGSReferenceCode,
			entry.RGSDescription,
			entry.RGSFlipCode,
		})
	}
	return output.Table(rt.Out, []string{"YUKI CODE", "ENABLED", "YUKI DESCRIPTION", "RGS CODE", "RGS DESCRIPTION", "RGS FLIP"}, rows)
}

type GLAccountsStartBalanceCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Bookyear       int    `name:"bookyear" required:"" help:"Bookyear to request."`
	FinancialMode  int    `name:"financial-mode" default:"1" help:"Yuki financial mode."`
}

func (c *GLAccountsStartBalanceCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.StartBalanceByGLAccount(rt.Context, sessionID, api.StartBalanceByGLAccountOptions{
		AdministrationID: administrationID,
		Bookyear:         c.Bookyear,
		FinancialMode:    c.FinancialMode,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	rows := make([][]string, 0, len(balances))
	for _, balance := range balances {
		rows = append(rows, []string{balance.AccountID, balance.StartBalance, balance.AccountDescription})
	}
	return output.Table(rt.Out, []string{"ACCOUNT", "START BALANCE", "DESCRIPTION"}, rows)
}

type GLAccountsBalanceCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Date           string `name:"date" required:"" help:"Transaction date, YYYY-MM-DD."`
}

func (c *GLAccountsBalanceCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.GLAccountBalance(rt.Context, sessionID, api.GLAccountBalanceOptions{
		AdministrationID: administrationID,
		TransactionDate:  c.Date,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	return renderGLAccountBalances(rt, balances)
}

type GLAccountsBalanceFiscalCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Date           string `name:"date" required:"" help:"Transaction date, YYYY-MM-DD."`
}

func (c *GLAccountsBalanceFiscalCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.GLAccountBalanceFiscal(rt.Context, sessionID, api.GLAccountBalanceOptions{
		AdministrationID: administrationID,
		TransactionDate:  c.Date,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	return renderGLAccountBalances(rt, balances)
}

type GLAccountsBalanceYearEndCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Date           string `name:"date" required:"" help:"Transaction date, YYYY-MM-DD."`
}

func (c *GLAccountsBalanceYearEndCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	balances, err := client.GLAccountBalanceYearEnd(rt.Context, sessionID, api.GLAccountBalanceOptions{
		AdministrationID: administrationID,
		TransactionDate:  c.Date,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, balances)
	}
	return renderGLAccountBalances(rt, balances)
}

func renderGLAccountBalances(rt *Runtime, balances []api.GLAccountBalanceItem) error {
	rows := make([][]string, 0, len(balances))
	for _, balance := range balances {
		rows = append(rows, []string{balance.Code, balance.BalanceType, balance.Amount, balance.Description})
	}
	return output.Table(rt.Out, []string{"CODE", "TYPE", "AMOUNT", "DESCRIPTION"}, rows)
}

type GLAccountsTransactionsCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	GLAccount      string `name:"gl-account" help:"GL account code. Pass an empty value to include all accounts."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *GLAccountsTransactionsCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	transactions, err := client.GLAccountTransactions(rt.Context, sessionID, api.GLAccountTransactionsOptions{
		AdministrationID: administrationID,
		GLAccountCode:    c.GLAccount,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, transactions)
	}
	return renderGLAccountTransactions(rt, transactions)
}

type GLAccountsTransactionsFiscalCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	GLAccount      string `name:"gl-account" help:"GL account code. Pass an empty value to include all accounts."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *GLAccountsTransactionsFiscalCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	transactions, err := client.GLAccountTransactionsFiscal(rt.Context, sessionID, api.GLAccountTransactionsOptions{
		AdministrationID: administrationID,
		GLAccountCode:    c.GLAccount,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, transactions)
	}
	return renderGLAccountTransactions(rt, transactions)
}

type GLAccountsTransactionsWithContactCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	GLAccount      string `name:"gl-account" help:"GL account code. Pass an empty value to include all accounts."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *GLAccountsTransactionsWithContactCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	transactions, err := client.GLAccountTransactionsAndContact(rt.Context, sessionID, api.GLAccountTransactionsOptions{
		AdministrationID: administrationID,
		GLAccountCode:    c.GLAccount,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, transactions)
	}
	return renderGLAccountTransactionsWithFile(rt, transactions)
}

func renderGLAccountTransactions(rt *Runtime, transactions []api.GLAccountTransaction) error {
	rows := make([][]string, 0, len(transactions))
	for _, tx := range transactions {
		project := tx.Project.Text
		if project == "" {
			project = tx.Project.Code
		}
		rows = append(rows, []string{tx.Date, tx.GLAccountCode, tx.Amount, tx.Contact, project, tx.Description})
	}
	return output.Table(rt.Out, []string{"DATE", "GL", "AMOUNT", "CONTACT", "PROJECT", "DESCRIPTION"}, rows)
}

func renderGLAccountTransactionsWithFile(rt *Runtime, transactions []api.GLAccountTransaction) error {
	rows := make([][]string, 0, len(transactions))
	for _, tx := range transactions {
		project := tx.Project.Text
		if project == "" {
			project = tx.Project.Code
		}
		rows = append(rows, []string{tx.Date, tx.GLAccountCode, tx.Amount, tx.Contact, project, tx.FileName, tx.Description})
	}
	return output.Table(rt.Out, []string{"DATE", "GL", "AMOUNT", "CONTACT", "PROJECT", "FILE", "DESCRIPTION"}, rows)
}

type JournalsCmd struct {
	Process JournalsProcessCmd `cmd:"" help:"Process a general journal XML file."`
}

type JournalsProcessCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	File           string `name:"file" required:"" help:"Journal XML file to process." type:"existingfile"`
	DryRun         bool   `name:"dry-run" help:"Validate and preview the journal without authenticating or sending it."`
}

type journalXMLDocument struct {
	Path    string
	Content string
	Bytes   int
	Root    string
}

type journalProcessDryRun struct {
	DryRun           bool   `json:"dry_run"`
	Operation        string `json:"operation"`
	AdministrationID string `json:"administration_id"`
	File             string `json:"file"`
	Bytes            int    `json:"bytes"`
	Root             string `json:"root"`
	Message          string `json:"message"`
}

func (c *JournalsProcessCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}
	doc, err := readJournalXML(c.File)
	if err != nil {
		return err
	}
	if c.DryRun {
		return renderJournalDryRun(rt, globals, journalProcessDryRun{
			DryRun:           true,
			Operation:        "ProcessJournal",
			AdministrationID: administrationID,
			File:             doc.Path,
			Bytes:            doc.Bytes,
			Root:             doc.Root,
			Message:          "dry run; no journal sent",
		})
	}
	if globals.Readonly {
		return errors.New("--readonly blocks mutating command: accounting journals process")
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.ProcessJournal(rt.Context, sessionID, api.JournalImportOptions{
		AdministrationID: administrationID,
		XMLDoc:           doc.Content,
	})
	if err != nil {
		return err
	}
	return renderJournalProcessResult(rt, globals, result)
}

func readJournalXML(path string) (journalXMLDocument, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return journalXMLDocument{}, fmt.Errorf("read %s: %w", path, err)
	}
	root, err := validateXMLDocument(data)
	if err != nil {
		return journalXMLDocument{}, fmt.Errorf("validate %s: %w", path, err)
	}
	if root != "Journal" {
		return journalXMLDocument{}, fmt.Errorf("validate %s: expected root element Journal, got %s", path, root)
	}
	data = stripXMLDeclaration(data)
	return journalXMLDocument{
		Path:    path,
		Content: string(data),
		Bytes:   len(data),
		Root:    root,
	}, nil
}

func renderJournalDryRun(rt *Runtime, globals *Globals, result journalProcessDryRun) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"OPERATION", "ADMINISTRATION", "FILE", "BYTES", "ROOT", "MESSAGE"}, [][]string{{
		result.Operation,
		result.AdministrationID,
		result.File,
		strconv.Itoa(result.Bytes),
		result.Root,
		result.Message,
	}})
}

func renderJournalProcessResult(rt *Runtime, globals *Globals, result api.JournalProcessResult) error {
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"ADMINISTRATION", "DOCUMENT"}, [][]string{{
		result.AdministrationID,
		result.DocumentID,
	}})
}

type PeriodsCmd struct {
	Table        PeriodsTableCmd        `cmd:"" help:"Get the start date table for a financial year."`
	ModifiedDate PeriodsModifiedDateCmd `cmd:"" name:"modified-date" help:"Show when a financial year was last modified."`
}

type PeriodsTableCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Year           int    `name:"year" required:"" help:"Financial year ID."`
}

func (c *PeriodsTableCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	period, err := client.PeriodDateTable(rt.Context, sessionID, api.PeriodDateTableOptions{
		AdministrationID: administrationID,
		YearID:           c.Year,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, period)
	}
	return output.Table(rt.Out, []string{"ADMINISTRATION", "YEAR", "NAME", "PERIOD", "WHOLE PERIOD", "ISO8601"}, [][]string{{
		period.AdministrationID,
		strconv.Itoa(period.YearID),
		period.Name,
		period.Period,
		period.WholePeriod,
		output.Bool(period.ISO8601Period),
	}})
}

type PeriodsModifiedDateCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Year           int    `name:"year" required:"" help:"Financial year ID."`
}

func (c *PeriodsModifiedDateCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	result, err := client.FinancialYearModifiedDate(rt.Context, sessionID, api.PeriodDateTableOptions{
		AdministrationID: administrationID,
		YearID:           c.Year,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"ADMINISTRATION", "YEAR", "MODIFIED"}, [][]string{{
		result.AdministrationID,
		strconv.Itoa(result.YearID),
		result.ModifiedDate,
	}})
}

type RevenueCmd struct {
	Net       RevenueNetCmd       `cmd:"" help:"Get net revenue for a date range."`
	NetFiscal RevenueNetFiscalCmd `cmd:"" name:"net-fiscal" help:"Get fiscal net revenue for a date range."`
}

type RevenueNetCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *RevenueNetCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	report, err := client.NetRevenue(rt.Context, sessionID, api.RevenueOptions{
		AdministrationID: administrationID,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, report)
	}
	return renderRevenueReport(rt, report)
}

type RevenueNetFiscalCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	From           string `name:"from" required:"" help:"Start date, YYYY-MM-DD."`
	To             string `name:"to" required:"" help:"End date, YYYY-MM-DD."`
}

func (c *RevenueNetFiscalCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	report, err := client.NetRevenueFiscal(rt.Context, sessionID, api.RevenueOptions{
		AdministrationID: administrationID,
		StartDate:        c.From,
		EndDate:          c.To,
	})
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, report)
	}
	return renderRevenueReport(rt, report)
}

func renderRevenueReport(rt *Runtime, report api.RevenueReport) error {
	return output.Table(rt.Out, []string{"FROM", "TO", "AMOUNT"}, [][]string{{report.StartDate, report.EndDate, report.Amount}})
}

type ContactDefaultsCmd struct {
	List ContactDefaultsListCmd `cmd:"" help:"List default accounting values for a contact."`
}

type ContactDefaultsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
	Contact        string `name:"contact" required:"" help:"Contact ID."`
}

func (c *ContactDefaultsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	defaults, err := client.ContactDefaultValues(rt.Context, sessionID, administrationID, c.Contact)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, defaults)
	}
	return renderContactDefaults(rt, defaults)
}

func renderContactDefaults(rt *Runtime, defaults []api.ContactDefaultValues) error {
	rows := [][]string{}
	for _, contactDefaults := range defaults {
		for _, defaultValue := range contactDefaults.DefaultValues {
			rows = append(rows, []string{
				contactDefaults.ContactName,
				contactDefaults.DefaultBankAccount,
				defaultValue.InputFields.DocumentType,
				strconv.Itoa(defaultValue.InputFields.Priority),
				defaultValue.InputFields.Amount,
				defaultValue.InputFields.Currency,
				defaultValue.OutputFields.GLAccount,
				defaultValue.OutputFields.VATCode,
				defaultValue.OutputFields.PaymentMethod,
				defaultValue.OutputFields.PaymentTerm,
				defaultValue.Created,
			})
		}
		if len(contactDefaults.DefaultValues) == 0 {
			rows = append(rows, []string{
				contactDefaults.ContactName,
				contactDefaults.DefaultBankAccount,
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
				"",
			})
		}
	}
	return output.Table(rt.Out, []string{"CONTACT", "BANK", "DOCUMENT", "PRIORITY", "AMOUNT", "CURRENCY", "GL", "VAT", "PAYMENT", "TERM", "CREATED"}, rows)
}

func resolveAdministrationID(globals *Globals, explicit string) (string, error) {
	profile, err := loadProfile(globals)
	if err != nil {
		return "", err
	}
	administrationID := explicit
	if administrationID == "" {
		administrationID = profile.AdministrationID
	}
	if administrationID == "" {
		return "", errors.New("missing --administration; set it in config or pass --administration <id>")
	}
	return administrationID, nil
}

type AccountingPaymentMethodsCmd struct {
	List AccountingPaymentMethodsListCmd `cmd:"" help:"List custom payment methods for an administration."`
}

type AccountingPaymentMethodsListCmd struct {
	Administration string `help:"Administration ID. Defaults to profile/global administration."`
}

func (c *AccountingPaymentMethodsListCmd) Run(rt *Runtime, globals *Globals) error {
	administrationID, err := resolveAdministrationID(globals, c.Administration)
	if err != nil {
		return err
	}

	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	methods, err := client.CustomPaymentMethods(rt.Context, sessionID, administrationID)
	if err != nil {
		return err
	}
	return renderPaymentMethods(rt, globals, methods)
}

func renderPaymentMethods(rt *Runtime, globals *Globals, methods []api.PaymentMethod) error {
	if globals.JSON {
		return output.JSON(rt.Out, methods)
	}

	rows := make([][]string, 0, len(methods))
	for _, method := range methods {
		rows = append(rows, []string{method.ID, method.Description})
	}
	return output.Table(rt.Out, []string{"ID", "DESCRIPTION"}, rows)
}
