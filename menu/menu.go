package menu

import (
	"fmt"
	"lab3/cmd/cmdUtils"
)

type Menu struct {
	Items       []Item
	Handler     func() error
	TableDrawer func() error
}

func (m *Menu) AddItem(item Item) {
	m.Items = append(m.Items, item)
}

func (m *Menu) CreateMenu(items []Item) {
	m.Items = items
}

func (m *Menu) Print() {
	fmt.Print("Доступные действия:\n")
	for i, item := range m.Items {
		fmt.Printf("%d -- %s\n", i+1, item.Name)
	}
	fmt.Print("0 -- выход\n")
}

func (m *Menu) validAction(action int) bool {
	return action >= 0 && action <= len(m.Items)
}

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
