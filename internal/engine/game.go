package engine

import (
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
	"gopherfish/internal/views"
	"log"
)

var GameState = types.GameState{
	CurrentView: types.GameViewMenu,
	Terminal:    types.TerminalInfo{},
}

func refreshGameState(g types.GameState) types.GameState {
	tSize, err := utils.TerminalGetSize()
	if err != nil {
		log.Println(err)
		return g
	}

	g.Terminal.Size = tSize

	return g
}

func RunGame() {
	for {
		var lines []string
		GameState = refreshGameState(GameState)
		switch GameState.CurrentView {
		case types.GameViewMenu:
			lines = (&views.MenuView{Game: GameState}).Draw()
		}

		RenderFrame(GameState, lines)
	}
}
