package ui

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	infoSymbol    string
	successSymbol string
	warningSymbol string
	errorSymbol   string
)

var (
	blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Blue)
	greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Green)
	yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Yellow)
	redStyle    = lipgloss.NewStyle().Foreground(lipgloss.Red)
	cyanBgStyle = lipgloss.NewStyle().Background(lipgloss.Cyan).Foreground(lipgloss.Black)
	greyStyle   = lipgloss.NewStyle().Foreground(lipgloss.BrightBlack)
)

var Logger = &logger{}

type logger struct{}

func init() {
	if isUnicodeSupported() {
		infoSymbol = "ℹ"
		successSymbol = "✔"
		warningSymbol = "⚠"
		errorSymbol = "✖"
	} else {
		infoSymbol = "i"
		successSymbol = "√"
		warningSymbol = "‼"
		errorSymbol = "×"
	}
}

func isUnicodeSupported() bool {
	if runtime.GOOS != "windows" {
		return true
	}
	if os.Getenv("CI") != "" ||
		os.Getenv("WT_SESSION") != "" ||
		os.Getenv("ConEmuTask") != "" ||
		os.Getenv("TERM_PROGRAM") == "vscode" ||
		strings.HasPrefix(os.Getenv("TERM"), "xterm") {
		return true
	}
	return false
}

func (l *logger) Log(args ...any) {
	fmt.Println(args...)
}

func (l *logger) Info(msg string, args ...any) {
	fmt.Println(blueStyle.Render(infoSymbol), fmt.Sprintf(msg, args...))
}

func (l *logger) Success(msg string, args ...any) {
	fmt.Println(greenStyle.Render(successSymbol), fmt.Sprintf(msg, args...))
}

func (l *logger) Warn(msg string, args ...any) {
	fmt.Fprintln(os.Stderr, yellowStyle.Render(warningSymbol), fmt.Sprintf(msg, args...))
}

func (l *logger) Error(msg string, args ...any) {
	fmt.Fprintln(os.Stderr, redStyle.Render(errorSymbol), fmt.Sprintf(msg, args...))
}

func (l *logger) Debug(msg string, args ...any) {
	if os.Getenv("DEBUG") == "" {
		return
	}
	fmt.Fprintln(os.Stderr, cyanBgStyle.Render(" DEBUG "), fmt.Sprintf(msg, args...))
}
