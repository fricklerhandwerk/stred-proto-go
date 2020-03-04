package ui

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

func Draw() {
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer func() { screen.Fini() }()

	style := tcell.StyleDefault.Foreground(tcell.ColorRed)
	writeLine(screen, "hello world - press any key to exit", 0, 0, style)
	screen.Show()

	quit := make(chan struct{})
	go func() {
		for {
			event := screen.PollEvent()
			switch event.(type) {
			case *tcell.EventKey:
				close(quit)
			case *tcell.EventResize:
				screen.Sync()
			}
		}
	}()
	for {
		select {
		case <-quit:
			return
		}
	}
}

func writeLine(screen tcell.Screen, text string, x, y int, style tcell.Style) {
	for i := 0; i < len(text); i++ {
		screen.SetContent(x+i, y, rune(text[i]), nil, style)
	}
}
