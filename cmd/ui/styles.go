package ui

import "github.com/charmbracelet/lipgloss"

var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(0, 1)

var MenuStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205"))

var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("170"))
