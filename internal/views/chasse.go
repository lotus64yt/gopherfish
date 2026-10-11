package views

import (
	"gopherfish/internal/sprites"
	"gopherfish/internal/types"
	"gopherfish/internal/utils"
	"strconv"
	"time"
)

type ChasseView struct {
	Game         *types.GameState
	Fishs        []sprites.Fish
	lastFishTick time.Time
	fishTime     time.Duration
	wasClicked   bool
}

func (m *ChasseView) Draw() []string {
	w := m.Game.Terminal.Size.Width
	h := m.Game.Terminal.Size.Height

	lines := utils.CreateBlankScreen(w, h)

	if len(m.Fishs) == 0 {
		for range 10 + m.Game.Level*5 {
			m.Fishs = sprites.SpawnRandomShip(*m.Game, m.Fishs)
		}
	}

	cursor := m.Game.Terminal.CursorPos
	justClicked := cursor.JustClicked || (cursor.Clicked && !m.wasClicked)
	m.wasClicked = cursor.Clicked

	if justClicked {
		for i := len(m.Fishs) - 1; i >= 0; i-- {
			if m.Fishs[i].IsClicked(cursor) {
				m.Game.Score++
				m.Fishs = append(m.Fishs[:i], m.Fishs[i+1:]...)
				if len(m.Fishs) == 0 {
					m.Game.Level++
				}
				break
			}
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

	lines = utils.DrawText(lines, m.Game.Terminal.Size.Width/10, m.Game.Terminal.Size.Height-3, "Score: "+strconv.Itoa(m.Game.Score)+" Level: "+strconv.Itoa(m.Game.Level))
	return lines
}
