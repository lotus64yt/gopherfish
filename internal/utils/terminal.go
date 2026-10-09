package utils

import (
	"gopherfish/internal/types"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

func TerminalGetSize() (types.TerminalSize, error) {
	sz, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return types.TerminalSize{}, err
	}

	return types.TerminalSize{
		Width:  int(sz.Col),
		Height: int(sz.Row),
	}, nil
}

func TerminalFitToScreen(lines []string, tSize types.TerminalSize) string {
	var sb strings.Builder

	for y := 0; y < tSize.Height; y++ {
		var currentLine string

		if y < len(lines) {
			currentLine = lines[y]
		} else {
			currentLine = ""
		}

		runes := []rune(currentLine)
		if len(runes) > tSize.Width {
			currentLine = string(runes[:tSize.Width])
		} else if len(runes) < tSize.Width {
			currentLine = currentLine + strings.Repeat(" ", tSize.Width-len(runes))
		}

		sb.WriteString(currentLine)

		if y < tSize.Height-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
