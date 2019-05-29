package main

import (
	"CHIP-8/system"
	e "github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const scale = 16

var chip8 *system.CHIP8

func run(screen *e.Image) error {
	checkKeyPress()
	gfx := <-chip8.ScreenState

	for x := 0; x < system.Width; x++ {
		for y := 0; y < system.Height; y++ {
			if gfx[y][x] == 1 {
				screen.Set(x, y, color.White)
			}
		}
	}

	return nil
}

func checkKeyPress() {
	if e.IsKeyPressed(e.Key1) {
		chip8.SendKeyPress(0x1) // 1
	}
	if e.IsKeyPressed(e.Key2) {
		chip8.SendKeyPress(0x2) // 2
	}
	if e.IsKeyPressed(e.Key3) {
		chip8.SendKeyPress(0x3) // 3
	}
	if e.IsKeyPressed(e.Key4) {
		chip8.SendKeyPress(0xC) // 12
	}
}

func main() {
	chip8 = new(system.CHIP8)

	if err := chip8.Run("roms/pong.ch8"); err != nil {
		log.Fatal(err)
	}

	if err := e.Run(run, system.Width, system.Height, scale, "CHIP-8"); err != nil {
		log.Fatal(err)
	}
}
