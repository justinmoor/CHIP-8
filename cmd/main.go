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

	if chip8.DrawFlag {
		for x := 0; x < system.Width; x++ {
			for y := 0; y < system.Height; y++ {
				if chip8.Gfx[y][x] == 1 {
					screen.Set(x, y, color.White)
				}
			}
		}
	}
	return nil
}

func main() {
	chip8 = new(system.CHIP8)

	if err := chip8.Run("roms/PONG1"); err != nil {
		log.Fatal(err)
	}

	if err := e.Run(run, system.Width, system.Height, scale, "CHIP-8"); err != nil {
		log.Fatal(err)
	}
}

func checkKeyPress() {
	if e.IsKeyPressed(e.Key1) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x1})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x1})
	}

	if e.IsKeyPressed(e.Key2) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x2})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x2})
	}

	if e.IsKeyPressed(e.Key3) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x3})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x3})
	}

	if e.IsKeyPressed(e.Key4) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x5})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x5})
	}

	if e.IsKeyPressed(e.KeyQ) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x4})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x4})
	}

	if e.IsKeyPressed(e.KeyW) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x5})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x5})
	}

	if e.IsKeyPressed(e.KeyE) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x6})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x6})
	}

	if e.IsKeyPressed(e.KeyR) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0xD})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0xD})
	}

	if e.IsKeyPressed(e.KeyA) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x7})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x7})
	}

	if e.IsKeyPressed(e.KeyS) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x8})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x8})
	}

	if e.IsKeyPressed(e.KeyD) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x9})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x9})
	}

	if e.IsKeyPressed(e.KeyF) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0xE})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0xE})
	}

	if e.IsKeyPressed(e.KeyZ) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0xA})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0xA})
	}

	if e.IsKeyPressed(e.KeyX) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0x0})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0x0})
	}

	if e.IsKeyPressed(e.KeyC) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0xB})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0xB})
	}

	if e.IsKeyPressed(e.KeyV) {
		chip8.SendKeyState(system.Key{Pressed: true, Hex: 0xF})
	} else {
		chip8.SendKeyState(system.Key{Pressed: false, Hex: 0xF})
	}
}
