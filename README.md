# cnfgr

A TUI for quickly opening config files. Browse categorized configs with keyboard navigation and open them in your editor of choice.

## Requirements

- Linux or macOS
- Go 1.21+

## Install

```sh
git clone https://github.com/therandomlance/cnfgr
cd cnfgr
./install.sh
```

## Configure

Config lives at `~/.config/cnfgr/config.toml`. cnfgr can open its own config file, the default editor is `nano`.

```sh
cnfgr -c
```
