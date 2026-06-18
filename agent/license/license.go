package license

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/usman8786/Siyaho-POS-Agent/agent/config"
)

const cacheTTL = 24 * time.Hour

type State struct {
	Valid           bool
	CompanyName     string
	AllowedOrigins  []string
	ExpiresAt       string
	LastValidatedAt time.Time
	LastError       string
}

type Manager struct {
	mu     sync.RWMutex
	cfg    config.Config
	state  State
	client *http.Client
}

func NewManager(cfg config.Config) *Manager {
	return &Manager{
		cfg: cfg,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (m *Manager) UpdateConfig(cfg config.Config) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cfg = cfg
}

func (m *Manager) Snapshot() State {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.state
}

func (m *Manager) Refresh() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.refreshLocked()
}

func (m *Manager) refreshLocked() {
	key := strings.TrimSpace(m.cfg.LicenseKey)
	if key == "" {
		m.state = State{
			Valid:           true,
			CompanyName:     "Siyaho",
			AllowedOrigins:  nil,
			LastValidatedAt: time.Now(),
			LastError:       "",
		}
		return
	}

	if !m.state.LastValidatedAt.IsZero() && time.Since(m.state.LastValidatedAt) < cacheTTL && m.state.Valid {
		return
	}

	apiBase := strings.TrimRight(strings.TrimSpace(m.cfg.APIBaseURL), "/")
	endpoint := fmt.Sprintf("%s/api/print-agent/validate?key=%s", apiBase, url.QueryEscape(key))
	resp, err := m.client.Get(endpoint)
	if err != nil {
		m.state.LastError = err.Error()
		m.state.LastValidatedAt = time.Now()
		m.state.Valid = false
		return
	}
	defer resp.Body.Close()

	var body struct {
		OK              bool     `json:"ok"`
		AllowedOrigins  []string `json:"allowedOrigins"`
		CompanyName     string   `json:"companyName"`
		ExpiresAt       string   `json:"expiresAt"`
		Error           string   `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		m.state.LastError = err.Error()
		m.state.Valid = false
		m.state.LastValidatedAt = time.Now()
		return
	}

	m.state = State{
		Valid:           body.OK && resp.StatusCode == http.StatusOK,
		CompanyName:     body.CompanyName,
		AllowedOrigins:  body.AllowedOrigins,
		ExpiresAt:       body.ExpiresAt,
		LastValidatedAt: time.Now(),
		LastError:       body.Error,
	}
}

func (m *Manager) IsOriginAllowed(origin string) bool {
	origin = strings.TrimSpace(origin)
	if origin == "" || origin == "null" {
		return true
	}

	if isSiyahoOrigin(origin) {
		return true
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, allowed := range builtinDevOrigins() {
		if strings.EqualFold(origin, allowed) {
			return true
		}
	}

	if strings.TrimSpace(m.cfg.LicenseKey) == "" {
		return isSiyahoOrigin(origin)
	}

	if !m.state.Valid {
		return false
	}

	for _, o := range m.state.AllowedOrigins {
		if strings.EqualFold(strings.TrimSpace(o), origin) {
			return true
		}
	}
	return false
}

func isSiyahoOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	host := strings.ToLower(u.Hostname())
	return host == "siyaho.com" || strings.HasSuffix(host, ".siyaho.com")
}

func builtinDevOrigins() []string {
	return []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:5000",
		"http://127.0.0.1:5000",
	}
}

func StartAutoRefresh(m *Manager, stop <-chan struct{}) {
	m.Refresh()
	ticker := time.NewTicker(cacheTTL)
	go func() {
		for {
			select {
			case <-ticker.C:
				m.Refresh()
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}
