package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/dedene/yuki-cli/internal/auth"
	"github.com/dedene/yuki-cli/internal/output"
)

type AuthCmd struct {
	Login   AuthLoginCmd   `cmd:"" help:"Store a Yuki webservice access key."`
	Session AuthSessionCmd `cmd:"" help:"Create one-shot Yuki session IDs."`
	Status  AuthStatusCmd  `cmd:"" help:"Show authentication status."`
	Logout  AuthLogoutCmd  `cmd:"" help:"Remove the stored access key."`
	Doctor  AuthDoctorCmd  `cmd:"" help:"Verify authentication with a read-only live check."`
}

type AuthLoginCmd struct {
	AccessKey string `name:"access-key" help:"Yuki WebserviceAccessKey. Omit to prompt securely."`
}

func (c *AuthLoginCmd) Run(rt *Runtime, globals *Globals) error {
	accessKey := strings.TrimSpace(c.AccessKey)
	if accessKey == "" {
		if globals.NoInput {
			return errors.New("missing --access-key in --no-input mode")
		}
		if !term.IsTerminal(int(os.Stdin.Fd())) {
			return errors.New("no TTY available; pass --access-key or set YUKI_ACCESS_KEY")
		}
		fmt.Fprint(rt.Err, "Yuki access key: ")
		data, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(rt.Err)
		if err != nil {
			return fmt.Errorf("read access key: %w", err)
		}
		accessKey = strings.TrimSpace(string(data))
	}

	store, err := rt.store()
	if err != nil {
		return err
	}
	if err := store.SetAccessKey(rt.Context, globals.Profile, accessKey); err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, map[string]any{"status": "saved", "profile": globals.Profile})
	}
	_, err = fmt.Fprintf(rt.Out, "Access key saved for profile %q.\n", globals.Profile)
	return err
}

type AuthSessionCmd struct {
	Client   AuthSessionClientCmd   `cmd:"" help:"Print a session ID using developer client credentials."`
	Username AuthSessionUsernameCmd `cmd:"" help:"Print a session ID using a Yuki username and password."`
}

type AuthSessionClientCmd struct {
	ClientID     string `name:"client-id" required:"" help:"Yuki developer ClientID."`
	ClientSecret string `name:"client-secret" help:"Yuki developer ClientSecret. Omit to prompt securely."`
	AccessKey    string `name:"access-key" help:"Yuki access key. Defaults to env/keyring when omitted."`
}

func (c *AuthSessionClientCmd) Run(rt *Runtime, globals *Globals) error {
	clientSecret, err := secretValue(c.ClientSecret, globals, rt, "client-secret", "Yuki client secret")
	if err != nil {
		return err
	}
	accessKey := strings.TrimSpace(c.AccessKey)
	if accessKey == "" {
		var sourceErr error
		accessKey, _, sourceErr = resolveAccessKey(rt.Context, rt, globals.Profile)
		if sourceErr != nil {
			if errors.Is(sourceErr, auth.ErrAccessKeyNotFound) {
				return fmt.Errorf("%w; pass --access-key, run 'yuki auth login --access-key <key>', or set %s", sourceErr, auth.AccessKeyEnv)
			}
			return sourceErr
		}
	}
	profile, err := loadProfile(globals)
	if err != nil {
		return err
	}
	client := rt.client(globals, profile)
	sessionID, err := client.AuthenticateClient(rt.Context, c.ClientID, clientSecret, accessKey)
	if err != nil {
		return err
	}
	return renderSessionID(rt, globals, sessionID)
}

type AuthSessionUsernameCmd struct {
	Username string `name:"username" required:"" help:"Yuki username or email address."`
	Password string `name:"password" help:"Yuki password. Omit to prompt securely."`
}

