package main

import (
	"CHIP-8/chip8"
	e "github.com/hajimehoshi/ebiten"
	"image/color"
	"log"
)

const scale = 16

var c *chip8.CHIP8

func run(screen *e.Image) error {
	c.Cycle()
	getKeyState()
	if c.DrawFlag {
		for x := 0; x < chip8.Width; x++ {
			for y := 0; y < chip8.Height; y++ {
				if c.Gfx[x][y] == 1 {
					screen.Set(x, y, color.White)
				}
			}
		}

		c.DrawFlag = false
	}

	return nil
}

func main() {
	c = chip8.New()

	if err := c.Load("roms/PONG"); err != nil {
		log.Fatal("Could not load ROM")
	}

	if err := e.Run(run, chip8.Width, chip8.Height, scale, "CHIP-8"); err != nil {
		log.Fatal(err)
	}
}

func getKeyState() {
	if e.IsKeyPressed(e.Key1) {
		c.SendKeyPress(chip8.KeyMap[e.Key1.String()])
	}
	if e.IsKeyPressed(e.Key2) {
		c.SendKeyPress(chip8.KeyMap[e.Key2.String()])
	}
	if e.IsKeyPressed(e.Key3) {
		c.SendKeyPress(chip8.KeyMap[e.Key3.String()])
	}
	if e.IsKeyPressed(e.Key4) {
		c.SendKeyPress(chip8.KeyMap[e.Key4.String()])
	}
	if e.IsKeyPressed(e.KeyQ) {
		c.SendKeyPress(chip8.KeyMap[e.KeyQ.String()])
	}
	if e.IsKeyPressed(e.KeyW) {
		c.SendKeyPress(chip8.KeyMap[e.KeyW.String()])
	}
	if e.IsKeyPressed(e.KeyE) {
		c.SendKeyPress(chip8.KeyMap[e.KeyE.String()])
	}
	if e.IsKeyPressed(e.KeyR) {
		c.SendKeyPress(chip8.KeyMap[e.KeyR.String()])
	}
	if e.IsKeyPressed(e.KeyA) {
		c.SendKeyPress(chip8.KeyMap[e.KeyA.String()])
	}
	if e.IsKeyPressed(e.KeyS) {
		c.SendKeyPress(chip8.KeyMap[e.KeyS.String()])
	}
	if e.IsKeyPressed(e.KeyD) {
		c.SendKeyPress(chip8.KeyMap[e.KeyD.String()])
	}
	if e.IsKeyPressed(e.KeyF) {
		c.SendKeyPress(chip8.KeyMap[e.KeyF.String()])
	}
	if e.IsKeyPressed(e.KeyZ) {
		c.SendKeyPress(chip8.KeyMap[e.KeyZ.String()])
	}
	if e.IsKeyPressed(e.KeyX) {
		c.SendKeyPress(chip8.KeyMap[e.KeyX.String()])
	}
	if e.IsKeyPressed(e.KeyC) {
		c.SendKeyPress(chip8.KeyMap[e.KeyC.String()])
	}
	if e.IsKeyPressed(e.KeyV) {
		c.SendKeyPress(chip8.KeyMap[e.KeyV.String()])
	}
}
