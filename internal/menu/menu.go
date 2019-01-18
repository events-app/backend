package menu

import (
	"fmt"
)

// Menu holds Name and wired card.
type Menu struct {
	Name string `json:"name"`
	Card string `json:"card"`
}

// menus stores all the menus
var menus = []Menu{
	Menu{Name: "Main menu", Card: "main"},
	Menu{Name: "Secured", Card: "secured"},
	Menu{Name: "Other menu", Card: "other"},
}

// GetAll returns all menus
func GetAll() (*[]Menu, error) {
	if len(menus) == 0 {
		return nil, fmt.Errorf("cannot find any menu object")
	}
	return &menus, nil
}

// Find returns menu object based on name
func Find(name string) (*Menu, error) {
	for _, c := range menus {
		if c.Name == name {
			return &c, nil
		}
	}
	return nil, fmt.Errorf("cannot find a menu named: %s", name)
}

// Add appends new Menu object
func Add(name, card string) error {
	// Name of a menu must be unique for specyfic user.
	if !ValidateName(name) {
		return fmt.Errorf("name should be 4-30 characters long and should consists of letters, numbers, -, _")
	}
	if m, _ := Find(name); m != nil {
		return fmt.Errorf("menu with the name already exists")
	}
	menus = append(menus, Menu{Name: name, Card: card})

	return nil
}

// Update changes Content object based on name
// Returns error if it was not found
func Update(name, card string) error {
	for i := range menus {
		if menus[i].Name == name {
			menus[i].Card = card
			return nil
		}
	}
	return fmt.Errorf("menu not found")
}

func Delete(name string) error {
	if len(menus) == 0 {
		return fmt.Errorf("no menus in database")
	}
	var index int
	var found bool
	for i := range menus {
		if menus[i].Name == name {
			index = i
			found = true
		}
	}
	if !found {
		return fmt.Errorf("menu not found")
	}
	menus = append(menus[:index], menus[index+1:]...)
	return nil
}
