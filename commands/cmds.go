package commands

import (
	"submoduleop/collector"
	pb "submoduleop/protos"

	tea "github.com/charmbracelet/bubbletea"
)

type SaveErrMsg struct {
	Err error
}

type LoadErrMsg struct {
	Err error
}

type LoadSucessMsg struct {
	Submodules *pb.SubmoduleList
}

func SaveSubmodulesCmd(submodules *pb.SubmoduleList) tea.Cmd {
	return func() tea.Msg {
		if err := collector.SaveSubmodulesToFile(submodules); err != nil {
			return SaveErrMsg{Err: err}
		}
		return SaveErrMsg{Err: nil}
	}
}

func LoadSubmoduleFromFileCmd() tea.Cmd {
	return func() tea.Msg {
		submodules, err := collector.LoadSubmodulesFromFile()
		if err != nil {
			return LoadErrMsg{Err: err}
		}
		return LoadSucessMsg{Submodules: submodules}
	}
}
