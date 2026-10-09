package engine

import (
	"fmt"
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
)

func RenderFrame(g types.GameState, lines []string) {
	out := utils.TerminalFitToScreen(lines, g.Terminal.Size)
	fmt.Print("\033[2J\033[H")
	fmt.Print(out)
}
