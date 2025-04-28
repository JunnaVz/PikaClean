// Package ui provides terminal user interface components for the PikaClean application.
// It includes various interactive elements like menus, forms, tables, and styled text
// components that create a consistent and user-friendly terminal experience.
package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableModel represents a tabular data display component for terminal interfaces.
// It implements the tea.Model interface from Bubble Tea and wraps the bubbles table
// component with consistent styling and behavior for the application.
type TableModel struct {
	table table.Model // The underlying bubbles table component being wrapped
}

// NewTable creates and initializes a new TableModel with the provided columns and rows data.
// It configures the table with application-specific styling and default settings.
//
// Parameters:
//   - columns: Column definitions including titles and formatting options
//   - rows: Data rows to display in the table
//
// Returns:
//   - TableModel: A ready-to-use table component
func NewTable(columns []table.Column, rows []table.Row) TableModel {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)

	t.SetStyles(style)

	return TableModel{table: t}
}

// Init is part of the tea.Model interface implementation.
// It defines the initial command to run when the table is first displayed.
//
// Returns:
//   - tea.Cmd: The command to execute on initialization (none for TableModel)
func (m TableModel) Init() tea.Cmd {
	return nil
}

// Update is part of the tea.Model interface implementation.
// It handles user input events, passing them to the underlying table component,
// and updates the table state accordingly. It also provides a way to quit the table
// view with keyboard shortcuts.
//
// Parameters:
//   - msg: The message containing user input or system event
//
// Returns:
//   - tea.Model: The updated model
//   - tea.Cmd: Command to execute after update
func (m TableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View is part of the tea.Model interface implementation.
// It renders the table as a formatted string for terminal display,
// applying the application's TableStyle to ensure consistent appearance.
//
// Returns:
//   - string: The formatted table display ready for terminal rendering
func (m TableModel) View() string {
	return TableStyle.Render(m.table.View())
}
