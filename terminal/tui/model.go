package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the application state
type Model struct {
	cursor      int    // Current menu selection
	currentPage string // Currently displayed page
	width       int    // Terminal width
	height      int    // Terminal height
	quitting    bool   // Whether user is quitting
}

// NewModel creates a new Model with default state
func NewModel() Model {
	return Model{
		cursor:      0,
		currentPage: "home",
		width:       80,
		height:      24,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.currentPage = MenuItems[m.cursor]
			}

		case "down", "j":
			if m.cursor < len(MenuItems)-1 {
				m.cursor++
				m.currentPage = MenuItems[m.cursor]
			}

		case "enter":
			m.currentPage = MenuItems[m.cursor]

		case "1":
			m.cursor = 0
			m.currentPage = MenuItems[0]

		case "2":
			m.cursor = 1
			m.currentPage = MenuItems[1]

		case "3":
			m.cursor = 2
			m.currentPage = MenuItems[2]

		case "4":
			m.cursor = 3
			m.currentPage = MenuItems[3]
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}
