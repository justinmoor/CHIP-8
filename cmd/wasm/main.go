package main

import (
	"CHIP-8/chip8"
	"syscall/js"
)

var (
	window, doc, body, canvas, ctx, beep js.Value
	width                                float64 = 64
	height                               float64 = 32
	c                                    *chip8.CHIP8
)

const scale = 16

func setup() {
	window = js.Global()
	doc = window.Get("document")
	body = doc.Get("body")

	height = window.Get("innerHeight").Float()
	width = window.Get("innerWidth").Float()

	canvas = doc.Call("createElement", "canvas")
	canvas.Set("height", height)
	canvas.Set("width", width)
	body.Call("appendChild", canvas)

	ctx = canvas.Call("getContext", "2d")
	ctx.Set("fillStyle", "black")
	ctx.Call("scale", scale, scale)

	doc.Call("addEventListener", "keydown", keyEvent)
}

var keyEvent = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	e := args[0]
	keyPress := e.Get("key")

	if key, ok := chip8.KeyMap[keyPress.String()]; ok {
		c.SendKeyPress(key)
	}

	return nil
})

func main() {
	setup()
	quit := make(chan struct{})
	c = chip8.New()
	c.Logging = true

	gfxBuffer := make(chan [chip8.Width][chip8.Height]byte, 10)
	var gfx [chip8.Width][chip8.Height]byte

	// load a ROM from the static files
	if err := c.LoadRomHTTP("http://localhost:8000/roms/PONG1"); err != nil {
		panic(err)
	}

	// start the emulator
	go func() {
		for range c.Timer.C {
			c.Cycle()
			if c.DrawFlag {
				select {
				case gfxBuffer <- c.Gfx:
				default:
				}
				c.DrawFlag = false
			}
		}
	}()

	// web renderer
	var renderer js.Func
	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("clearRect", 0, 0, width, height)

		select {
		case gfx = <-gfxBuffer:
		default:
		}

		for x := 0; x < chip8.Width; x++ {
			for y := 0; y < chip8.Height; y++ {
				if gfx[x][y] == 1 {
					ctx.Call("fillRect", x, y, 1, 1)
				}
			}
		}

		window.Call("requestAnimationFrame", renderer)
		return nil
	})

	window.Call("requestAnimationFrame", renderer)

	<-quit
}
