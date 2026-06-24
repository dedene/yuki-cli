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
	Login  AuthLoginCmd  `cmd:"" help:"Store a Yuki webservice access key."`
	Status AuthStatusCmd `cmd:"" help:"Show authentication status."`
	Logout AuthLogoutCmd `cmd:"" help:"Remove the stored access key."`
	Doctor AuthDoctorCmd `cmd:"" help:"Verify authentication with a read-only live check."`
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
