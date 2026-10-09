package views

import (
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
)

type MenuView struct {
	Game types.GameState
}

func (m *MenuView) Draw() []string {
	w := m.Game.Terminal.Size.Width
	h := m.Game.Terminal.Size.Height

	lines := utils.CreateBlankScreen(w, h)

	lines = utils.DrawCenterText(lines, h/6, w, "GopherFish")

	return lines
}
