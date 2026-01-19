package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"capc0dotcom/terminal/server"
)

func main() {
	// Parse command line flags
	sshPort := flag.Int("ssh-port", 2222, "SSH server port")
	webPort := flag.Int("web-port", 8080, "Web server port")
	host := flag.String("host", "0.0.0.0", "Host to bind to")
	flag.Parse()

	log.Println("Starting CapC0 Terminal...")

	// Create .ssh directory for host keys if it doesn't exist
	if err := os.MkdirAll(".ssh", 0700); err != nil {
		log.Fatalf("Failed to create .ssh directory: %v", err)
	}

	// Start SSH server
	sshServer := server.NewSSHServer(*host, *sshPort)
	if err := sshServer.StartAsync(); err != nil {
		log.Fatalf("Failed to start SSH server: %v", err)
	}

	// Start Web server
	webServer := server.NewWebServer(*host, *webPort)
	if err := webServer.StartAsync(); err != nil {
		log.Fatalf("Failed to start web server: %v", err)
	}

	log.Printf("SSH server running on %s:%d", *host, *sshPort)
	log.Printf("Web server running on http://%s:%d", *host, *webPort)
	log.Println("Press Ctrl+C to exit")

	// Wait for interrupt signal
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Println("Shutting down...")
}
