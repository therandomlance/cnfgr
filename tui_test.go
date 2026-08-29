package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func catItems(cats ...string) []listItem {
	var items []listItem
	for _, c := range cats {
		items = append(items, listItem{category: c})
	}
	return items
}

func TestJumpCategoryWraps(t *testing.T) {
	cases := []struct {
		items []listItem
		start int
		dir   int
		want  int
	}{
		{catItems("a", "a", "b", "c"), 0, -1, 3},  // up from first cat -> last cat
		{catItems("a", "b", "c", "c"), 2, 1, 0},   // down from last cat -> first cat
		{catItems("a", "a", "b"), 1, 1, 2},        // normal forward jump
		{catItems("a", "b", "b"), 1, -1, 0},       // backward jump lands on block start
		{catItems("a", "a"), 0, 1, 0},             // single category: no-op
		{catItems("a", "b"), 0, -1, 1},            // two cats, wrap up
	}
	for _, c := range cases {
		m := model{items: c.items, selected: c.start}
		m.jumpCategory(c.dir)
		if m.selected != c.want {
			t.Errorf("start=%d dir=%d: got %d, want %d", c.start, c.dir, m.selected, c.want)
		}
	}
}

func TestUpDownWraps(t *testing.T) {
	m := model{items: catItems("a", "b", "c")}
	next, _ := m.updateList(tea.KeyMsg{Type: tea.KeyUp})
	if m := next.(model); m.selected != 2 {
		t.Errorf("up from top: got %d, want 2", m.selected)
	}
	next, _ = next.(model).updateList(tea.KeyMsg{Type: tea.KeyDown})
	if m := next.(model); m.selected != 0 {
		t.Errorf("down from bottom: got %d, want 0", m.selected)
	}
}
