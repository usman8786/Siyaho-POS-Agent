package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runtime"
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
		Type string `json:"type"`
		Name string `json:"name"`
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"printer"`
	Payload string `json:"payload"`
}

type v1PrintBody struct {
	Target struct {
		Type string `json:"type"`
		Name string `json:"name"`
		Host string `json:"host"`
		IP   string `json:"ip"`
		Port int    `json:"port"`
	} `json:"target"`
	Printer struct {
		Type string `json:"type"`
		Name string `json:"name"`
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
	mux.HandleFunc("/printers", s.handlePrinters)
	mux.HandleFunc("/v1/printers", s.handlePrinters)
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
		"ok":       true,
		"service":  meta.ServiceName,
		"version":  meta.Version,
		"port":     s.Port,
		"platform": runtime.GOOS,
		"license": map[string]any{
			"valid":       lic.Valid,
			"companyName": lic.CompanyName,
			"expiresAt":   lic.ExpiresAt,
		},
	})
}

func (s *Server) handlePrinters(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"ok": false, "error": "Method not allowed"})
		return
	}
	if runtime.GOOS != "windows" {
		writeJSON(w, http.StatusNotImplemented, map[string]any{
			"ok":    false,
			"error": "windows printer listing is only supported on Windows",
		})
		return
	}

	printers, err := print.ListWindowsPrinters()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"ok": false, "error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"ok":       true,
		"printers": printers,
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

	if err := validatePrintTarget(target, format, payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"ok": false, "error": err.Error()})
		return
	}

	var printErr error
	switch format {
	case "escpos":
		printErr = print.PrintBase64ToTarget(target, payload, print.DefaultTimeout)
	case "dantsu":
		printErr = print.PrintDantsuToTarget(target, payload, print.DefaultTimeout)
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

func validatePrintTarget(target print.Target, format, payload string) error {
	if target.IsWindows() {
		if strings.TrimSpace(target.Name) == "" {
			return errf("printer.name is required for windows printers")
		}
	} else if strings.TrimSpace(target.Host) == "" {
		return errf("printer.ip is required for network printers")
	}
	if format == "dantsu" && payload == "" {
		return errf("payload is required")
	}
	return nil
}

func errf(msg string) error {
	return &parseError{msg: msg}
}

type parseError struct{ msg string }

func (e *parseError) Error() string { return e.msg }

func parsePrintRequest(body v1PrintBody, raw []byte) (print.Target, string, string, error) {
	target := resolveTarget(body.Target, body.Printer)

	var legacy legacyPrintBody
	_ = json.Unmarshal(raw, &legacy)
	if legacy.Printer.Type != "" || legacy.Printer.Name != "" || legacy.Printer.IP != "" {
		target = resolveTarget(struct {
			Type string `json:"type"`
			Name string `json:"name"`
			Host string `json:"host"`
			IP   string `json:"ip"`
			Port int    `json:"port"`
		}{
			Type: legacy.Printer.Type,
			Name: legacy.Printer.Name,
			Host: legacy.Printer.IP,
			IP:   legacy.Printer.IP,
			Port: legacy.Printer.Port,
		}, body.Printer)
	}

	if body.Payload != "" || (body.Data.Content == "" && (target.Host != "" || target.Name != "")) {
		payload := body.Payload
		if payload == "" {
			payload = legacy.Payload
		}
		return target, payload, "dantsu", nil
	}

	format := strings.ToLower(strings.TrimSpace(body.Data.Format))
	if format == "" {
		format = "dantsu"
	}
	return target, body.Data.Content, format, nil
}

func resolveTarget(targetFields struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Host string `json:"host"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}, printerFields struct {
	Type string `json:"type"`
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}) print.Target {
	printerType := strings.ToLower(strings.TrimSpace(targetFields.Type))
	name := strings.TrimSpace(targetFields.Name)
	host := strings.TrimSpace(targetFields.Host)
	if host == "" {
		host = strings.TrimSpace(targetFields.IP)
	}
	port := targetFields.Port

	if printerType == "" {
		printerType = strings.ToLower(strings.TrimSpace(printerFields.Type))
	}
	if name == "" {
		name = strings.TrimSpace(printerFields.Name)
	}
	if host == "" {
		host = strings.TrimSpace(printerFields.IP)
	}
	if port <= 0 {
		port = printerFields.Port
	}

	if printerType == "windows" {
		return print.Target{Type: print.TypeWindows, Name: name}
	}
	if host != "" {
		return print.Target{Type: print.TypeNetwork, Host: host, Port: port}
	}
	if name != "" {
		return print.Target{Type: print.TypeWindows, Name: name}
	}

	return print.Target{Type: print.TypeNetwork, Host: host, Port: port}
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
	log.Printf("Siyaho Printer Agent listening on http://%s", ListenAddress(port))
	log.Printf("Version %s", meta.Version)
	log.Printf("Started at %s", time.Now().Format(time.RFC3339))
}
