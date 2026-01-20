package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Fixed dimensions for stable layout
const (
	contentWidth  = 45
	contentHeight = 16
	menuWidth     = 18
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginBottom(1).
			Align(lipgloss.Center)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			MarginRight(2).
			Width(menuWidth)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			Width(contentWidth).
			Height(contentHeight)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1).
			Align(lipgloss.Center)
)

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	// Build menu
	var menuItems []string
	for i, item := range MenuItems {
		page := Pages[item]
		label := fmt.Sprintf("%d. %s", i+1, page.Title)

		if i == m.cursor {
			menuItems = append(menuItems, selectedStyle.Render("> "+label))
		} else {
			menuItems = append(menuItems, normalStyle.Render("  "+label))
		}
	}
	// Set menu height to match content height for alignment
	menu := menuStyle.Height(contentHeight).Render(strings.Join(menuItems, "\n"))

	// Build content
	page := Pages[m.currentPage]
	content := contentStyle.Render(page.Content)

	// Layout: menu on left, content on right
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, menu, content)

	// Get the width of the main content for centering title and help
	contentWidth := lipgloss.Width(mainContent)

	// Title centered above content
	title := titleStyle.Width(contentWidth).Render("CapC0's Terminal")

	// Help text centered below content
	help := helpStyle.Width(contentWidth).Render("↑/↓ or j/k: navigate • enter: select • 1-4: quick access • q: quit")

	// Combine all elements vertically
	fullUI := lipgloss.JoinVertical(lipgloss.Center, title, mainContent, help)

	// Center the entire UI in the terminal
	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		fullUI,
	)
}
