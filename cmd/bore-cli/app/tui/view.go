package tui

import (
	"encoding/json"
	"io"

	"github.com/rivo/tview"
)

type Manager struct {
	application *tview.Application
}

func NewViewManager() *Manager {
	return &Manager{application: tview.NewApplication()}
}

func (m *Manager) RenderJSON(writer io.Writer, data any) error {
	var (
		serialized []byte
		err        error
	)

	switch d := data.(type) {
	case string:
		serialized = []byte(d)
	case []byte:
		serialized = d
	default:
		serialized, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			return err
		}
	}

	_, err = writer.Write(append(serialized, '\n'))
	return err
}
