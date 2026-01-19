package export

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"portscango/internal/scanner"
)

// WriteCSV exports results to CSV file
func WriteCSV(filename string, target string, results []scanner.Result) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Target", "Port", "State", "Service", "Banner", "Timestamp"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	for _, r := range results {
		row := []string{
			target,
			fmt.Sprintf("%d", r.Port),
			r.State,
			r.Service,
			r.Banner,
			timestamp,
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}
