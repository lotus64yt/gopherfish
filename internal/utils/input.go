package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/charmbracelet/x/term"
)

var mouseSeqRegex = regexp.MustCompile(`\x1b\[<(\d+);(\d+);(\d+)([Mm])`)

func StartInputLoop(onMouse func(x, y int, pressed bool)) (func(), error) {
	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		return nil, err
	}

	fmt.Print("\033[?1003h\033[?1006h")

	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				return
			}
			seq := string(buf[:n])

			matches := mouseSeqRegex.FindAllStringSubmatch(seq, -1)
			if len(matches) > 0 {
				for _, match := range matches {
					b, _ := strconv.Atoi(match[1])
					x, _ := strconv.Atoi(match[2])
					y, _ := strconv.Atoi(match[3])
					isRelease := match[4] == "m"

					isPressed := false
					if !isRelease && (b&3) != 3 {
						isPressed = true
					}

					onMouse(x-1, y-1, isPressed)
				}
			} else if seq == "\x03" || seq == "\x1b" {
				fmt.Print("\033[?1006l\033[?1003l")
				term.Restore(os.Stdin.Fd(), oldState)
				fmt.Print("\033[2J\033[H")
				os.Exit(0)
			}
		}
	}()

	cleanup := func() {
		fmt.Print("\033[?1006l\033[?1003l")
		term.Restore(os.Stdin.Fd(), oldState)
	}

	return cleanup, nil
}
