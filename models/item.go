package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
)

type Item struct {
	title, desc string
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	return i.desc
}

func (i Item) FilterValue() string {
	return i.title
}

func ItemsFromSubmodules(submodules []Submodule) []list.Item {
	items := make([]list.Item, len(submodules))
	for i, sub := range submodules {
		items[i] = Item{
			title: sub.Name,
			desc:  fmt.Sprintf("Branch: %s\nURL: %s", sub.Branch, sub.Url),
		}
	}
	return items
}
