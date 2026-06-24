package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	AppName        = "yuki"
	DefaultBaseURL = "https://api.yukiworks.be/ws"
)

type Config struct {
	DefaultProfile string             `yaml:"default_profile"`
	Profiles       map[string]Profile `yaml:"profiles"`
}

type Profile struct {
	BaseURL          string `yaml:"base_url,omitempty"`
	AdministrationID string `yaml:"administration_id,omitempty"`
	DomainID         string `yaml:"domain_id,omitempty"`
}

func Default() Config {
	return Config{
		DefaultProfile: "default",
		Profiles: map[string]Profile{
			"default": {BaseURL: DefaultBaseURL},
		},
	}
}

func Load() (Config, error) {
	path, err := ConfigFile()
	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Default(), nil
		}
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	if len(data) == 0 {
		return Default(), nil
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	cfg.normalize()
	return cfg, nil
}

func Save(cfg Config) error {
	cfg.normalize()
	dir, err := ConfigDir()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("ensure config dir: %w", err)
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func (c Config) Profile(name string) Profile {
	if name == "" {
		name = c.DefaultProfile
	}
	profile, ok := c.Profiles[name]
	if !ok {
		profile = Profile{}
	}
	if profile.BaseURL == "" {
		profile.BaseURL = DefaultBaseURL
	}
	return profile
}

func ConfigDir() (string, error) {
	if dir := os.Getenv("YUKI_CONFIG_DIR"); dir != "" {
		return dir, nil
	}
	userDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("determine user config dir: %w", err)
	}
	return filepath.Join(userDir, AppName), nil
}

func ConfigFile() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func KeyringDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "keyring"), nil
}

func (c *Config) normalize() {
	if c.DefaultProfile == "" {
		c.DefaultProfile = "default"
	}
	if c.Profiles == nil {
		c.Profiles = map[string]Profile{}
	}
	profile := c.Profiles[c.DefaultProfile]
	if profile.BaseURL == "" {
		profile.BaseURL = DefaultBaseURL
	}
	c.Profiles[c.DefaultProfile] = profile
}
