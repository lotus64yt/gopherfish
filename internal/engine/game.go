package engine

import (
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
	"gopherfish/internal/views"
	"log"
	"time"
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
	cleanup, err := utils.StartInputLoop(func(x, y int, pressed bool) {
		GameState.Terminal.CursorPos.X = x
		GameState.Terminal.CursorPos.Y = y
		if pressed && !GameState.Terminal.CursorPos.Clicked {
			GameState.Terminal.CursorPos.JustClicked = true
		}
		GameState.Terminal.CursorPos.Clicked = pressed
	})
	if err != nil {
		log.Println("cursor err:", err)
	} else {
		defer cleanup()
	}

	menuView := &views.MenuView{}
	chasseView := &views.ChasseView{}

	for {
		var lines []string
		GameState = refreshGameState(GameState)
		switch GameState.CurrentView {
		case types.GameViewMenu:
			menuView.Game = &GameState
			lines = menuView.Draw()
		case types.GameViewChasing:
			chasseView.Game = &GameState
			lines = chasseView.Draw()
		}

		GameState.Terminal.CursorPos.JustClicked = false

		RenderFrame(GameState, lines)
		time.Sleep(32 * time.Millisecond)
	}
}
