package main

import (
	"CHIP-8/chip8"
	"fmt"
	"github.com/nsf/termbox-go"
	"log"
	"os"
)

var c *chip8.CHIP8

func main() {
	err := termbox.Init()

	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	rom := os.Args[1]
	c = chip8.New()
	c.Logging = false
	err = c.Load(fmt.Sprintf("../../static/roms/%v", rom))
	if err != nil {
		log.Fatalf("Could not find rom %v", rom)
	}

	keyEvents := make(chan termbox.Event)

	go func() {
		for {
			keyEvents <- termbox.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-keyEvents:
			if ev.Type == termbox.EventKey {
				if ev.Key == termbox.KeyEsc {
					os.Exit(0)
				}

				if k, ok := chip8.KeyMap[string(ev.Ch)]; ok {
					c.SendKeyPress(k)
				}
			}
		case <-c.Timer.C:
			if err := c.Cycle(); err != nil {
				log.Fatal(err)
			}
			if err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault); err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
		}
	}
}
