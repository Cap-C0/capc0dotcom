package server

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"

	"capc0dotcom/terminal/tui"
)

//go:embed static/*
var staticFiles embed.FS

// WebServer handles HTTP and WebSocket connections
type WebServer struct {
	Host     string
	Port     int
	upgrader websocket.Upgrader
}

// NewWebServer creates a new web server
func NewWebServer(host string, port int) *WebServer {
	return &WebServer{
		Host: host,
		Port: port,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

// wsWriter wraps a WebSocket connection to implement io.Writer
type wsWriter struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func (w *wsWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	err = w.conn.WriteMessage(websocket.TextMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// wsReader wraps a WebSocket connection to implement io.Reader
type wsReader struct {
	conn    *websocket.Conn
	buf     []byte
	bufPos  int
	inputCh chan []byte
}

func newWSReader(conn *websocket.Conn) *wsReader {
	return &wsReader{
		conn:    conn,
		inputCh: make(chan []byte, 100),
	}
}

func (r *wsReader) Read(p []byte) (n int, err error) {
	// If we have buffered data, return it first
	if r.bufPos < len(r.buf) {
		n = copy(p, r.buf[r.bufPos:])
		r.bufPos += n
		return n, nil
	}

	// Wait for new data from channel
	data, ok := <-r.inputCh
	if !ok {
		return 0, io.EOF
	}

	r.buf = data
	r.bufPos = 0
	n = copy(p, r.buf)
	r.bufPos = n
	return n, nil
}

func (r *wsReader) receiveLoop() {
	defer close(r.inputCh)
	for {
		_, message, err := r.conn.ReadMessage()
		if err != nil {
			return
		}
		r.inputCh <- message
	}
}

// handleWebSocket handles WebSocket connections
func (s *WebServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection from %s", r.RemoteAddr)

	// Create writer and reader for the terminal
	writer := &wsWriter{conn: conn}
	reader := newWSReader(conn)

	// Start reading from WebSocket
	go reader.receiveLoop()

	// Create and run the Bubble Tea program
	model := tui.NewModel()
	p := tea.NewProgram(
		model,
		tea.WithInput(reader),
		tea.WithOutput(writer),
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Printf("Bubble Tea error: %v", err)
	}

	log.Printf("WebSocket connection closed from %s", r.RemoteAddr)
}

// serveStaticFile serves the embedded index.html
func (s *WebServer) serveStaticFile(w http.ResponseWriter, r *http.Request) {
	// Try to read from embedded files first
	content, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		// Fallback to reading from disk
		content, err = os.ReadFile("static/index.html")
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(content)
}

// Start starts the web server
func (s *WebServer) Start() error {
	http.HandleFunc("/", s.serveStaticFile)
	http.HandleFunc("/ws", s.handleWebSocket)

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	log.Printf("Web server starting on http://%s", addr)

	return http.ListenAndServe(addr, nil)
}

// StartAsync starts the web server in a goroutine
func (s *WebServer) StartAsync() error {
	http.HandleFunc("/", s.serveStaticFile)
	http.HandleFunc("/ws", s.handleWebSocket)

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	log.Printf("Web server starting on http://%s", addr)

	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Printf("Web server error: %v", err)
		}
	}()

	return nil
}
