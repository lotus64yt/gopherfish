package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/x/term"
)

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

			if strings.HasPrefix(seq, "\033[<") {
				isRelease := strings.HasSuffix(seq, "m")
				isPressOrMove := strings.HasSuffix(seq, "M")

				if isRelease || isPressOrMove {
					trimmed := seq[3 : len(seq)-1]
					parts := strings.Split(trimmed, ";")
					if len(parts) == 3 {
						b, _ := strconv.Atoi(parts[0])
						x, _ := strconv.Atoi(parts[1])
						y, _ := strconv.Atoi(parts[2])

						isPressed := false
						if !isRelease && (b&3) != 3 {
							isPressed = true
						}

						onMouse(x-1, y-1, isPressed)
					}
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
