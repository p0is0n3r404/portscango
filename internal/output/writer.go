package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"portscango/internal/scanner"
)

// ScanOutput tarama çıktısını temsil eder
type ScanOutput struct {
	Target     string           `json:"target"`
	TotalPorts int              `json:"total_ports"`
	OpenPorts  int              `json:"open_ports"`
	ScanTime   string           `json:"scan_time"`
	Results    []scanner.Result `json:"results"`
}

// WriteJSON sonuçları JSON dosyasına yazar
func WriteJSON(filename string, output ScanOutput) error {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// WriteTXT sonuçları TXT dosyasına yazar
func WriteTXT(filename string, output ScanOutput) error {
	var sb strings.Builder

	sb.WriteString("PortScanGO Scan Report\n")
	sb.WriteString("======================\n\n")
	sb.WriteString(fmt.Sprintf("Target: %s\n", output.Target))
	sb.WriteString(fmt.Sprintf("Scanned Ports: %d\n", output.TotalPorts))
	sb.WriteString(fmt.Sprintf("Open Ports: %d\n", output.OpenPorts))
	sb.WriteString(fmt.Sprintf("Scan Time: %s\n\n", output.ScanTime))
	sb.WriteString(fmt.Sprintf("%-10s %-10s %-15s %s\n", "PORT", "STATE", "SERVICE", "BANNER"))
	sb.WriteString(strings.Repeat("-", 60) + "\n")

	for _, r := range output.Results {
		banner := r.Banner
		if len(banner) > 30 {
			banner = banner[:30] + "..."
		}
		sb.WriteString(fmt.Sprintf("%-10d %-10s %-15s %s\n", r.Port, r.State, r.Service, banner))
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

// PrintTable sonuçları tablo formatında konsola yazdırır
func PrintTable(results []scanner.Result, useColor bool) {
	if len(results) == 0 {
		return
	}

	// Header
	fmt.Printf("\n%-10s %-10s %-15s %s\n", "PORT", "STATE", "SERVICE", "BANNER")
	fmt.Println(strings.Repeat("─", 65))

	for _, r := range results {
		banner := r.Banner
		if len(banner) > 35 {
			banner = banner[:35] + "..."
		}

		if useColor {
			// Renkli çıktı (ANSI)
			fmt.Printf("\033[36m%-10d\033[0m \033[32m%-10s\033[0m \033[33m%-15s\033[0m %s\n",
				r.Port, r.State, r.Service, banner)
		} else {
			fmt.Printf("%-10d %-10s %-15s %s\n", r.Port, r.State, r.Service, banner)
		}
	}
	fmt.Println()
}
