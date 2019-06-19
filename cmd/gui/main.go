package main

import (
	"CHIP-8/system"
	e "github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
	"time"
)

const scale = 16

var chip8 *system.CHIP8

func run(screen *e.Image) error {
	getKeyState()

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
	chip8 = system.New()

	if err := chip8.Load("roms/TICTAC"); err != nil {
		log.Fatal("Could not load ROM")
	}

	go func() {
		for range time.Tick(16 * time.Nanosecond) {
			chip8.Cycle()
		}
	}()

	if err := e.Run(run, system.Width, system.Height, scale, "CHIP-8"); err != nil {
		log.Fatal(err)
	}
}

func getKeyState() {
	chip8.ResetKeys()

	if e.IsKeyPressed(e.Key1) {
		chip8.SendKeyPress(system.KeyMap[e.Key1.String()])
	}
	if e.IsKeyPressed(e.Key2) {
		chip8.SendKeyPress(system.KeyMap[e.Key2.String()])
	}
	if e.IsKeyPressed(e.Key3) {
		chip8.SendKeyPress(system.KeyMap[e.Key3.String()])
	}
	if e.IsKeyPressed(e.Key4) {
		chip8.SendKeyPress(system.KeyMap[e.Key4.String()])
	}
	if e.IsKeyPressed(e.KeyQ) {
		chip8.SendKeyPress(system.KeyMap[e.KeyQ.String()])
	}
	if e.IsKeyPressed(e.KeyW) {
		chip8.SendKeyPress(system.KeyMap[e.KeyW.String()])
	}
	if e.IsKeyPressed(e.KeyE) {
		chip8.SendKeyPress(system.KeyMap[e.KeyE.String()])
	}
	if e.IsKeyPressed(e.KeyR) {
		chip8.SendKeyPress(system.KeyMap[e.KeyR.String()])
	}
	if e.IsKeyPressed(e.KeyA) {
		chip8.SendKeyPress(system.KeyMap[e.KeyA.String()])
	}
	if e.IsKeyPressed(e.KeyS) {
		chip8.SendKeyPress(system.KeyMap[e.KeyS.String()])
	}
	if e.IsKeyPressed(e.KeyD) {
		chip8.SendKeyPress(system.KeyMap[e.KeyD.String()])
	}
	if e.IsKeyPressed(e.KeyF) {
		chip8.SendKeyPress(system.KeyMap[e.KeyF.String()])
	}
	if e.IsKeyPressed(e.KeyZ) {
		chip8.SendKeyPress(system.KeyMap[e.KeyZ.String()])
	}
	if e.IsKeyPressed(e.KeyX) {
		chip8.SendKeyPress(system.KeyMap[e.KeyX.String()])
	}
	if e.IsKeyPressed(e.KeyC) {
		chip8.SendKeyPress(system.KeyMap[e.KeyC.String()])
	}
	if e.IsKeyPressed(e.KeyV) {
		chip8.SendKeyPress(system.KeyMap[e.KeyV.String()])
	}
}
