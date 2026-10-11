package views

import (
	"gopherfish/internal/sprites"
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
	"time"
)

type ChasseView struct {
	Game         *types.GameState
	Fishs        []sprites.Fish
	lastFishTick time.Time
	fishTime     time.Duration
}

func (m *ChasseView) Draw() []string {
	w := m.Game.Terminal.Size.Width
	h := m.Game.Terminal.Size.Height

	lines := utils.CreateBlankScreen(w, h)

	if len(m.Fishs) == 0 {
		for range 10 {
			m.Fishs = sprites.SpawnRandomShip(*m.Game, m.Fishs)
		}
	}

	if m.lastFishTick.IsZero() {
		m.lastFishTick = time.Now()
	} else {
		m.fishTime += time.Since(m.lastFishTick)
		m.lastFishTick = time.Now()
	}

	const fishInterval = time.Second / 10
	for m.fishTime >= fishInterval {
		for i := range m.Fishs {
			m.Fishs[i].Tick(*m.Game)
		}
		m.fishTime -= fishInterval
	}

	for i := range m.Fishs {
		lines = m.Fishs[i].Draw(lines)
	}
	return lines
}
