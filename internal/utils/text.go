package utils

import (
	"strings"
)

func CreateBlankScreen(width, height int) []string {
	screen := make([]string, height)
	blankLine := strings.Repeat(" ", width)
	for i := 0; i < height; i++ {
		screen[i] = blankLine
	}
	return screen
}

func DrawText(screen []string, x, y int, text string) []string {
	if y < 0 || y >= len(screen) {
		return screen
	}

	line := []rune(screen[y])
	textRunes := []rune(text)

	for i := 0; i < len(textRunes); i++ {
		targetX := x + i
		if targetX >= 0 && targetX < len(line) {
			line[targetX] = textRunes[i]
		}
	}

	screen[y] = string(line)

	return screen
}

func DrawCenterText(screen []string, y int, width int, text string) []string {
	textRunes := []rune(text)
	x := (width - len(textRunes)) / 2
	if x < 0 {
		x = 0
	}
	return DrawText(screen, x, y, text)
}
