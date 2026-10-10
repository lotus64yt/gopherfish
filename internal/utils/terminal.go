package utils

import (
	"fmt"
	"gopherfish/internal/types"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

func TerminalGetSize() (types.TerminalSize, error) {
	w, h, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		w, h, err = term.GetSize(os.Stdin.Fd())
		if err != nil {
			w, h, err = term.GetSize(os.Stderr.Fd())
		}
	}

	if err != nil || w == 0 || h == 0 {
		return types.TerminalSize{Width: 80, Height: 24}, err
	}

	return types.TerminalSize{
		Width:  w,
		Height: h,
	}, nil
}

func TerminalGetCursorPosition() (types.TerminalCursorPos, error) {
	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		return types.TerminalCursorPos{}, err
	}
	defer term.Restore(os.Stdin.Fd(), oldState)

	_, err = fmt.Fprint(os.Stdout, "\033[6n")
	if err != nil {
		return types.TerminalCursorPos{}, err
	}

	buf := make([]byte, 32)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return types.TerminalCursorPos{}, err
	}

	response := string(buf[:n])

	if !strings.HasPrefix(response, "\033[") || !strings.HasSuffix(response, "R") {
		return types.TerminalCursorPos{}, fmt.Errorf("invalid terminal res %q", response)
	}

	trimmed := strings.TrimSuffix(strings.TrimPrefix(response, "\033["), "R")
	parts := strings.Split(trimmed, ";")
	if len(parts) != 2 {
		return types.TerminalCursorPos{}, fmt.Errorf("pos33s format incorect %q", trimmed)
	}

	row, err1 := strconv.Atoi(parts[0])
	col, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return types.TerminalCursorPos{}, fmt.Errorf("cannot parse coords")
	}

	return types.TerminalCursorPos{
		X: col - 1,
		Y: row - 1,
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
		}

		currentLine = strings.TrimRight(currentLine, " ")

		sb.WriteString(currentLine)

		if y < tSize.Height-1 {
			sb.WriteString("\r\n")
		}
	}

	return sb.String()
}
