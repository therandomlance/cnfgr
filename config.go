package main

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Program struct {
	Name        string `toml:"name"`
	Path        string `toml:"path"`
	Description string `toml:"description"`
	Root        bool   `toml:"root"`
	Editor      string `toml:"editor"`
}

type Category struct {
	Name     string    `toml:"name"`
	Programs []Program `toml:"program"`
}

type Config struct {
	DefaultEditor    string     `toml:"default_editor"`    // editor used when no per-program override is set
	AvailableEditors []string   `toml:"available_editors"` // editors shown in the TUI picker (press 'e')
	Categories       []Category `toml:"category"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "cnfgr", "config.toml"), nil
}

func loadConfig() (*Config, string, error) {
	p, err := configPath()
	if err != nil {
		return nil, "", err
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		if err := writeDefaultConfig(p); err != nil {
			return nil, p, err
		}
	}
	var cfg Config
	if _, err := toml.DecodeFile(p, &cfg); err != nil {
		return nil, p, err
	}
	return &cfg, p, nil
}

func writeDefaultConfig(p string) error {
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	return os.WriteFile(p, []byte(defaultConfig), 0o644)
}

const defaultConfig = `# cnfgr configuration
# Edit freely. Run "cnfgr config" or press 'c' in the TUI to reopen.

# default_editor: editor used when no per-program override is set
default_editor = "nano"

# available_editors: editors shown in the TUI picker (press 'e')
available_editors = ["nano", "vim"]

# Each [[category]] groups its [[category.program]] entries.
# Fields per program:
#   name        (required) label shown in the TUI
#   path        (required) file to open; ~ and $ENV are expanded
#   description (optional) short note shown next to the name
#   root        (optional) true => opened with sudo. Default false
#   editor      (optional) overrides the top-level editor for this entry

[[category]]
name = "System"

[[category.program]]
name = "fstab"
path = "/etc/fstab"
root = true
description = "filesystem mount table"

[[category.program]]
name = "hosts"
path = "/etc/hosts"
root = true

[[category]]
name = "Shell"

[[category.program]]
name = "bashrc"
path = "~/.bashrc"
description = "bash configuration"

[[category.program]]
name = "profile"
path = "~/.profile"
editor = "vim"
`
