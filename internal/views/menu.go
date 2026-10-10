package views

import (
	"fmt"
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
)

type MenuView struct {
	Game    types.GameState
	PlayBtn *utils.Button
	Message string
}

func (m *MenuView) Draw() []string {
	w := m.Game.Terminal.Size.Width
	h := m.Game.Terminal.Size.Height

	lines := utils.CreateBlankScreen(w, h)

	utils.DrawCenterText(lines, h/6, w, "GopherFish")

	utils.DrawCenterText(lines, h-2, w, fmt.Sprintf("Terminal: (W: %d, H: %d)", w, h))
	utils.DrawCenterText(lines, h-1, w, fmt.Sprintf("Souris: (X: %d, Y: %d) Clic: %v", m.Game.Terminal.CursorPos.X, m.Game.Terminal.CursorPos.Y, m.Game.Terminal.CursorPos.Clicked))

	if m.PlayBtn == nil {
		text := "JOUER"
		btnX := (w - (len(text) + 4)) / 2
		m.PlayBtn = utils.NewButton(text, btnX, h/2)
	}

	if m.PlayBtn.UpdateAndDraw(lines, m.Game.Terminal.CursorPos) {
		m.Message = "Oui ?"
	}

	if m.Message != "" {
		utils.DrawCenterText(lines, h/2+2, w, m.Message)
	}

	return lines
}
