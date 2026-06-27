package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const VERSION = "v0.1.0"

func main() {
	var (
		editCfg  bool
		editName string
		showVer  bool
		initCfg  bool
	)
	flag.BoolVar(&editCfg, "c", false, "edit cnfgr's own config")
	flag.StringVar(&editName, "e", "", "open a specific config file")
	flag.BoolVar(&showVer, "v", false, "show version")
	flag.BoolVar(&initCfg, "init", false, "create the default config if it does not exist")
	flag.Parse()

	switch {
	case showVer:
		fmt.Println(VERSION)
		return
	case initCfg:
		initConfig()
		return
	case editCfg:
		editConfig()
		return
	case editName != "":
		editByName(editName)
		return
	}
	runTUI()
}

func editByName(name string) {
	cfg, _, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading config:", err)
		os.Exit(1)
	}
	for _, cat := range cfg.Categories {
		for _, prog := range cat.Programs {
			if prog.Name == name {
				if err := openFile(prog.Path, prog.Root, editorFor(prog, cfg.DefaultEditor)); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
					os.Exit(1)
				}
				return
			}
		}
	}
	fmt.Fprintln(os.Stderr, "no config named:", name)
	os.Exit(1)
}

func initConfig() {
	p, err := configPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if _, err := os.Stat(p); err == nil {
		fmt.Println("config already exists:", p)
		return
	}
	if err := writeDefaultConfig(p); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println("created config:", p)
}

func editConfig() {
	cfg, p, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading config:", err)
		os.Exit(1)
	}
	if err := openFile(p, false, resolveEditor(cfg.DefaultEditor)); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func runTUI() {
	cfg, cfgPath, err := loadConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading config:", err)
		os.Exit(1)
	}
	p := tea.NewProgram(initialModel(cfg, cfgPath), tea.WithAltScreen())
	final, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	m, _ := final.(model)
	switch {
	case m.editCfg:
		err = openFile(cfgPath, false, resolveEditor(cfg.DefaultEditor))
	case m.picked != nil:
		ed := editorFor(m.picked.program, cfg.DefaultEditor)
		if m.editorOverride != "" {
			ed = m.editorOverride
		}
		err = openFile(m.picked.program.Path, m.picked.program.Root, ed)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func editorFor(prog Program, cfgEditor string) string {
	if prog.Editor != "" {
		return prog.Editor
	}
	return resolveEditor(cfgEditor)
}

func resolveEditor(cfgEditor string) string {
	if cfgEditor != "" {
		return cfgEditor
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if _, err := exec.LookPath("nano"); err == nil {
		return "nano"
	}
	return "vi"
}

func expandPath(p string) (string, error) {
	if strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, p[2:]), nil
	}
	if p == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return home, nil
	}
	return os.ExpandEnv(p), nil
}

func openFile(path string, root bool, editor string) error {
	var err error
	path, err = expandPath(path)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	if root {
		cmd = exec.Command("sudo", editor, path)
	} else {
		cmd = exec.Command(editor, path)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
