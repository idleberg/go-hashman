package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/idleberg/go-hashman/internal/hasher"
)

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
			Logger.Error("error computing %s: %v", r.Algorithm.Display, r.Err)
			continue
		}
		name := fmt.Sprintf("%-*s", maxDisplayLen, r.Algorithm.Display)
		dur := fmt.Sprintf("%.2fms", float64(r.Duration.Microseconds())/1000.0)
		Logger.Info("%s %s %s", name, blueStyle.Render(r.Hash), greyStyle.Render(dur))
	}

	fmt.Println()
	dur := fmt.Sprintf("%.2fms", float64(totalDuration.Microseconds())/1000.0)
	Logger.Success("Completed in %s", dur)
}

func PrintSeparator() {
	fmt.Println()
	fmt.Println("───")
}

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
