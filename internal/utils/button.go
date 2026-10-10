package utils

import (
	"gopherfish/internal/types"
)

type Button struct {
	Text       string
	X          int
	Y          int
	Width      int
	WasClicked bool
}

func NewButton(text string, x, y int) *Button {
	return &Button{
		Text:       text,
		X:          x,
		Y:          y,
		Width:      len(text) + 4,
		WasClicked: false,
	}
}

func (b *Button) UpdateAndDraw(screen []string, cursor types.TerminalCursorPos) bool {
	isHovered := cursor.Y == b.Y && cursor.X >= b.X && cursor.X < b.X+b.Width

	isActionTriggered := false

	if isHovered {
		if b.WasClicked && !cursor.Clicked {
			isActionTriggered = true
		}
	}

	var displayText string
	if isHovered {
		if cursor.Clicked {
			displayText = "={" + b.Text + "}="
		} else {
			displayText = "> " + b.Text + " <"
		}
	} else {
		displayText = "[ " + b.Text + " ]"
	}

	DrawText(screen, b.X, b.Y, displayText)

	if isHovered {
		b.WasClicked = cursor.Clicked
	} else {
		b.WasClicked = false
	}

	return isActionTriggered
}
