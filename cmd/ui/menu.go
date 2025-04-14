package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

type MenuModel struct {
	choices  []string
	Cursor   int
	selected map[int]struct{}
}

func NewMenu(choices []string) MenuModel {
	return MenuModel{
		choices:  choices,
		selected: make(map[int]struct{}),
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.choices)-1 {
				m.Cursor++
			}
		case "enter", " ":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString(HeaderStyle.Render("Выберите действие:") + "\n\n")

	for i, choice := range m.choices {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}
		s.WriteString(MenuStyle.Render(fmt.Sprintf("%s%s\n", cursor, choice)))
	}

	s.WriteString("\n" + BaseStyle.Render("(↑/↓) для навигации • (enter) для выбора"))
	return s.String()
}
