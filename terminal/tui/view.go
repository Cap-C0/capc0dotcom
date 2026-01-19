package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("212")).
			MarginBottom(1)

	menuStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			MarginRight(2)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("212")).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	contentStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)
)

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("CapC0's Terminal"))
	b.WriteString("\n\n")

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
	menu := menuStyle.Render(strings.Join(menuItems, "\n"))

	// Build content
	page := Pages[m.currentPage]
	content := contentStyle.Render(page.Content)

	// Layout: menu on left, content on right
	layout := lipgloss.JoinHorizontal(lipgloss.Top, menu, content)
	b.WriteString(layout)

	// Help text
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("↑/↓ or j/k: navigate • enter: select • 1-4: quick access • q: quit"))

	return b.String()
}
