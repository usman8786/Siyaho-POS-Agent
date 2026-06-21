package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

const DefaultPort = 17890

type Config struct {
	Port int `json:"port"`
}

func Default() Config {
	return Config{
		Port: DefaultPort,
	}
}

func ConfigDir() string {
	if runtime.GOOS == "windows" {
		programData := os.Getenv("PROGRAMDATA")
		if programData == "" {
			programData = `C:\ProgramData`
		}
		return filepath.Join(programData, "Siyaho", "PrinterAgent")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".siyaho", "printer-agent")
}

func ConfigPath() string {
	return filepath.Join(ConfigDir(), "config.json")
}

func Load() (Config, error) {
	cfg := Default()
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}
	if cfg.Port <= 0 {
		cfg.Port = DefaultPort
	}
	return cfg, nil
}

func Save(cfg Config) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigPath(), data, 0o644)
}
