package types

type GameViewPage int

const (
	GameViewMenu    GameViewPage = 0
	GameViewChasing GameViewPage = 1
	GameViewGallery GameViewPage = 2
)

type GameState struct {
	CurrentView GameViewPage
	Terminal    TerminalInfo
	Score       int
}
