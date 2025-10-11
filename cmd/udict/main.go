package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"regexp"
	"strings"

	"github.com/93aakash/urban_dict/internal/udict"
	"github.com/urfave/cli/v3"
	"golang.org/x/term"
)

const (
	defaultWidth   = 79
	termWidthRatio = 0.5
)

func main() {
	cmd := &cli.Command{
		Name: "udict",
		Flags: []cli.Flag{
			&cli.BoolFlag{Name: "delete", Aliases: []string{"d"}},
			&cli.BoolFlag{Name: "debug", Aliases: []string{"x"}},
		},
		Usage:     "Look up terms using Urban Dictionary",
		ArgsUsage: "term",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() == 0 {
				cli.ShowAppHelpAndExit(cmd, 1)
			}

			if cmd.Bool("debug") {
				setLogLevel(slog.LevelDebug)
			}

			d, err := udict.New()
			if err != nil {
				return cli.Exit(fmt.Sprintf("failed to initialize dictionary: %v", err), 1)
			}

			query := strings.ToLower(strings.Join(cmd.Args().Slice(), " "))

			// Handle delete flag.
			if cmd.Bool("delete") {
				if err := d.Delete(query); err != nil {
					return cli.Exit(fmt.Sprintf("failed to delete definition: %v", err), 1)
				}
			}

			// Get definition.
			definitions, err := d.Get(query)
			if err != nil || len(definitions) == 0 {
				errMessage := "definition not found"
				if err != udict.ErrDefNotFound {
					errMessage = fmt.Sprintf("failed to get definition: %v", err)
				}
				return cli.Exit(errMessage, 1)
			}

			// Format and print definitions.
			fmt.Printf("%d definitions found\n", len(definitions))
			for i, d := range definitions {
				fmt.Println()
				printDefinition(i+1, d.Definition)
			}

			return nil
		},
		Version: "1.0.0",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func getTermWidth() int {
	if !term.IsTerminal(int(os.Stdout.Fd())) {
		return defaultWidth
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		slog.Debug("Error getting terminal size", "error", err)
		return defaultWidth
	}

	return int(math.Round(float64(width) * termWidthRatio))
}

func wrapText(text string, width int) string {
	var result strings.Builder
	words := strings.Fields(text) // Split into words

	currentLineLen := 0
	for _, word := range words {
		if currentLineLen+len(word)+1 > width { // +1 for space
			result.WriteString("\n     ") // New line
			currentLineLen = 0
		}
		if currentLineLen > 0 {
			result.WriteString(" ") // Add space between words
			currentLineLen++
		}
		result.WriteString(word)
		currentLineLen += len(word)
	}
	return result.String()
}

func printDefinition(index int, definition string) {
	for i, line := range splitOnNewlines(definition) {
		if i == 0 {
			fmt.Printf("%-5s", fmt.Sprintf("%d.", index))
		} else {
			fmt.Printf("%-5s", "")
		}
		// Add padding for subsequent lines
		fmt.Printf("%-5s\n", wrapText(line, getTermWidth()))
	}
}

func splitOnNewlines(text string) []string {
	// Regular expression to match 1 or more consecutive \r\n or \n sequences
	re := regexp.MustCompile(`(\r\n|\n)+`)
	// Split the string based on these sequences of newlines
	return re.Split(text, -1)
}

func setLogLevel(level slog.Level) {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
