//go:build !windows

package print

import "fmt"

type WindowsPrinter struct {
	Name       string `json:"name"`
	PortName   string `json:"portName"`
	DriverName string `json:"driverName"`
	IsDefault  bool   `json:"isDefault"`
}

func sendWindowsRaw(_ string, _ []byte) error {
	return fmt.Errorf("windows spooler printing is only supported on Windows")
}

func ListWindowsPrinters() ([]WindowsPrinter, error) {
	return nil, fmt.Errorf("windows printer listing is only supported on Windows")
}
