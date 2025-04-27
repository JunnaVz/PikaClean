// Package menu provides a command-line menu system for the PikaClean application.
// It implements an interactive menu with flexible options, handlers, and display
// capabilities to create consistent user interfaces across the application.
package menu

// Item represents a selectable option in a command-line menu.
// Each menu item consists of a display name and an associated handler function
// that will be executed when the item is selected by the user.
//
// The handler function should encapsulate all logic needed for the operation,
// returning an error if the operation fails, or nil on success.
type Item struct {
	Name    string       // Displayed text for the menu option
	Handler func() error // Function to execute when this item is selected
}
