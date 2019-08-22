package main

import (
	"CHIP-8/chip8"
	"github.com/nsf/termbox-go"
	"os"
)

var c *chip8.CHIP8

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	c = chip8.New()
	c.Logging = false
	err = c.Load("../../static/roms/TICTAC")
	if err != nil {
		panic(err)
	}

	eventQueue := make(chan termbox.Event)

	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventQueue:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					os.Exit(0)
				}

				if k, ok := chip8.KeyMap[string(ev.Ch)]; ok {
					c.SendKeyPress(k)
				}
			}
		case <-c.Timer.C:
			c.Cycle()
			if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
				panic(err)
			}
			if c.DrawFlag {
				for x := 0; x < chip8.Width; x++ {
					for y := 0; y < chip8.Height; y++ {
						if c.Gfx[x][y] == 1 {
							termbox.SetCell(x, y, 'x', termbox.ColorWhite, termbox.ColorWhite)
						}
					}
				}
			}

			if err := termbox.Flush(); err != nil {
				panic(err)
			}
		}
	}
}
