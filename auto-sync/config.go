package main

import (
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

// Config represents on-disk YAML configuration.
// VaultPath: absolute path to vault; IntervalMinutes: sync interval; RepoURL: HTTPS remote (without token)
// Username: GitHub username (for informative purposes)
// RemoteName defaults to "origin"
type Config struct {
    VaultPath      string `yaml:"vault_path"`
    IntervalMinutes int    `yaml:"interval_minutes"`
    RepoURL        string `yaml:"repo_url"`
    Username       string `yaml:"username"`
    RemoteName     string `yaml:"remote_name"`
}

func defaultConfigDir() (string, error) {
    dir, err := os.UserConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(dir, "obsidian-auto-sync"), nil
}

func configFilePath() (string, error) {
    d, err := defaultConfigDir()
    if err != nil {
        return "", err
    }
    return filepath.Join(d, "config.yaml"), nil
}

func loadConfig() (*Config, error) {
    fp, err := configFilePath()
    if err != nil {
        return nil, err
    }
    b, err := os.ReadFile(fp)
    if err != nil {
        return nil, err
    }
    var c Config
    if err := yaml.Unmarshal(b, &c); err != nil {
        return nil, err
    }
    if c.RemoteName == "" {
        c.RemoteName = "origin"
    }
    if c.IntervalMinutes == 0 {
        c.IntervalMinutes = 60
    }
    return &c, nil
}

func saveConfig(c *Config) error {
    fp, err := configFilePath()
    if err != nil {
        return err
    }
    if err := os.MkdirAll(filepath.Dir(fp), 0o755); err != nil {
        return err
    }
    data, err := yaml.Marshal(c)
    if err != nil {
        return err
    }
    return os.WriteFile(fp, data, 0o600)
} 