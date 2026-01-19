package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"

	"capc0dotcom/terminal/tui"
)

// SSHServer handles SSH connections
type SSHServer struct {
	Host string
	Port int
}

// NewSSHServer creates a new SSH server
func NewSSHServer(host string, port int) *SSHServer {
	return &SSHServer{
		Host: host,
		Port: port,
	}
}

// teaHandler returns the Bubble Tea middleware handler
func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	pty, _, _ := s.Pty()

	model := tui.NewModel()
	// Set initial window size from PTY
	updatedModel, _ := model.Update(tea.WindowSizeMsg{
		Width:  pty.Window.Width,
		Height: pty.Window.Height,
	})

	return updatedModel, []tea.ProgramOption{tea.WithAltScreen()}
}

// Start starts the SSH server
func (s *SSHServer) Start() error {
	srv, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", s.Host, s.Port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create SSH server: %w", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("SSH server starting on %s:%d", s.Host, s.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Printf("SSH server error: %v", err)
			done <- nil
		}
	}()

	<-done

	log.Println("Shutting down SSH server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}

// StartAsync starts the SSH server in a goroutine
func (s *SSHServer) StartAsync() error {
	srv, err := wish.NewServer(
		wish.WithAddress(fmt.Sprintf("%s:%d", s.Host, s.Port)),
		wish.WithHostKeyPath(".ssh/term_info_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create SSH server: %w", err)
	}

	log.Printf("SSH server starting on %s:%d", s.Host, s.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != ssh.ErrServerClosed {
			log.Printf("SSH server error: %v", err)
		}
	}()

	return nil
}
