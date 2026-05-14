package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	"github.com/idleberg/go-hashman/internal/algo"
	"github.com/idleberg/go-hashman/internal/hasher"
	"github.com/idleberg/go-hashman/internal/ui"
)

var (
	allFlag   bool
	algoFlags = make(map[string]*bool)
)

func main() {
	rootCmd := &cobra.Command{
		Use:           "hashman [flags] <file> [file...]",
		Short:         "Calculate multiple hashes for files concurrently",
		Args:          cobra.MinimumNArgs(1),
		RunE:          run,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.Flags().BoolVarP(&allFlag, "all", "A", false, "use all supported hashes")

	for _, a := range algo.Registry {
		b := false
		algoFlags[a.ID] = &b
		desc := fmt.Sprintf("create %s hash", a.Display)
		if a.Deprecated {
			desc += " (deprecated)"
		}
		rootCmd.Flags().BoolVar(algoFlags[a.ID], a.Flag, false, desc)
	}

	if err := rootCmd.Execute(); err != nil {
		ui.Logger.Error("%s", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	selected := resolveAlgorithms()
	if len(selected) == 0 {
		ui.Logger.Error("no hashing algorithm provided")
		fmt.Println()
		return cmd.Help()
	}

	maxDisplayLen := 0
	for _, a := range selected {
		if len(a.Display) > maxDisplayLen {
			maxDisplayLen = len(a.Display)
		}
	}

	maxWorkers := runtime.NumCPU()
	isTTY := isTerminal()

	for i, filePath := range args {
		startTime := time.Now()

		var results []hasher.Result

		if isTTY {
			model := ui.NewSpinnerModel(filePath, selected, maxWorkers)
			p := tea.NewProgram(model)
			finalModel, err := p.Run()
			if err != nil {
				return fmt.Errorf("error processing %s: %w", filePath, err)
			}
			results = finalModel.(ui.SpinnerModel).Results()
		} else {
			results = hasher.HashFile(filePath, selected, maxWorkers)
		}

		totalDuration := time.Since(startTime)
		ui.PrintResults(filePath, results, totalDuration, maxDisplayLen)

		if i < len(args)-1 {
			ui.PrintSeparator()
		}
	}

	return nil
}

func resolveAlgorithms() []algo.Algorithm {
	var selected []algo.Algorithm
	for _, a := range algo.Registry {
		if allFlag {
			selected = append(selected, a)
		} else if b, ok := algoFlags[a.ID]; ok && *b {
			selected = append(selected, a)
		}
	}
	return selected
}

func isTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd())
}
