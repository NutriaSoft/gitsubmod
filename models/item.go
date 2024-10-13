package models

import (
	"fmt"

	pb "submoduleop/protos"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var descriptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))

type Item struct {
	title, branch, path, url string
}

func (i Item) Title() string {
	return i.title
}

func (i Item) Description() string {
	desc := fmt.Sprintf("Branch: %s\nPath: %s\nURL: %s", i.branch, i.path, i.url)
	return descriptionStyle.Render(desc)
}

func (i Item) FilterValue() string {
	return i.title
}

func ItemsFromSubmodules(submodules *pb.SubmoduleList) []list.Item {
	items := make([]list.Item, len(submodules.Submodules))
	for i, sub := range submodules.Submodules {
		items[i] = Item{
			title:  sub.Name,
			branch: sub.Branch,
			url:    sub.Url,
			path:   sub.Path,
		}
	}
	return items
}
