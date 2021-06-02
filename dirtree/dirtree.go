package dirtree

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/knipferrc/fm/icons"
	"github.com/knipferrc/fm/utils"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type directoryMsg []fs.FileInfo

type Model struct {
	Files               []fs.FileInfo
	Cursor              int
	ShowIcons           bool
	ShowHidden          bool
	SelectedItemColor   string
	UnselectedItemColor string
}

func NewModel(files []fs.FileInfo, showIcons bool, selectedItemColor, unselectedItemColor string) Model {
	return Model{
		Files:               files,
		Cursor:              0,
		ShowIcons:           showIcons,
		ShowHidden:          true,
		SelectedItemColor:   selectedItemColor,
		UnselectedItemColor: unselectedItemColor,
	}
}

func (m *Model) SetContent(files []fs.FileInfo) {
	m.Files = files
}

func (m *Model) GotoTop() {
	m.Cursor = 0
}

func (m *Model) GotoBottom() {
	m.Cursor = len(m.Files) - 1
}

func (m Model) GetSelectedFile() fs.FileInfo {
	return m.Files[m.Cursor]
}

func (m Model) GetCursor() int {
	return m.Cursor
}

func (m *Model) GoDown() {
	m.Cursor++
}

func (m *Model) GoUp() {
	m.Cursor--
}

func (m Model) GetTotalFiles() int {
	return len(m.Files)
}

func (m *Model) ToggleHidden() {
	m.ShowHidden = !m.ShowHidden
}

func (m Model) Init() tea.Cmd {
	return nil
}

func updateDirectoryListing(dir string, showHidden bool) tea.Cmd {
	return func() tea.Msg {
		files := utils.GetDirectoryListing(dir, showHidden)

		return directoryMsg(files)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case directoryMsg:
		m.Files = msg
		m.Cursor = 0

	case tea.KeyMsg:
		switch msg.String() {
		case ".":
			m.ToggleHidden()

			return m, updateDirectoryListing(".", m.ShowHidden)
		}
	}

	return m, cmd
}

func (m Model) dirItem(selected bool, file fs.FileInfo) string {
	if !m.ShowIcons && !selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name())
	} else if !m.ShowIcons && selected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name())
	} else if selected && file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))

		return lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(listing)
	} else if !selected && file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))

		return listing
	} else if selected && !file.IsDir() {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.SelectedItemColor)).
			Render(file.Name()))

		return listing
	} else {
		icon, color := icons.GetIcon(file.Name(), filepath.Ext(file.Name()), icons.GetIndicator(file.Mode()))
		fileIcon := fmt.Sprintf("%s%s", color, icon)
		listing := fmt.Sprintf("%s\033[0m %s", fileIcon, lipgloss.NewStyle().
			Foreground(lipgloss.Color(m.UnselectedItemColor)).
			Render(file.Name()))

		return listing
	}
}

func (m Model) View() string {
	doc := strings.Builder{}
	curFiles := ""

	for i, file := range m.Files {
		curFiles += fmt.Sprintf("%s\n", m.dirItem(m.Cursor == i, file))
	}

	doc.WriteString(curFiles)

	return doc.String()
}
