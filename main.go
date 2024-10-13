package main

import (
	"fmt"
	"log"
	"submoduleop/collector"
	"submoduleop/commands"
	"submoduleop/models"

	pb "submoduleop/protos"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type screenState int

var (
	appStyle           = lipgloss.NewStyle().Padding(1, 2)
	titleStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFDF5")).Background(lipgloss.Color("#25A065")).Padding(0, 1)
	statusMessageStyle = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).Render
)

const (
	screenList screenState = iota
	screenInput
)

type model struct {
	list        list.Model
	screen      screenState
	inputURL    string
	inputBranch string
	inputName   string
	cursor      int
	err         error
	submodules  *pb.SubmoduleList
}

func initialModel() model {
	submodules, err := collector.LoadSubmodulesFromFile()
	if err != nil {
		// TODO: Handle error
		log.Println(err)
	}
	itemList := models.ItemsFromSubmodules(submodules)
	delegate := list.NewDefaultDelegate()

	submoduleList := list.New(itemList, delegate, 20, 10)
	submoduleList.Title = "Submodules"
	submoduleList.SetShowStatusBar(false)
	submoduleList.SetShowFilter(false)
	submoduleList.SetShowHelp(false)
	submoduleList.Styles.Title = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))

	return model{
		list:   submoduleList,
		screen: screenList,
		cursor: 0,
	}
}

func (m model) Init() tea.Cmd {
	return commands.LoadSubmoduleFromFileCmd()
}

func (m model) renderInputView() string {
	var urlField, branchField, nameField string
	urlField = "URL: " + m.inputURL
	branchField = "Branch: " + m.inputBranch
	nameField = "Name: " + m.inputName

	if m.cursor == 0 {
		urlField = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(urlField)
	}
	if m.cursor == 1 {
		branchField = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(branchField)
	}
	if m.cursor == 2 {
		nameField = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).Render(nameField)
	}

	return fmt.Sprintf(
		"Add a new submodule:\n\n%s\n%s\n%s\n\nPress Enter to save or 'q' to cancel.",
		urlField, branchField, nameField,
	)
}

func (m model) View() string {
	if m.screen == screenInput {
		return m.renderInputView()
	}

	return m.list.View() + "\nPress 'a' to add a new submodule or 'q' to quit.\n"
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case commands.LoadSucessMsg:
		m.submodules = msg.Submodules
		m.list.SetItems(models.ItemsFromSubmodules(m.submodules))
		m.err = nil
		return m, nil
	case commands.LoadErrMsg:
		m.err = msg.Err
		return m, nil
	case commands.SaveErrMsg:
		if msg.Err != nil {
			m.err = msg.Err
			return m, nil
		}
		m.err = nil
		if m.screen == screenInput {
			m.screen = screenList
			m.inputURL = ""
			m.inputBranch = ""
			m.inputName = ""
		}
		m.list.SetItems(models.ItemsFromSubmodules(m.submodules))
		return m, nil
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		if m.screen == screenList {
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit

			case "a":
				m.screen = screenInput
				return m, nil

			case "delete":
				if len(m.submodules.Submodules) > 0 {
					index := m.list.Index()
					name := m.submodules.Submodules[index].Name
					collector.DeleteSubmodule(m.submodules, name)
					return m, commands.SaveSubmodulesCmd(m.submodules)
				}
				return m, nil

			case "enter":
				return m, nil
			}
		}

		if m.screen == screenInput {
			switch msg.String() {
			case "enter":
				newModule := &pb.Submodule{
					Url:    m.inputURL,
					Branch: m.inputBranch,
					Name:   m.inputName,
				}
				collector.AddSubmodule(newModule, m.submodules)
				return m, commands.SaveSubmodulesCmd(m.submodules)

			case "up", "down":
				if msg.String() == "up" && m.cursor > 0 {
					m.cursor--
				} else if msg.String() == "down" && m.cursor < 2 {
					m.cursor++
				}

			case "backspace":
				if m.cursor == 0 && len(m.inputURL) > 0 {
					m.inputURL = m.inputURL[:len(m.inputURL)-1]
				} else if m.cursor == 1 && len(m.inputBranch) > 0 {
					m.inputBranch = m.inputBranch[:len(m.inputBranch)-1]
				} else if m.cursor == 2 && len(m.inputName) > 0 {
					m.inputName = m.inputName[:len(m.inputName)-1]
				}

			default:
				if msg.Type == tea.KeyRunes {
					switch m.cursor {
					case 0: // URL input
						m.inputURL += msg.String()
					case 1: // Branch input
						m.inputBranch += msg.String()
					case 2: // Name input
						m.inputName += msg.String()
					}
				}
			}
		}
	}

	if m.screen == screenList {
		m.list, cmd = m.list.Update(msg)
	}
	return m, cmd
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
