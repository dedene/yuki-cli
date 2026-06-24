package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadReturnsDefaultsWhenConfigMissing(t *testing.T) {
	t.Setenv("YUKI_CONFIG_DIR", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.DefaultProfile != "default" {
		t.Fatalf("DefaultProfile = %q", cfg.DefaultProfile)
	}
	profile := cfg.Profile("default")
	if profile.BaseURL != DefaultBaseURL {
		t.Fatalf("BaseURL = %q, want %q", profile.BaseURL, DefaultBaseURL)
	}
}

func TestSaveAndLoadRoundTripProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("YUKI_CONFIG_DIR", dir)

	cfg := Default()
	cfg.DefaultProfile = "zenjoy"
	cfg.Profiles["zenjoy"] = Profile{
		BaseURL:          "https://api.yukiworks.nl/ws",
		AdministrationID: "admin-1",
		DomainID:         "domain-1",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "config.yaml")); err != nil {
		t.Fatalf("config file missing: %v", err)
	}

	got, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got.DefaultProfile != "zenjoy" {
		t.Fatalf("DefaultProfile = %q", got.DefaultProfile)
	}
	profile := got.Profile("zenjoy")
	if profile.BaseURL != "https://api.yukiworks.nl/ws" ||
		profile.AdministrationID != "admin-1" ||
		profile.DomainID != "domain-1" {
		t.Fatalf("profile = %#v", profile)
	}
}
