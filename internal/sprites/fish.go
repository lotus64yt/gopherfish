package sprites

import (
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
	"math"
	"math/rand"
)

type Fish struct {
	X         int
	Y         int
	DirX      int
	DirY      int
	DirYStep  int
	DirYTick  int
	AsciiFish string
}

func (f *Fish) GetAscii() string {
	if f.DirX == -1 {
		return "@><"
	}
	return "><@"
}

func SpawnRandomShip(g types.GameState, fishs []Fish) []Fish {
	dirX := -1
	if rand.Intn(2) == 1 {
		dirX = 1
	}

	w := g.Terminal.Size.Width
	if w <= 10 {
		w = 80
	}
	minX := 4
	maxX := w - 4
	if maxX <= minX {
		maxX = minX + 1
	}

	minY := 2
	maxY := g.Terminal.Size.Height - 4
	if maxY <= minY {
		maxY = minY + 1
	}

	cols := 6
	rows := 4
	if (maxX-minX)/cols < 6 {
		cols = (maxX - minX) / 6
		if cols < 1 {
			cols = 1
		}
	}
	if (maxY-minY)/rows < 2 {
		rows = (maxY - minY) / 2
		if rows < 1 {
			rows = 1
		}
	}

	cellW := (maxX - minX) / cols
	cellH := (maxY - minY) / rows
	if cellW < 1 {
		cellW = 1
	}
	if cellH < 1 {
		cellH = 1
	}

	counts := make([][]int, cols)
	for c := 0; c < cols; c++ {
		counts[c] = make([]int, rows)
	}

	for _, f := range fishs {
		c := (f.X - minX) / cellW
		r := (f.Y - minY) / cellH
		if c >= 0 && c < cols && r >= 0 && r < rows {
			counts[c][r]++
		}
	}

	minCount := 999999
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			if counts[c][r] < minCount {
				minCount = counts[c][r]
			}
		}
	}

	type cellCoord struct{ c, r int }
	var bestCells []cellCoord
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			if counts[c][r] == minCount {
				bestCells = append(bestCells, cellCoord{c, r})
			}
		}
	}

	chosen := bestCells[rand.Intn(len(bestCells))]
	cLeft := minX + chosen.c*cellW
	rTop := minY + chosen.r*cellH

	bestX := cLeft + rand.Intn(cellW)
	bestY := rTop + rand.Intn(cellH)
	minDistX := 6
	minDistY := 2

	for attempt := 0; attempt < 20; attempt++ {
		candX := cLeft + rand.Intn(cellW)
		candY := rTop + rand.Intn(cellH)
		tooClose := false
		for _, other := range fishs {
			dx := int(math.Abs(float64(candX - other.X)))
			dy := int(math.Abs(float64(candY - other.Y)))
			if dx < minDistX && dy < minDistY {
				tooClose = true
				break
			}
		}
		if !tooClose {
			bestX = candX
			bestY = candY
			break
		}
		if attempt == 10 {
			minDistX = 4
			minDistY = 1
		}
	}

	newFish := Fish{
		DirX:     dirX,
		X:        bestX,
		Y:        bestY,
		DirY:     rand.Intn(3) - 1,
		DirYStep: rand.Intn(2)*2 - 1,
		DirYTick: rand.Intn(5),
	}
	newFish.AsciiFish = newFish.GetAscii()

	return append(fishs, newFish)
}

func (f *Fish) Tick(g types.GameState) {
	f.X += f.DirX

	f.DirYTick++
	if f.DirYTick%5 == 0 {
		if f.DirYStep == 0 {
			f.DirYStep = 1
		}
		f.DirY += f.DirYStep
		if f.DirY >= 1 {
			f.DirY = 1
			f.DirYStep = -1
		} else if f.DirY <= -1 {
			f.DirY = -1
			f.DirYStep = 1
		}
		f.Y += f.DirY
		if f.Y < 1 {
			f.Y = 1
			f.DirY = 1
		} else if g.Terminal.Size.Height > 4 && f.Y >= g.Terminal.Size.Height-4 {
			f.Y = g.Terminal.Size.Height - 4
			f.DirY = -1
		}
	}

	if utils.TerminalHitScreenBorder(g.Terminal.Size, types.TerminalCursorPos{
		X: f.X,
		Y: f.Y,
	}) {
		f.DirX = -f.DirX
	}
}

func (f *Fish) Draw(screen []string) []string {
	f.AsciiFish = f.GetAscii()

	if f.Y < 0 || f.Y >= len(screen) {
		return screen
	}

	line := []rune(screen[f.Y])
	fish := []rune(f.AsciiFish)
	start := f.X - len(fish)/2

	for i, r := range fish {
		pos := start + i
		if pos >= 0 && pos < len(line) {
			line[pos] = r
		}
	}

	screen[f.Y] = string(line)

	return screen
}

func (f *Fish) IsClicked(cPos types.TerminalCursorPos) bool {
	if cPos.Y == f.Y {
		fishStr := []rune(f.GetAscii())
		start := f.X - len(fishStr)/2
		if cPos.X >= start-1 && cPos.X <= start+len(fishStr) {
			return true
		}
	}

	return false
}
