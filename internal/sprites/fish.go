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
	h := g.Terminal.Size.Height - 10
	if h <= 0 {
		h = 14
	}

	newFish := Fish{
		DirX:     dirX,
		X:        rand.Intn(w),
		Y:        rand.Intn(int(math.Abs(float64(h+5)))) + 2,
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
