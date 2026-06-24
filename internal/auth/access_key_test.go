package auth

import (
	"context"
	"errors"
	"testing"
)

func TestResolveAccessKeyPrefersEnvironment(t *testing.T) {
	t.Setenv("YUKI_ACCESS_KEY", "env-key")
	store := &fakeStore{key: "stored-key"}

	key, source, err := ResolveAccessKey(context.Background(), store, "default")
	if err != nil {
		t.Fatalf("ResolveAccessKey: %v", err)
	}
	if key != "env-key" || source != SourceEnv {
		t.Fatalf("key/source = %q/%s", key, source)
	}
}

func TestResolveAccessKeyUsesStoreWhenEnvironmentEmpty(t *testing.T) {
	store := &fakeStore{key: "stored-key"}

	key, source, err := ResolveAccessKey(context.Background(), store, "default")
	if err != nil {
		t.Fatalf("ResolveAccessKey: %v", err)
	}
	if key != "stored-key" || source != SourceKeyring {
		t.Fatalf("key/source = %q/%s", key, source)
	}
}

func TestResolveAccessKeyReportsMissingCredentials(t *testing.T) {
	store := &fakeStore{err: ErrAccessKeyNotFound}

	_, _, err := ResolveAccessKey(context.Background(), store, "default")
	if !errors.Is(err, ErrAccessKeyNotFound) {
		t.Fatalf("err = %v, want ErrAccessKeyNotFound", err)
	}
}

func TestResolveAccessKeyRejectsEmptyStoredValue(t *testing.T) {
	store := &fakeStore{key: "   "}

	_, _, err := ResolveAccessKey(context.Background(), store, "default")
	if !errors.Is(err, ErrAccessKeyNotFound) {
		t.Fatalf("err = %v, want ErrAccessKeyNotFound", err)
	}
}

type fakeStore struct {
	key string
	err error
}

func (s *fakeStore) SetAccessKey(context.Context, string, string) error {
	return nil
}

func (s *fakeStore) AccessKey(_ context.Context, _ string) (string, error) {
	if s.err != nil {
		return "", s.err
	}
	return s.key, nil
}

func (s *fakeStore) DeleteAccessKey(context.Context, string) error {
	return nil
}
