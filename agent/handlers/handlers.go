package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/usman8786/Siyaho-POS-Agent/agent/config"
	"github.com/usman8786/Siyaho-POS-Agent/agent/license"
	"github.com/usman8786/Siyaho-POS-Agent/agent/meta"
	"github.com/usman8786/Siyaho-POS-Agent/agent/print"
)

type Server struct {
	Port     int
	Config   config.Config
	Licenses *license.Manager
}

type legacyPrintBody struct {
	Printer struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"printer"`
	Payload string `json:"payload"`
}

type v1PrintBody struct {
	Target struct {
		Type string `json:"type"`
		Host string `json:"host"`
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"target"`
	Printer struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"printer"`
	Data struct {
		Format   string `json:"format"`
		Encoding string `json:"encoding"`
		Content  string `json:"content"`
	} `json:"data"`
	Payload string `json:"payload"`
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/v1/health", s.handleHealth)
	mux.HandleFunc("/print", s.handlePrint)
	mux.HandleFunc("/v1/print", s.handlePrint)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "error": "Method not allowed"})
		return
	}
	lic := s.Licenses.Snapshot()
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"service": meta.ServiceName,
		"version": meta.Version,
		"port":    s.Port,
		"license": map[string]any{
			"valid":       lic.Valid,
			"companyName": lic.CompanyName,
			"expiresAt":   lic.ExpiresAt,
		},
	})
}

func (s *Server) handlePrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "error": "Method not allowed"})
		return
	}

	raw, err := io.ReadAll(io.LimitReader(r.Body, 4<<20))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "Failed to read body"})
		return
	}

	var body v1PrintBody
	if err := json.Unmarshal(raw, &body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "Invalid JSON"})
		return
	}

	target, payload, format, err := parsePrintRequest(body, raw)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": err.Error()})
		return
	}

	if target.Host == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "printer.ip is required"})
		return
	}
	if format == "dantsu" && payload == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "payload is required"})
		return
	}

	var printErr error
	switch format {
	case "escpos":
		printErr = print.PrintBase64(print.NetworkTarget{Host: target.Host, Port: target.Port}, payload, print.DefaultTimeout)
	case "dantsu":
		printErr = print.PrintDantsu(print.NetworkTarget{Host: target.Host, Port: target.Port}, payload, print.DefaultTimeout)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": "unsupported data format"})
		return
	}
	if printErr != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": printErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func parsePrintRequest(body v1PrintBody, raw []byte) (print.NetworkTarget, string, string, error) {
	host := strings.TrimSpace(body.Printer.IP)
	port := body.Printer.Port
	if host == "" {
		host = strings.TrimSpace(body.Target.Host)
	}
	if host == "" {
		host = strings.TrimSpace(body.Target.IP)
	}
	if port <= 0 {
		port = body.Target.Port
	}

	// Legacy DantSu body shape.
	if body.Payload != "" || (body.Data.Content == "" && host != "") {
		var legacy legacyPrintBody
		_ = json.Unmarshal(raw, &legacy)
		if strings.TrimSpace(legacy.Printer.IP) != "" {
			host = strings.TrimSpace(legacy.Printer.IP)
		}
		if legacy.Printer.Port > 0 {
			port = legacy.Printer.Port
		}
		payload := body.Payload
		if payload == "" {
			payload = legacy.Payload
		}
		return print.NetworkTarget{Host: host, Port: port}, payload, "dantsu", nil
	}

	format := strings.ToLower(strings.TrimSpace(body.Data.Format))
	if format == "" {
		format = "dantsu"
	}
	return print.NetworkTarget{Host: host, Port: port}, body.Data.Content, format, nil
}

func writeJSON(w http.ResponseWriter, status int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	_ = enc.Encode(payload)
}

func ListenAddress(port int) string {
	if port <= 0 {
		port = config.DefaultPort
	}
	return "127.0.0.1:" + strconv.Itoa(port)
}

func LogStartup(port int) {
	fmt.Printf("Siyaho Printer Agent listening on http://%s\n", ListenAddress(port))
	fmt.Printf("Version %s\n", meta.Version)
	fmt.Printf("Started at %s\n", time.Now().Format(time.RFC3339))
}
