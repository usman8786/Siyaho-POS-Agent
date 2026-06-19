package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/usman8786/Siyaho-POS-Agent/agent/config"
	"github.com/usman8786/Siyaho-POS-Agent/agent/cors"
	"github.com/usman8786/Siyaho-POS-Agent/agent/handlers"
	"github.com/usman8786/Siyaho-POS-Agent/agent/license"
	"github.com/usman8786/Siyaho-POS-Agent/agent/logging"
)

func main() {
	if err := logging.Setup(); err != nil {
		log.Fatalf("logging: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	if portEnv := os.Getenv("SIYAHO_PRINT_AGENT_PORT"); portEnv != "" {
		if p, err := parsePort(portEnv); err == nil {
			cfg.Port = p
		}
	}

	licenses := license.NewManager(cfg)
	stop := make(chan struct{})
	license.StartAutoRefresh(licenses, stop)

	srv := &handlers.Server{
		Port:     cfg.Port,
		Config:   cfg,
		Licenses: licenses,
	}

	mux := cors.New(licenses).Wrap(srv.Handler())
	addr := handlers.ListenAddress(cfg.Port)
	handlers.LogStartup(cfg.Port)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		close(stop)
		_ = httpServer.Close()
	}()

	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		if isAddrInUse(err) {
			log.Printf("agent already running on %s", addr)
			return
		}
		log.Fatalf("server: %v", err)
	}
}

func isAddrInUse(err error) bool {
	var opErr *net.OpError
	if !errors.As(err, &opErr) || opErr.Op != "listen" {
		return false
	}
	msg := strings.ToLower(opErr.Err.Error())
	return strings.Contains(msg, "address already in use") || strings.Contains(msg, "only one usage")
}

func parsePort(raw string) (int, error) {
	var p int
	_, err := fmt.Sscanf(raw, "%d", &p)
	return p, err
}
