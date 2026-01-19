package tui

// Page represents a page in the terminal UI
type Page struct {
	Title   string
	Content string
}

// ASCII art logo
const Logo = `
   ____            ____  ___
  / ___|__ _ _ __ / ___|/ _ \
 | |   / _` + "`" + ` | '_ \| |   | | | |
 | |__| (_| | |_) | |___| |_| |
  \____\__,_| .__/ \____|\___/
            |_|
`

// Pages contains all static page content
var Pages = map[string]Page{
	"home": {
		Title: "Home",
		Content: Logo + `
Welcome to my terminal!

Navigate using arrow keys or j/k
Press Enter to select
Press 1-4 for quick access
Press q to quit
`,
	},
	"about": {
		Title: "About",
		Content: `
╔═══════════════════════════════════════╗
║              ABOUT ME                 ║
╚═══════════════════════════════════════╝

I figured I ought to have a website.

This is a terminal-based personal site
accessible via both SSH and web browser.

Built with Go and the Charm stack:
  • Bubble Tea for the TUI
  • Wish for SSH server
  • Lipgloss for styling
`,
	},
	"interests": {
		Title: "Interests",
		Content: `
╔═══════════════════════════════════════╗
║          THINGS I LIKE                ║
╚═══════════════════════════════════════╝

  ┌─────────────────────────────────┐
  │  Rust: It's type safe!          │
  └─────────────────────────────────┘

`,
	},
	"contact": {
		Title: "Contact",
		Content: `
╔═══════════════════════════════════════╗
║             CONTACT                   ║
╚═══════════════════════════════════════╝

Say hi on GitHub or drop a message.
I'm always open to interesting collaborations.

  GitHub:   https://github.com/sdm9252
  LinkedIn: https://linkedin.com/in/your-handle
  Email:    your@email.com

`,
	},
}

// MenuItems defines the navigation menu order
var MenuItems = []string{"home", "about", "interests", "contact"}
