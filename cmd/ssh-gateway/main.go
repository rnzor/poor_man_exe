package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/gliderlabs/ssh"
	"github.com/rnzor/poor_man_exe/internal/auth"
	"github.com/rnzor/poor_man_exe/internal/caddy"
	"github.com/rnzor/poor_man_exe/internal/config"
	"github.com/rnzor/poor_man_exe/internal/db"
	"github.com/rnzor/poor_man_exe/internal/router"
	"github.com/rnzor/poor_man_exe/internal/runner"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	cfg := config.Load()

	// Connect to DB
	database, err := db.Connect(cfg.DBPath)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()

	// Init DB schema (embedded)
	if err := database.Init(); err != nil {
		log.Printf("Warning: Failed to init schema (may already exist): %v", err)
	}

	// Init Docker runner
	dockerRunner, err := runner.NewDockerRunner()
	if err != nil {
		log.Fatalf("Failed to init Docker runner: %v", err)
	}

	// Init Authenticator, Caddy, and Router
	authenticator := auth.NewAuthenticator(database)
	caddyClient := caddy.NewClient(cfg.CaddyURL)
	rtr := router.NewRouter(database, dockerRunner, cfg, caddyClient)

	// Start rate limiter cleanup goroutine
	go func() {
		for range time.Tick(5 * time.Minute) {
			authenticator.Limiter.Cleanup()
		}
	}()

	// Start health check server
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		log.Printf("Starting health check server on :8080...")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Printf("Health check server failed: %v", err)
		}
	}()

	// Setup SSH server
	server := &ssh.Server{
		Addr: fmt.Sprintf(":%d", cfg.SSHPort),
		Handler: func(sess ssh.Session) {
			defer authenticator.OnSessionClose(sess)
			rtr.HandleSession(sess)
		},
		PublicKeyHandler: authenticator.PublicKeyHandler,
		BannerHandler:    authenticator.BannerHandler,
		ConnCallback:     authenticator.ConnCallback,
		IdleTimeout:      time.Duration(cfg.IdleTimeout) * time.Second,
		MaxTimeout:       time.Duration(cfg.IdleTimeout*2) * time.Second,
	}

	// Load or generate host key
	if _, err := os.Stat(cfg.HostKeyPath); os.IsNotExist(err) {
		log.Printf("Generating host key...")
		if err := generateHostKey(cfg.HostKeyPath); err != nil {
			log.Fatalf("Failed to generate host key: %v", err)
		}
	}

	keyData, err := ioutil.ReadFile(cfg.HostKeyPath)
	if err != nil {
		log.Fatalf("Failed to read host key: %v", err)
	}

	key, err := gossh.ParsePrivateKey(keyData)
	if err != nil {
		log.Fatalf("Failed to parse host key: %v", err)
	}
	server.AddHostKey(key)

	// Start server
	log.Printf("Starting SSH gateway on port %d...", cfg.SSHPort)
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Printf("SSH server stopped: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Printf("Shutting down...")
}

func generateHostKey(path string) error {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	keyFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer keyFile.Close()

	pemBlock, err := gossh.MarshalPrivateKey(privateKey, "")
	if err != nil {
		return err
	}

	if err := pem.Encode(keyFile, pemBlock); err != nil {
		return err
	}

	return nil
}
