package main

import (
	"CHIP-8/system"
	"fmt"
	e "github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const scale = 16

var chip8 *system.CHIP8

func run(screen *e.Image) error {
	chip8.Cycle()
	var y int
	if chip8.DrawFlag {
		for x := 0; x < len(chip8.Gfx); x++ {
			if x%64 == 0 {
				y += 1
			}
			if chip8.Gfx[x] == 1 {
				screen.Set(x, y, color.White)
			}
		}
		fmt.Println(9)
	}

	return nil
}

func main() {
	chip8 = new(system.CHIP8)
	chip8.Initialize()
	err := chip8.Load("roms/pong.ch8")
	if err != nil {
		log.Fatal(err)
	}

	err = e.Run(run, system.Width, system.Height, scale, "CHIP-8")

	if err != nil {
		log.Fatal(err)
	}
}
