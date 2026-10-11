package views

import (
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
)

type MenuView struct {
	Game    *types.GameState
	PlayBtn *utils.Button
	Message string
}

func (m *MenuView) Draw() []string {
	w := m.Game.Terminal.Size.Width
	h := m.Game.Terminal.Size.Height

	lines := utils.CreateBlankScreen(w, h)

	utils.DrawCenterText(lines, h/6, w, "GopherFish")

	if m.PlayBtn == nil {
		text := "PLAY"
		btnX := (w - (len(text) + 4)) / 2
		m.PlayBtn = utils.NewButton(text, btnX, h/2)
	}

	if m.PlayBtn.UpdateAndDraw(lines, m.Game.Terminal.CursorPos) {
		m.Game.CurrentView = types.GameViewChasing
	}

	return lines
}
