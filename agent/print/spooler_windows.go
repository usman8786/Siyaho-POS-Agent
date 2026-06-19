//go:build windows

package print

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type WindowsPrinter struct {
	Name       string `json:"name"`
	PortName   string `json:"portName"`
	DriverName string `json:"driverName"`
	IsDefault  bool   `json:"isDefault"`
}

type docInfo1 struct {
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
}

var (
	modWinspool          = windows.NewLazySystemDLL("winspool.drv")
	procOpenPrinterW     = modWinspool.NewProc("OpenPrinterW")
	procClosePrinter     = modWinspool.NewProc("ClosePrinter")
	procStartDocPrinterW = modWinspool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = modWinspool.NewProc("EndDocPrinter")
	procStartPagePrinter = modWinspool.NewProc("StartPagePrinter")
	procEndPagePrinter   = modWinspool.NewProc("EndPagePrinter")
	procWritePrinter     = modWinspool.NewProc("WritePrinter")
)

func sendWindowsRaw(printerName string, data []byte) error {
	if strings.TrimSpace(printerName) == "" {
		return fmt.Errorf("windows printer name is required")
	}
	if len(data) == 0 {
		return fmt.Errorf("empty print data")
	}

	namePtr, err := windows.UTF16PtrFromString(printerName)
	if err != nil {
		return err
	}

	var handle windows.Handle
	ret, _, callErr := procOpenPrinterW.Call(
		uintptr(unsafe.Pointer(namePtr)),
		uintptr(unsafe.Pointer(&handle)),
		0,
	)
	if ret == 0 {
		return fmt.Errorf("open printer %q: %w", printerName, callErr)
	}
	defer procClosePrinter.Call(uintptr(handle))

	docName, _ := windows.UTF16PtrFromString("Siyaho POS")
	dataType, _ := windows.UTF16PtrFromString("RAW")
	di := docInfo1{
		DocName:  docName,
		Datatype: dataType,
	}

	docID, _, callErr := procStartDocPrinterW.Call(
		uintptr(handle),
		1,
		uintptr(unsafe.Pointer(&di)),
	)
	if docID == 0 {
		return fmt.Errorf("start document: %w", callErr)
	}
	defer procEndDocPrinter.Call(uintptr(handle))

	ret, _, callErr = procStartPagePrinter.Call(uintptr(handle))
	if ret == 0 {
		return fmt.Errorf("start page: %w", callErr)
	}

	var written uint32
	ret, _, callErr = procWritePrinter.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)),
		uintptr(unsafe.Pointer(&written)),
	)
	if ret == 0 {
		return fmt.Errorf("write printer: %w", callErr)
	}

	ret, _, callErr = procEndPagePrinter.Call(uintptr(handle))
	if ret == 0 {
		return fmt.Errorf("end page: %w", callErr)
	}

	return nil
}

func ListWindowsPrinters() ([]WindowsPrinter, error) {
	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-NonInteractive",
		"-Command",
		`$default = (Get-Printer | Where-Object Default -eq $true | Select-Object -First 1 -ExpandProperty Name); Get-Printer | Select-Object @{n='name';e={$_.Name}}, @{n='portName';e={$_.PortName}}, @{n='driverName';e={$_.DriverName}}, @{n='isDefault';e={$_.Name -eq $default}} | ConvertTo-Json -Compress`,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("list printers: %w", err)
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "null" {
		return []WindowsPrinter{}, nil
	}

	var printers []WindowsPrinter
	if raw[0] == '[' {
		if err := json.Unmarshal([]byte(raw), &printers); err != nil {
			return nil, err
		}
	} else {
		var one WindowsPrinter
		if err := json.Unmarshal([]byte(raw), &one); err != nil {
			return nil, err
		}
		printers = []WindowsPrinter{one}
	}

	return printers, nil
}
