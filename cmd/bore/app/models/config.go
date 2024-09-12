package models

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"go.trulyao.dev/bore/cmd/bore/app/styles"
)

type configDumpModel struct {
	table table.Model
}

func NewConfigDumpModel(t table.Model) *configDumpModel {
	return &configDumpModel{table: t}
}

func (m *configDumpModel) Init() tea.Cmd { return nil }

func (m *configDumpModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "esc", "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			row := m.table.SelectedRow()

			return m, tea.Batch(
				tea.Printf("Property: %s\nValue: %s\nDescription: %s\n", row[0], row[1], row[2]),
			)
		}
	}

	m.table, cmd = m.table.Update(message)
	return m, cmd
}

func (m *configDumpModel) View() string {
	return styles.BaseTableStyle.Render(m.table.View()) + "\n  " + m.table.HelpView() + "\n"
}
