package print

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/usman8786/Siyaho-POS-Agent/agent/dantsu"
)

type TargetType string

const (
	TypeNetwork TargetType = "network"
	TypeWindows TargetType = "windows"
)

type Target struct {
	Type TargetType
	Host string
	Port int
	Name string
}

func (t Target) IsWindows() bool {
	return t.Type == TypeWindows
}

func PrintDantsuToTarget(target Target, payload string, timeout time.Duration) error {
	data := dantsu.ToEscPos(payload)
	return sendToTarget(target, data, timeout)
}

func PrintBase64ToTarget(target Target, b64 string, timeout time.Duration) error {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return fmt.Errorf("invalid base64 content: %w", err)
	}
	return sendToTarget(target, raw, timeout)
}

func sendToTarget(target Target, data []byte, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	switch target.Type {
	case TypeWindows:
		if target.Name == "" {
			return fmt.Errorf("windows printer name is required")
		}
		return sendWindowsRaw(target.Name, data)
	default:
		if target.Host == "" {
			return fmt.Errorf("printer host is required")
		}
		return sendTCP(NetworkTarget{Host: target.Host, Port: target.Port}, data, timeout)
	}
}