func (c *AuthSessionUsernameCmd) Run(rt *Runtime, globals *Globals) error {
	password, err := secretValue(c.Password, globals, rt, "password", "Yuki password")
	if err != nil {
		return err
	}
	profile, err := loadProfile(globals)
	if err != nil {
		return err
	}
	client := rt.client(globals, profile)
	sessionID, err := client.AuthenticateByUserName(rt.Context, c.Username, password)
	if err != nil {
		return err
	}
	return renderSessionID(rt, globals, sessionID)
}

func secretValue(value string, globals *Globals, rt *Runtime, flagName, prompt string) (string, error) {
	secret := strings.TrimSpace(value)
	if secret != "" {
		return secret, nil
	}
	if globals.NoInput {
		return "", fmt.Errorf("missing --%s in --no-input mode", flagName)
	}
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", fmt.Errorf("no TTY available; pass --%s", flagName)
	}
	fmt.Fprintf(rt.Err, "%s: ", prompt)
	data, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(rt.Err)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", strings.ToLower(prompt), err)
	}
	secret = strings.TrimSpace(string(data))
	if secret == "" {
		return "", fmt.Errorf("missing %s", strings.ToLower(prompt))
	}
	return secret, nil
}

func renderSessionID(rt *Runtime, globals *Globals, sessionID string) error {
	payload := map[string]string{"session_id": sessionID}
	if globals.JSON {
		return output.JSON(rt.Out, payload)
	}
	return output.Table(rt.Out, []string{"SESSION ID"}, [][]string{{sessionID}})
}

type AuthStatusCmd struct{}

func (c *AuthStatusCmd) Run(rt *Runtime, globals *Globals) error {
	_, source, err := resolveAccessKey(rt.Context, rt, globals.Profile)
	authenticated := err == nil
	if errors.Is(err, auth.ErrAccessKeyNotFound) {
		err = nil
	}
	if err != nil {
		return err
	}

	payload := map[string]any{
		"authenticated": authenticated,
		"profile":       globals.Profile,
	}
	if authenticated {
		payload["source"] = string(source)
	}
	if globals.JSON {
		return output.JSON(rt.Out, payload)
	}
	if authenticated {
		_, err = fmt.Fprintf(rt.Out, "Authenticated via %s for profile %q.\n", source, globals.Profile)
		return err
	}
	_, err = fmt.Fprintf(rt.Out, "Not authenticated. Run 'yuki auth login --access-key <key>' or set %s.\n", auth.AccessKeyEnv)
	return err
}

type AuthLogoutCmd struct{}

func (c *AuthLogoutCmd) Run(rt *Runtime, globals *Globals) error {
	store, err := rt.store()
	if err != nil {
		return err
	}
	err = store.DeleteAccessKey(rt.Context, globals.Profile)
	if err != nil {
		return err
	}
	if globals.JSON {
		return output.JSON(rt.Out, map[string]any{"status": "deleted", "profile": globals.Profile})
	}
	_, err = fmt.Fprintf(rt.Out, "Access key removed for profile %q.\n", globals.Profile)
	return err
}

type AuthDoctorCmd struct{}

func (c *AuthDoctorCmd) Run(rt *Runtime, globals *Globals) error {
	client, sessionID, err := authenticatedClient(rt.Context, rt, globals)
	if err != nil {
		return err
	}
	domain, domainErr := client.CurrentDomain(rt.Context, sessionID)
	payload := map[string]any{
		"authenticated": true,
		"profile":       globals.Profile,
	}
	if domainErr == nil {
		payload["current_domain"] = domain
	}
	if globals.JSON {
		if domainErr != nil {
			payload["warning"] = domainErr.Error()
		}
		return output.JSON(rt.Out, payload)
	}
	if _, err := fmt.Fprintln(rt.Out, "Authentication OK."); err != nil {
		return err
	}
	if domainErr == nil {
		_, err = fmt.Fprintf(rt.Out, "Current domain: %s (%s)\n", domain.Name, domain.ID)
		return err
	}
	_, err = fmt.Fprintf(rt.Out, "Current domain check skipped: %v\n", domainErr)
	return err
}
