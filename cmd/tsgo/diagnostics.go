package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// tableRow represents a single row in the diagnostic table
type tableRow struct {
	name  string
	value string
}

// table collects and formats diagnostic data
type table struct {
	rows []tableRow
}

// add adds a statistic to the table, automatically formatting durations
func (t *table) add(name string, value any) {
	if d, ok := value.(time.Duration); ok {
		value = formatDuration(d)
	}
	t.rows = append(t.rows, tableRow{name, fmt.Sprint(value)})
}

// toJSON converts the table to a map suitable for JSON serialization
func (t *table) toJSON() map[string]string {
	result := make(map[string]string)
	for _, row := range t.rows {
		result[row.name] = row.value
	}
	return result
}

// print outputs the table in a formatted way to stdout
func (t *table) print() {
	nameWidth := 0
	valueWidth := 0
	for _, r := range t.rows {
		nameWidth = max(nameWidth, len(r.name))
		valueWidth = max(valueWidth, len(r.value))
	}

	for _, r := range t.rows {
		fmt.Printf("%-*s %*s\n", nameWidth+1, r.name+":", valueWidth, r.value)
	}
}

// formatDuration formats a duration in seconds with 3 decimal places
func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.3fs", d.Seconds())
}

// OutputStats handles the output of diagnostic statistics based on the extendedDiagnostics option
// Returns any error encountered during the output process
func OutputStats(stats *table, extendedDiagnostics string) error {
	if extendedDiagnostics == "" {
		return nil
	}

	if extendedDiagnostics == "inline" {
		stats.print()
		return nil
	}

	jsonData := stats.toJSON()

	if strings.HasSuffix(extendedDiagnostics, ".json") {
		// Ensure directory exists
		dir := filepath.Dir(extendedDiagnostics)
		if dir != "." {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("error creating directory for stats JSON file: %w", err)
			}
		}

		// Write to file
		file, err := os.Create(extendedDiagnostics)
		if err != nil {
			return fmt.Errorf("error creating stats JSON file: %w", err)
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ")
		if err := enc.Encode(jsonData); err != nil {
			return fmt.Errorf("error writing stats to JSON file: %w", err)
		}
		return nil
	}

	return errors.New("invalid value for --extendedDiagnostics. Use 'inline' or a path ending with '.json'")
}
