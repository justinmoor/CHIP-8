package main

import (
	"CHIP-8/chip8"
	"syscall/js"
	"time"
)

var (
	window, doc, body, canvas, ctx, beep js.Value
	width                                float64 = 64
	height                               float64 = 32
	c                                    *chip8.CHIP8
	then                                 time.Time
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
	c = chip8.New()

	go func() {
		for range time.After(16 * time.Microsecond) {
			c.Cycle()
		}
	}()

	if err := c.LoadRomHTTP("http://localhost:8000/roms/PONG1"); err != nil {
		panic(err)
	}

	var renderer js.Func
	renderer = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ctx.Call("clearRect", 0, 0, width, height)

		if c.DrawFlag {
			for x := 0; x < chip8.Width; x++ {
				for y := 0; y < chip8.Height; y++ {
					if c.Gfx[x][y] == 1 {
						ctx.Call("fillRect", x, y, 1, 1)
					}
				}
			}
		}

		window.Call("requestAnimationFrame", renderer)
		return nil
	})

	window.Call("requestAnimationFrame", renderer)

	select {}
}
