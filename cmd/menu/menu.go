// Package menu provides a command-line menu system for the PikaClean application.
// It implements an interactive menu with flexible options, handlers, and display
// capabilities to create consistent user interfaces across the application.
package menu

import (
	"fmt"
	"teamdev/cmd/cmdUtils"
)

// Menu represents an interactive command-line menu with items and handlers.
// It provides methods for displaying options, capturing user input, and
// invoking appropriate handlers based on selection.
type Menu struct {
	Items       []Item       // Collection of menu items to be displayed
	Handler     func() error // Optional function executed when menu is shown
	TableDrawer func() error // Optional function to draw tabular data
}

// AddItem appends a new item to the menu's list of available options.
//
// Parameters:
//   - item: The menu item to be added
func (m *Menu) AddItem(item Item) {
	m.Items = append(m.Items, item)
}

// CreateMenu initializes the menu with a predefined list of items,
// replacing any existing items.
//
// Parameters:
//   - items: Slice of menu items to populate the menu
func (m *Menu) CreateMenu(items []Item) {
	m.Items = items
}

// Print displays the menu items with sequential numbering on the console.
// Each item is shown with its assigned number and name, followed by
// a "0 -- exit" option at the end.
func (m *Menu) Print() {
	fmt.Print("Доступные действия:\n")
	for i, item := range m.Items {
		fmt.Printf("%d -- %s\n", i+1, item.Name)
	}
	fmt.Print("0 -- выход\n")
}

// validAction checks if the user-selected action is within the valid range.
//
// Parameters:
//   - action: The numeric selection made by the user
//
// Returns:
//   - bool: True if the action is valid (0 to number of items), false otherwise
func (m *Menu) validAction(action int) bool {
	return action >= 0 && action <= len(m.Items)
}

// Menu runs the interactive menu loop, displaying options and processing
// user selections until the exit option (0) is chosen. It validates inputs
// and executes the handler associated with the selected menu item.
//
// Returns:
//   - error: Any error that might occur during menu operation or nil on
//     successful exit
func (m *Menu) Menu() error {
	for {
		m.Print()

		action := cmdUtils.EndlessReadInt("Выберите действие")

		if action == 0 {
			return nil // exit action
		}

		if !m.validAction(action) {
			fmt.Println("Неверный номер")
			continue
		}

		err := m.Items[action-1].Handler()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
