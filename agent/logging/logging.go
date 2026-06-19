package logging

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/usman8786/Siyaho-POS-Agent/agent/config"
)

func LogPath() string {
	return filepath.Join(config.ConfigDir(), "agent.log")
}

// Setup writes agent logs to ProgramData. When stderr is a TTY (e.g. go run), also mirrors there.
func Setup() error {
	if err := os.MkdirAll(config.ConfigDir(), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(LogPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	writers := []io.Writer{f}
	if fi, statErr := os.Stderr.Stat(); statErr == nil && (fi.Mode()&os.ModeCharDevice) != 0 {
		writers = append(writers, os.Stderr)
	}

	log.SetOutput(io.MultiWriter(writers...))
	log.SetFlags(log.Ldate | log.Ltime)
	return nil
}
