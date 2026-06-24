package cmd

import "github.com/dedene/yuki-cli/internal/output"

type LanguageCmd struct {
	Current   LanguageCurrentCmd   `cmd:"" help:"Show the current session language."`
	Supported LanguageSupportedCmd `cmd:"" help:"List supported session languages."`
	Set       LanguageSetCmd       `cmd:"" help:"Set the language for this session."`
}

type LanguageCurrentCmd struct{}

func (c *LanguageCurrentCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	language, err := client.Language(rt.Context, sessionID)
	if err != nil {
		return err
	}
	result := languageResult{Language: language}
	if globals.JSON {
		return output.JSON(rt.Out, result)
	}
	return output.Table(rt.Out, []string{"LANGUAGE"}, [][]string{{result.Language}})
}

type languageResult struct {
	Language string `json:"language"`
}

type LanguageSupportedCmd struct{}

func (c *LanguageSupportedCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	languages, err := client.SupportedLanguages(rt.Context, sessionID)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, languages)
	}

	rows := make([][]string, 0, len(languages))
	for _, language := range languages {
		rows = append(rows, []string{language.Code, language.Description, language.NativeDescription})
	}
	return output.Table(rt.Out, []string{"CODE", "DESCRIPTION", "NATIVE DESCRIPTION"}, rows)
}

type LanguageSetCmd struct {
	Language string `name:"language" required:"" help:"Session language code, e.g. en or nl-be."`
	DryRun   bool   `name:"dry-run" help:"Print the planned session update without authenticating or sending it."`
}

func (c *LanguageSetCmd) Run(rt *Runtime, globals *Globals) error {
	result := sessionSettingResult{
		Language: c.Language,
		Message:  "language set for this session",
	}
	if c.DryRun {
		result.DryRun = true
		result.Message = "dry run; language not sent"
		return renderSessionSetting(rt, globals, result)
	}
	client, sessionID, _, err := authenticatedSession(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	if err := client.SetLanguage(rt.Context, sessionID, c.Language); err != nil {
		return err
	}
	return renderSessionSetting(rt, globals, result)
}
