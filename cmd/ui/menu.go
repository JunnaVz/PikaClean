// Package ui provides terminal user interface components for the PikaClean application.
// It includes various interactive elements like menus, forms, tables, and styled text
// components that create a consistent and user-friendly terminal experience.
package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

// MenuModel represents an interactive menu component for terminal interfaces.
// It implements the tea.Model interface from Bubble Tea and provides a
// navigable list of choices with cursor-based selection.
type MenuModel struct {
	choices  []string         // Available menu options displayed to the user
	Cursor   int              // Current cursor position (highlighted menu item)
	selected map[int]struct{} // Keeps track of which items are selected (for multi-select menus)
}

// NewMenu creates and initializes a new MenuModel with the provided choices.
// It sets up the selection tracking map and prepares the component for rendering.
//
// Parameters:
//   - choices: String slice containing the menu options to display
//
// Returns:
//   - MenuModel: A ready-to-use menu component
func NewMenu(choices []string) MenuModel {
	return MenuModel{
		choices:  choices,
		selected: make(map[int]struct{}),
	}
}

// Init is part of the tea.Model interface implementation.
// It defines the initial command to run when the menu is first displayed.
//
// Returns:
//   - tea.Cmd: The command to execute on initialization (none for MenuModel)
func (m MenuModel) Init() tea.Cmd {
	return nil
}

// Update is part of the tea.Model interface implementation.
// It handles user input events and updates the menu state accordingly.
// Supports keyboard navigation with arrows or vim-like keys (j/k)
// and selection with enter or space.
//
// Parameters:
//   - msg: The message containing user input or system event
//
// Returns:
//   - tea.Model: The updated model (potentially with changed cursor position)
//   - tea.Cmd: Command to execute after update (quit on selection)
func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			// Move cursor up, with boundary check
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			// Move cursor down, with boundary check
			if m.Cursor < len(m.choices)-1 {
				m.Cursor++
			}
		case "enter", " ":
			// Selection made, exit the menu
			return m, tea.Quit
		}
	}
	return m, nil
}

// View is part of the tea.Model interface implementation.
// It renders the menu as a formatted string for terminal display.
// The current cursor position is indicated with a ">" character,
// and includes usage instructions at the bottom.
//
// Returns:
//   - string: The formatted menu display ready for terminal rendering
func (m MenuModel) View() string {
	s := strings.Builder{}
	s.WriteString(HeaderStyle.Render("Выберите действие:") + "\n\n")

	// Render each menu item, highlighting the active selection
	for i, choice := range m.choices {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}
		s.WriteString(MenuStyle.Render(fmt.Sprintf("%s%s\n", cursor, choice)))
	}

	// Add usage instructions at the bottom
	s.WriteString("\n" + BaseStyle.Render("(↑/↓) для навигации • (enter) для выбора"))
	return s.String()
}
