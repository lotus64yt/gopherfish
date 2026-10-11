package types

type TerminalSize struct {
	Width  int
	Height int
}

type TerminalCursorPos struct {
	X           int
	Y           int
	Clicked     bool
	JustClicked bool
}

type TerminalInfo struct {
	Size      TerminalSize
	CursorPos TerminalCursorPos
}
