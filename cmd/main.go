package main

import (
	"CHIP8/system"
	e "github.com/hajimehoshi/ebiten"
	"log"
)

const scale = 16

var chip8 *system.CHIP8

func run(screen *e.Image) error {
	chip8.Cycle()

	if chip8.DrawFlag {
		//screen.renderstuff()
	}

	return nil
}

func main() {
	chip8 = new(system.CHIP8)
	chip8.Initialize()
	err := chip8.Load("roms/BC_test.ch8")
	if err != nil {
		log.Fatal(err)
	}

	err = e.Run(run, system.Width, system.Height, scale, "CHIP-8")

	if err != nil {
		log.Fatal(err)
	}
}
