package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/99designs/keyring"
	"golang.org/x/term"

	"github.com/dedene/yuki-cli/internal/config"
)

const (
	keyringBackendEnv  = "YUKI_KEYRING_BACKEND"
	keyringPasswordEnv = "YUKI_KEYRING_PASSWORD" //nolint:gosec // env var name, not a secret value
)

type KeyringStore struct {
	ring keyring.Keyring
}

func OpenDefault() (*KeyringStore, error) {
	ring, err := openKeyring()
	if err != nil {
		return nil, err
	}
	return &KeyringStore{ring: ring}, nil
}

func (s *KeyringStore) SetAccessKey(_ context.Context, profile string, accessKey string) error {
	accessKey = strings.TrimSpace(accessKey)
	if accessKey == "" {
		return errors.New("access key cannot be empty")
	}
	item := keyring.Item{
		Key:  accessKeyName(profile),
		Data: []byte(accessKey),
	}
	if err := s.ring.Set(item); err != nil {
		return fmt.Errorf("store access key: %w", err)
	}
	return nil
}

func (s *KeyringStore) AccessKey(_ context.Context, profile string) (string, error) {
	item, err := s.ring.Get(accessKeyName(profile))
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return "", ErrAccessKeyNotFound
		}
		return "", fmt.Errorf("read access key: %w", err)
	}
	accessKey := strings.TrimSpace(string(item.Data))
	if accessKey == "" {
		return "", ErrAccessKeyNotFound
	}
	return accessKey, nil
}

func (s *KeyringStore) DeleteAccessKey(_ context.Context, profile string) error {
	if err := s.ring.Remove(accessKeyName(profile)); err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return ErrAccessKeyNotFound
		}
		return fmt.Errorf("delete access key: %w", err)
	}
	return nil
}

func accessKeyName(profile string) string {
	if profile == "" {
		profile = "default"
	}
	return "access-key:" + profile
}

func openKeyring() (keyring.Keyring, error) {
	backend := strings.ToLower(strings.TrimSpace(os.Getenv(keyringBackendEnv)))
	backends, err := allowedBackends(backend)
	if err != nil {
		return nil, err
	}
	if runtime.GOOS == "darwin" && (backend == "" || backend == "auto") {
		backends = []keyring.BackendType{keyring.KeychainBackend}
	}
	if runtime.GOOS == "linux" && (backend == "" || backend == "auto") && os.Getenv("DBUS_SESSION_BUS_ADDRESS") == "" {
		backends = []keyring.BackendType{keyring.FileBackend}
	}

	keyringDir, err := config.KeyringDir()
	if err != nil {
		return nil, err
	}

	cfg := keyring.Config{
		ServiceName:              config.AppName,
		KeychainTrustApplication: true,
		AllowedBackends:          backends,
		FileDir:                  keyringDir,
		FilePasswordFunc:         filePasswordFunc(),
	}

	if usesFileBackend(backends) {
		if err := os.MkdirAll(keyringDir, 0o700); err != nil {
			return nil, fmt.Errorf("ensure keyring dir: %w", err)
		}
	}

	ring, err := keyring.Open(cfg)
	if err != nil {
		return nil, fmt.Errorf("open keyring: %w", err)
	}
	return ring, nil
}

func allowedBackends(backend string) ([]keyring.BackendType, error) {
	switch backend {
	case "", "auto":
		return nil, nil
	case "keychain":
		return []keyring.BackendType{keyring.KeychainBackend}, nil
	case "file":
		return []keyring.BackendType{keyring.FileBackend}, nil
	case "secret-service":
		return []keyring.BackendType{keyring.SecretServiceBackend}, nil
	case "wincred":
		return []keyring.BackendType{keyring.WinCredBackend}, nil
	default:
		return nil, fmt.Errorf("invalid keyring backend %q", backend)
	}
}

func usesFileBackend(backends []keyring.BackendType) bool {
	for _, backend := range backends {
		if backend == keyring.FileBackend {
			return true
		}
	}
	return false
}

func filePasswordFunc() keyring.PromptFunc {
	if password := os.Getenv(keyringPasswordEnv); password != "" {
		return keyring.FixedStringPrompt(password)
	}
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return keyring.TerminalPrompt
	}
	return func(string) (string, error) {
		return "", fmt.Errorf("no TTY available for file keyring password; set %s", keyringPasswordEnv)
	}
}
