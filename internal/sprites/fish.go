package sprites

import (
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
	"math"
	"math/rand"
)

type Fish struct {
	X        int
	Y        int
	DirX     int
	DirY     int
	DirYStep int
	DirYTick int
}

func SpawnRandomShip(g types.GameState, fishs []Fish) []Fish {
	dirX := -1
	if rand.Intn(2) == 1 {
		dirX = 1
	}

	newFish := Fish{
		DirX:     dirX,
		X:        rand.Intn(g.Terminal.Size.Width*2) / 2,
		Y:        rand.Intn(int(math.Abs(float64(g.Terminal.Size.Height-10)))*2) / 2,
		DirY:     rand.Intn(3) - 1,
		DirYStep: rand.Intn(2)*2 - 1,
		DirYTick: rand.Intn(5),
	}

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
	}

	if utils.TerminalHitScreenBorder(g.Terminal.Size, types.TerminalCursorPos{
		X: f.X,
		Y: f.Y,
	}) {
		f.DirX = -f.DirX
	}
}

func (f *Fish) Draw(screen []string) []string {
	fishStr := "><@"
	if f.DirX == -1 {
		runes := []rune(fishStr)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		for i, r := range runes {
			switch r {
			case '>':
				runes[i] = '<'
			case '<':
				runes[i] = '>'
			}
		}
		fishStr = string(runes)
	}

	if f.Y < 0 || f.Y >= len(screen) {
		return screen
	}

	line := []rune(screen[f.Y])
	fish := []rune(fishStr)
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
