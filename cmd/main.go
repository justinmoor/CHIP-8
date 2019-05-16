package main

import (
	"CHIP-8/system"
	e "github.com/hajimehoshi/ebiten"
	"log"
)

const (
	width  = 64
	height = 32
	scale  = 16
)

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
	err := chip8.Load("roms/pong.ch8")
	if err != nil {
		log.Fatal(err)
	}

	err = e.Run(run, width, height, scale, "CHIP-8")

	if err != nil {
		log.Fatal(err)
	}
}
