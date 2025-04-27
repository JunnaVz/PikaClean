// Package ui provides terminal user interface components for the PikaClean application.
// It includes various interactive elements like menus, forms, tables, and styled text
// components that create a consistent and user-friendly terminal experience.
package ui

import "github.com/charmbracelet/lipgloss"

// BaseStyle defines the default styling for standard UI elements.
// It adds a normal border with a subtle gray color (240) to create
// visual separation between different UI components.
var BaseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// TableStyle defines the styling for tabular data displays.
// It inherits the border style from BaseStyle but adds horizontal
// and vertical padding to improve readability of tabular content.
var TableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(0, 1)

// MenuStyle defines the styling for interactive menu items.
// It uses a vibrant pink color (205) to make selectable options
// stand out from other terminal text.
var MenuStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("205"))

// HeaderStyle defines the styling for section headers and titles.
// It combines bold text with a purple color (170) to create visual
// hierarchy and help users navigate through different sections.
var HeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("170"))
