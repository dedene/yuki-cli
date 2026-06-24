package auth

import (
	"context"
	"errors"
	"os"
	"strings"
)

const AccessKeyEnv = "YUKI_ACCESS_KEY"

var ErrAccessKeyNotFound = errors.New("yuki access key not found")

type Source string

const (
	SourceEnv     Source = "env"
	SourceKeyring Source = "keyring"
)

type Store interface {
	SetAccessKey(context.Context, string, string) error
	AccessKey(context.Context, string) (string, error)
	DeleteAccessKey(context.Context, string) error
}

func ResolveAccessKey(ctx context.Context, store Store, profile string) (string, Source, error) {
	if key, ok := EnvAccessKey(); ok {
		return key, SourceEnv, nil
	}
	if store == nil {
		return "", "", ErrAccessKeyNotFound
	}
	key, err := store.AccessKey(ctx, profile)
	if err != nil {
		return "", "", err
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return "", "", ErrAccessKeyNotFound
	}
	return key, SourceKeyring, nil
}

func EnvAccessKey() (string, bool) {
	key := strings.TrimSpace(os.Getenv(AccessKeyEnv))
	return key, key != ""
}
