package print

import (
	"encoding/base64"
	"fmt"
	"net"
	"time"

	"github.com/usman8786/Siyaho-POS-Agent/agent/dantsu"
)

const DefaultTimeout = 30 * time.Second

type NetworkTarget struct {
	Host string
	Port int
}

func PrintDantsu(target NetworkTarget, payload string, timeout time.Duration) error {
	data := dantsu.ToEscPos(payload)
	return sendTCP(target, data, timeout)
}

func PrintRaw(target NetworkTarget, raw []byte, timeout time.Duration) error {
	return sendTCP(target, raw, timeout)
}

func PrintBase64(target NetworkTarget, b64 string, timeout time.Duration) error {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return fmt.Errorf("invalid base64 content: %w", err)
	}
	return sendTCP(target, raw, timeout)
}

func sendTCP(target NetworkTarget, data []byte, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	port := target.Port
	if port <= 0 {
		port = 9100
	}
	addr := net.JoinHostPort(target.Host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	_ = conn.SetDeadline(time.Now().Add(timeout))
	_, err = conn.Write(data)
	return err
}
