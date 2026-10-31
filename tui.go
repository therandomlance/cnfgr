package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	catStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).Bold(true)
	selStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	rootStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("203"))
	titleStyl = lipgloss.NewStyle().Foreground(lipgloss.Color("213")).Bold(true)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	modeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("221")).Bold(true)
)

type listItem struct {
	category string
	program  Program
}

type mode int

const (
	modeList mode = iota
	modeEditor
)

type model struct {
	cfg            *Config
	cfgPath        string
	items          []listItem
	selected       int
	picked         *listItem
	editCfg        bool
	editorOverride string

	mode      mode
	editors   []string
	editorSel int
}

func initialModel(cfg *Config, cfgPath string) model {
	var items []listItem
	for _, c := range cfg.Categories {
		for _, p := range c.Programs {
			items = append(items, listItem{category: c.Name, program: p})
		}
	}
	return model{cfg: cfg, cfgPath: cfgPath, items: items}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	k, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if m.mode == modeEditor {
		return m.updateEditor(k)
	}
	return m.updateList(k)
}

func (m model) updateList(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "ctrl+c", "q", "esc":
		return m, tea.Quit
	case "up", "k":
		if len(m.items) > 0 {
			m.selected = (m.selected - 1 + len(m.items)) % len(m.items)
		}
	case "down", "j":
		if len(m.items) > 0 {
			m.selected = (m.selected + 1) % len(m.items)
		}
	case "home", "g":
		m.selected = 0
	case "end", "G":
		if len(m.items) > 0 {
			m.selected = len(m.items) - 1
		}
	case "ctrl+up":
		m.jumpCategory(-1)
	case "ctrl+down":
		m.jumpCategory(1)
	case "enter":
		if len(m.items) > 0 {
			p := m.items[m.selected]
			m.picked = &p
		}
		return m, tea.Quit
	case "ctrl+enter", "e":
		if len(m.items) > 0 && len(m.cfg.AvailableEditors) > 0 {
			m.editors = m.cfg.AvailableEditors
			m.editorSel = 0
			m.mode = modeEditor
			return m, tea.ClearScreen
		}
	case "c":
		m.editCfg = true
		return m, tea.Quit
	}
	return m, nil
}

func (m model) updateEditor(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "esc", "ctrl+c":
		m.mode = modeList
		return m, tea.ClearScreen
	case "up", "k":
		if m.editorSel > 0 {
			m.editorSel--
		}
	case "down", "j":
		if m.editorSel < len(m.editors)-1 {
			m.editorSel++
		}
	case "enter":
		if len(m.editors) > 0 {
			p := m.items[m.selected]
			m.picked = &p
			m.editorOverride = m.editors[m.editorSel]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m *model) jumpCategory(dir int) {
	n := len(m.items)
	if n == 0 {
		return
	}
	curCat := m.items[m.selected].category
	i := m.selected
	for {
		i = (i + dir + n) % n
		if m.items[i].category != curCat {
			break
		}
		if i == m.selected { // single category, nowhere to go
			return
		}
	}
	if dir < 0 {
		for j := i; j > 0 && m.items[j-1].category == m.items[i].category; j-- {
			i = j - 1
		}
	}
	m.selected = i
}

func (m model) View() string {
	if m.mode == modeEditor {
		return m.viewEditor()
	}
	return m.viewList()
}

func (m model) viewList() string {
	var b strings.Builder
	b.WriteString(titleStyl.Render(" cnfgr ") + "\n\n")
	if len(m.items) == 0 {
		fmt.Fprintf(&b, "No files configured.\n\nEdit config: %s\n(press 'c' to edit, 'q' to quit)", m.cfgPath)
		return b.String()
	}
	curCat := ""
	for i, it := range m.items {
		if it.category != curCat {
			if curCat != "" {
				b.WriteString("\n")
			}
			b.WriteString(catStyle.Render(it.category) + "\n")
			curCat = it.category
		}
		marker := "  "
		name := it.program.Name
		if i == m.selected {
			marker = "▶ "
			name = selStyle.Render(name)
		}
		line := marker + name
		if it.program.Root {
			line += " " + rootStyle.Render("[root]")
		}
		if it.program.Description != "" {
			line += "  " + descStyle.Render(it.program.Description)
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ navigate · ctrl+↑/↓ jump category · enter open · e pick editor · c edit config · q quit"))
	return b.String()
}

func (m model) viewEditor() string {
	var b strings.Builder
	b.WriteString(titleStyl.Render(" cnfgr ") + " " + modeStyle.Render("[ pick editor ]") + "\n\n")
	b.WriteString(helpStyle.Render(fmt.Sprintf("Opening: %s", m.items[m.selected].program.Name)) + "\n\n")
	for i, e := range m.editors {
		marker := "  "
		name := e
		if i == m.editorSel {
			marker = "▶ "
			name = selStyle.Render(name)
		}
		b.WriteString(marker + name + "\n")
	}
	b.WriteString("\n" + helpStyle.Render("↑/↓ select · enter confirm · esc cancel"))
	return b.String()
}
