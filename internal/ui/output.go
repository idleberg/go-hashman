package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"

	"github.com/idleberg/hashman/internal/hasher"
)

var (
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	greyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)

var Logger *log.Logger

func init() {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: false,
	})

	styles := log.DefaultStyles()
	styles.Levels[log.InfoLevel] = lipgloss.NewStyle().
		SetString("ℹ").
		Foreground(lipgloss.Color("39"))
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("✖").
		Foreground(lipgloss.Color("196"))
	styles.Levels[log.DebugLevel] = lipgloss.NewStyle().
		SetString("DEBUG").
		Bold(true).
		Background(lipgloss.Color("6")).
		Foreground(lipgloss.Color("0"))
	styles.Message = lipgloss.NewStyle()
	styles.Key = lipgloss.NewStyle()
	styles.Value = lipgloss.NewStyle()

	Logger.SetStyles(styles)
}

// PrintResults renders the hash results for a single file.
func PrintResults(filePath string, results []hasher.Result, totalDuration time.Duration, maxDisplayLen int) {
	fmt.Println()
	fmt.Println(yellowStyle.Render(filepath.Base(filePath)))

	info, err := os.Stat(filePath)
	if err == nil {
		size := info.Size()
		fmt.Printf("%s (%d bytes)\n", formatBytes(size), size)
	}

	fmt.Println()

	for _, r := range results {
		if r.Err != nil {
			Logger.Error(fmt.Sprintf("error computing %s: %v", r.Algorithm.Display, r.Err))
			continue
		}
		name := fmt.Sprintf("%-*s", maxDisplayLen, r.Algorithm.Display)
		dur := fmt.Sprintf("%.2fms", float64(r.Duration.Microseconds())/1000.0)
		Logger.Infof("%s %s %s", name, blueStyle.Render(r.Hash), greyStyle.Render(dur))
	}

	fmt.Println()
	dur := fmt.Sprintf("%.2fms", float64(totalDuration.Microseconds())/1000.0)
	fmt.Printf("%s Completed in %s\n", greenStyle.Render("✔"), dur)
}

// PrintSeparator renders the separator between multiple file outputs.
func PrintSeparator() {
	fmt.Println()
	fmt.Println("───")
}

// formatBytes converts bytes to a human-readable string using SI units (1000-based),
// matching the behavior of the pretty-bytes npm package.
func formatBytes(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	suffixes := []string{"kB", "MB", "GB", "TB", "PB", "EB"}
	val := float64(b) / float64(div)
	if val >= 100 {
		return fmt.Sprintf("%.0f %s", val, suffixes[exp])
	}
	if val >= 10 {
		return fmt.Sprintf("%.1f %s", val, suffixes[exp])
	}
	return fmt.Sprintf("%.2f %s", val, suffixes[exp])
}
