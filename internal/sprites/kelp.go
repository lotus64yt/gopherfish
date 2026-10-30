package sprites

import (
	"gopherfish/internal/types"
	"math"
	"math/rand"
	"strings"
)

type Kelp struct {
	X          int
	Height     int
	BaseHeight int
	HTick      int
	WiggleTick int
	startY     int
	totalLines int
	actualX    int
}

func (k *Kelp) GetAscii() string {
	lines := make([]string, 0, k.Height+2)
	wiggleState := (k.WiggleTick / 4) % 6

	switch wiggleState {
	case 2:
		lines = append(lines, "    ^  ", "  /   \\")
	case 1:
		lines = append(lines, "   ^   ", "  /   \\")
	case 4:
		lines = append(lines, "   ^   ", "/   \\  ")
	case 5:
		lines = append(lines, "  ^    ", "/   \\  ")
	default:
		lines = append(lines, "   ^   ", " /   \\ ")
	}

	for i := 0; i < k.Height; i++ {
		isBottom := (i == k.Height-1)
		isNearBottom := (i == k.Height-2)

		var seg string
		switch wiggleState {
		case 2:
			if isBottom {
				seg = " |   | "
			} else if isNearBottom {
				seg = " /   / "
			} else {
				seg = "  /   /"
			}
		case 1:
			if isBottom {
				seg = " |   | "
			} else {
				seg = " /   / "
			}
		case 4:
			if isBottom {
				seg = " |   | "
			} else {
				seg = " \\   \\ "
			}
		case 5:
			if isBottom {
				seg = " |   | "
			} else if isNearBottom {
				seg = " \\   \\ "
			} else {
				seg = "\\   \\  "
			}
		default:
			seg = " |   | "
		}
		lines = append(lines, seg)
	}

	return strings.Join(lines, "\n")
}

func SpawnRandomKelp(g types.GameState, kelps []Kelp) []Kelp {
	w := g.Terminal.Size.Width
	if w <= 10 {
		w = 80
	}
	h := g.Terminal.Size.Height
	if h <= 10 {
		h = 24
	}

	minHeight := 10
	maxHeight := h / 2
	if maxHeight < minHeight {
		maxHeight = minHeight + 1
	}
	kelpHeight := rand.Intn(maxHeight-minHeight+1) + minHeight

	minX := 6
	maxX := w - 6
	if maxX <= minX {
		maxX = minX + 1
	}

	bestX := rand.Intn(maxX-minX) + minX
	minDistKelp := 15
	for attempt := 0; attempt < 30; attempt++ {
		candX := rand.Intn(maxX-minX) + minX
		tooClose := false
		for _, other := range kelps {
			dx := int(math.Abs(float64(candX - other.X)))
			if dx < minDistKelp {
				tooClose = true
				break
			}
		}
		if !tooClose {
			bestX = candX
			break
		}
		if attempt == 15 {
			minDistKelp = 8
		}
	}

	totalLines := kelpHeight + 2
	startY := h - totalLines

	newKelp := Kelp{
		X:          bestX,
		Height:     kelpHeight,
		BaseHeight: kelpHeight,
		actualX:    bestX,
		startY:     startY,
		totalLines: totalLines,
		HTick:      rand.Intn(40),
		WiggleTick: rand.Intn(24),
	}

	return append(kelps, newKelp)
}

func (k *Kelp) Tick(g types.GameState) {
	if k.BaseHeight <= 0 {
		k.BaseHeight = k.Height
	}
	if k.BaseHeight <= 0 {
		k.BaseHeight = 5
	}

	k.HTick++
	k.WiggleTick++

	offsetCycle := (k.HTick / 20) % 4
	var offset int
	switch offsetCycle {
	case 1:
		offset = 1
	case 3:
		offset = -1
	default:
		offset = 0
	}

	newHeight := k.BaseHeight + offset
	if newHeight < 2 {
		newHeight = 2
	}
	k.Height = newHeight
}

func (k *Kelp) Draw(screen []string) []string {
	if len(screen) == 0 {
		return screen
	}

	lines := strings.Split(strings.TrimPrefix(k.GetAscii(), "\n"), "\n")
	k.totalLines = len(lines)
	k.startY = len(screen) - k.totalLines

	for i, spriteLine := range lines {
		y := k.startY + i
		if y < 0 || y >= len(screen) {
			continue
		}

		line := []rune(screen[y])
		kelpLine := []rune(spriteLine)
		startX := k.X - len(kelpLine)/2

		if startX+len(kelpLine) > len(line) {
			startX = len(line) - len(kelpLine)
		}
		if startX < 0 {
			startX = 0
		}
		k.actualX = startX + len(kelpLine)/2

		firstNonSpace := -1
		lastNonSpace := -1
		for j, r := range kelpLine {
			if r != ' ' {
				if firstNonSpace == -1 {
					firstNonSpace = j
				}
				lastNonSpace = j
			}
		}

		if firstNonSpace == -1 {
			continue
		}

		for j := firstNonSpace; j <= lastNonSpace; j++ {
			pos := startX + j
			if pos >= 0 && pos < len(line) {
				line[pos] = kelpLine[j]
			}
		}
		screen[y] = string(line)
	}

	return screen
}

func (k *Kelp) IsClicked(cPos types.TerminalCursorPos) bool {
	if cPos.Y >= k.startY && cPos.Y < k.startY+k.totalLines {
		if cPos.X >= k.actualX-3 && cPos.X <= k.actualX+3 {
			return true
		}
	}
	return false
}

func (k *Kelp) CoversFish(f Fish) bool {
	if f.Y < k.startY || f.Y >= k.startY+k.totalLines {
		return false
	}

	if f.X+1 >= k.actualX-2 && f.X-1 <= k.actualX+2 {
		return true
	}
	return false
}
